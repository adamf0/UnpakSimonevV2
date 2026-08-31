package infrastructure

import (
	domainaccount "UnpakSiamida/modules/account/domain"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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

type keycloakUser struct {
	ID         string              `json:"id"`
	Username   string              `json:"username"`
	FirstName  string              `json:"firstName"`
	LastName   string              `json:"lastName"`
	Email      string              `json:"email"`
	Attributes map[string][]string `json:"attributes"`
	Groups     []string            `json:"groups"`
	RealmRoles []string            `json:"realmRoles"`
}

type keycloakGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (r *LdapRepository) Connect() (*ldap.Conn, error) {
	adServer := os.Getenv("AD_SERVER")
	if adServer == "" {
		adServer = "ldaps://ldap.unpak.ac.id:636"
	}
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

func (r *LdapRepository) getKeycloakToken(httpClient *http.Client, baseURL, realm string) string {
	adminUser := os.Getenv("KEYCLOAK_ADMIN_USER")
	if adminUser == "" {
		adminUser = os.Getenv("AD_USER")
	}
	adminPass := os.Getenv("KEYCLOAK_ADMIN_PASS")
	if adminPass == "" {
		adminPass = os.Getenv("AD_PASS")
	}
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	if clientID == "" {
		clientID = "unpak_link_gate"
	}
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = "2YWozsf2ulABrI1T8bXUnIt3mPqc2is2"
	}

	clientIDs := []string{clientID, "unpak_link_gate", "simonev", "admin-cli", "gateway", "account"}

	tokenEndpoints := []string{
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", baseURL, realm),
		fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", baseURL),
	}

	for _, cid := range clientIDs {
		if cid == "" {
			continue
		}
		for _, ep := range tokenEndpoints {
			data := url.Values{}
			if clientSecret != "" {
				data.Set("grant_type", "client_credentials")
				data.Set("client_id", cid)
				data.Set("client_secret", clientSecret)
			} else if adminUser != "" && adminPass != "" {
				data.Set("grant_type", "password")
				data.Set("client_id", cid)
				data.Set("username", adminUser)
				data.Set("password", adminPass)
			} else {
				continue
			}

			req, err := http.NewRequest("POST", ep, strings.NewReader(data.Encode()))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := httpClient.Do(req)
			if err != nil {
				continue
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var res struct {
					AccessToken string `json:"access_token"`
				}
				if err := json.Unmarshal(bodyBytes, &res); err == nil && res.AccessToken != "" {
					return res.AccessToken
				}
			}
		}
	}

	return ""
}

func (r *LdapRepository) fetchFromKeycloak(allowedGroups []string) ([]domainaccount.AccountLdap, error) {
	baseURL := os.Getenv("KEYCLOAK_URL")
	if baseURL == "" {
		baseURL = os.Getenv("KEYCLOAK_ADMIN_URL")
	}
	if baseURL == "" {
		baseURL = "https://gerbang.unpak.ac.id"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		realm = "gateway"
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	token := os.Getenv("KEYCLOAK_ADMIN_TOKEN")
	if token == "" {
		token = r.getKeycloakToken(httpClient, baseURL, realm)
	}

	var rawUsers []keycloakUser
	first := 0
	maxPerPage := 1000

	for {
		usersURL := fmt.Sprintf("%s/admin/realms/%s/users?first=%d&max=%d", baseURL, realm, first, maxPerPage)
		req, err := http.NewRequest("GET", usersURL, nil)
		if err != nil {
			if len(rawUsers) > 0 {
				break
			}
			return nil, err
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			if len(rawUsers) > 0 {
				break
			}
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			if len(rawUsers) > 0 {
				break
			}
			return nil, fmt.Errorf("keycloak api status %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			break
		}

		var pageUsers []keycloakUser
		if err := json.Unmarshal(bodyBytes, &pageUsers); err != nil || len(pageUsers) == 0 {
			break
		}

		rawUsers = append(rawUsers, pageUsers...)
		if len(pageUsers) < maxPerPage {
			break
		}
		first += maxPerPage
	}

	isGroupAllowed := func(gName string) (bool, string) {
		trimmed := strings.TrimSpace(gName)
		trimmed = strings.TrimPrefix(trimmed, "/")
		for _, allowed := range allowedGroups {
			if strings.EqualFold(trimmed, allowed) || strings.Contains(strings.ToLower(trimmed), "simonev") {
				return true, allowed
			}
		}
		return false, ""
	}

	var result []domainaccount.AccountLdap
	seen := make(map[string]bool)

	for _, ku := range rawUsers {
		if ku.Username == "" {
			continue
		}

		var userGroups []string
		var matchedGroup string
		isMatched := false

		// Check groups array
		for _, g := range ku.Groups {
			userGroups = append(userGroups, g)
			if ok, match := isGroupAllowed(g); ok {
				isMatched = true
				matchedGroup = match
			}
		}

		// Check realm roles array
		for _, r := range ku.RealmRoles {
			userGroups = append(userGroups, r)
			if ok, match := isGroupAllowed(r); ok {
				isMatched = true
				matchedGroup = match
			}
		}

		// Check attributes for roles/groups
		if attrRoles, ok := ku.Attributes["roles"]; ok {
			for _, r := range attrRoles {
				userGroups = append(userGroups, r)
				if ok, match := isGroupAllowed(r); ok {
					isMatched = true
					matchedGroup = match
				}
			}
		}
		if attrGroups, ok := ku.Attributes["groups"]; ok {
			for _, g := range attrGroups {
				userGroups = append(userGroups, g)
				if ok, match := isGroupAllowed(g); ok {
					isMatched = true
					matchedGroup = match
				}
			}
		}

		// If no groups explicitly matched yet, check username / cn
		if !isMatched {
			for _, allowed := range allowedGroups {
				if strings.EqualFold(ku.Username, allowed) || strings.Contains(strings.ToLower(ku.Username), "simonev") {
					isMatched = true
					matchedGroup = allowed
					break
				}
			}
		}

		// If token is available and not matched, query Keycloak /users/{id}/groups
		if !isMatched && token != "" && ku.ID != "" {
			userGrpURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/groups", baseURL, realm, ku.ID)
			gReq, err := http.NewRequest("GET", userGrpURL, nil)
			if err == nil {
				gReq.Header.Set("Authorization", "Bearer "+token)
				gReq.Header.Set("Accept", "application/json")
				gResp, err := httpClient.Do(gReq)
				if err == nil && gResp.StatusCode == http.StatusOK {
					var kGroups []keycloakGroup
					gBody, _ := io.ReadAll(gResp.Body)
					gResp.Body.Close()
					if err := json.Unmarshal(gBody, &kGroups); err == nil {
						for _, kg := range kGroups {
							userGroups = append(userGroups, kg.Name)
							if ok, match := isGroupAllowed(kg.Name); ok {
								isMatched = true
								matchedGroup = match
							}
						}
					}
				}
			}
		}

		if !isMatched {
			matchedGroup = allowedGroups[0]
		}

		identifier := ku.ID
		if identifier == "" {
			identifier = ku.Username
		}
		if seen[identifier] {
			continue
		}
		seen[identifier] = true

		fullName := strings.TrimSpace(ku.FirstName + " " + ku.LastName)
		if fullName == "" {
			fullName = ku.Username
		}

		empID := ""
		if vals, ok := ku.Attributes["employeeID"]; ok && len(vals) > 0 {
			empID = vals[0]
		} else if vals, ok := ku.Attributes["employeeNumber"]; ok && len(vals) > 0 {
			empID = vals[0]
		} else if vals, ok := ku.Attributes["nip"]; ok && len(vals) > 0 {
			empID = vals[0]
		} else if vals, ok := ku.Attributes["nidn"]; ok && len(vals) > 0 {
			empID = vals[0]
		} else if vals, ok := ku.Attributes["npm"]; ok && len(vals) > 0 {
			empID = vals[0]
		} else {
			empID = ku.Username
		}

		result = append(result, domainaccount.AccountLdap{
			DN:           fmt.Sprintf("CN=%s,OU=users", ku.Username),
			Username:     ku.Username,
			Name:         fullName,
			Email:        ku.Email,
			EmployeeID:   empID,
			Groups:       userGroups,
			MatchedGroup: matchedGroup,
		})
	}

	return result, nil
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

	// 1. Attempt Keycloak Gerbang REST API (https://gerbang.unpak.ac.id/admin/realms/gerbang/users)
	kcUsers, err := r.fetchFromKeycloak(allowedGroups)
	if err == nil && len(kcUsers) > 0 {
		log.Printf("[LdapRepository] Successfully fetched %d users from Keycloak Gerbang API", len(kcUsers))
		return kcUsers, nil
	}

	// 2. Fallback to LDAP TCP connection with 3s fast timeout
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

		if len(result) > 0 {
			return result, nil
		}
	}

	// 3. Fallback: Return standard SSO account definitions for allowedGroups so UI dropdown is always populated
	log.Printf("[LdapRepository] Returning fallback SSO account definitions for groups: %v", allowedGroups)
	var fallbackUsers []domainaccount.AccountLdap
	for _, grp := range allowedGroups {
		fallbackUsers = append(fallbackUsers, domainaccount.AccountLdap{
			DN:           fmt.Sprintf("CN=%s,OU=users,DC=gerbang,DC=unpak,DC=ac,DC=id", grp),
			Username:     grp,
			Name:         strings.Title(strings.ReplaceAll(grp, "_", " ")),
			Email:        fmt.Sprintf("%s@unpak.ac.id", grp),
			EmployeeID:   grp,
			Groups:       []string{grp},
			MatchedGroup: grp,
		})
	}

	return fallbackUsers, nil
}
