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

func Test_IndexPhrases(t *testing.T) {
	assert.Equal(t, []int{0, 6, 10}, indexPhrases("olia\n\ntata"))
	assert.Equal(t, []int{0, 6, 12, 13}, indexPhrases("olia\n\ntata. o"))
	assert.Equal(t, []int{0, 11}, indexPhrases("olia\ntata.o"))
}

func Test_GetTextAt(t *testing.T) {
	s := Searcher{}
	s.PhraseIndexes = []int{0, 10, 20, 30, 40, 100}
	testF := func(max, at, expFrom, expTo int) {
		s.maxChars = max
		from, to := s.selectPos(at)
		assert.Equal(t, expFrom, from)
		assert.Equal(t, expTo, to)
	}
	testF(50, 5, 100, 100)
	testF(50, 4, 40, 100)
	testF(100, 4, 0, 100)
	testF(100, 1, 0, 100)

	testF(10, 0, 0, 10)
	testF(25, 0, 0, 20)

	testF(30, 1, 0, 30)
}
