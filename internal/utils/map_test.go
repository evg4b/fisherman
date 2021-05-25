package utils_test

import (
	"fisherman/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMapToInterfaceMap(t *testing.T) {
	data := map[string]string{
		"1": "11111",
		"2": "11111",
	}

	actual := utils.StringMapToInterfaceMap(data)

	assert.IsType(t, map[string]interface{}{}, actual)
	assert.Equal(t, len(data), len(actual))
	for key, value := range data {
		assert.Equal(t, value, actual[key])
	}
}
