package main

import (
	"example/services"
	"example/utils"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/time/rate"

	httptransport "github.com/go-kit/kit/transport/http" // 第三方
	mymux "github.com/gorilla/mux"                       // 第三方
)

func main() {

	name := flag.String("n", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")

	flag.Parse()

	if *name == "" {
		fmt.Println("服务名称未填写")
		os.Exit(0)
	}

	if *port == 0 {
		fmt.Println("服务端口未填写")
		os.Exit(0)
	}

	// 设置服务名与端口
	utils.SetServiceNameAndPort(*name, *port)

	s := &services.UserService{}
	// 限流
	l := rate.NewLimiter(1, 5)
	ep := services.RateLimitMiddleWare(l)(services.GenerateEndpoint(s))

	// 自定义错误
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(services.MyErrorEncoder),
	}
	serverHandler := httptransport.NewServer(ep, services.DecodeUserRequest, services.EncodeUserResponse, options...)
	r := mymux.NewRouter()
	// r.Handle(`/user/{uid:\d+}`, serverHandler)
	r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
	r.Methods("GET").Path(`/health`).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// 将运行和信号监听分成两个goroutine

	ss := make(chan error)

	// 运行
	go func() {
		// 将服务注册至consul
		utils.RegisterService()
		// 运行服务
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r)
		if err != nil {
			fmt.Println(err.Error())
			ss <- err
		}
	}()

	// 信号监听
	go func() {
		// 通道过滤系统信号
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		sign := <-signalChannel
		ss <- fmt.Errorf("%s", sign)
	}()

	syserr := <-ss
	utils.UnregService()
	fmt.Println(syserr)

}
