package bot

import (
	"net/http"
)

func Exec(c *ServerConfig) error {
	return http.ListenAndServe(":80", nil)
}
