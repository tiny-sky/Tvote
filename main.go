package main

import (
	"context"
	"flag"

	"github.com/tiny-sky/Tvote/conf/file"
	"github.com/tiny-sky/Tvote/core"
	"github.com/tiny-sky/Tvote/core/server"
	"github.com/tiny-sky/Tvote/log"
)

var filepath = flag.String("f", "/conf.yml", "configuration file")

func main() {
	flag.Parse()
	c := file.NewFile(*filepath)

	settings, err := c.Load()

	if err != nil {
		log.Fatalf("%s", err)
	}
	settings.Init()

	var servers []core.Server

	// 生成票服务
	genticketSrv := server.New(settings.Server)
	servers = append(servers, genticketSrv)

	// TODO : 增加其他服务
	var opts []core.Option
	opts = append(opts, core.WithServers(servers...))

	newCore := core.New(opts...)
	if err := newCore.Run(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}
	log.Infof("server is stopped")
}
