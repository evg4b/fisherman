package utils_test

import (
	"testing"

	. "github.com/evg4b/fisherman/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestStringMapToInterfaceMap(t *testing.T) {
	data := map[string]string{
		"1": "11111",
		"2": "11111",
	}

	actual := StringMapToInterfaceMap(data)

	assert.IsType(t, map[string]any{}, actual)
	assert.Equal(t, len(data), len(actual))
	for key, value := range data {
		assert.Equal(t, value, actual[key])
	}
}
