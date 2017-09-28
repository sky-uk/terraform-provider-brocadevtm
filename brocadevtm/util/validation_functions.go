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

func ValidatePortNumber(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) < 1 || v.(int) > 65535 {
		errors = append(errors, fmt.Errorf("Port has to be between 1 and 65535", k))
	}
	return
}
