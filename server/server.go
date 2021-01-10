package server

import (
	"github.com/jake-hansen/followrs/config"
)

func Init(env string) {
	r := NewRouter(env)
	var address string = config.GetConfig().GetString("server.address")
	r.Run(address)
}
