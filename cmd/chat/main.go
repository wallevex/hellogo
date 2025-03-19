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
	"syscall"

	grpc_logging "hellogo/pkg/grpc-middleware/logging"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	go StartServices(conf)

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

func StartServices(conf *Config) {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_logging.PayloadUnaryServerInterceptor(log.Named("server-requests")),
		),
		grpc.MaxRecvMsgSize(math.MaxInt32-1),
		grpc.MaxSendMsgSize(math.MaxInt32-1),
	)
	chatImpl := chatsrv.NewChatServiceImpl()
	chatpb.RegisterChatServiceServer(grpcServer, chatImpl)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := chatpb.RegisterChatServiceHandlerServer(ctx, mux, chatImpl)
	if err != nil {
		panic(err)
	}

	h2s := &http2.Server{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Cookie,Grpc-Timeout,X-Grpc-Web")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}

		if r.ProtoMajor == 2 {
			grpcServer.ServeHTTP(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})
	log.Infof("chat server listen on %s", conf.Listen)
	if err := http.ListenAndServe(conf.Listen, h2c.NewHandler(handler, h2s)); err != nil {
		panic(err)
	}
}
