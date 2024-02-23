package main

import (
	"flag"
	"flygoose/configs"
	"flygoose/datasource"
	"fmt"
)

func main() {
	//解析命令行参数
	configPath := flag.String("c", "cmd/admin/flygoose-config.yaml", "指定配置文件路径")
	flag.Parse()

	//根据启动指定的配置文件生成对应的 struct
	cfg, err := configs.NewConfig(*configPath)
	if err != nil {
		fmt.Errorf("生成配置文件出错. err: %w\n", err)
		return
	}
	//生成数据库,注意使用admin目录下的配置文件
	datasource.CheckAndCreateDatabase(cfg)
}
