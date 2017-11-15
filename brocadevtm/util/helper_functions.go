package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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

// AddSimpleGetAttributesToMap : wrapper for d.Get
func AddSimpleGetAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		attributeValue := d.Get(attributeName)
		switch attributeValue.(type) {
		case bool:
			mapItem[item] = attributeValue.(bool)
		case string:
			mapItem[item] = attributeValue.(string)
		case int:
			mapItem[item] = attributeValue.(int)
		case *schema.Set:
			mapItem[item] = attributeValue.(*schema.Set).List()
		default:
		}
	}
	return mapItem
}

// AddSimpleGetOkAttributesToMap : wrapper for d.GetOk
func AddSimpleGetOkAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
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

// AddChangedSimpleAttributesToMap : wrapper for d.HasChange & d.Get
func AddChangedSimpleAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		if d.HasChange(attributeName) {
			attributeValue := d.Get(attributeName)

			switch attributeValue.(type) {
			case bool:
				mapItem[item] = attributeValue.(bool)
			case string:
				mapItem[item] = attributeValue.(string)
			case int:
				mapItem[item] = attributeValue.(int)
			case *schema.Set:
				mapItem[item] = attributeValue.(*schema.Set).List()
			default:
			}
		}
	}
	return mapItem
}

// SetSimpleAttributesFromMap : wrapper for d.Set
func SetSimpleAttributesFromMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		d.Set(attributeName, mapItem[item])
	}
}

// BuildListMaps : builds a list of maps from a list of interfaces using a list of attributes
func BuildListMaps(itemList *schema.Set, attributeNames []string) []map[string]interface{} {

	listOfMaps := make([]map[string]interface{}, 0)

	for _, item := range itemList.List() {

		definedItem := item.(map[string]interface{})
		newMap := make(map[string]interface{})

		for _, attributeName := range attributeNames {

			if attributeValue, ok := definedItem[attributeName]; ok {
				switch attributeValue.(type) {
				case bool:
					newMap[attributeName] = attributeValue.(bool)
				case string:
					newMap[attributeName] = attributeValue.(string)
				case int:
					newMap[attributeName] = attributeValue.(int)
				case []interface{}:
					newMap[attributeName] = attributeValue
				case *schema.Set:
					newMap[attributeName] = attributeValue.(*schema.Set).List()
				default:
				}
			}
		}
		listOfMaps = append(listOfMaps, newMap)
	}
	return listOfMaps
}

// SetListMaps : sets a list of maps in the state
func SetListMaps(mapItems []interface{}, attributeNames []string) []map[string]interface{} {

	sectionMapList := make([]map[string]interface{}, 0)

	for _, mapItem := range mapItems {

		definedMapItem := mapItem.(map[string]interface{})
		sectionMapListItem := make(map[string]interface{})

		for _, attributeName := range attributeNames {

			if mapItemValue, ok := definedMapItem[attributeName]; ok {
				switch mapItemValue.(type) {
				case bool:
					sectionMapListItem[attributeName] = mapItemValue.(bool)
				case string:
					sectionMapListItem[attributeName] = mapItemValue.(string)
				case int:
					sectionMapListItem[attributeName] = mapItemValue.(int)
				case []interface{}:
					sectionMapListItem[attributeName] = mapItemValue
				case *schema.Set:
					sectionMapListItem[attributeName] = mapItemValue.(*schema.Set).List()
				default:
				}
			} else {
				log.Printf(fmt.Sprintf("[DEBUG] mapItemValue not found"))
			}
		}
		sectionMapList = append(sectionMapList, sectionMapListItem)
	}
	return sectionMapList
}
