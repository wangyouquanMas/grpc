package day3

import (
	"errors"
	"strings"
	"sync"
)

//服务定义： 用结构体映射
type Server struct{
	serviceMap sync.Map
}

//服务注册
func (server *Server) Register(rcvr interface{}) error  {
	 s := newService(rcvr)
	 if _,dup := server.serviceMap.LoadOrStore(s.name,s); dup{
		 return errors.New("rpc: service already defined: " + s.name)
	 }
	 return nil
}


//服务发现
func (server *Server) findService(serviceMethod string) (svc *service,mtype *methodType,err error){
	//1 根据 服务名称 获取服务实例
	dot :=strings.LastIndex(serviceMethod,".")
	if dot<0{
		errors.New("rpc server:service/method request ill-formed:"+serviceMethod)
		return
	}

	serviceName,methodName := serviceMethod[:dot],serviceMethod[dot+1:]
	svci,ok :=server.serviceMap.Load(serviceName)
	if !ok{
		err = errors.New("rpc server: can't find service"+serviceName)
		return
	}
	//2 根据 方法名称 获取方法类型
	 svc = svci.(*service)
	 mtype = svc.method[methodName]
	 if mtype == nil{
	 	err = errors.New("rpc server: can't find method"+methodName)
	 }
	 return
}


