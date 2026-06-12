package domain_test

import (
	"UnpakSiamida/common/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_ErrorMethod(t *testing.T) {
	err := domain.Error{
		Code:        "Domain.TestError",
		Description: "Something went wrong",
		Type:        domain.ErrorFailure,
	}

	assert.Equal(t, "Domain.TestError: Something went wrong", err.Error())
}

func TestError_FactoryMethods(t *testing.T) {
	// FailureError
	errFail := domain.FailureError("Fail.Code", "Fail description")
	assert.Equal(t, "Fail.Code", errFail.Code)
	assert.Equal(t, "Fail description", errFail.Description)
	assert.Equal(t, domain.ErrorFailure, errFail.Type)

	// NotFoundError
	errNotFound := domain.NotFoundError("NotFound.Code", "NotFound description")
	assert.Equal(t, "NotFound.Code", errNotFound.Code)
	assert.Equal(t, "NotFound description", errNotFound.Description)
	assert.Equal(t, domain.NotFound, errNotFound.Type)

	// ProblemError
	errProblem := domain.ProblemError("Problem.Code", "Problem description")
	assert.Equal(t, "Problem.Code", errProblem.Code)
	assert.Equal(t, "Problem description", errProblem.Description)
	assert.Equal(t, domain.Problem, errProblem.Type)

	// ConflictError
	errConflict := domain.ConflictError("Conflict.Code", "Conflict description")
	assert.Equal(t, "Conflict.Code", errConflict.Code)
	assert.Equal(t, "Conflict description", errConflict.Description)
	assert.Equal(t, domain.Conflict, errConflict.Type)
}
