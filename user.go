package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kilicmu/user-service/internal/database"

	"github.com/kilicmu/user-service/github.com/kilicmu/user-service"
	"github.com/kilicmu/user-service/internal/config"
	"github.com/kilicmu/user-service/internal/server"
	"github.com/kilicmu/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	godotenv.Load(".env.test")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_service.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	database.InitRedisUtils()
	database.InitDB()
	defer database.DestroyDB()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
