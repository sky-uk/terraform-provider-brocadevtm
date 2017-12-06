package util

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/schema"
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

// ReorderTablesInSection - Reorders the elements of a nested table to match the order within the state file.
func ReorderTablesInSection(mapToTraverse map[string]interface{}, tableNames map[string]string, sectionName string, d *schema.ResourceData) map[string]interface{} {
	for key, value := range mapToTraverse[sectionName].(map[string]interface{}) {
		tableUniqueKey, ok := tableNames[key]
		if ok {
			// We create a list of maps from the value
			valueAsListOfMaps := make([]map[string]interface{}, 0)
			for _, element := range value.([]interface{}) {
				valueAsListOfMaps = append(valueAsListOfMaps, element.(map[string]interface{}))
			}

			orderedTableMap := make([]map[string]interface{}, 0)

			//We Loop over the current key (value within tableNames) list within the given section of the resource in the state file
			for _, stateTableValue := range d.Get(sectionName + ".0." + key).([]interface{}) {
				//For each occurance of the key (value within tableNames) in the statefile, We Loop Over the list of that key within the given section of the response from the API
				for i, responseTableValue := range valueAsListOfMaps {
					// We compare the name of the key (value within tableNames) block in the state file to that of the API response
					if stateTableValue.(map[string]interface{})[tableUniqueKey] == responseTableValue[tableUniqueKey] {
						//We append the ifList with the correct value as per state file order
						orderedTableMap = append(orderedTableMap, responseTableValue)
						// We remove the value we just appended onto orderedTableMap from our valueAsListOfMaps we got from brocade
						valueAsListOfMaps = append(valueAsListOfMaps[:i], valueAsListOfMaps[i+1:]...)
					}
				}
			}
			orderedTableMap = append(orderedTableMap, valueAsListOfMaps...)
			// As the config in the statefile is a list of interfaces, we need to turn the list of maps into a list of interfaces
			mapSliceAsInterfaceSlice := make([]interface{}, len(orderedTableMap))

			for i, j := range orderedTableMap {
				mapSliceAsInterfaceSlice[i] = j
			}

			mapToTraverse[sectionName].(map[string]interface{})[key] = mapSliceAsInterfaceSlice
		}
	}
	return mapToTraverse[sectionName].(map[string]interface{})
}

// ReorderLists - reorders list taking from granted what already present in the state file
// basically this to emulate an unordered list data type and having T not complain for
// items being in the wrong order...
func ReorderLists(attr map[string]interface{}, attrName string, unorderedListNames map[string]bool, d *schema.ResourceData) map[string]interface{} {
	log.Println("[ReorderLists]: attr: \n", spew.Sdump(attr))
	// we traverse all attribute keys...
	for key, value := range attr {
		newValue := value
		attrProperName := attrName + "." + key
		log.Println("[ReorderLists]: attr proper name: ", attrProperName)
		log.Println("[ReorderLists]: key: ", key)
		// if key type is list of maps or list of interfaces and key is set to be an unordered list we
		// we take state value for granted...
		t := reflect.TypeOf(value)
		log.Println("[ReorderLists]: type of value: ", t.String())
		log.Println("[ReorderLists]: Value: ", spew.Sdump(value))
		log.Println("[ReorderLists]: T value: ", d.Get(attrProperName))
		switch t.String() {
		case "[]map[string]interface {}", "[]interface {}":
			if unorderedListNames[key] {
				log.Println("[ReorderLists] Key should be an unordered list! ", key)
				// we have to browse the list in the state file and the list from brocade and generate a new list to set into state
				// for each item in blist, bitem:
				//	for each item in slist, sitem:
				//		if sitem == bitem
				//			newList[spos] = sitem
				//			delete item in brocade list
				// newList = append(newList, blist)
				blist := value.([]interface{})
				slist := d.Get(attrProperName).([]interface{})
				remainingItems := make([]interface{}, 0)
				log.Println("[ReorderLists] blist: ", spew.Sdump(blist))
				log.Println("[ReorderLists] slist: ", spew.Sdump(slist))
				for bpos, bitem := range blist {
					log.Println("[ReorderLists] bpos: ", bpos)
					found := false
					for spos, sitem := range slist {
						log.Println("[ReorderLists] spos: ", spos)
						if bitem == sitem {
							found = true
							newValue.([]interface{})[spos] = sitem
						}
					}
					if found == false {
						remainingItems = append(remainingItems, bitem)
					}
				}
				newValue = append(newValue.([]interface{}), remainingItems)
				log.Println("[ReorderLists] NewValue: ", spew.Sdump(newValue))
			} else {
				for kk, vv := range value.([]interface{}) {
					if isMap, ok := vv.(map[string]interface{}); ok {
						index := strconv.Itoa(kk)
						attrPPName := attrProperName + "." + index
						log.Println("[ReorderLists]: attr proper name: ", attrPPName)
						newValue.([]interface{})[kk] = ReorderLists(isMap, attrPPName, unorderedListNames, d)
					}
				}
			}
		case "map[string]interface {}":
			log.Printf("key %s is a map, keep traversing...", key)
			newValue = ReorderLists(value.(map[string]interface{}), attrProperName, unorderedListNames, d)
		}
		// if item is a map, we keep traversing...
		attr[key] = newValue
	}
	return attr
}
