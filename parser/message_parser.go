package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/messageparser/clients"
	"golang.org/x/net/html"
)

var re = regexp.MustCompile("^[a-zA-Z0-9]*$")

// Link symbolizes a simple mapping of a URL and the related
// page title.
type Link struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// MessageContent provides information about message.
type MessageContent struct {
	Mentions []string `json:"mentions,omitempty"`
	Emojis   []string `json:"emoticons,omitempty"`
	Links    []Link   `json:"links,omitempty"`
}

// The parsing delimiters. Could be config driven.
const (
	mentionPrefix = byte('@')
	space         = byte(' ')
	emojiStart    = byte('(')
	emojiStop     = byte(')')
	h             = byte('h')
	urlStart      = "http"
)

// Attempts to parse a URL from byte array from a starting index
// Returns nil if a link was not found and the last index examined
func parseURL(bites []byte, start *int) *Link {
	fmt.Println("Attempting to parse URL.")
	end := *start + 4 // account for 'ttp' since start is the index of 'h'
	fmt.Println("URL PREFIX:", string(bites[*start:end]))
	if urlStart == string(bites[*start:end]) {

		link := ParseSection(bites, start, space, -1, true, false)
		if "" != link {
			return &Link{URL: link}
		}
	}
	fmt.Println("Not a valid URL Link.")
	return nil
}

// Utility to match the global regex.
func isAlphaNum(bites []byte) bool {
	return re.Match(bites)
}

// ParseSection will parse out a section of bytes from a given buffer starting with the given
// 'start' index and end at the first matching byte denoted by 'end'
// If you wish to restrict the size of a section returned provide a 'maxSize' or -1 for unlimited size.
// The 'inclusive' flag will ensure the start and end delimiters are part of the section returned.
// The 'an' flag will enforce the word is alphanumeric.
func ParseSection(data []byte, start *int, end byte, maxSize int, inclusive bool, an bool) string {
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
			break
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
	var b []byte

	if !inclusive {
		b = data[*start+1 : tmp]

	} else {
		// Inclusive delimiter parse
		b = data[*start : tmp+1]
	}

	// If this is the end of the string there is a null terminator at the
	// end potentially since we are using bytes.
	if tmp == len(data) {
		fmt.Println("Trimming off the null.")
		b = bytes.Trim(b, "\x00")
	}
	// Check to make sure the string is alpha numeric. (if flag set)
	word := ""
	if an {
		fmt.Println("Checking alpha numeric:", string(b))
		fmt.Println("Checking for Alpha numeric match.")
		if isAlphaNum(b) {
			fmt.Println("Alphanumeric filter failed. Will return empty string.")
			word = string(b)
		}
	} else {
		word = string(b)
	}

	fmt.Println("Parsed section:", word)
	*start = tmp // Start where the slice ended for the outer loop
	return strings.TrimSpace(word)

}

func getWebTitle(url *string) (string, error) {
	fmt.Println("Parsing url:", *url)
	r, err := clients.Get(url)
	if nil != err {
		fmt.Println("Could not retrieve url:", *url, err)
		return "", err
	}

	// Make sure the header starts with text/html
	if !strings.Contains(r.Header.Get("Content-Type"), "text/html") {
		fmt.Println("URL did not return valid HTML, title will not be found.")
		return "", nil
	}

	t := html.NewTokenizer(r.Body)

	// A large assumption is the web page being crawled is W3C compliant
	// So the title tag is near the top.
	// Submitting a page to this method that does not adhere to this assumption
	// will take longer to parse.
	for {
		token := t.Next()
		switch {
		case token == html.ErrorToken:
			return "", nil
		case token == html.StartTagToken && "title" == t.Token().Data:
			// We are at the title token, the next token will be the contents.
			t.Next()
			return t.Token().Data, nil
		}
	}
}

// Appends a message to a slice. A slice will be created
// if the one passed in is nil.
func appendString(s *[]string, message *string) {
	if nil == s {
		s = &[]string{*message}
		return
	}
	*s = append(*s, *message)
}

// Appends a Link to a slice. A slice will be created
// if the one passed in is nil.
func appendLink(s *[]Link, link *Link) {
	if nil == s {
		s = &[]Link{*link}
		return
	}
	*s = append(*s, *link)
}

// ParseMessageContents parses a string represnted by a byte Buffer
// for a given message.
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
		b := bites[current]
		switch {

		case mentionPrefix == b:
			m := ParseSection(bites, &current, space, -1, false, false)
			if "" != m {
				appendString(&mentions, &m)
			}
			continue

		case emojiStart == b:
			// Emojis cannot be longer than 15 (not including the '()' )
			e := ParseSection(bites, &current, emojiStop, 15, false, true)
			if "" != e {
				appendString(&emojis, &e)
			}
			continue

		case h == b:
			// We MAY be dealing with a URL
			l := parseURL(bites, &current)
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

	// Speed up the processing of web links.
	// This will suck when some cool individual decides to send a large of N of links in their message.
	//size := len(links)
	//ch := make(chan Link, size)
	var wg sync.WaitGroup
	wg.Add(len(links))
	for i := range links {
		go func(l *Link) {
			defer wg.Done() // Tell the wait group were done after this go routine.
			t, _ := getWebTitle(&l.URL)
			l.Title = t
		}(&links[i])
	}
	//close(ch)

	wg.Wait()

	return &MessageContent{Mentions: mentions, Emojis: emojis, Links: links}
}
