package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"reflect"
)

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

// BuildReadMap : used by a read to build a list of maps which contain bools, strings, ints, float64s and lists of strings
func BuildReadMap(inputMap map[string]interface{}) (map[string]interface{}, error) {

	builtMap := make(map[string]interface{})

	for key, value := range inputMap {

		switch value.(type) {
		case bool:
			builtMap[key] = value.(bool)
		case string:
			builtMap[key] = value.(string)
		case float64:
			builtMap[key] = value.(float64)
		// []interface{} only configured / tested for a list of strings
		case []interface{}:
			builtMap[key] = schema.NewSet(schema.HashString, value.([]interface{}))
		default:
			return builtMap, fmt.Errorf("[ERROR] util.BuildReadListMaps doesn't understand type for %+v", value)
		}
	}
	return builtMap, nil
}

// GetSection : used to build a section in the schema into a map
func GetSection(d *schema.ResourceData, sectionName string, properties map[string]interface{}, keys []string) error {
	m, err := GetAttributesToMap(d, keys)
	if err != nil {
		log.Println("[ERROR] Error getting section ", sectionName, err)
		return err
	}
	properties[sectionName] = m
	return nil
}

// GetAttributesToMap : wrapper for d.Get
func GetAttributesToMap(d *schema.ResourceData, attributeNames []string) (map[string]interface{}, error) {

	m := make(map[string]interface{})

	for _, item := range attributeNames {
		v := d.Get(item)
		switch v.(type) {
		case bool:
			m[item] = v.(bool)
		case string:
			m[item] = v.(string)
		case int:
			m[item] = v.(int)
		case float64:
			m[item] = v.(float64)
		case []byte:
			m[item] = v.([]byte)
		case map[string]interface{}:
			m[item] = v.(map[string]interface{})
		case []map[string]interface{}:
			m[item] = v.([]map[string]interface{})
		case []interface{}:
			m[item] = v.([]interface{})
		case []string:
			m[item] = v.([]string)
		case *schema.Set:
			m[item] = v.(*schema.Set).List()
		default:
			return nil, fmt.Errorf("[ERROR] error, key %s of not valid type", item)
		}
	}
	return m, nil
}

// TraverseMapTypes - traverses the map fixing attr types accordingly
// Any *schema.Set attr is encoded into a list of maps
func TraverseMapTypes(m map[string]interface{}) {

	for attr := range m {
		t := reflect.TypeOf(m[attr])

		switch t.String() {
		case "*schema.Set":
			m[attr] = m[attr].(*schema.Set).List()
			for _, item := range m[attr].([]interface{}) {
				if v, ok := item.(map[string]interface{}); ok {
					TraverseMapTypes(v)
				}
			}
		case "[]interface {}":
			for _, item := range m[attr].([]interface{}) {
				if v, ok := item.(map[string]interface{}); ok {
					TraverseMapTypes(v)
				}
			}
		case "map[string]interface {}":
			TraverseMapTypes(m[attr].(map[string]interface{}))
		}
	}
}
