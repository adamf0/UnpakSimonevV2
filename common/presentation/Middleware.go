package presentation

import (
	"context"
	"errors"
	"html"
	"io"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"UnpakSiamida/common/helper"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainaccount "UnpakSiamida/modules/account/domain"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/text/unicode/norm"
)

// =======================
// CONFIG
// =======================
var jwtSecret = []byte("secret")

type HeaderSecurityConfig struct {
	BlacklistedHeaderNames map[string]bool
	AllowDomains           []string
	MaxHeaderLen           int
	ResolveAndCheck        bool
	LookupTimeout          time.Duration
	BlockedCIDRs           []string
}

func DefaultBlacklistedHeaderNames() map[string]bool {
	names := []string{
		"x-forwarded-for", "x-forwarded-host", "forwarded", "forwarded-host",
		"x-forwarded-proto", "x-forwarded-port", "x-forwarded-scheme",
		"x-real-ip", "client-ip", "true-client-ip", "cf-connecting-ip",
		"x-remote-ip", "x-originating-ip",
		"x-original-host", "via", "x-via",
		"host", "x-host", "x-rewrite-url", "x-original-url",
		"x-request-url", "x-request-uri", "redirect", "location",
		"authorization", "proxy-authorization", "x-api-key",
		"metadata", "x-aws-ec2-metadata", "referer",
	}
	out := map[string]bool{}
	for _, v := range names {
		out[strings.ToLower(v)] = true
	}
	return out
}

func DefaultHeaderSecurityConfig() *HeaderSecurityConfig {
	domains := []string{"simonev.unpak.ac.id", "gerbang.unpak.ac.id", "localhost", "localhost:3000", "localhost:4000", "127.0.0.1", "127.0.0.1:3000", "127.0.0.1:4000", "thunderclient.com", "postman"}
	if envDomains := os.Getenv("ALLOWED_HOSTS"); envDomains != "" {
		parts := strings.Split(envDomains, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				domains = append(domains, p)
			}
		}
	}

	return &HeaderSecurityConfig{
		BlacklistedHeaderNames: DefaultBlacklistedHeaderNames(),
		AllowDomains:           domains,
		MaxHeaderLen:           8192,
		ResolveAndCheck:        false,
		LookupTimeout:          1 * time.Second,
		BlockedCIDRs:           []string{},
	}
}

// =======================
// REGEX
// =======================

var (
	crlfRe      = regexp.MustCompile(`[\r\n]`)
	nullRe      = regexp.MustCompile(`\x00`)
	protoRe     = regexp.MustCompile(`(?i)^(javascript|data|vbscript|file|view-source):`)
	punyRe      = regexp.MustCompile(`(?i)xn--[a-z0-9-]+`)
	zeroWidthRe = regexp.MustCompile(`[\x{200B}\x{200C}\x{200D}\x{2060}\x{FEFF}]`)
	hostExtract = regexp.MustCompile(`(?i)(?:https?://)?([a-z0-9\.\-]+\.[a-z]{2,})(:\d+)?`)
)

type Account struct {
	// UUID         string      `json:"uuid"`
	// NidnUsername string      `json:"nidn_username"`
	// Password     string      `json:"password"`
	// Level        string      `json:"level"`
	// Name         string      `json:"name"`
	// Email        string      `json:"email"`
	// FakultasUnit string      `json:"fakultas_unit"`
	// ExtraRole    []ExtraRole `gorm:"-" json:"extrarole,omitempty"`
	ID          string  `json:"ID"`
	UUID        *string `json:"UUID"`
	Username    *string `json:"Username"`
	Password    *string `json:"-"`
	Level       *string `json:"Level"`
	Name        *string `json:"Name"`
	Email       *string `json:"Email"`
	RefFakultas *string `json:"RefFakultas"`
	Fakultas    *string `json:"Fakultas"`
	RefProdi    *string `json:"RefProdi"`
	Prodi       *string `json:"Prodi"`
	Unit        *string `json:"Unit"`
	Resource    *string `json:"Resource"`
	CodeCtx     *string `json:"CodeCtx"` //untuk membedakan dosen & mahasiswa
}
type ExtraRole struct {
	Tahun string `json:"tahun"`
	Role  string `json:"role"`
}

