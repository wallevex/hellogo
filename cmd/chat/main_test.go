package main

import (
	"context"
	"math"
	"net/http"
	"testing"

	grpc_logging "hellogo/pkg/grpc-middleware/logging"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	chatpb "hellogo/api/chat"
	chatsrv "hellogo/internal/chat/service"
	"hellogo/pkg/log"
)

const Listen = "0.0.0.0:2000"

func TestServices(t *testing.T) {
	if err := log.SetLogger("",
		"./log",
		"foobar.log",
		5,
		20,
		"MB",
		"info",
		1); err != nil {
		t.Fatal(err)
	}
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
		t.Fatal(err)
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
	t.Logf("chat server listen on %s", Listen)
	if err := http.ListenAndServe(Listen, h2c.NewHandler(handler, h2s)); err != nil {
		t.Fatal(err)
	}
}
