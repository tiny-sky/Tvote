package conf

import (
	"log"

	"github.com/tiny-sky/Tvote/core/dao"
	"github.com/tiny-sky/Tvote/core/dao/mysqlx"
	"github.com/tiny-sky/Tvote/core/dao/redisx"
	"github.com/tiny-sky/Tvote/core/server"

	"net/http"
	_ "net/http/pprof"
)

type (
	Settings struct {
		Mysql  mysqlx.Settings `yaml:"mysql"`
		Redis  redisx.Settings `yaml:"redis"`
		Server server.Settings `yaml:"server"`
	}
)

func (s *Settings) Init() {
	s.Mysql.Init()
	s.Redis.Init()

	// 性能分析
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	dao.NewDB()
}
