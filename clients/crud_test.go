package clients

import (
	"testing"
)

func TestGetURLSuccess(t *testing.T) {
	url := "http://www.google.com"
	_, err := Get(&url)
	if nil != err {
		t.Error("Error performing GET. This shouldn't have happened.")
	}
}

func TestGetURLError(t *testing.T) {
	url := "http://localhost:7777"
	_, err := Get(&url)
	if nil == err {
		t.Error("Expected error did not occur when performing GET")
	}
}
