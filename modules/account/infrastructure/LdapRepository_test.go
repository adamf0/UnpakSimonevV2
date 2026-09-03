package infrastructure

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMatchedGroup_PriorityOrder(t *testing.T) {
	priorityOrder := []string{
		"adm_simonev_prodi",
		"adm_simonev_fakultas",
		"adm_simonev",
		"adm_pusat",
	}

	allowedGroups := []string{
		"adm_simonev_prodi",
		"adm_simonev_fakultas",
		"adm_simonev",
		"adm_pusat",
		"admin",
		"superadmin",
	}

	getMatchedGroup := func(userGroups []string) string {
		var targets []string
		seenTarget := make(map[string]bool)

		for _, p := range priorityOrder {
			targets = append(targets, p)
			seenTarget[strings.ToLower(p)] = true
		}

		for _, grp := range allowedGroups {
			gClean := strings.TrimSpace(grp)
			gLower := strings.ToLower(gClean)
			if gLower != "" && !seenTarget[gLower] {
				targets = append(targets, gClean)
				seenTarget[gLower] = true
			}
		}

		for _, target := range targets {
			targetLower := strings.ToLower(strings.TrimSpace(target))
			for _, ug := range userGroups {
				ugClean := strings.TrimSpace(ug)
				ugLower := strings.ToLower(ugClean)
				if ugLower == targetLower || strings.HasPrefix(ugLower, targetLower) {
					return target
				}
			}
		}

		for _, ug := range userGroups {
			ugLower := strings.ToLower(strings.TrimSpace(ug))
			if strings.Contains(ugLower, "simonev") {
				return ug
			}
		}

		if len(allowedGroups) > 0 {
			return allowedGroups[0]
		}
		return ""
	}

	// Test 1: Siti Maimunah with adm_simonev_fakultas
	groups1 := []string{"adm_simonev_fakultas", "adm_tu", "adm_fakultas", "FEB", "Dosen"}
	assert.Equal(t, "adm_simonev_fakultas", getMatchedGroup(groups1))

	// Test 2: User with both prodi and fakultas groups
	groups2 := []string{"adm_simonev_fakultas", "adm_simonev_prodi", "FEB"}
	assert.Equal(t, "adm_simonev_prodi", getMatchedGroup(groups2))

	// Test 3: User with adm_simonev and adm_pusat
	groups3 := []string{"adm_simonev", "adm_pusat"}
	assert.Equal(t, "adm_simonev", getMatchedGroup(groups3))

	// Test 4: User with adm_pusat
	groups4 := []string{"adm_pusat", "staff"}
	assert.Equal(t, "adm_pusat", getMatchedGroup(groups4))
}
