package util

import "fmt"

// ValidateUnsignedInteger : check integer is unsigned
func ValidateUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

// ValidatePortNumber : check port number is valid
func ValidatePortNumber(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 1 || v.(int) > 65535 {
		errors = append(errors, fmt.Errorf("Port has to be between 1 and 65535"))
	}
	return
}

// ValidateUDPSize : check UDP size is valid
func ValidateUDPSize(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 512 || v.(int) > 4096 {
		errors = append(errors, fmt.Errorf("%q must be a value within 512-4096", k))
	}
	return
}
