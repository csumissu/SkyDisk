package main

import (
	"fmt"
	"github.com/csumissu/SkyDisk/conf"
	"github.com/csumissu/SkyDisk/graph"
	"github.com/csumissu/SkyDisk/util/logger"
)

func main() {
	logger.InitLogger(logger.LevelDebug)

	r := graph.InitRouters()
	err := r.Run(fmt.Sprintf(":%d", conf.ServerCfg.Port))
	if err != nil {
		logger.Fatal("cannot start the server, %s", err)
	}
}
