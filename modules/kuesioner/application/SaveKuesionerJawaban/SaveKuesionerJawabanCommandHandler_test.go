package application

import (
	"context"
	"errors"
	"testing"

	helper "UnpakSiamida/common/helper"
	mockkuesioner "UnpakSiamida/modules/kuesioner/application/mock"
	domainkuesioner "UnpakSiamida/modules/kuesioner/domain"
	mockjawaban "UnpakSiamida/modules/templatejawaban/application/mock"
	domainjawaban "UnpakSiamida/modules/templatejawaban/domain"
	mockpertanyaan "UnpakSiamida/modules/templatepertanyaan/application/mock"
	domainpertanyaan "UnpakSiamida/modules/templatepertanyaan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSaveKuesionerJawabanCommandHandler(t *testing.T) {
	kuesionerUUID := uuid.New()
	pertanyaanUUID := uuid.New()
	jawabanUUID1 := uuid.New()


	baseCmd := SaveKuesionerJawabanCommand{
		UuidKuesioner:  kuesionerUUID.String(),
		UuidPertanyaan: pertanyaanUUID.String(),
		Jawaban:        `[{"uuid":"` + jawabanUUID1.String() + `","freetext":""}]`,
		SID:            "user-123",
		Resource:       "simak",
	}

	t.Run("Fail BeginTx Error", func(t *testing.T) {
		beginErr := errors.New("begin tx error")
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			BeginTxFunc: func(ctx context.Context) (*gorm.DB, error) {
				return nil, beginErr
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, beginErr)
	})

	t.Run("Fail Invalid Kuesioner UUID", func(t *testing.T) {
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}
		handler := &SaveKuesionerJawabanCommandHandler{
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		cmd := baseCmd
		cmd.UuidKuesioner = "invalid-uuid"

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.InvalidUuid(), err)
	})

	t.Run("Fail Invalid Pertanyaan UUID", func(t *testing.T) {
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}
		handler := &SaveKuesionerJawabanCommandHandler{
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		cmd := baseCmd
		cmd.UuidPertanyaan = "invalid-uuid"

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.InvalidPertanyaan(), err)
	})

	t.Run("Fail Get Kuesioner By Uuid Error", func(t *testing.T) {
		getErr := errors.New("get kuesioner error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return nil, getErr
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Fail Get Pertanyaan By Uuid Error", func(t *testing.T) {
		getErr := errors.New("get pertanyaan error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return nil, getErr
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Fail Invalid Jawaban JSON", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		cmd := baseCmd
		cmd.Jawaban = "invalid-json"

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, domainkuesioner.InvalidJawaban(), err)
	})

	t.Run("Fail Invalid Payload UUID", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		cmd := baseCmd
		cmd.Jawaban = `[{"uuid":"invalid-uuid"}]`

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, "uuid jawaban tidak valid", err.Error())
	})

	t.Run("Fail GetByUUIDs Error", func(t *testing.T) {
		getErr := errors.New("get by uuids error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				return nil, getErr
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Fail GetFreeTextByPertanyaan Error", func(t *testing.T) {
		getErr := errors.New("get freetext error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				return []domainjawaban.TemplateJawaban{
					{ID: 10, UUID: jawabanUUID1},
				}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				return nil, getErr
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Fail GetByPertanyaanAndUser Error", func(t *testing.T) {
		getErr := errors.New("get existing error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				return []domainjawaban.TemplateJawaban{
					{ID: 10, UUID: jawabanUUID1},
				}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				return nil, nil
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			GetByPertanyaanAndUserFunc: func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domainkuesioner.KuesionerJawaban, error) {
				return nil, getErr
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Fail Delete Existing Error", func(t *testing.T) {
		deleteErr := errors.New("delete existing error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				// BaseCmd requests jawabanUUID1 (ID: 10).
				// We return only this, so jawabanUUID2 (ID: 20) is NOT selected.
				return []domainjawaban.TemplateJawaban{
					{ID: 10, UUID: jawabanUUID1},
				}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				return nil, nil
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			GetByPertanyaanAndUserFunc: func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domainkuesioner.KuesionerJawaban, error) {
				// Existing is jawabanUUID2 (ID: 20). Since it's not selected, handler will try to DELETE it.
				optId := uint(20)
				return []domainkuesioner.KuesionerJawaban{
					{
						ID:                99,
						IdTemplateJawaban: &optId,
					},
				}, nil
			},
			DeleteFunc: func(ctx context.Context, id uint) error {
				assert.Equal(t, uint(99), id)
				return deleteErr
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, deleteErr)
	})

	t.Run("Fail Create New Answer Error", func(t *testing.T) {
		createErr := errors.New("create answer error")
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				return []domainjawaban.TemplateJawaban{
					{ID: 10, UUID: jawabanUUID1},
				}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				return nil, nil
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			GetByPertanyaanAndUserFunc: func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domainkuesioner.KuesionerJawaban, error) {
				return nil, nil // No existing
			},
			CreateFunc: func(ctx context.Context, data *domainkuesioner.KuesionerJawaban) error {
				return createErr
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		_, err := handler.Handle(context.Background(), baseCmd)
		assert.ErrorIs(t, err, createErr)
	})

	t.Run("Fail FreeText Unsupported", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				return []domainjawaban.TemplateJawaban{}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				return nil, nil // Unsupported
			},
		}
		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			GetByPertanyaanAndUserFunc: func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domainkuesioner.KuesionerJawaban, error) {
				return nil, nil
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		cmd := SaveKuesionerJawabanCommand{
			UuidKuesioner:  kuesionerUUID.String(),
			UuidPertanyaan: pertanyaanUUID.String(),
			Jawaban:        `[{"uuid":"","freetext":"Some free text answer"}]`,
			SID:            "user-123",
			Resource:       "simak",
		}

		_, err := handler.Handle(context.Background(), cmd)
		assert.Error(t, err)
		assert.Equal(t, "pertanyaan ini tidak mendukung free text", err.Error())
	})

	t.Run("Success Complex Scenario (Insert, Delete, Update FreeText)", func(t *testing.T) {
		repo := &mockkuesioner.MockKuesionerRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainkuesioner.Kuesioner, error) {
				return &domainkuesioner.Kuesioner{ID: 1, UUID: kuesionerUUID}, nil
			},
		}
		repoPertanyaan := &mockpertanyaan.MockTemplatePertanyaanRepository{
			GetByUuidFunc: func(ctx context.Context, uid uuid.UUID) (*domainpertanyaan.TemplatePertanyaan, error) {
				return &domainpertanyaan.TemplatePertanyaan{ID: 2, UUID: pertanyaanUUID}, nil
			},
		}
		repoJawaban := &mockjawaban.MockTemplateJawabanRepository{
			GetByUUIDsFunc: func(ctx context.Context, uuids []string) ([]domainjawaban.TemplateJawaban, error) {
				// Input contains:
				// - Jawaban 1 (selected choice): ID: 10, UUID: jawabanUUID1
				return []domainjawaban.TemplateJawaban{
					{ID: 10, UUID: jawabanUUID1},
				}, nil
			},
			GetFreeTextByPertanyaanFunc: func(ctx context.Context, pertanyaanID uint) (*domainjawaban.TemplateJawaban, error) {
				// Supports free text via template ID: 30
				return &domainjawaban.TemplateJawaban{ID: 30}, nil
			},
		}

		deletedCount := 0
		createdCount := 0

		repoJawabanKuesioner := &mockkuesioner.MockKuesionerJawabanRepository{
			GetByPertanyaanAndUserFunc: func(ctx context.Context, pertanyaanID uint, sid string, resource string) ([]domainkuesioner.KuesionerJawaban, error) {
				optIdExisting1 := uint(20) // Selected choice that should be DELETED
				optIdExisting2 := uint(30) // Existing free text answer that should be UPDATED
				return []domainkuesioner.KuesionerJawaban{
					{
						ID:                99,
						IdTemplateJawaban: &optIdExisting1,
					},
					{
						ID:                100,
						IdTemplateJawaban: &optIdExisting2,
						FreeText:          helper.StrPtr("Old free text"),
					},
				}, nil
			},
			DeleteFunc: func(ctx context.Context, id uint) error {
				if id == 99 {
					deletedCount++
					return nil
				}
				return errors.New("unexpected delete")
			},
			CreateFunc: func(ctx context.Context, data *domainkuesioner.KuesionerJawaban) error {
				if *data.IdTemplateJawaban == 10 {
					// Insert new selected choice
					createdCount++
					return nil
				}
				if *data.IdTemplateJawaban == 30 {
					// Update existing free text or insert new free text
					assert.Equal(t, uint(100), data.ID)
					assert.Equal(t, "New free text", *data.FreeText)
					createdCount++
					return nil
				}
				return errors.New("unexpected create")
			},
		}

		handler := &SaveKuesionerJawabanCommandHandler{
			Repo:                 repo,
			RepoPertanyaan:       repoPertanyaan,
			RepoJawaban:          repoJawaban,
			RepoJawabanKuesioner: repoJawabanKuesioner,
		}

		// Input:
		// - choice: jawabanUUID1
		// - free text: "New free text"
		cmd := SaveKuesionerJawabanCommand{
			UuidKuesioner:  kuesionerUUID.String(),
			UuidPertanyaan: pertanyaanUUID.String(),
			Jawaban:        `[{"uuid":"` + jawabanUUID1.String() + `","freetext":""},{"uuid":"","freetext":"New free text"}]`,
			SID:            "user-123",
			Resource:       "simak",
		}

		res, err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, kuesionerUUID.String(), res)
		assert.Equal(t, 1, deletedCount)
		assert.Equal(t, 2, createdCount)
	})
}
