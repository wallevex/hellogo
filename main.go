package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	// 加载服务端的自签名证书
	caCert, err := ioutil.ReadFile("./ssl/server.crt")
	if err != nil {
		panic(err)
	}

	// 创建 CA 池并添加自签名证书
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	// TLS 配置使用自定义 CA 池
	tlsConfig := &tls.Config{
		RootCAs: caPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// 请求
	resp, err := client.Get("https://192.168.1.24:2000/chat/users/100")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
