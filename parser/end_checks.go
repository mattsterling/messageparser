package parser

// Check function to determine if section parsing
// should continue.
type stopCheck func(b *byte) bool

// Stop when you hit a closin parenthesis
func stopForEmojiEnd(b *byte) bool {
	return *b == 41 // )
}

// Checks a byte value to see if it's a non-word.
func stopForNonWord(b *byte) bool {
	// Non word checks cover Capital, Lowercase, and Numbers
	if (*b >= 48 && *b <= 57) || (*b >= 65 && *b <= 90) || (*b >= 97 && *b <= 122) {
		return false
	} else if *b == 95 {
		// Check the _
		return false
	}
	return true
}

// Checks for non-word allowing for
// special characters that could exist in a URL.
func stopForLinkEnd(b *byte) bool {
	// Checks for reserver characters first then alpha numeric
	v := *b
	// Kind of hate my life for this but whatevs
	// You can find these values here: http://ascii.cl/
	switch {
	// &, ', (, ), *, +, ',' , -, . , /
	case (v >= 38 && v <= 47):
		return false
	// !, #, $
	case (v == 33 || v == 35 || v == 36):
		return false
	// : , ;
	case (v == 58 || v == 59):
		return false
	// =, ?, @
	case (v == 61 || v == 63 || v == 64):
		return false
	// [, ], _, `,
	case (v == 91 || v == 93 || v == 95 || v == 96 || v == 126):
		return false
	default:
		if stopForNonWord(b) {
			return true
		}
		return false
	}
}
