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

// BuildStringListFromSet : take an interface and convert it into an array of strings
func BuildStringListFromSet(strings *schema.Set) []string {
	stringList := make([]string, 0)
	for _, stringValue := range strings.List() {
		stringList = append(stringList, stringValue.(string))
	}
	return stringList
}

func AddBooleansToMap(d *schema.ResourceData, mapItem map[string]interface{}, boolOptions []string) map[string]interface{} {

	for _, item := range boolOptions {
		mapItem[item] = d.Get(item).(bool)
	}
	return mapItem
}

func AddIntegersToMap(d *schema.ResourceData, mapItem map[string]interface{}, integerOptions []string) map[string]interface{} {

	for _, item := range integerOptions {
		mapItem[item] = d.Get(item).(int)
	}
	return mapItem
}

func AddStringsToMap(d *schema.ResourceData, mapItem map[string]interface{}, stringOptions []string) map[string]interface{} {

	for _, item := range stringOptions {
		mapItem[item] = d.Get(item).(string)
	}
	return mapItem
}
