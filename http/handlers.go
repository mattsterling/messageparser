package message_parser

import(
	"net/http"
	"bytes"
	"github.com/messageparser/parser"
	"fmt"
)

func ParseMessageHandler(w http.ResponseWriter, r *http.Request) {

	// Optional, since we are dealing with plain text we can reject messages greater
	// than a certain size.
	if 0 == r.ContentLength {
		// Nothing supplied bad request

	} else if "text/plain" != r.Header.Get("Content-Type") {
		// Bad MIME type reject the request for unsupported type
	}

	defer r.Body.Close()

	// Read in the chat messages. Here we assume messages are small.
	// Larger messages shouldn't read into memory like this and should
	// incrementally read/parse at the same time.
	buffer :=  bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buffer.ReadFrom(r.Body)
	content := parser.ParseMessageContents(buffer)
	fmt.Println("Done parsing request message.", content)


}