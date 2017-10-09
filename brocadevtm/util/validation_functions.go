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
		errors = append(errors, fmt.Errorf("%q must be a valid IP. i.e 10.0.0.1", k))
	}
	return
}

// ValidateTCPPort : check TCP port is within range
func ValidateTCPPort(v interface{}, k string) (ws []string, errors []error) {
	port := v.(int)
	if port < 1 || port > 65535 {
		errors = append(errors, fmt.Errorf("%q must be a valid port number in the range 1 to 65535", k))
	}
	return
}
