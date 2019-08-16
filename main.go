package main

import (
	"flag"
	"github.com/wailovet/osmanthuswine"
	"github.com/wailovet/osmanthuswine/src/core"
	"github.com/wailovet/pigo-web/face"
	"github.com/wailovet/pigo-web/pigoutil"
	"log"
)

func main() {

	var crossDomain string
	var port string
	var dev bool
	flag.StringVar(&port, "port", "80", "端口号:默认[80]")
	flag.StringVar(&crossDomain, "cross", "*", "跨域命令:默认[*]")
	flag.BoolVar(&dev, "dev", false, "开发模式,启用后动态读取本地静态文件")
	flag.Parse()

	core.SetConfig(&core.Config{
		Port:        port,
		CrossDomain: crossDomain,
		UpdatePath:  "pigo_update_linux",
	})

	data, _ := Asset("static/facefinder")
	pigoutil.InitCascade(data)

	//go-bindata-assetfs static/...
	if dev {
		log.Println("开发模式")
	} else {
		log.Println("生产模式")
		core.GetInstanceConfig().StaticFileSystem = assetFS()
	}
	core.GetInstanceRouterManage().Registered(&face.Detect{})
	osmanthuswine.Run()
}
