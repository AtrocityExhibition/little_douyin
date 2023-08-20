package main

import (
	"DouYin/repository"
	"DouYin/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// 如果启动有问题，大概是你的IP地址出现变化，需要在项目依赖的服务器中配置安全组
func main() {
	repository.Init()
	util.InitSFTP()
	util.InitFfmpeg()
	r := gin.Default()
	initRouter(r)
	pprof.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
