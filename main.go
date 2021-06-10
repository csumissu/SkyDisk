package main

import (
	"fmt"

	"github.com/csumissu/SkyDisk/conf"
	"github.com/csumissu/SkyDisk/graph"
)

func main() {
	r := graph.InitRouters()
	r.Run(fmt.Sprintf(":%d", conf.ServerCfg.Port))
}