// =======================
// MIDDLEWARE
// =======================
const logCommonTokenLabel = "common.token"
const logCommonRbac = "common.rbac"

func HeaderSecurityMiddleware(cfg *HeaderSecurityConfig) fiber.Handler {
	if cfg == nil {
		cfg = DefaultHeaderSecurityConfig()
	}

	blocked := parseBlockedCIDRs(cfg.BlockedCIDRs)

	return func(c *fiber.Ctx) error {

		for name, vals := range c.GetReqHeaders() {
			for _, val := range vals {

				if err := validateHeaderLength(name, val, cfg); err != nil {
					return badRequest(c, err)
				}

				if err := validateControlChars(name, val); err != nil {
					return badRequest(c, err)
				}

				decoded := normalizeHeader(val)

				if err := validateProtocol(name, decoded); err != nil {
					return badRequest(c, err)
				}

				if err := validatePunycode(name, decoded); err != nil {
					return badRequest(c, err)
				}

				if err := validateZeroWidth(name, val); err != nil {
					return badRequest(c, err)
				}

				if err := validateURLDomain(decoded, cfg, blocked); err != nil {
					return badRequest(c, err)
				}

				if err := validateHostHeader(name, decoded, cfg); err != nil {
					return badRequest(c, err)
				}
			}
		}

		if err := validateEmbeddedDomains(c, cfg); err != nil {
			return badRequest(c, err)
		}

		return c.Next()
	}
}

// =======================
// VALIDATION HELPERS
// =======================

func validateHeaderLength(name, val string, cfg *HeaderSecurityConfig) error {
	if len(val) > cfg.MaxHeaderLen {
		return commoninfra.NewResponseError("common.check[A+1]", "header too long: "+name)
	}
	return nil
}

func validateControlChars(name, val string) error {
	if crlfRe.MatchString(val) || nullRe.MatchString(val) {
		return commoninfra.NewResponseError("common.check[A+2]", "header ctrl char: "+name)
	}
	return nil
}

func normalizeHeader(val string) string {
	decoded := multiUnescape(html.UnescapeString(val), 3)
	decoded = zeroWidthRe.ReplaceAllString(decoded, "")
	return norm.NFKC.String(decoded)
}

func validateProtocol(name, decoded string) error {
	if protoRe.MatchString(decoded) {
		return commoninfra.NewResponseError("common.check[A+3]", "protocol attack: "+name)
	}
	return nil
}

func validatePunycode(name, decoded string) error {
	if punyRe.MatchString(decoded) {
		return commoninfra.NewResponseError("common.check[A+4]", "punycode forbidden: "+name)
	}
	return nil
}

func validateZeroWidth(name, val string) error {
	if zeroWidthRe.MatchString(val) {
		return commoninfra.NewResponseError("common.check[A+5]", "zero width attack: "+name)
	}
	return nil
}

func validateURLDomain(decoded string, cfg *HeaderSecurityConfig, blocked []*net.IPNet) error {
	u, err := url.Parse(decoded)
	if err != nil || u.Host == "" {
		return nil
	}

	host := u.Hostname()

	if !domainAllowed(host, cfg.AllowDomains) {
		return commoninfra.NewResponseError("common.check[A+6]", "domain not allowed: "+host)
	}

	if cfg.ResolveAndCheck {
		if err := resolveAndCheckIP(host, cfg.LookupTimeout, blocked); err != nil {
			return err
		}
	}

	return nil
}

func resolveAndCheckIP(host string, timeout time.Duration, blocked []*net.IPNet) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ips, _ := net.DefaultResolver.LookupIP(ctx, "ip", host)

	for _, ip := range ips {
		if ipInNets(ip, blocked) {
			return commoninfra.NewResponseError("common.check[A+7]", "domain resolves to forbidden IP: "+host)
		}
	}

	return nil
}

