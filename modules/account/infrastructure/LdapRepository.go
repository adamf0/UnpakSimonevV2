package infrastructure

import (
	domainaccount "UnpakSiamida/modules/account/domain"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type LdapRepository struct{}

func NewLdapRepository() *LdapRepository {
	return &LdapRepository{}
}

func (r *LdapRepository) Connect() (*ldap.Conn, error) {
	adServer := os.Getenv("AD_SERVER")
	adUser := os.Getenv("AD_USER")
	adPass := os.Getenv("AD_PASS")

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	dialer := &net.Dialer{Timeout: 3 * time.Second}

	var conn *ldap.Conn
	var err error

	if strings.HasPrefix(adServer, "ldaps://") {
		conn, err = ldap.DialURL(adServer, ldap.DialWithTLSConfig(tlsConfig), ldap.DialWithDialer(dialer))
	} else {
		conn, err = ldap.DialURL(adServer, ldap.DialWithDialer(dialer))
		if err == nil {
			_ = conn.StartTLS(tlsConfig)
		}
	}

	if err != nil {
		log.Printf("[LDAP Connect Notice] Server %s unreachable (timeout 3s): %v", adServer, err)
		return nil, fmt.Errorf("failed to connect to LDAP server (%s): %w", adServer, err)
	}

	if adUser != "" && adPass != "" {
		err = conn.Bind(adUser, adPass)
		if err != nil {
			conn.Close()
			log.Printf("[LDAP Bind Notice] User: %s, Error: %v", adUser, err)
			return nil, fmt.Errorf("failed to bind to LDAP: %w", err)
		}
	}

	return conn, nil
}

func (r *LdapRepository) GetAccountsByGroups(allowedGroups []string) ([]domainaccount.AccountLdap, error) {
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

	// Direct Active Directory LDAP Search (ldaps://ldap.unpak.ac.id:636)
	conn, err := r.Connect()
	if err == nil && conn != nil {
		defer conn.Close()

		ldapDN := os.Getenv("LDAP_DN")
		if ldapDN == "" {
			ldapDN = "DC=unpak,DC=ac,DC=id"
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

		// Search users directly in Active Directory
		var userFilterParts []string
		for _, grp := range allowedGroups {
			userFilterParts = append(userFilterParts, fmt.Sprintf("(memberOf=CN=%s,OU=Groups,%s)", ldap.EscapeFilter(grp), ldapDN))
			userFilterParts = append(userFilterParts, fmt.Sprintf("(memberOf=CN=%s,%s)", ldap.EscapeFilter(grp), ldapDN))
			userFilterParts = append(userFilterParts, fmt.Sprintf("(memberOf=*CN=%s*)", ldap.EscapeFilter(grp)))
		}
		userFilterParts = append(userFilterParts, "(memberOf=*simonev*)")
		userFilterParts = append(userFilterParts, "(sAMAccountName=*simonev*)")

		userFilter := fmt.Sprintf("(&(objectClass=user)(|%s))", strings.Join(userFilterParts, ""))

		userSearchReq := ldap.NewSearchRequest(
			ldapDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, 0, false,
			userFilter,
			[]string{"dn", "sAMAccountName", "displayName", "cn", "mail", "employeeID", "employeeNumber", "memberOf"},
			nil,
		)

		uRes, err := conn.Search(userSearchReq)
		if err == nil && len(uRes.Entries) > 0 {
			for _, uEntry := range uRes.Entries {
				username := uEntry.GetAttributeValue("sAMAccountName")
				if username == "" {
					username = uEntry.GetAttributeValue("cn")
				}
				if username == "" || seen[username] {
					continue
				}
				seen[username] = true

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
				matchedGroup := allowedGroups[0]

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

		// Also search groups in Active Directory
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

					uEntry := uRes.Entries[0]
					username := uEntry.GetAttributeValue("sAMAccountName")
					if username == "" {
						username = uEntry.GetAttributeValue("cn")
					}
					if username == "" || seen[username] {
						continue
					}
					seen[username] = true
					seen[mDN] = true

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

		if len(result) > 0 {
			log.Printf("[LdapRepository] Successfully fetched %d users from Active Directory LDAP (%s)", len(result), os.Getenv("AD_SERVER"))
			return result, nil
		}
	}

	// Fallback: Return standard SSO account definitions for allowedGroups so UI dropdown is always populated
	// log.Printf("[LdapRepository] Returning fallback SSO account definitions for groups: %v", allowedGroups)
	var fallbackUsers []domainaccount.AccountLdap
	// for _, grp := range allowedGroups {
	// 	fallbackUsers = append(fallbackUsers, domainaccount.AccountLdap{
	// 		DN:           fmt.Sprintf("CN=%s,OU=users,DC=unpak,DC=ac,DC=id", grp),
	// 		Username:     grp,
	// 		Name:         strings.Title(strings.ReplaceAll(grp, "_", " ")),
	// 		Email:        fmt.Sprintf("%s@unpak.ac.id", grp),
	// 		EmployeeID:   grp,
	// 		Groups:       []string{grp},
	// 		MatchedGroup: grp,
	// 	})
	// }

	return fallbackUsers, nil
}
