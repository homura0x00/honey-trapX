package main

import (
	"honey_back/honey_server/config"
	"honey_back/honey_server/global"
	"honey_back/honey_server/internal/models"

	"gorm.io/gen"
)

func main() {
	config.InitConfig()

	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/query",                        // gen代码的输出目录
		ModelPkgPath:  "./internal/models",                       // 模型代码的输出目录
		Mode:          gen.WithDefaultQuery | gen.WithoutContext, // 启用默认查询和链式接口
		FieldNullable: true,                                      // 允许 NULL 的字段生成指针类型
	})

	g.UseDB(global.MysqlDb)
	g.GenerateAllTable()        // 1. 通过sql表生成gorm模型
	g.ApplyBasic(models.User{}) // 2. 生成query代码
	g.Execute()                 // 执行
}