func validateHostHeader(name, decoded string, cfg *HeaderSecurityConfig) error {
	if strings.ToLower(name) != "host" {
		return nil
	}

	// godump.Dd(decoded, cfg.AllowDomains)
	if !domainAllowed(decoded, cfg.AllowDomains) {
		return commoninfra.NewResponseError("common.check[A+8]", "host header spoof: "+decoded)
	}
	return nil
}

func validateEmbeddedDomains(c *fiber.Ctx, cfg *HeaderSecurityConfig) error {

	urlHeaders := []string{
		"referer", "origin", "location", "refferer",
		"referrer", "redirect", "url", "http-url",
		"x-rewrite-url", "x-http-destinationurl",
		"x-http-host-override", "x-forwarded-host",
	}

	for _, h := range urlHeaders {
		val := c.Get(h)
		if val == "" {
			continue
		}

		decoded := normalizeHeader(val)
		hosts := extractHostsFromText(decoded)

		for _, host := range hosts {
			if !domainAllowed(host, cfg.AllowDomains) {
				return commoninfra.NewResponseError(
					"common.check[A+9]",
					"embedded domain not allowed: "+host,
				)
			}
		}
	}

	return nil
}

// =======================
// UTILITIES
// =======================

func badRequest(c *fiber.Ctx, err error) error {
	return c.Status(400).JSON(err)
}

func parseBlockedCIDRs(cidrs []string) []*net.IPNet {
	var blocked []*net.IPNet
	for _, c := range cidrs {
		_, ipnet, err := net.ParseCIDR(c)
		if err == nil {
			blocked = append(blocked, ipnet)
		}
	}
	return blocked
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr, err := extractBearerToken(c)
		if err != nil {
			return err
		}

		token, err := parseJWT(tokenStr)
		if err != nil {
			return err
		}

		claims, err := validateClaims(token)
		if err != nil {
			return err
		}

		injectRequestValues(c, claims, tokenStr)

		return c.Next()
	}
}

func extractBearerToken(c *fiber.Ctx) (string, error) {
	authQuery := c.Query("ctxtoken")
	authHeader := c.Get("Authorization")
	log.Printf("Authorization header: %s", authHeader)

	if authHeader == "" && authQuery == "" {
		log.Println("Authorization header missing")
		return "", c.Status(400).
			JSON(commoninfra.NewResponseError(logCommonRbac, "authorization header missing"))
	}

	var jwt = authHeader
	if authQuery != "" {
		jwt = strings.TrimPrefix(authQuery, "Bearer ")
		jwt = strings.TrimSpace(jwt)
		jwt = "Bearer " + jwt
	}

	parts := strings.Split(jwt, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		log.Println("Invalid authorization header format")
		return "", c.Status(400).
			JSON(commoninfra.NewResponseError(logCommonRbac, "authorization header format must be Bearer token"))
	}

	token := parts[1]
	log.Printf("Token: %s", token)
	return token, nil
}

func parseJWT(tokenStr string) (*jwt.Token, error) {
	if tokenStr == "" || tokenStr == "undefined" || tokenStr == "null" {
		return nil, fiber.NewError(400, "token authorization header missing or empty")
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fiber.NewError(400, "token format invalid: expected JWT with 3 segments")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err == nil && token.Valid {
		return token, nil
	}

	// Fallback: parse unverified for Keycloak RSA tokens (RS256/RS512)
	token, _, err = new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, fiber.NewError(400, "failed to parse token: "+err.Error())
	}

	return token, nil
}

func validateClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.NewError(400, "invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fiber.NewError(400, "token expired")
		}
	}

	return claims, nil
}

