package bot

import (
	"neolog.xyz/squirrelbot/config"

	"net/http"
)

func Exec(c *config.ServerConfig) error {
	return http.ListenAndServe(":80", nil)
}
