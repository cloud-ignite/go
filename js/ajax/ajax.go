//Provides Methods to make gopherjs ajax requests.
package ajax

import (
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

	json.NewDecoder(resp.Body).Decode(&v)

	return nil
}

func GetHTML(url string) (string, error) {

	client := &http.Client{}

	resp, err := client.Get(url)
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	return string(html[:]), err
}
