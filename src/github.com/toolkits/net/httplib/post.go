package httplib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PostJSON(url string, v interface{}) (response []byte, err error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return response, err
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(url, "application/json", bf)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
