package util

import (
	"fmt"
	"regexp"
)

// ValidateUnsignedInteger : check integer is unsigned
func ValidateUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

// ValidateIP : check valid IP address
func ValidateIP(v interface{}, k string) (ws []string, errors []error) {
	ip := v.(string)
	validateIP := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]$`)
	if !validateIP.MatchString(ip) {
		errors = append(errors, fmt.Errorf("Must be a valid IP. i.e 10.0.0.1"))
	}
	return
}
