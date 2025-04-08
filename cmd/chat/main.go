package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	grpc_logging "hellogo/pkg/grpc-middleware/logging"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	chatpb "hellogo/api/chat"
	chatsrv "hellogo/internal/chat/service"
	"hellogo/pkg/log"
)

var (
	c = flag.String("c", "../../configs/chat.json", "config file")
)

func main() {
	flag.Parse()

	buf, err := os.ReadFile(*c)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadFile", *c, err)
		os.Exit(1)
	}
	conf := &Config{}
	err = json.Unmarshal(buf, conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal", string(buf), err)
		os.Exit(3)
	}
	if conf.Logger.Dir == "" {
		conf.Logger.Dir = os.TempDir()
	}

	go StartServices(conf)
	go Debug(conf.Debug)

	if err := log.SetLogger("",
		conf.Logger.Dir,
		conf.Logger.File,
		conf.Logger.Count,
		conf.Logger.Size,
		conf.Logger.Unit,
		conf.Logger.Level,
		conf.Logger.Compress); err != nil {
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Warn("receive signal", (<-ch).String())
}

// TODO: 端口复用换成cmux
func StartServices(conf *Config) {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_logging.PayloadUnaryServerInterceptor(log.Named("server-requests")),
		),
		grpc.MaxRecvMsgSize(math.MaxInt32-1),
		grpc.MaxSendMsgSize(math.MaxInt32-1),
	)
	chatServer := chatsrv.NewChatServiceImpl()
	chatpb.RegisterChatServiceServer(grpcServer, chatServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	err := chatpb.RegisterChatServiceHandlerServer(ctx, gwMux, chatServer)
	if err != nil {
		panic(err)
	}

	// HTTP和gRPC服务端口复用
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Cookie,Grpc-Timeout,X-Grpc-Web")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}

		if r.ProtoMajor == 2 && strings.HasPrefix(
			r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r) // gRPC请求
		} else {
			gwMux.ServeHTTP(w, r) // HTTP请求
		}
	})
	log.Infof("chat server listen on %s", conf.Listen)
	if err := http.ListenAndServe(conf.Listen, handler); err != nil {
		panic(err)
	}
}
