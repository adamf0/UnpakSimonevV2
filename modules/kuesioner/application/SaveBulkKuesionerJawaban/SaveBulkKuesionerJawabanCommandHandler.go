package SaveBulkKuesionerJawaban

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	helper "UnpakSiamida/common/helper"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	domainjawaban "UnpakSiamida/modules/templatejawaban/domain"
	domainpertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
)

type JawabanPayload struct {
	UUID     string `json:"uuid"`
	FreeText string `json:"freetext"`
}

type SaveBulkKuesionerJawabanCommandHandler struct {
	Repo                 domainkuesioner.IKuesionerRepository
	RepoPertanyaan       domainpertanyaan.ITemplatePertanyaanRepository
	RepoJawaban          domainjawaban.ITemplateJawabanRepository
	RepoJawabanKuesioner domainkuesioner.IKuesionerJawabanRepository
}

func (h *SaveBulkKuesionerJawabanCommandHandler) Handle(
	ctx context.Context,
	cmd SaveBulkKuesionerJawabanCommand,
) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if len(cmd.Items) == 0 {
		return "", errors.New("tidak ada jawaban yang dikirim")
	}

	// 1. Validasi UUID Kuesioner
	kuesionerUUID, err := uuid.Parse(cmd.UuidKuesioner)
	if err != nil {
		return "", domainkuesioner.InvalidUuid()
	}

	kuesioner, err := h.Repo.GetByUuid(ctx, kuesionerUUID)
	if err != nil {
		return "", err
	}

	// 2. Pre-process items & validate queries OUTSIDE transaction
	type ProcessedItem struct {
		Pertanyaan   *domainpertanyaan.TemplatePertanyaan
		JawabanList  []domainjawaban.TemplateJawaban
		FreeTemplate *domainjawaban.TemplateJawaban
		FreeTexts    []string
		FreeTextMap  map[string]string
		SelectedMap  map[uint]bool
	}

	processedItems := make([]ProcessedItem, 0, len(cmd.Items))

	for _, item := range cmd.Items {
		pertanyaanUUID, err := uuid.Parse(item.UuidPertanyaan)
		if err != nil {
			continue // skip invalid pertanyaan UUID
		}

		pertanyaan, err := h.RepoPertanyaan.GetByUuid(ctx, pertanyaanUUID)
		if err != nil {
			continue
		}

		var raw []JawabanPayload
		if err := json.Unmarshal([]byte(item.Jawaban), &raw); err != nil {
			continue
		}

		var selectedUUIDs []string
		var freeTexts []string
		freeTextMap := map[string]string{}

		for _, p := range raw {
			if p.UUID != "" {
				if _, err := uuid.Parse(p.UUID); err == nil {
					selectedUUIDs = append(selectedUUIDs, p.UUID)
				}
			}
			if p.FreeText != "" {
				freeTexts = append(freeTexts, p.FreeText)
				if p.UUID != "" {
					freeTextMap[p.UUID] = p.FreeText
				}
			}
		}

		jawabanList, err := h.RepoJawaban.GetByUUIDs(ctx, selectedUUIDs)
		if err != nil {
			continue
		}

		selectedMap := map[uint]bool{}
		for _, j := range jawabanList {
			selectedMap[j.ID] = true
		}

		freeTemplate, _ := h.RepoJawaban.GetFreeTextByPertanyaan(ctx, pertanyaan.ID)

		processedItems = append(processedItems, ProcessedItem{
			Pertanyaan:   pertanyaan,
			JawabanList:  jawabanList,
			FreeTemplate: freeTemplate,
			FreeTexts:    freeTexts,
			FreeTextMap:  freeTextMap,
			SelectedMap:  selectedMap,
		})
	}

	if len(processedItems) == 0 {
		return "", errors.New("tidak ada item jawaban valid untuk disimpan")
	}

	// 3. TRANSACTION - 1 Single DB Transaction for all bulk items
	tx, err := h.RepoJawabanKuesioner.BeginTx(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	commit := false
	defer func() {
		if !commit {
			_ = tx.Rollback()
		}
	}()

	repoWithTx := h.RepoJawabanKuesioner.WithTx(tx)

	for _, item := range processedItems {
		existing, err := repoWithTx.GetByPertanyaanAndUser(ctx, item.Pertanyaan.ID, cmd.SID, cmd.Resource)
		if err != nil {
			return "", err
		}

		existingMap := map[uint]domainkuesioner.KuesionerJawaban{}
		var existingFree []domainkuesioner.KuesionerJawaban

		for _, v := range existing {
			if v.IdTemplateJawaban != nil {
				if item.FreeTemplate != nil && *v.IdTemplateJawaban == item.FreeTemplate.ID {
					existingFree = append(existingFree, v)
				} else {
					existingMap[*v.IdTemplateJawaban] = v
				}
			}
		}

		// Delete unselected
		for id, data := range existingMap {
			if !item.SelectedMap[id] {
				if err := repoWithTx.Delete(ctx, data.ID); err != nil {
					return "", err
				}
			}
		}

		// Insert choice answers
		for _, j := range item.JawabanList {
			if _, exist := existingMap[j.ID]; !exist {
				var freeText *string
				if val, ok := item.FreeTextMap[j.UUID.String()]; ok {
					freeText = &val
				}

				newData := domainkuesioner.KuesionerJawaban{
					UUID:                 helper.StrPtr(uuid.New().String()),
					IdKuesioner:          kuesioner.ID,
					IdTemplatePertanyaan: item.Pertanyaan.ID,
					IdTemplateJawaban:    &j.ID,
					FreeText:             freeText,
					CreatedBy:            &cmd.Resource,
					CreatedByRef:         &cmd.SID,
				}

				if err := repoWithTx.Create(ctx, &newData); err != nil {
					return "", err
				}
			}
		}

		// Free text answers
		existingCount := len(existingFree)
		inputCount := len(item.FreeTexts)

		min := existingCount
		if inputCount < min {
			min = inputCount
		}

		for i := 0; i < min; i++ {
			text := item.FreeTexts[i]
			existingFree[i].FreeText = &text
			if err := repoWithTx.Create(ctx, &existingFree[i]); err != nil {
				return "", err
			}
		}

		if inputCount > existingCount && item.FreeTemplate != nil {
			for i := existingCount; i < inputCount; i++ {
				text := item.FreeTexts[i]
				newData := domainkuesioner.KuesionerJawaban{
					UUID:                 helper.StrPtr(uuid.New().String()),
					IdKuesioner:          kuesioner.ID,
					IdTemplatePertanyaan: item.Pertanyaan.ID,
					IdTemplateJawaban:    &item.FreeTemplate.ID,
					FreeText:             &text,
					CreatedBy:            &cmd.Resource,
					CreatedByRef:         &cmd.SID,
				}
				if err := repoWithTx.Create(ctx, &newData); err != nil {
					return "", err
				}
			}
		}

		if existingCount > inputCount {
			for i := inputCount; i < existingCount; i++ {
				if err := repoWithTx.Delete(ctx, existingFree[i].ID); err != nil {
					return "", err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	commit = true

	return kuesioner.UUID.String(), nil
}
