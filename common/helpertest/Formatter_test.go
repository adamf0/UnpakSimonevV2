package commontest

import (
	"UnpakSiamida/common/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateChain_Success(t *testing.T) {
	// Setup parsing strategy
	parser := helper.DateLayoutFirstFactory{}.CreateParser()
	formatter := helper.DateLayoutFirstFactory{}.CreateFormatter()

	// Parse valid RFC3339 date
	dc := helper.NewDateChain("2024-01-02T10:20:30Z").
		UseParseStrategy(parser).
		UseFormatStrategy(formatter).
		Parse()

	assert.NotEmpty(t, dc.FormatString())
	assert.NotEqual(t, "-", dc.FormatString())
	
	// Parse plain date
	dc2 := helper.NewDateChain("2024-01-02").
		UseParseStrategy(parser).
		UseFormatStrategy(formatter).
		Parse()

	assert.Contains(t, dc2.FormatString(), "Januari")
}

func TestDateChain_Truncate(t *testing.T) {
	parser := helper.DateLayoutFirstFactory{}.CreateParser()
	dc := helper.NewDateChain("2024-01-02T10:20:30Z").
		UseParseStrategy(parser).
		Parse().
		Truncate(24 * time.Hour)

	tm, err := dc.Ptr()
	require.NoError(t, err)
	require.NotNil(t, tm)
	assert.Equal(t, 0, tm.Hour())
	assert.Equal(t, 0, tm.Minute())
}

func TestDateChain_NoStrategyError(t *testing.T) {
	dc := helper.NewDateChain("2024-01-02").Parse()
	_, err := dc.Ptr()
	assert.Error(t, err)
	assert.Equal(t, "parse strategy not set", err.Error())
}

func TestDateChain_FallbackOnEmptyOrInvalid(t *testing.T) {
	parser := helper.DateLayoutFirstFactory{}.CreateParser()

	// Empty string fallback
	dcEmpty := helper.NewDateChain("  ").
		UseParseStrategy(parser).
		Parse()
	assert.Equal(t, "-", dcEmpty.FormatString())

	// Invalid format fallback
	dcInvalid := helper.NewDateChain("invalid-date-format").
		UseParseStrategy(parser).
		Parse()
	assert.Equal(t, "invalid-date-format", dcInvalid.FormatString())
}

func TestDateLayoutSecondFactory(t *testing.T) {
	parser := helper.DateLayoutSecondFactory{}.CreateParser()
	formatter := helper.DateLayoutSecondFactory{}.CreateFormatter()

	dc := helper.NewDateChain("2024-01-02").
		UseParseStrategy(parser).
		UseFormatStrategy(formatter).
		Parse()

	assert.Contains(t, dc.FormatString(), "Januari")
}
