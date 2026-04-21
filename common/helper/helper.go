package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

var (
	reUnpakEmail = regexp.MustCompile(
		`^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?@unpak\.ac\.id$`,
	)

	reUUIDv4 = regexp.MustCompile(
		`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`,
	)

	rePlus        = regexp.MustCompile(`\+`)
	reDoubleSep   = regexp.MustCompile(`(\.\.|__|--)`)
	reWhitespace  = regexp.MustCompile(`\s`)
	reURLEncoded  = regexp.MustCompile(`%[0-9A-Fa-f]{2}`)
	reURLEncoded2 = regexp.MustCompile(`%25[0-9A-Fa-f]{2}`)
	reNonASCII    = regexp.MustCompile(`[^\x20-\x7F]`)
)

func IsValidUnpakEmail(email string) bool {
	if len(email) > 254 { //[note] dalam pemantauan
		return false
	}

	// 1. Base pattern
	reg := regexp.MustCompile(reUnpakEmail.String())

	if !reg.MatchString(email) {
		return false
	}

	// 2. No plus (+)
	if regexp.MustCompile(rePlus.String()).MatchString(email) {
		return false
	}

	// 3. Double separator
	if regexp.MustCompile(reDoubleSep.String()).MatchString(email) {
		return false
	}

	// 4. No whitespace
	if regexp.MustCompile(reWhitespace.String()).MatchString(email) {
		return false
	}

	// 5. No URL-encoded chars
	if regexp.MustCompile(reURLEncoded.String()).MatchString(email) {
		return false
	}
	if regexp.MustCompile(reURLEncoded2.String()).MatchString(email) {
		return false
	}

	// 6. No non-ASCII
	if regexp.MustCompile(reNonASCII.String()).MatchString(email) {
		return false
	}

	return true
}

func ValidateJSONArray(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return validation.NewError("validation_json", "must be a JSON string")
	}

	var arr []interface{}
	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return validation.NewError("validation_json", "invalid JSON format")
	}

	if len(arr) == 0 {
		return validation.NewError("validation_json", "cannot be empty")
	}

	for i, v := range arr {

		// =====================================
		// CASE 1 : STRING UUID
		// =====================================
		if val, ok := v.(string); ok {

			if strings.TrimSpace(val) == "" {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("string at index %d cannot be empty", i),
				)
			}

			if _, err := uuid.Parse(val); err != nil {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("invalid uuid at index %d", i),
				)
			}

			continue
		}

		// =====================================
		// CASE 2 : OBJECT {uuid:"", freetext:""}
		// =====================================
		if obj, ok := v.(map[string]interface{}); ok {

			// key uuid wajib ada
			rawUUIDValue, existsUUID := obj["uuid"]
			if !existsUUID {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("missing uuid key at index %d", i),
				)
			}

			rawUUID, okUUID := rawUUIDValue.(string)
			if !okUUID || strings.TrimSpace(rawUUID) == "" {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("uuid required at index %d", i),
				)
			}

			if _, err := uuid.Parse(rawUUID); err != nil {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("invalid uuid at index %d", i),
				)
			}

			// key freetext wajib ada
			rawTextValue, existsText := obj["freetext"]
			if !existsText {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("missing freetext key at index %d", i),
				)
			}

			_, okText := rawTextValue.(string)
			if !okText {
				return validation.NewError(
					"validation_json",
					fmt.Sprintf("freetext must be string at index %d", i),
				)
			}

			continue
		}

		return validation.NewError(
			"validation_json",
			fmt.Sprintf("invalid type at index %d", i),
		)
	}

	return nil
}

func ValidateUnpakEmail(value interface{}) error {
	if value == nil {
		return fmt.Errorf("Email cannot be blank")
	}

	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("Email invalid type")
	}

	if !IsValidUnpakEmail(email) {
		return fmt.Errorf("Email is not valid unpak.ac.id")
	}

	return nil
}

func ValidateUUIDv4(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("UUID invalid type")
	}

	s = strings.TrimSpace(s)

	// Cek null padding ASCII ( \x00 )
	if strings.Contains(s, "\x00") {
		return fmt.Errorf("UUID contains invalid null padding")
	}

	if len(s) != 36 {
		return fmt.Errorf("UUID must be a valid UUIDv4 format")
	}

	// format regex UUID v4
	matched := regexp.MustCompile(reUUIDv4.String()).MatchString(s)
	if !matched {
		return fmt.Errorf("UUID must be a valid UUIDv4 format")
	}

	return nil
}

