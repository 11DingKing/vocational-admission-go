package domain

import (
	"strings"
	"unicode"
)

func ValidProvince(v string) bool {
	v = strings.TrimSpace(v)
	if len(v) < 2 || len(v) > 32 {
		return false
	}
	for _, r := range v {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func ValidStudentNo(v string) bool {
	if len(v) < 6 || len(v) > 24 {
		return false
	}
	for _, r := range v {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
			return false
		}
	}
	return true
}
func ValidYear(y int) bool     { return y >= 2020 && y <= 2100 }
func ValidCapacity(c int) bool { return c > 0 && c < 100000 }
func ValidatePlan(p AdmissionPlan) error {
	if !ValidYear(p.Year) || !ValidProvince(p.Province) || !ValidCapacity(p.TotalCapacity) {
		return ErrConflict
	}
	return nil
}
