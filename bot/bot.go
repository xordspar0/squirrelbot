package bot

import (
	"fmt"
	"net/http"
)

func Start(c *ServerConfig) error {
	listenAddr := fmt.Sprintf("%d:%s", c.Port, c.Endpoint)
	return http.ListenAndServe(listenAddr, nil)
}
