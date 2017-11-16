package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
)

// AccTestCheckValueInKeyPattern : Function is used by acceptance tests where we want to check for a value in a set
func AccTestCheckValueInKeyPattern(resourceName string, keyPattern *regexp.Regexp, checkValue string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if ok {
			for attributeKey, attributeValue := range rs.Primary.Attributes {
				if keyPattern.MatchString(attributeKey) {
					if attributeValue == checkValue {
						return nil
					}
				}
			}
		}
		return fmt.Errorf("value %s not found in resource %s", checkValue, resourceName)
	}
}

// AccTestCreateRegexPatternForSet : generates a pattern for a list of items in a set
func AccTestCreateRegexPatternForSet(basePattern string) *regexp.Regexp {
	return regexp.MustCompile(basePattern + `\.[0-9]+`)
}

// AccTestCreateRegexPatternForSetItems : generates a pattern for a set with items
func AccTestCreateRegexPatternForSetItems(basePattern, appendPattern string) *regexp.Regexp {
	return regexp.MustCompile(basePattern + `\.[0-9]+\.` + appendPattern)
}

// AccTestCreateRegexPatternForNestedSets : generates a pattern for a set within a set
func AccTestCreateRegexPatternForNestedSets(basePattern, appendPattern string) *regexp.Regexp {
	return regexp.MustCompile(basePattern + `\.[0-9]+\.` + appendPattern + `\.[0-9]+`)
}
