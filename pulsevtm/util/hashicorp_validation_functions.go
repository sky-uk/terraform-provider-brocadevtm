package util

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

/*
THIS FILE WAS COPIED FROM https://github.com/hashicorp/terraform/blob/master/helper/validation/validation.go FOR TEMPOARY USE
*/

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func IntBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("[ERROR] expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}

// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least min (inclusive)
func IntAtLeast(min int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be int", k))
			return
		}

		if v < min {
			es = append(es, fmt.Errorf("[ERROR] expected %s to be at least (%d), got %d", k, min, v))
			return
		}

		return
	}
}

// IntAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at most max (inclusive)
func IntAtMost(max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be int", k))
			return
		}

		if v > max {
			es = append(es, fmt.Errorf("[ERROR] expected %s to be at most (%d), got %d", k, max, v))
			return
		}

		return
	}
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be string", k))
			return
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return
			}
		}

		es = append(es, fmt.Errorf("[ERROR] expected %s to be one of %v, got %s", k, valid, v))
		return
	}
}

// StringLenBetween returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length between min and max (inclusive)
func StringLenBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be string", k))
			return
		}
		if len(v) < min || len(v) > max {
			es = append(es, fmt.Errorf("[ERROR] expected length of %s to be in the range (%d - %d), got %s", k, min, max, v))
		}
		return
	}
}

// NoZeroValues is a SchemaValidateFunc which tests if the provided value is
// not a zero value. It's useful in situations where you want to catch
// explicit zero values on things like required fields during validation.
func NoZeroValues(i interface{}, k string) (s []string, es []error) {
	if reflect.ValueOf(i).Interface() == reflect.Zero(reflect.TypeOf(i)).Interface() {
		switch reflect.TypeOf(i).Kind() {
		case reflect.String:
			es = append(es, fmt.Errorf("[ERROR] %s must not be empty", k))
		case reflect.Int, reflect.Float64:
			es = append(es, fmt.Errorf("[ERROR] %s must not be zero", k))
		default:
			// this validator should only ever be applied to TypeString, TypeInt and TypeFloat
			panic(fmt.Errorf("[ERROR] can't use NoZeroValues with %T attribute %s", i, k))
		}
	}
	return
}

// CIDRNetwork returns a SchemaValidateFunc which tests if the provided value
// is of type string, is in valid CIDR network notation, and has significant bits between min and max (inclusive)
func CIDRNetwork(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("[ERROR] expected type of %s to be string", k))
			return
		}

		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			es = append(es, fmt.Errorf(
				"[ERROR] expected %s to contain a valid CIDR, got: %s with err: %s", k, v, err))
			return
		}

		if ipnet == nil || v != ipnet.String() {
			es = append(es, fmt.Errorf(
				"[ERROR] expected %s to contain a valid network CIDR, expected %s, got %s",
				k, ipnet, v))
		}

		sigbits, _ := ipnet.Mask.Size()
		if sigbits < min || sigbits > max {
			es = append(es, fmt.Errorf(
				"[ERROR] expected %q to contain a network CIDR with between %d and %d significant bits, got: %d",
				k, min, max, sigbits))
		}

		return
	}
}

// ValidateListUniqueStrings is a ValidateFunc that ensures a list has no
// duplicate items in it. It's useful for when a list is needed over a set
// because order matters, yet the items still need to be unique.
func ValidateListUniqueStrings(v interface{}, k string) (ws []string, errors []error) {
	for n1, v1 := range v.([]interface{}) {
		for n2, v2 := range v.([]interface{}) {
			if v1.(string) == v2.(string) && n1 != n2 {
				errors = append(errors, fmt.Errorf("[ERROR] %q: duplicate entry - %s", k, v1.(string)))
			}
		}
	}
	return
}

// ValidateRegexp returns a SchemaValidateFunc which tests to make sure the
// supplied string is a valid regular expression.
func ValidateRegexp(v interface{}, k string) (ws []string, errors []error) {
	if _, err := regexp.Compile(v.(string)); err != nil {
		errors = append(errors, fmt.Errorf("[ERROR] %q: %s", k, err))
	}
	return
}
