package main

import (
	"github.com/tiny-sky/Tvote/core/dao/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	outPath := "core/dao/mysqlx/query"

	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		OutFile: outPath + "/query.go",
	})

	db, err := gorm.Open(mysql.Open("root:vote@(127.0.0.1:3306)/Vote?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("gorm err")
	}

	//复用已有的SQL连接配置db(*gorm.DB)
	g.UseDB(db)

	// 生成对应 CRUD 代码
	g.ApplyBasic(entity.Ticket{}, entity.User{})

	g.Execute()
}
