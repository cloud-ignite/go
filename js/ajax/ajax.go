//Provides Methods to make gopherjs ajax requests.
package ajax

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetJSON(url string, v interface{}) error {

	client := &http.Client{}

	resp, err := client.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

func PostJSON(url string, data []byte, v interface{}) error {
	client := &http.Client{}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

func GetHTML(url string) (string, error) {

	client := &http.Client{}

	resp, err := client.Get(url)
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	return string(html[:]), err
}
