package domaintest

import (
	"UnpakSiamida/modules/templatejawaban/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateJawabanErrors(t *testing.T) {
	assert.Equal(t, "TemplateJawaban.EmptyData", domain.EmptyData().Code)
	assert.Equal(t, "TemplateJawaban.InvalidUuid", domain.InvalidUuid().Code)
	assert.Equal(t, "TemplateJawaban.InvalidTemplatePertanyaan", domain.InvalidTemplatePertanyaan().Code)
	assert.Equal(t, "TemplateJawaban.NotFoundTemplatePertanyaan", domain.NotFoundTemplatePertanyaan().Code)
	assert.Equal(t, "TemplateJawaban.InvalidData", domain.InvalidData().Code)
	assert.Equal(t, "TemplateJawaban.NotFound", domain.NotFound("123").Code)
	assert.Equal(t, "TemplateJawaban.InvalidOwner", domain.InvalidOwner().Code)
}
