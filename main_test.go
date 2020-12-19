package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Highlight(t *testing.T) {
	assert.Equal(t, "", highlightText("", "olia"))
	assert.Equal(t, "oli", highlightText("oli", "olia"))
	assert.Equal(t, "<b>olia</b>", highlightText("olia", "olia"))
	assert.Equal(t, "<b>olia</b> hi <b>olia</b>", highlightText("olia hi olia", "olia"))
	assert.Equal(t, "oh <b>olia</b> hi", highlightText("oh olia hi", "olia"))
	assert.Equal(t, "oh<b>olia</b>hi", highlightText("oholiahi", "olia"))
}
