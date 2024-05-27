package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/proto"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//1. 初始化logger
	initialize.InitLogger()

	//2. 初始化配置文件
	initialize.Viper()

	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//服务注册
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic("grpc dial error")
	}
	defer conn.Close()
	global.UserSrvClient = proto.NewUserClient(conn)

	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
