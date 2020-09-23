package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LastChars(t *testing.T) {
	a := Word("Hello!")

	assert.Equal(t, "!", a.LastChars(1))
	assert.Equal(t, "o!", a.LastChars(2))
	assert.Equal(t, "Hello!", a.LastChars(100))
}

func Test_Chars(t *testing.T) {
	a := Word("12345")

	assert.Equal(t, "4", a.Chars(-2, -1))
	assert.Equal(t, "", a.Chars(-2, -100))
	assert.Equal(t, "1234", a.Chars(0, -1))
	assert.Equal(t, "12345", a.Chars(0, 1000))
}
