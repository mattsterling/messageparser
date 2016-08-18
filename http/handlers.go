package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	conf "github.com/messageparser/config"
	"github.com/messageparser/parser"
)

// ParseMessageHandler will accepted a HTTP requet with a
// text/plain body and return the marshalled JSON of the information
// that was parsed from the request.
func ParseMessageHandler(w http.ResponseWriter, r *http.Request) {

	// Validate the request.
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Provided operation is not allowed."))
		return
	}

	// In a future iteration is probably ideal to cap the message size
	if 0 == r.ContentLength {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty messages are not allowed."))
		return
	} else if conf.GlobalConfig.MsgSize != 0 && conf.GlobalConfig.MsgSize < r.ContentLength {
		// It's too big.
		w.WriteHeader(http.StatusBadRequest)
		format := "Message too large, must be smaller than %d chars."
		w.Write([]byte(fmt.Sprintf(format, conf.GlobalConfig.MsgSize)))
		return
	} else if "text/plain" != r.Header.Get("Content-Type") {
		// Bad MIME type reject the request for unsupported type
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("A raw message should be supplied via text/plain."))
		return
	}

	// Read in the chat messages. Here we assume messages are small.
	// Larger messages shouldn't read into memory like this and should
	// incrementally read/parse at the same time.
	defer r.Body.Close()
	buffer := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	buffer.ReadFrom(r.Body)

	// Get the message contents
	content := parser.ParseMessageContents(buffer)

	// Create the response
	b, err := json.Marshal(*content)
	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating JSON resp. Please contact the Admin."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// HealthCheck HTTP Handler.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if "GET" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed for Health Check"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "online"}`))
}