func injectRequestValues(c *fiber.Ctx, claims jwt.MapClaims, tokenStr string) {
	iss, _ := claims["iss"].(string)
	employeeId, _ := claims["employeeid"].(string)

	isKeycloak := strings.Contains(iss, "gerbang.unpak.ac.id") || employeeId != ""

	if isKeycloak {
		if employeeId == "" {
			if pref, ok := claims["preferred_username"].(string); ok && pref != "" {
				employeeId = pref
			} else if sub, ok := claims["sub"].(string); ok {
				employeeId = sub
			} else if uName, ok := claims["username"].(string); ok {
				employeeId = uName
			}
		}

		c.Request().PostArgs().Set("sid", employeeId)
		c.Locals("sid", employeeId)

		var allGroups []string
		if groupRaw, ok := claims["group"].([]interface{}); ok {
			for _, g := range groupRaw {
				if gStr, ok := g.(string); ok {
					allGroups = append(allGroups, strings.ToLower(gStr))
				}
			}
		}
		if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
			if rolesRaw, ok := realmAccess["roles"].([]interface{}); ok {
				for _, r := range rolesRaw {
					if rStr, ok := r.(string); ok {
						allGroups = append(allGroups, strings.ToLower(rStr))
					}
				}
			}
		}

		source := "simpeg" // default tendik
		codectx := ""
		role := "tendik"

		hasGroup := func(targets ...string) bool {
			for _, g := range allGroups {
				for _, t := range targets {
					if g == t {
						return true
					}
				}
			}
			return false
		}

		if hasGroup("adm_simonev", "admin", "superadmin", "adm_pusat") {
			role = "admin"
		} else if hasGroup("adm_simonev_fakultas") {
			role = "fakultas"
		} else if hasGroup("adm_simonev_prodi") {
			role = "prodi"
		}

		isDosen := false
		isMahasiswa := false
		for _, g := range allGroups {
			if g == "dosen" {
				isDosen = true
				break
			} else if g == "mahasiswa" || g == "mhs" {
				isMahasiswa = true
				break
			}
		}

		if isDosen {
			source = "simak"
			codectx = domainaccount.CtxDosen // "dosen"
			if role == "tendik" {
				role = "dosen"
			}
		} else if isMahasiswa {
			source = "simak"
			codectx = domainaccount.CtxMahasiswa // "mahasiswa"
			if role == "tendik" {
				role = "mahasiswa"
			}
		} else {
			// tendik / pegawai
			source = "simpeg"
			codectx = ""
		}

		c.Request().PostArgs().Set("source", source)
		c.Request().PostArgs().Set("resource", source)
		c.Request().PostArgs().Set("codectx", codectx)
		c.Request().PostArgs().Set("role", role)
		c.Locals("source", source)
		c.Locals("resource", source)
		c.Locals("codectx", codectx)
		c.Locals("role", role)
	} else {
		// Local login token
		if sid, ok := claims["sid"].(string); ok {
			c.Request().PostArgs().Set("sid", sid)
			c.Locals("sid", sid)
		}
		if resource, ok := claims["resource"].(string); ok {
			c.Request().PostArgs().Set("resource", resource)
			c.Locals("resource", resource)
		}
		if codectx, ok := claims["codectx"].(string); ok {
			c.Request().PostArgs().Set("codectx", codectx)
			c.Locals("codectx", codectx)
		} else if codetx, ok := claims["codetx"].(string); ok {
			c.Request().PostArgs().Set("codectx", codetx)
			c.Request().PostArgs().Set("codetx", codetx)
			c.Locals("codectx", codetx)
			c.Locals("codetx", codetx)
		}
		if role, ok := claims["role"].(string); ok {
			c.Request().PostArgs().Set("role", role)
		} else if level, ok := claims["level"].(string); ok {
			c.Request().PostArgs().Set("role", level)
		}
	}

	c.Locals("claims", claims)
	c.Request().PostArgs().Set("token", tokenStr)
}

func getTahun(c *fiber.Ctx) string {
	if t := c.Params("tahun"); t != "" {
		return t
	}
	return c.Query("ctxtahun")
}