func ValidateFakultasUnit(value interface{}, level interface{}) error {
	levelStr, ok := level.(string)
	if !ok {
		return fmt.Errorf("level invalid type")
	}

	var s string
	switch v := value.(type) {
	case string:
		s = strings.TrimSpace(v)
	case *string:
		if v != nil {
			s = strings.TrimSpace(*v)
		}
	case nil:
		s = ""
	default:
		return fmt.Errorf("FakultasUnit invalid type")
	}

	if levelStr == "fakultas" && s == "" {
		return fmt.Errorf("FakultasUnit cannot be blank")
	}
	if (levelStr == "admin" || levelStr == "user") && s != "" {
		return fmt.Errorf("FakultasUnit required to be blank")
	}

	if s == "" {
		return nil
	}

	return ValidateUUIDv4(s)
}

func ValidateParent(value interface{}) error {
	var s string

	switch v := value.(type) {
	case string:
		s = strings.TrimSpace(v)

	case *string:
		if v != nil {
			s = strings.TrimSpace(*v)
		}

	case nil:
		return nil

	default:
		return fmt.Errorf("Parent invalid type")
	}

	if s == "" {
		return nil
	}

	return ValidateUUIDv4(value)
}

func ValidateLevel(value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return fmt.Errorf("level invalid type")
	}

	validLevels := map[string]struct{}{
		"admin":    {},
		"user":     {},
		"fakultas": {},
	}

	if _, exists := validLevels[val]; !exists {
		return fmt.Errorf("level not exist")
	}

	return nil
}

func ParseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			switch numErr.Err {
			case strconv.ErrRange:
				return 0, fmt.Errorf("Number out of range")
			case strconv.ErrSyntax:
				return 0, fmt.Errorf("Must be a number")
			}
		}
		return 0, fmt.Errorf("Invalid number")
	}
	return val, nil
}

func ParseUint(s string) (uint, error) {
	s = strings.TrimSpace(s)
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			switch numErr.Err {
			case strconv.ErrRange:
				return 0, fmt.Errorf("Number out of range")
			case strconv.ErrSyntax:
				return 0, fmt.Errorf("Must be a positive number")
			}
		}
		return 0, fmt.Errorf("Invalid number")
	}

	return uint(val), nil
}

func IsValidTugas(tugas string) bool {
	switch tugas {
	case
		"auditor1",
		"auditor2":
		return true
	}
	return false
}

func IsValidTypeGenerate(tipe string) bool {
	switch tipe {
	case
		"renstra",
		"dokumen_tambahan":
		return true
	}
	return false
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func StringHtmlValue(v *string) string {
	if v == nil || *v == "" {
		return "-"
	}
	return html.EscapeString(*v)
}
func Status(v *uint) string {
	if v == nil {
		return "-"
	}
	if *v == 1 {
		return "Ya"
	}
	return "Tidak"
}

func Contains(list []string, value string) bool {
	return slices.Contains(list, value)
}

func GrantedContains(audit []string, tahun string, granted string, isother bool) bool {
	entries := strings.Split(granted, ",")
	for _, e := range entries {
		parts := strings.Split(e, "#")
		if len(parts) != 2 {
			return false
		}

		year := strings.TrimSpace(parts[0])
		level := strings.TrimSpace(parts[1])

		if isother {
			if Contains(audit, level) {
				return true
			}
		} else {
			if year == tahun && Contains(audit, level) {
				return true
			}
		}
	}
	return false
}

func TimeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func StrPtr(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func UUIDString(u *uuid.UUID) string {
	if u == nil {
		return ""
	}
	return u.String()
}

func NullableString(u *string) string {
	if u == nil {
		return ""
	}
	return *u
}

func UintPtr(v uint) *uint { return &v }

func GenerateQRBase64(content string, size int) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

func CheckFilesExist(paths ...string) error {
	for _, f := range paths {
		if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
			return errors.New("file not found: " + f)
		}
	}
	return nil
}
