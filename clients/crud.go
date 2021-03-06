package clients

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Duration could be configuration driven
var client = &http.Client{
	Timeout: time.Duration(2 * time.Second),
}

// Get is a simple encapsulation of http.Get.
func Get(url *string) (*http.Response, error) {

	req, err := http.NewRequest("GET", *url, nil)
	if nil != err {
		log.Error("Error occured creating HTTP request for url: ", *url, err)
		return nil, err

	}

	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	return resp, nil
}
