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

	return getRequest(getUrl.String())
}

func (n NodeProxy) GetAll(collection string) (json.RawMessage, error) {
	getUrl := &url.URL{
		Scheme: "http",
		Host:   n.HostAddr,
		Path:   "/doc",
	}
	query := getUrl.Query()
	query.Add("collection", collection)
	query.Add("all", "true")

	getUrl.RawQuery = query.Encode()

	return getRequest(getUrl.String())
}

func getRequest(url string) (json.RawMessage, error) {
	resp, err := http.Get(url)
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
