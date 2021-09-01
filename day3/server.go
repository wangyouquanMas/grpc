package day3

import (
	"errors"
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