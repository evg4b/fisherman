package utils_test

import (
	"fmt"
	"github.com/evg4b/fisherman/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	collection []int
	min        int
	max        int
}{
	{
		collection: []int{0},
		min:        0,
		max:        0,
	},
	{
		collection: []int{-1, 0},
		min:        -1,
		max:        0,
	},
	{
		collection: []int{1, 2, 3, 4, 5},
		min:        1,
		max:        5,
	},
	{
		collection: []int{5, 4, 3, 2, 1},
		min:        1,
		max:        5,
	},
	{
		collection: []int{0, -100, 100, -100, 100, 0},
		min:        -100,
		max:        100,
	},
}

func TestMin(t *testing.T) {
	for _, tt := range testCases {
		t.Run(fmt.Sprintf("returns %d for %v", tt.min, tt.collection), func(t *testing.T) {
			actual := utils.Min(tt.collection...)

			assert.Equal(t, tt.min, actual)
		})
	}

	t.Run("panic for where no arguments", func(t *testing.T) {
		assert.PanicsWithError(t, "min: no arguments", func() {
			_ = utils.Min()
		})
	})
}

func TestMax(t *testing.T) {
	for _, tt := range testCases {
		t.Run(fmt.Sprintf("returns %d for %v", tt.max, tt.collection), func(t *testing.T) {
			actual := utils.Max(tt.collection...)

			assert.Equal(t, tt.max, actual)
		})
	}

	t.Run("panic for where no arguments", func(t *testing.T) {
		assert.PanicsWithError(t, "max: no arguments", func() {
			_ = utils.Max()
		})
	})
}
