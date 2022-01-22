package client

import (
	"client/internal/config"
	"errors"
	"net/http"
)

func healthRequest() error {
	url := "http://" + config.GetApiAddr() + "/health"

	resp, err := http.DefaultClient.Get(url)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = errors.New("status code not OK")
	}

	return err

}
