package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	helper "UnpakSiamida/common/helper"
	domainkategori "UnpakSiamida/modules/kategori/domain"
	"time"

	"github.com/bytedance/sonic"
	"github.com/goforj/godump"
)

type UpdateKategoriOrderCommandHandler struct {
	Repo domainkategori.IKategoriRepository
}

func (h *UpdateKategoriOrderCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateKategoriOrderCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var payload []domainkategori.KategoriPayload

	if err := DecodeJSONArrayStrict(strings.NewReader(cmd.Payload), &payload); err != nil {
		return "", err
	}

	uuidList := make([]string, 0)
	uuidSet := make(map[string]struct{})

	for _, item := range payload {

		if _, exists := uuidSet[item.UUID]; !exists {
			uuidSet[item.UUID] = struct{}{}
			uuidList = append(uuidList, item.UUID)
		}

		if item.UUIDSub != nil && *item.UUIDSub != "" {
			if _, exists := uuidSet[*item.UUIDSub]; !exists {
				uuidSet[*item.UUIDSub] = struct{}{}
				uuidList = append(uuidList, *item.UUIDSub)
			}
		}
	}

	kategori, _, err := h.Repo.GetAll(
		ctx,
		"",
		nil,
		nil,
		nil,
		false,
	)

	if err != nil {
		return "", err
	}

	uuidToID := make(map[string]uint)

	for _, item := range kategori {
		uuidToID[item.UUID.String()] = item.ID
	}

	updates := make([]domainkategori.UpdateRow, 0, len(payload))

	for _, item := range payload {

		id, ok := uuidToID[item.UUID]

		if !ok {
			continue
		}

		var subKategori *uint

		if item.UUIDSub != nil {

			if parentID, ok := uuidToID[*item.UUIDSub]; ok {
				subKategori = &parentID
			}
		}

		updates = append(updates, domainkategori.UpdateRow{
			ID:          id,
			SubKategori: subKategori,
			FullText:    item.FullText,
		})
	}

	godump.Dump(updates)
	err = h.Repo.UpdateParentBatch(ctx, updates)
	if err != nil {

	}

	return "ok", nil
}

func DecodeJSONArrayStrict(r io.Reader, dst *[]domainkategori.KategoriPayload) error {
	limited := io.LimitReader(r, helper.MaxPayloadSize+1)

	raw, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if len(raw) > helper.MaxPayloadSize {
		return helper.ErrPayloadTooLarge
	}

	if !utf8.Valid(raw) {
		return helper.ErrInvalidUTF8
	}

	trimmed := bytes.TrimSpace(raw)

	// ensure root array
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return helper.ErrInvalidRoot
	}

	// reject duplicate fields
	if err := helper.RejectDuplicateKeys(trimmed); err != nil {
		return err
	}

	cfg := sonic.Config{
		UseInt64:                true,
		NoValidateJSONMarshaler: true,
	}.Froze()

	dec := cfg.NewDecoder(bytes.NewReader(trimmed))
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	// reject trailing garbage
	var extra any

	if err := dec.Decode(&extra); err != io.EOF {
		return helper.ErrTrailingData
	}

	if len(*dst) > helper.MaxArrayItems {
		return helper.ErrTooManyItems
	}

	// validate every item
	for i := range *dst {
		if err := validateItem(&(*dst)[i]); err != nil {
			return fmt.Errorf("item[%d]: %w", i, err)
		}
	}

	return nil
}

func validateItem(v *domainkategori.KategoriPayload) error {
	err := helper.ValidateUUIDv4(v.UUID)
	if err != nil {
		return helper.ErrInvalidUUID
	}

	v.FullText = strings.TrimSpace(v.FullText)

	if v.FullText == "" {
		return helper.ErrEmptyName
	}

	// if len(v.FullText) > helper.MaxNameLength {
	// 	return helper.ErrNameTooLong
	// }

	if !utf8.ValidString(v.FullText) {
		return helper.ErrInvalidUTF8
	}

	return nil
}
