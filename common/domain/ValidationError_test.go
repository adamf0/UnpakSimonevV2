package domain_test

import (
	"UnpakSiamida/common/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	err1 := domain.FailureError("Code.1", "Desc 1")
	err2 := domain.FailureError("Code.2", "Desc 2")

	ve := domain.NewValidationError([]domain.Error{err1, err2})
	assert.Len(t, ve.Errors, 2)
	assert.Equal(t, err1, ve.Errors[0])

	// Test FromResults
	results := []domain.Result{
		domain.Success(),
		domain.Failure(err1),
		domain.Success(),
		domain.Failure(err2),
	}

	ve2 := domain.FromResults(results)
	assert.Len(t, ve2.Errors, 2)
	assert.Equal(t, err1, ve2.Errors[0])
	assert.Equal(t, err2, ve2.Errors[1])
}
