package nodeproxy

import (
	"errors"
	"net/http"
)

func (n NodeProxy) Probe() error {
	url := "http://" + n.HostAddr + "/health"

	resp, err := http.DefaultClient.Get(url)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = errors.New("status code not OK")
	}

	return err
}
