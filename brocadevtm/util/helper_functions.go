package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

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

// AddSimpleAttributeToMap : wrapper for d.Get
func AddSimpleAttributeToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributes []string) map[string]interface{} {

	for _, item := range attributes {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		attributeValue := d.Get(attributeName)
		switch attributeValue.(type) {
		case bool:
			mapItem[item] = attributeValue.(bool)
		case string:
			mapItem[item] = attributeValue.(string)
		case int:
			mapItem[item] = attributeValue.(int)
		default:
		}
	}
	return mapItem
}

// AddSimpleAttributeOkToMap : wrapper for d.GetOk
func AddSimpleAttributeOkToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributes []string) map[string]interface{} {

	for _, item := range attributes {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		if attributeValue, ok := d.GetOk(attributeName); ok {
			switch attributeValue.(type) {
			case bool:
				mapItem[item] = attributeValue.(bool)
			case string:
				mapItem[item] = attributeValue.(string)

			case int:
				mapItem[item] = attributeValue.(int)
			default:
			}
		}
	}
	return mapItem
}
