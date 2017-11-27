package util

import (
	"fmt"
	"regexp"
)

// ValidateUnsignedInteger : check integer is unsigned
func ValidateUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("[ERROR] %q can't be negative", k))
	}
	return
}

// ValidateIP : check valid IP address
func ValidateIP(v interface{}, k string) (ws []string, errors []error) {
	ip := v.(string)
	validateIP := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`)
	if !validateIP.MatchString(ip) {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be a valid IP. i.e 10.0.0.1", k))
	}
	return
}

// ValidatePortNumber : check port number is valid
func ValidatePortNumber(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 1 || v.(int) > 65535 {
		errors = append(errors, fmt.Errorf("[ERROR] Port has to be between 1 and 65535"))
	}
	return
}

// ValidateTCPPort : check TCP port is within range
func ValidateTCPPort(v interface{}, k string) (ws []string, errors []error) {
	port := v.(int)
	if port < 1 || port > 65535 {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be a valid port number in the range 1 to 65535", k))
	}
	return
}

// ValidateUDPSize : check UDP size is valid
func ValidateUDPSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 512 || v.(int) > 4096 {
		errors = append(errors, fmt.Errorf("[ERROR] %q must be a value within 512-4096", k))
	}
	return
}
