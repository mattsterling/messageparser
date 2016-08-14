package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/messageparser/parser"
)

func ParseMessageHandler(w http.ResponseWriter, r *http.Request) {

	// Validate the request.
	if "POST" != r.Method {
		// TODO: 405
	}

	// In a future iteration is probably ideal to cap the message size
	if 0 == r.ContentLength {
		// TODO: 400

	} else if "text/plain" != r.Header.Get("Content-Type") {
		// Bad MIME type reject the request for unsupported type
		// TODO: 415
	}

	// Read in the chat messages. Here we assume messages are small.
	// Larger messages shouldn't read into memory like this and should
	// incrementally read/parse at the same time.
	defer r.Body.Close()
	buffer := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buffer.ReadFrom(r.Body)
	content := parser.ParseMessageContents(buffer)
	b, err := json.Marshal(*content)
	if nil != err {
		//TODO: 500 internal server error
	}
	fmt.Println("Bytes marshalled:", string(b))

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
