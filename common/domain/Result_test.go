package domain_test

import (
	"UnpakSiamida/common/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_Success(t *testing.T) {
	res := domain.Success()
	assert.True(t, res.IsSuccess)
	assert.Equal(t, domain.None, res.Error)
}

func TestResult_Failure(t *testing.T) {
	testErr := domain.FailureError("Test.Code", "Test description")
	res := domain.Failure(testErr)
	assert.False(t, res.IsSuccess)
	assert.Equal(t, testErr, res.Error)
}

func TestResultValue_Success(t *testing.T) {
	val := "hello"
	resVal := domain.SuccessValue(val)
	assert.True(t, resVal.IsSuccess)
	assert.Equal(t, val, resVal.Value)

	got, err := resVal.GetValue()
	assert.NoError(t, err)
	assert.Equal(t, val, got)
}

func TestResultValue_Failure(t *testing.T) {
	testErr := domain.FailureError("Test.Code", "Test description")
	resVal := domain.FailureValue[string](testErr)
	assert.False(t, resVal.IsSuccess)
	assert.Empty(t, resVal.Value)

	_, err := resVal.GetValue()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot get value of failed result: Test description")
}
