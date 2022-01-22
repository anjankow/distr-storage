package nodeproxy

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type NodeProxy struct {
	HostAddr string
	Logger   *zap.Logger
}

const (
	waitTimeout = 400 * time.Millisecond
)

func (n NodeProxy) WaitReady() {
	url := "http://" + n.HostAddr + "/health"

	for {
		time.Sleep(waitTimeout)
		resp, err := http.DefaultClient.Get(url)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			continue
		}

		break
	}

}
