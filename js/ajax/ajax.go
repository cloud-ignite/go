//Provides Methods to make gopherjs ajax requests.
package ajax

import (
	"encoding/json"
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
