package util

import "github.com/hashicorp/terraform/helper/schema"

// BuildStringArrayFromInterface : take an interface and convert it into an array of strings
func BuildStringArrayFromInterface(strings interface{}) []string {
	stringList := make([]string, len(strings.([]interface{})))
	for idx, stringValue := range strings.([]interface{}) {
		stringList[idx] = stringValue.(string)
	}
	return stringList
}

// BuildStringSetFromInterface : take an interface and convert it into an array of strings
func BuildStringListFromSet(strings *schema.Set) []string {
	stringList := make([]string, 0)
	for _, stringValue := range strings.List() {
		stringList = append(stringList, stringValue.(string))
	}
	return stringList
}
