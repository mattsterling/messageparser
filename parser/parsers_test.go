package parser

import (
	"bytes"
	"testing"
)

var goodURL = "http://www.test.com"
var badURL = "htp://www.test.com"
var emoticon = "(dumplings)" // 12 AM and I am hungry

func TestParseURLEndOfString(t *testing.T) {
	b := []byte(goodURL)
	s := 0
	l := parseURL(b, &s)
	if l.URL != goodURL && l.Title != "" {
		t.Error("Valid URL Test failed. URL:", l.URL)
	}
}

func TestParseURLInString(t *testing.T) {
	b := []byte(goodURL)
	s := 0
	l := parseURL(b, &s)
	if l.URL != goodURL && l.Title != "" {
		t.Error("Valid URL Test failed. URL:", l.URL)
	}
}

func TestParseBadURL(t *testing.T) {
	b := []byte(badURL)
	s := 0
	l := parseURL(b, &s)
	if nil != l {
		t.Error("Invalid URL parsed unexpectedly. URL:", l.URL)
	}
}

func TestParseSectionNilSlice(t *testing.T) {
	var b []byte
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, false, false)
	if "" != section {
		t.Error("Nil section parsed to a non-empty string. Section:", section)
	}
}

func TestParseSectionInclusive(t *testing.T) {
	b := []byte(emoticon)
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, true, false)
	if emoticon != section {
		t.Error("Inclusive section parsed is malformed. Section:", section)
	}
}

func TestParseSectionNonInclusive(t *testing.T) {
	b := []byte(emoticon)
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, false, false)
	if "dumplings" != section {
		t.Error("Non-Inclusive section parse is malformed. Section:", section)
	}
}

func TestParseSectionIsAlphaNumeric(t *testing.T) {
	b := []byte(emoticon)
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, false, true)
	if "dumplings" != section {
		t.Error("AlphaNumeric parsed failed and shouldn't have. Section:", section)
	}
}

func TestParseSectionIsNotAlphaNumeric(t *testing.T) {
	b := []byte("($hi$@#)")
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, false, true)
	if "" != section {
		t.Error("AlphaNumeric parse passed unexpectantly. Section:", section)
	}
}

func TestParseSectionIsTooLong(t *testing.T) {
	b := []byte("(adfasdfasdfadfasdfadfaf)")
	start := 0
	section := ParseSection(b, &start, stopForNonWord, 15, true, true)
	if "" != section {
		t.Error("Length filter parse passed unexpectantly. Section:", section)
	}
}

func TestParseSectionSpaceTrimmed(t *testing.T) {
	b := []byte("thinger ")
	start := 0
	section := ParseSection(b, &start, stopForNonWord, -1, true, false)
	if "thinger" != section {
		t.Error("Expected trimmed string. Section:", section)
	}
}

func TestParseMessageHappyPath(t *testing.T) {
	data := bytes.NewBuffer([]byte("@chris you around? Good morning! (megusta) (coffee) Olympics are starting soon; http://www.nbcolympics.com"))
	c := ParseMessageContents(data)
	if "chris" != c.Mentions[0] {
		t.Error("Did not parse mentions correctly")
	}
	if "megusta" != c.Emojis[0] && "coffee" != c.Emojis[1] {
		t.Error("Did not parse emojis correctly")
	}

	if "http://www.nbcolympics.com" != c.Links[0].URL && "2016 Rio Olympic Games | NBC Olympics" != c.Links[0].Title {
		t.Error("Did not parse links correctly")
	}
}

func TestParseMessageNoLinks(t *testing.T) {
	data := bytes.NewBuffer([]byte("@chris you around? Good morning! (megusta) (coffee)"))
	c := ParseMessageContents(data)
	if "chris" != c.Mentions[0] {
		t.Error("Did not parse mentions correctly")
	}
	if "megusta" != c.Emojis[0] && "coffee" != c.Emojis[1] {
		t.Error("Did not parse emojis correctly")
	}
	if nil != c.Links {
		t.Error("Unexpected link slice found.")
	}
}

func TestParseMessageNoEmoticons(t *testing.T) {
	data := bytes.NewBuffer([]byte("@chris you around?"))
	c := ParseMessageContents(data)
	if "chris" != c.Mentions[0] {
		t.Error("Did not parse mentions correctly")
	}
	if nil != c.Emojis {
		t.Error("Unexpected slice found or emojis")
	}
	if nil != c.Links {
		t.Error("Unexpected link slice found.")
	}
}

func TestMessageNoData(t *testing.T) {
	data := bytes.NewBuffer([]byte("chris you around?"))
	c := ParseMessageContents(data)
	if nil != c.Mentions {
		t.Error("Unexpected slice found for mentions.")
	}
	if nil != c.Emojis {
		t.Error("Unexpected slice foundf or emojis")
	}
	if nil != c.Links {
		t.Error("Unexpected link slice found.")
	}
}

// Because I didn't get too fancy :(
func TestParseMessageEmbeddedLink(t *testing.T) {
	data := bytes.NewBuffer([]byte("http://www.nbcolympics.comadfasfhttp://www.nbcolympics.com"))
	c := ParseMessageContents(data)

	if nil != c.Links {
		t.Error("Somehow a link was magically parsed and it shouldn't have been.")
	}
}

func TestParseMessageNonAlphanumericEmoticon(t *testing.T) {
	data := bytes.NewBuffer([]byte("(adfasdf41$$$$)"))
	c := ParseMessageContents(data)
	if nil != c.Emojis {
		t.Error("Non alphnumeric link was returned, expected nil")
	}
}

func TestParseMessageEmoticonTooLong(t *testing.T) {
	data := bytes.NewBuffer([]byte("(123456789123456789adfa)"))
	c := ParseMessageContents(data)
	if nil != c.Emojis {
		t.Error("Long emoji was parsed but expected nil")
	}
}

func TestAppendLinkNil(t *testing.T) {
	var links []Link
	l := &Link{}
	appendLink(&links, l)
	if len(links) != 1 && nil == links {
		t.Error("Append link failed to create a new link slice.")
	}
}

func TestAppendLink(t *testing.T) {
	links := []Link{Link{}}
	l := &Link{}
	appendLink(&links, l)
	if len(links) != 2 && nil == links {
		t.Error("Append link failed to add a new link to the slice.")
	}
}

func TestAppendStringNil(t *testing.T) {
	var s []string
	data := "blah"
	appendString(&s, &data)
	if nil == s || len(s) != 1 {
		t.Error("Append string did not create a new slice.")
	}
}

func TestAppendString(t *testing.T) {
	s := []string{"blardy"}
	data := "blah"
	appendString(&s, &data)
	if len(s) != 2 {
		t.Error("Append string did not add string entry")
	}
}