func RBACMiddleware(whitelist []string, whoamiURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, err := extractBearerToken(c)
		if err != nil {
			return err
		}

		user, err := fetchWhoAmI(token, whoamiURL, c)
		if err != nil {
			return err
		}
		nidn := ""
		nip := ""
		npm := ""

		codeCtx := strings.ToLower(helper.NullableString(user.CodeCtx))
		resource := strings.ToLower(helper.NullableString(user.Resource))
		level := strings.ToLower(helper.NullableString(user.Level))

		if (codeCtx == domainaccount.CtxDosen || level == "dosen") && (resource == "simak" || resource == "") {
			nidn = helper.NullableString(&user.ID)
		}
		if (codeCtx == domainaccount.CtxMahasiswa || level == "mahasiswa") && (resource == "simak" || resource == "") {
			npm = helper.NullableString(&user.ID)
		}
		if resource == "simpeg" || level == "tendik" || level == "pegawai" {
			nip = helper.NullableString(&user.ID)
		}

		c.Request().PostArgs().Set("nidn", nidn)
		c.Request().PostArgs().Set("nip", nip)
		c.Request().PostArgs().Set("npm", npm)

		c.Request().PostArgs().Set("fakultas", helper.NullableString(user.RefFakultas))
		c.Request().PostArgs().Set("prodi", helper.NullableString(user.RefProdi))
		c.Request().PostArgs().Set("unit", helper.NullableString(user.Unit))
		c.Request().PostArgs().Set("nama", helper.NullableString(user.Name))
		c.Request().PostArgs().Set("level", helper.NullableString(user.Level))
		c.Request().PostArgs().Set("source", helper.NullableString(user.Resource))
		c.Request().PostArgs().Set("sid", user.ID)

		log.Println("[RBAC] Middleware passed, continue to handler")

		return c.Next()
	}
}

//
// ========================
// WHOAMI CALL
// ========================
//

func fetchWhoAmI(token, whoamiURL string, c *fiber.Ctx) (*Account, error) {
	req, err := http.NewRequest("GET", whoamiURL, nil)
	if err != nil {
		log.Printf("[RBAC] Failed to create request: %v", err)
		return nil, c.Status(500).
			JSON(commoninfra.NewResponseError(logCommonRbac, "Failed to create request: "+err.Error()))
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[RBAC] Failed to call whoami: %v", err)
		return nil, c.Status(500).
			JSON(commoninfra.NewResponseError(logCommonRbac, "Failed to call whoami: "+err.Error()))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[RBAC] Whoami response status: %d, body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return nil, handleWhoAmIError(body, c)
	}

	var user Account
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("[RBAC] Failed to parse whoami response: %v", err)
		return nil, c.Status(400).
			JSON(commoninfra.NewResponseError(logCommonRbac, "Failed to parse whoami response"))
	}

	log.Printf("[RBAC] Whoami user: %+v", user)
	return &user, nil
}

func handleWhoAmIError(body []byte, c *fiber.Ctx) error {
	var errResp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
		log.Printf("[RBAC] Whoami error code: %s, message: %s", errResp.Code, errResp.Message)
		return c.Status(400).
			JSON(commoninfra.NewResponseError(errResp.Code, errResp.Message))
	}

	log.Println("[RBAC] Whoami response not JSON or invalid format")
	return c.Status(401).
		JSON(commoninfra.NewResponseError(logCommonRbac, "Invalid format response"))
}

//
// ========================
// VALIDATION
// ========================
//

func isAdmin(user *Account) bool {
	return helper.NullableString(user.Level) == "admin"
}

