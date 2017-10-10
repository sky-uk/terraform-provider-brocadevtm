package api

import (
	"math/rand"
	"strconv"
	"time"
)

// SetTestResourceName : used by resource tests to generate a random unique name.
func SetTestResourceName(prefix string) string {
	rand.Seed(time.Now().Unix())
	return prefix + strconv.Itoa(rand.Int())
}
