package parser

import(
	"bytes"
	"fmt"
	"strings"
)


// Link symbolizes a simple mapping of a URL and the related
// page title.
type Link struct {
	URL string
	Title string
}

// MessageContent provides information about message.
type MessageContent struct {
	mentions []string
	emojis []string
	links []Link
}


// The parsing delimiters. Could be config driven.
const  (
	mentionPrefix = byte('@')
	space = byte(' ')
	emojiStart = byte('(')
	emojiStop = byte(')')
	h = byte('h')
	urlStart = "http"
)

// Attempts to parse a URL from byte array from a starting index
// Returns nil if a link was not found and the last index examined
func parseUrl(bites []byte, start *int) *Link {
	fmt.Println("Attempting to parse url.")
	end := *start + 4 // account for 'ttp' since start is the index of 'h'
	fmt.Println("URL PREFIX:", string(bites[*start: end]))
	if urlStart == string(bites[*start: end]) {

		link := ParseSection(bites, start, space, -1)
		if "" != link {
			return &Link{URL: link}
		}
	}
	fmt.Println("Not a valid URL Link.")
	return nil
}

func ParseSection(data []byte, start *int, end byte, maxSize int) string {
	fmt.Println(fmt.Sprintf("Parsing section from index %d ending with character '%s'", *start, string(end)))
	// Loop until we reach the end or we reach the end of the buffer.
	tmp := *start
	for tmp < len(data) {
		fmt.Println("ParseSection index:", tmp)
		if end == data[tmp] {
			fmt.Println("Found delimeter at index: ", tmp)
			tmp++
			break // We found our stopping point (Could be the end of the buffer)
		} else if tmp == len(data) {
			// The end is nigh
			break;
		}
		tmp++ // Increment to the next

	}

	// Reject finding sections that are bigger than an allowable size.
	// Here we subtract 2 to not include the start index and ending index.
	size := (tmp - *start) - 2
	if -1 != maxSize && size > maxSize {
		fmt.Println(fmt.Sprintf("Section is too big between delimeters, will skip. Size = %d", size))
		*start = tmp - 1 // Skip outer loop from double checking indices we just touched
		return ""
	}

	// Update the start index for the outer loop to continue where this one left off
	// and return the word as a string
	word := strings.TrimSpace(string(data[*start: tmp]))
	fmt.Println("Found section:", word)
	*start = tmp - 1 // Start where the slice ended
	return word

}


func ParseMessageContents(data *bytes.Buffer) *MessageContent {

	bites := data.Bytes()

	// Pre-allocate some slices for the metadata we find.
	// Could be memory inefficient, an area worth checking out.
	mentions := make([]string, 5)
	emojis := make([]string, 5)
	links := make([]Link, 5)

	// N iteration loop
	fmt.Println("Buffer size: ", len(bites))
	for current := 0; current < len(bites); current++ {
		fmt.Println("Current is: ", current)
		b:= bites[current]
		fmt.Println("Current letter: ", string(b))
		switch {

		case mentionPrefix == b:
			m := ParseSection(bites, &current, space, -1)
			if "" != m {
				mentions = append(mentions, m)
			}
			continue

		case emojiStart == b:
			// Emojis cannot be longer than 15 (not including the '()' )
			e := ParseSection(bites, &current, emojiStop, 15)
			if "" != e {
				emojis = append(emojis, e)
			}
			continue

		case h == b:
			// We MAY be dealing with a URL
			l := parseUrl(bites, &current)
			if nil != l {
				// Tell the another go routine to process the link
				links = append(links, *l)
			}
			continue

		default:
			// Keep moving forward
			continue
		}
	}

	return &MessageContent{mentions, emojis, links}
}

