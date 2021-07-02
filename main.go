package main

import (
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/graph"
	"github.com/csumissu/SkyDisk/util"
)

func init() {
	util.InitLogger(util.LevelDebug)
}

func main() {
	r := graph.InitRouters()
	err := r.Run(fmt.Sprintf(":%d", config.ServerCfg.Port))
	if err != nil {
		util.Log().Panic("cannot start the server, %s", err)
	}
}