func validateTahun(c *fiber.Ctx) (string, error) {
	tahun := getTahun(c)
	log.Printf("[RBAC] Tahun: %s", tahun)

	if tahun == "" {
		return "", c.Status(400).
			JSON(commoninfra.NewResponseError(logCommonRbac, "Query parameter 'tahun' is required"))
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil || tahunInt < 2000 {
		return "", c.Status(400).
			JSON(commoninfra.NewResponseError(logCommonRbac, "Query parameter 'tahun' invalid"))
	}

	return tahun, nil
}

//
// ========================
// ROLE CHECK
// ========================
//

func checkRoleAccess(user *Account, tahun string, whitelist []string) (bool, []string) {

	grantedAccess := []string{}

	// for _, r := range user.ExtraRole {
	// 	key := r.Tahun + "#" + r.Role
	// 	grantedAccess = append(grantedAccess, key)

	// 	if r.Tahun != tahun {
	// 		continue
	// 	}

	// 	role := strings.ToLower(strings.TrimSpace(r.Role))
	// 	if roleInWhitelist(role, whitelist) {
	// 		log.Printf("[RBAC] User has role '%s' for tahun %s, access granted", r.Role, r.Tahun)
	// 		return true, grantedAccess
	// 	}
	// }

	return false, grantedAccess
}

func roleInWhitelist(role string, whitelist []string) bool {
	for _, w := range whitelist {
		if role == strings.ToLower(w) {
			return true
		}
	}
	return false
}

func WSError(conn *websocket.Conn, code string, msg string) error {

	conn.WriteJSON(map[string]interface{}{
		"code":        code,
		"description": msg,
	})

	conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.ClosePolicyViolation,
			msg,
		),
		time.Now().Add(time.Second),
	)

	conn.Close()
	return errors.New(msg)
}

type WSSession struct {
	Token         string
	SID           string
	User          *Account
	GrantedAccess []string
}

// =======================
// HELPERS
// =======================

//	func domainAllowed(host string, allow []string) bool {
//		host = strings.ToLower(host)
//		for _, a := range allow {
//			if strings.HasSuffix(host, strings.ToLower(a)) {
//				return true
//			}
//		}
//		return false
//	}
func domainAllowed(host string, allow []string) bool {
	host = strings.ToLower(host)
	for _, a := range allow {
		u, err := url.Parse(a)
		var domain string
		if err == nil && u.Host != "" {
			domain = u.Hostname()
		} else {
			domain = strings.ToLower(a)
		}
		if strings.HasSuffix(host, domain) {
			return true
		}
	}
	return false
}

//	func extractHosts(s string) []string {
//		out := []string{}
//		words := strings.Fields(s)
//		for _, w := range words {
//			u, err := url.Parse(w)
//			if err == nil && u.Host != "" {
//				out = append(out, u.Hostname())
//				continue
//			}
//			m := hostExtract.FindStringSubmatch(w)
//			if len(m) > 1 {
//				out = append(out, m[1])
//			}
//		}
//		return out
//	}
func extractHostsFromText(s string) []string {
	out := []string{}

	// Split string berdasarkan spasi, koma, titik koma, newline
	words := regexp.MustCompile(`[ \t\r\n,;]+`).Split(s, -1)

	for _, w := range words {
		if w == "" {
			continue
		}

		// Hanya parsing kata yang terlihat seperti URL
		if strings.Contains(w, "://") {
			u, err := url.Parse(w)
			if err == nil && u.Host != "" {
				out = append(out, u.Hostname())
				continue
			}
		}

		// Fallback regex: cocokkan host sederhana (domain.tld)
		m := hostExtract.FindStringSubmatch(w)
		if len(m) > 1 {
			out = append(out, m[1])
		}
	}

	return out
}

func multiUnescape(s string, n int) string {
	cur := s
	for i := 0; i < n; i++ {
		u, err := url.QueryUnescape(cur)
		if err != nil || u == cur {
			break
		}
		cur = u
	}
	return cur
}

func ipInNets(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func SmartCompress() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ct := string(c.Response().Header.ContentType())

		// Jangan compress streaming
		if strings.Contains(ct, "text/event-stream") ||
			strings.Contains(ct, "application/x-ndjson") {
			return c.Next()
		}

		return compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		})(c)
	}
}
