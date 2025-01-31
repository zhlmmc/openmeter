package slicesx_test

import (
	"testing"

	"github.com/openmeterio/openmeter/pkg/slicesx"
	"github.com/stretchr/testify/assert"
)

func TestUniqueGroupBy(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var input []string
		result, ok := slicesx.UniqueGroupBy(input, func(s string) string {
			return s
		})
		assert.True(t, ok)
		assert.Empty(t, result)
	})

	t.Run("unique values", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		result, ok := slicesx.UniqueGroupBy(input, func(s string) string {
			return s
		})
		assert.True(t, ok)
		assert.Equal(t, map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		}, result)
	})

	t.Run("non-unique values", func(t *testing.T) {
		input := []string{"a", "a", "b"}
		result, ok := slicesx.UniqueGroupBy(input, func(s string) string {
			return s
		})
		assert.False(t, ok)
		assert.Nil(t, result)
	})

	t.Run("custom key function", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		input := []Person{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		result, ok := slicesx.UniqueGroupBy(input, func(p Person) int {
			return p.ID
		})
		assert.True(t, ok)
		assert.Equal(t, map[int]Person{
			1: {ID: 1, Name: "Alice"},
			2: {ID: 2, Name: "Bob"},
		}, result)
	})

	t.Run("custom key with duplicates", func(t *testing.T) {
		type Person struct {
			ID   int
			Name string
		}
		input := []Person{
			{ID: 1, Name: "Alice"},
			{ID: 1, Name: "Bob"},
		}
		result, ok := slicesx.UniqueGroupBy(input, func(p Person) int {
			return p.ID
		})
		assert.False(t, ok)
		assert.Nil(t, result)
	})
}
