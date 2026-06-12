package commontest

import (
	"UnpakSiamida/common/helper"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSanitizer(t *testing.T) {
	sanitizer := helper.NewDefaultSanitizer()

	// Safe HTML tags
	inputSafe := "<p>Hello <b>World</b></p>"
	outputSafe := sanitizer.Sanitize(inputSafe)
	assert.Equal(t, template.HTML(inputSafe), outputSafe)

	// Dangerous HTML tags should be stripped or escaped
	inputDangerous := "<script>alert(1)</script><div>Normal text</div>"
	outputDangerous := sanitizer.Sanitize(inputDangerous)
	assert.NotContains(t, string(outputDangerous), "<script>")
	assert.NotContains(t, string(outputDangerous), "<div>")
	assert.Contains(t, string(outputDangerous), "Normal text")
}

func TestSanitizerContext(t *testing.T) {
	ctx := &helper.SanitizerContext{}

	// When strategy is not set, fallback behavior returns original input
	input := "<script>unsafe</script>"
	assert.Equal(t, template.HTML(input), ctx.Sanitize(input))

	// Set strategy
	ctx.SetStrategy(helper.NewDefaultSanitizer())
	output := ctx.Sanitize(input)
	assert.NotContains(t, string(output), "<script>")
}
