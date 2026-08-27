package infrastructure

import (
	domainaccount "UnpakSiamida/modules/account/domain"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type LdapRepository struct{}

func NewLdapRepository() *LdapRepository {
	return &LdapRepository{}
}

func (r *LdapRepository) Connect() (*ldap.Conn, error) {
	adServer := os.Getenv("AD_SERVER")
	if adServer == "" {
		adServer = "ldaps://ldap.unpak.ac.id:636"
	}
	adUser := os.Getenv("AD_USER")
	adPass := os.Getenv("AD_PASS")

	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	var conn *ldap.Conn
	var err error

	if strings.HasPrefix(adServer, "ldaps://") {
		conn, err = ldap.DialURL(adServer, ldap.DialWithTLSConfig(tlsConfig))
	} else {
		conn, err = ldap.DialURL(adServer)
		if err == nil {
			_ = conn.StartTLS(tlsConfig)
		}
	}

	if err != nil {
		log.Printf("[LDAP Connect Error] Server: %s, Error: %v", adServer, err)
		return nil, fmt.Errorf("failed to connect to LDAP server (%s): %w", adServer, err)
	}

	if adUser != "" && adPass != "" {
		err = conn.Bind(adUser, adPass)
		if err != nil {
			conn.Close()
			log.Printf("[LDAP Bind Error] User: %s, Error: %v", adUser, err)
			return nil, fmt.Errorf("failed to bind to LDAP: %w", err)
		}
	}

	return conn, nil
}

func (r *LdapRepository) GetAccountsByGroups(allowedGroups []string) ([]domainaccount.AccountLdap, error) {
	conn, err := r.Connect()
	if err != nil {
		log.Printf("[LDAP Error] %v", err)
		return []domainaccount.AccountLdap{}, err
	}
	defer conn.Close()

	ldapDN := os.Getenv("LDAP_DN")
	if ldapDN == "" {
		ldapDN = "DC=unpak,DC=ac,DC=id"
	}

	if len(allowedGroups) == 0 {
		allowedGroups = []string{
			"adm_simonev",
			"admin",
			"superadmin",
			"adm_pusat",
			"adm_simonev_fakultas",
			"adm_simonev_prodi",
		}
	}

	result := []domainaccount.AccountLdap{}
	seen := make(map[string]bool)
	cnRegex := regexp.MustCompile(`(?i)CN=([^,]+)`)

	isGroupAllowed := func(groupName string) (bool, string) {
		trimmed := strings.TrimSpace(groupName)
		for _, allowed := range allowedGroups {
			if strings.EqualFold(trimmed, allowed) || strings.Contains(strings.ToLower(trimmed), "simonev") {
				return true, allowed
			}
		}
		return false, ""
	}

	// ------------------------------------------------------------------
	// STRATEGY 1: Search AD Group Objects and fetch member User DNs
	// ------------------------------------------------------------------
	var groupFilterParts []string
	for _, grp := range allowedGroups {
		groupFilterParts = append(groupFilterParts, fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(grp)))
		groupFilterParts = append(groupFilterParts, fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(grp)))
	}
	groupFilterParts = append(groupFilterParts, "(cn=*simonev*)")

	groupFilter := fmt.Sprintf("(&(objectClass=group)(|%s))", strings.Join(groupFilterParts, ""))

	groupSearchReq := ldap.NewSearchRequest(
		ldapDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		groupFilter,
		[]string{"cn", "sAMAccountName", "member"},
		nil,
	)

	groupRes, err := conn.Search(groupSearchReq)
	if err == nil && len(groupRes.Entries) > 0 {
		for _, gEntry := range groupRes.Entries {
			grpCN := gEntry.GetAttributeValue("cn")
			if grpCN == "" {
				grpCN = gEntry.GetAttributeValue("sAMAccountName")
			}

			members := gEntry.GetAttributeValues("member")
			for _, mDN := range members {
				if seen[mDN] || strings.TrimSpace(mDN) == "" {
					continue
				}

				userReq := ldap.NewSearchRequest(
					mDN,
					ldap.ScopeBaseObject,
					ldap.NeverDerefAliases,
					0, 0, false,
					"(objectClass=*)",
					[]string{"dn", "sAMAccountName", "displayName", "cn", "mail", "employeeID", "employeeNumber", "memberOf"},
					nil,
				)

				uRes, err := conn.Search(userReq)
				if err != nil || len(uRes.Entries) == 0 {
					continue
				}

				seen[mDN] = true
				uEntry := uRes.Entries[0]
				username := uEntry.GetAttributeValue("sAMAccountName")
				if username == "" {
					username = uEntry.GetAttributeValue("cn")
				}
				name := uEntry.GetAttributeValue("displayName")
				if name == "" {
					name = uEntry.GetAttributeValue("cn")
				}
				email := uEntry.GetAttributeValue("mail")
				empID := uEntry.GetAttributeValue("employeeID")
				if empID == "" {
					empID = uEntry.GetAttributeValue("employeeNumber")
				}

				rawMemberOf := uEntry.GetAttributeValues("memberOf")
				var userGroups []string
				matchedGroup := grpCN

				for _, m := range rawMemberOf {
					matches := cnRegex.FindStringSubmatch(m)
					if len(matches) > 1 {
						gName := matches[1]
						userGroups = append(userGroups, gName)
						if ok, match := isGroupAllowed(gName); ok {
							matchedGroup = match
						}
					}
				}

				result = append(result, domainaccount.AccountLdap{
					DN:           uEntry.DN,
					Username:     username,
					Name:         name,
					Email:        email,
					EmployeeID:   empID,
					Groups:       userGroups,
					MatchedGroup: matchedGroup,
				})
			}
		}
	}

	// ------------------------------------------------------------------
	// STRATEGY 2: Search active users under ldapDN and check memberOf in memory
	// ------------------------------------------------------------------
	userSearchReq := ldap.NewSearchRequest(
		ldapDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(&(objectClass=user)(sAMAccountName=*))",
		[]string{"dn", "sAMAccountName", "displayName", "cn", "mail", "employeeID", "employeeNumber", "memberOf"},
		nil,
	)

	userRes, err := conn.Search(userSearchReq)
	if err == nil && len(userRes.Entries) > 0 {
		for _, uEntry := range userRes.Entries {
			dn := uEntry.DN
			if seen[dn] {
				continue
			}

			rawMemberOf := uEntry.GetAttributeValues("memberOf")
			var userGroups []string
			var matchedGroup string
			isMatched := false

			for _, m := range rawMemberOf {
				matches := cnRegex.FindStringSubmatch(m)
				if len(matches) > 1 {
					gName := matches[1]
					userGroups = append(userGroups, gName)
					if ok, match := isGroupAllowed(gName); ok {
						isMatched = true
						matchedGroup = match
					}
				}
			}

			if !isMatched {
				continue
			}

			seen[dn] = true
			username := uEntry.GetAttributeValue("sAMAccountName")
			if username == "" {
				username = uEntry.GetAttributeValue("cn")
			}
			name := uEntry.GetAttributeValue("displayName")
			if name == "" {
				name = uEntry.GetAttributeValue("cn")
			}
			email := uEntry.GetAttributeValue("mail")
			empID := uEntry.GetAttributeValue("employeeID")
			if empID == "" {
				empID = uEntry.GetAttributeValue("employeeNumber")
			}

			result = append(result, domainaccount.AccountLdap{
				DN:           dn,
				Username:     username,
				Name:         name,
				Email:        email,
				EmployeeID:   empID,
				Groups:       userGroups,
				MatchedGroup: matchedGroup,
			})
		}
	}

	return result, nil
}
