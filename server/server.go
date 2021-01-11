package server

import (
	"time"

	"github.com/jake-hansen/followrs/config"
)

func Init(env string, startTime time.Time) {
	r := NewRouter(env, startTime)
	var address string = config.GetConfig().GetString("server.address")
	r.Run(address)
}
