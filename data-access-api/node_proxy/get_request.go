package nodeproxy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (n NodeProxy) Get(collection string, id string) (json.RawMessage, error) {

	getUrl := &url.URL{
		Scheme: "http",
		Host:   n.HostAddr,
		Path:   "/doc",
	}
	query := getUrl.Query()
	query.Add("collection", collection)
	query.Add("id", id)

	getUrl.RawQuery = query.Encode()

	resp, err := http.Get(getUrl.String())
	if err != nil {
		return nil, errors.New("get request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get return code: " + resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read the response: " + err.Error())
	}

	var responseBody json.RawMessage
	if err := json.Unmarshal(respBytes, &responseBody); err != nil {
		return nil, errors.New("failed to unmarshal the response: " + err.Error())
	}

	return responseBody, nil
}
