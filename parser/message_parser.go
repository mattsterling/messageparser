package parser

import(
	"bytes"
	"fmt"
	"strings"
)


// Link symbolizes a simple mapping of a URL and the related
// page title.
type Link struct {
	URL string `json: "url"`
	Title string `json: "title, omitempty"`
}

// MessageContent provides information about message.
type MessageContent struct {
	Mentions []string `json: "mentions, omitempty"`
	Emojis []string `json: "emoticons, omitempty"`
	Links []Link `json: "links, omitempty"`
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

		link := ParseSection(bites, start, space, -1, true)
		if "" != link {
			return &Link{URL: link}
		}
	}
	fmt.Println("Not a valid URL Link.")
	return nil
}

func ParseSection(data []byte, start *int, end byte, maxSize int, inclusive bool) string {
	fmt.Println(fmt.Sprintf("Parsing section from index %d ending with character '%s'", *start, string(end)))
	// Loop until we reach the end or we reach the end of the buffer.
	tmp := *start
	for tmp < len(data) {
		fmt.Println("ParseSection index:", tmp)
		if end == data[tmp] {
			fmt.Println("Found delimeter at index: ", tmp)
			break // We found our stopping point (Could be the end of the buffer)
		} else if tmp == len(data) {
			// The end of the buffer
			break;
		}
		tmp++ // Increment to the next

	}

	// Reject finding sections that are bigger than an allowable size.
	// Here we subtract 2 to not include the start index and ending index.
	size := (tmp - *start) - 2
	if -1 != maxSize && size > maxSize {
		fmt.Println(fmt.Sprintf("Section is too big between delimeters, will skip. Size = %d", size))
		*start = tmp // Skip outer loop from double checking indices we just touched
		return ""
	}



	// If the parse is not inclusive include exclude the start/stop delimiters
	word := ""
	if !inclusive {
		word = strings.TrimSpace(string(data[*start + 1: tmp]))

	} else {
		// Inclusive delimiter parse
		word = strings.TrimSpace(string(data[*start: tmp + 1]))
	}

	fmt.Println("Found section:", word)
	*start = tmp // Start where the slice ended for the outer loop
	return word


}



// Appends a message to a slice. A slice will be created
// if the one passed in is nil.
func appendString(s *[]string, message *string){
	if nil == s {
		fmt.Println("Appending string.")
		s = &[]string{*message}
		return
	}
	*s = append(*s, *message)
}


// Appends a Link to a slice. A slice will be created
// if the one passed in is nil.
func appendLink(s *[]Link, link *Link) {
	fmt.Println("Link slice", s)
	if nil == s {
		fmt.Println("Appending link.")
		s = &[]Link{*link}
		return
	}
	*s = append(*s, *link)
}


func ParseMessageContents(data *bytes.Buffer) *MessageContent {

	bites := data.Bytes()

	// Set up some slice references.
	var mentions []string
	var emojis []string
	var links []Link

	// N iteration loop
	fmt.Println("Buffer size: ", len(bites))
	for current := 0; current < len(bites); current++ {
		fmt.Println("Current is: ", current)
		b:= bites[current]
		fmt.Println("Current letter: ", string(b))
		switch {

		case mentionPrefix == b:
			m := ParseSection(bites, &current, space, -1, false)
			if "" != m {
				appendString(&mentions, &m)
			}
			continue

		case emojiStart == b:
			// Emojis cannot be longer than 15 (not including the '()' )
			e := ParseSection(bites, &current, emojiStop, 15, false)
			if "" != e {
				appendString(&emojis, &e)
			}
			continue

		case h == b:
			// We MAY be dealing with a URL
			l := parseUrl(bites, &current)
			if nil != l {
				// Tell the another go routine to process the link
				appendLink(&links, l)
			}
			continue

		default:
			// Keep moving forward
			continue
		}
	}

	fmt.Println("Pointers: ", mentions, emojis, links)
	return &MessageContent{Mentions: mentions, Emojis: emojis, Links: links}
}

