package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"sample/services"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consul "github.com/go-kit/kit/sd/consul" // sd包创建第二个client，用于服务的注册与发现
	"github.com/go-kit/kit/sd/lb"            // 负载均衡器
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api" // 创建第一个client
)

var logger log.Logger

func main() {

	/*

		// 客户端直接访问

		// 1.创建一个直连client
		target, _ := url.Parse("http://172.31.50.219:8080")
		client := httptransport.NewClient("GET", target, services.GetUserInfoEncodeReqFunc, services.GetUserInfoDecodeResFunc)
		// 2.暴露出一个endpoint
		GetUser := client.Endpoint()
		// 3.创建contenxt与request
		ctx := context.Background()
		request := services.UserRequest{ID: 100, Method: "GET"}
		// 4.执行
		res, err := GetUser(ctx, request)
		if err != nil {
			fmt.Println(err.Error())
		}
		// 5.断言得到返回值
		res_, _ := res.(services.UserResponse)
		fmt.Println(res_.Result)

	*/

	// 服务的发现

	// 1.使用consulapi创建一个client

	config := consulapi.DefaultConfig()
	config.Address = "172.31.50.219:8500" // consul地址
	apiClient, _ := consulapi.NewClient(config)

	// 2.sd包创建第二个client，用于服务的注册与发现
	client := consul.NewClient(apiClient)

	// 3.创建一个instancer，查询服务实例状态信息，发现服务
	logger = log.NewLogfmtLogger(os.Stdout)
	tags := []string{"primary"}
	// 可实时查询服务状态的实例
	instancer := consul.NewInstancer(client, logger, "test1", tags, true) // 此处带入sd包创建的client
	factory := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
		tart, _ := url.Parse("http://" + serviceUrl) // 真实服务的地址 172.31.50.219:8080
		return httptransport.NewClient("GET", tart, services.GetUserInfoEncodeReqFunc, services.GetUserInfoDecodeResFunc).Endpoint(), nil, nil
	}
	endpointer := sd.NewEndpointer(instancer, factory, logger) // 关联
	// endpoints, _ := endpointer.Endpoints()
	// fmt.Println("服务有", len(endpoints), "条")

	// 发现并获取服务

	// GetUser := endpoints[0]

	// 使用kit包的负载均衡器
	// mylb := lb.NewRoundRobin(endpointer)                        // 轮询
	mylb := lb.NewRoundRobin(endpointer, time.Now().UnixNano()) // 随机
	for {
		GetUser, err := mylb.Endpoint()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		// 4.创建contenxt与request
		ctx := context.Background()
		request := services.UserRequest{ID: 100, Method: "GET"}
		// 5.执行
		res, err := GetUser(ctx, request)
		if err != nil {
			fmt.Println(err.Error())
		}
		// 6.断言得到返回值
		reS, _ := res.(services.UserResponse)
		fmt.Println(reS.Result)
		time.Sleep(time.Second * 3)
	}

}
