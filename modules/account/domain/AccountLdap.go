package domain

type AccountLdap struct {
	DN           string   `json:"dn"`
	Username     string   `json:"username"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	EmployeeID   string   `json:"employee_id"`
	Groups       []string `json:"groups"`
	MatchedGroup string   `json:"matched_group,omitempty"`
}
