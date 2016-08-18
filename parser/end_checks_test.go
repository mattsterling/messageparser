package parser

import (
	"strings"
	"testing"
)

func TestEmojiStop(t *testing.T) {
	b := byte(')')
	result := stopForEmojiEnd(&b)
	if !result {
		t.Error("Stop for emoji ending with parenthesis should be true.")
	}
}

func TestNonWordStop(t *testing.T) {
	numbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	letters := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "r",
		"x", "z",
	}

	for _, v := range letters {
		char := []byte(v)
		result := stopForNonWord(&char[0])
		if result {
			t.Error("Lower case letter was determined to be a non word.")
		}

		char = []byte(strings.ToUpper(v))
		result = stopForNonWord(&char[0])
		if result {
			t.Error("Upper case letter was determined to be a non word.")
		}
	}

	for _, v := range numbers {
		char := []byte(v)
		result := stopForNonWord(&char[0])
		if result {
			t.Error("Number was determined to be a non word.")
		}
	}

	b := byte('_')
	result := stopForNonWord(&b)
	if result {
		t.Error("Underscore was marked as non word character.")
	}
	// Make sure a non word fails
	b = byte('/')
	result = stopForNonWord(&b)
	if !result {
		t.Error("Non word character was marked as word character.")
	}
}

func TestLinkStop(t *testing.T) {
	// &, ', (, ), *, +, ',' , -, . , / !, #, $ : , ; =, ?, @ [, ], _, `,
	symbols := []byte{
		'&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', '!',
		'#', '$', ':', ';', '=', '?', '@', '[', ']', '_', '`',
	}
	for _, v := range symbols {
		result := stopForLinkEnd(&v)
		if result {
			t.Error("Valid URL character was marked invalid, char: ", v)
		}
	}

	b := byte(' ')
	result := stopForLinkEnd(&b)
	if !result {
		t.Error("Invalid URL character was determined to be valid.")
	}
}
