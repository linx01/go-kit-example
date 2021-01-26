package utils

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
)

var client *api.Client
var err error
var serviceID string
var serviceName string
var ServicePort int

// init ... 导入时执行
func init() {
	// consul配置
	config := api.DefaultConfig()
	config.Address = "localhost:8500"
	// 注册
	client, err = api.NewClient(config)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// SetServiceNameAndPort ...
func SetServiceNameAndPort(name string, port int) {
	serviceName = name
	ServicePort = port
}

// RegisterService ...
func RegisterService() {

	// 注册信息
	uuID := uuid.NewV4()
	serviceID = serviceName + uuID.String() // 通过serviceid来区分相同服务的不同实例
	reg := api.AgentServiceRegistration{}
	reg.ID = serviceID
	reg.Name = serviceName
	reg.Address = "localhost"
	reg.Port = ServicePort
	reg.Tags = []string{"primary"}
	// 健康检查信息
	check := api.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = "http://172.31.50.219:" + strconv.Itoa(ServicePort) + "/health"

	// 将健康信息绑定至注册信息
	reg.Check = &check

	// 使用注册信息注册服务
	err := client.Agent().ServiceRegister(&reg)
	if err != nil {
		fmt.Println(err.Error())
	}

}

// UnregService ...
func UnregService() {
	client.Agent().ServiceDeregister(serviceID)
}
