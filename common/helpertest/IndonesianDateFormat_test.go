package commontest

import (
	"UnpakSiamida/common/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIndonesianDateFormatter(t *testing.T) {
	formatter := helper.IndonesianDateFormatter{}

	tm := time.Date(2024, time.January, 15, 10, 20, 30, 0, time.UTC)

	assert.Equal(t, "Senin", formatter.NameDay(tm))
	assert.Equal(t, "15", formatter.Day(tm))
	assert.Equal(t, "Januari", formatter.Month(tm))
	assert.Equal(t, "2024", formatter.Year(tm))
	assert.Equal(t, "15 Januari 2024", formatter.Format(tm))
	assert.Equal(t, "Senin, 15 Januari 2024", formatter.FormatWithDay(tm))
	assert.Equal(t, "15 Januari 2024 10:20:30 WIB", formatter.FormatWithTime(tm))

	// Zero values
	var zero time.Time
	assert.Empty(t, formatter.Format(zero))
	assert.Empty(t, formatter.FormatWithDay(zero))
	assert.Empty(t, formatter.FormatWithTime(zero))

	// Pointer tests
	assert.Equal(t, "15 Januari 2024", formatter.FormatDefault(&tm))
	assert.Equal(t, "Senin, 15 Januari 2024", formatter.FormatWithDayDefault(&tm))
	assert.Equal(t, "15 Januari 2024 10:20:30 WIB", formatter.FormatWithTimeDefault(&tm))

	// Nil pointer tests
	assert.Empty(t, formatter.FormatDefault(nil))
	assert.Empty(t, formatter.FormatWithDayDefault(nil))
	assert.Empty(t, formatter.FormatWithTimeDefault(nil))
}

func TestDateContext(t *testing.T) {
	formatter := helper.IndonesianDateFormatter{}
	ctx := &helper.DateContext{}
	ctx.SetStrategy(formatter)

	tm := time.Date(2024, time.January, 15, 10, 20, 30, 0, time.UTC)

	assert.Equal(t, "Senin", ctx.NameDay(tm))
	assert.Equal(t, "15", ctx.Day(tm))
	assert.Equal(t, "Januari", ctx.Month(tm))
	assert.Equal(t, "2024", ctx.Year(tm))
	assert.Equal(t, "15 Januari 2024", ctx.Format(tm))
	assert.Equal(t, "Senin, 15 Januari 2024", ctx.FormatWithDay(tm))
	assert.Equal(t, "15 Januari 2024 10:20:30 WIB", ctx.FormatWithTime(tm))
	assert.Equal(t, "15 Januari 2024", ctx.FormatDefault(&tm))
	assert.Equal(t, "Senin, 15 Januari 2024", ctx.FormatWithDayDefault(&tm))
	assert.Equal(t, "15 Januari 2024 10:20:30 WIB", ctx.FormatWithTimeDefault(&tm))
}
