package nodeproxy

import (
	"errors"
	"net/http"
	"net/url"
)

func (n NodeProxy) Delete(collection string, id string) error {

	deleteUrl := &url.URL{
		Scheme: "http",
		Host:   n.HostAddr,
		Path:   "/doc",
	}
	query := deleteUrl.Query()
	query.Add("collection", collection)
	query.Add("id", id)

	deleteUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodDelete, deleteUrl.String(), nil)
	if err != nil {
		return errors.New("failed to create delete request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("delete request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("delete return code: " + resp.Status)
	}

	return nil
}
