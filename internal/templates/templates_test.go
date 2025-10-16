package templates

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplatesBuilder_Process(t *testing.T) {
	t.Run("Test randNum", func(t *testing.T) {
		tb := NewTemplateBuilder()
		str, err := tb.Process(`{{randNum 50 100}}`)
		require.NoError(t, err)
		num, err := strconv.Atoi(str)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, num, 50)
		assert.LessOrEqual(t, num, 100)
	})
}
