package application

import (
	domainAccount "UnpakSiamida/modules/account/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockLdapRepo struct {
	GetAccountsByGroupsFunc func(groups []string) ([]domainAccount.AccountLdap, error)
}

func (m *MockLdapRepo) GetAccountsByGroups(allowedGroups []string) ([]domainAccount.AccountLdap, error) {
	if m.GetAccountsByGroupsFunc != nil {
		return m.GetAccountsByGroupsFunc(allowedGroups)
	}
	return nil, nil
}

func TestGetAllAccountsLdapQueryHandler_Handle(t *testing.T) {
	mockRepo := &MockLdapRepo{
		GetAccountsByGroupsFunc: func(groups []string) ([]domainAccount.AccountLdap, error) {
			assert.Contains(t, groups, "adm_simonev")
			return []domainAccount.AccountLdap{
				{
					DN:           "CN=Adam,OU=Accounts,DC=unpak,DC=ac,DC=id",
					Username:     "adam",
					Name:         "Adam Furqon",
					Email:        "adam@unpak.ac.id",
					MatchedGroup: "adm_simonev",
				},
			}, nil
		},
	}

	handler := &GetAllAccountsLdapQueryHandler{
		LdapRepo: mockRepo,
	}

	res, err := handler.Handle(context.Background(), GetAllAccountsLdapQuery{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "adam", res[0].Username)
}
