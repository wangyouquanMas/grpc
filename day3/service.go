package day3

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

type methodType struct{
	method reflect.Method
	ArgType reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}
//NumCalls(): 计算方法调用次数
func (m *methodType) NumCalls() uint64{
	return atomic.LoadUint64(&m.numCalls)
}
//
func (m *methodType) newArgv() reflect.Value{
	var argv reflect.Value
	//arg 可能是指针，也可能是值类型
	if m.ArgType.Kind() == reflect.Ptr{
		//reflect.New: 返回一个指向0值的指针
		argv = reflect.New(m.ArgType.Elem())
	}else{
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

func (m *methodType) newReplyv() reflect.Value{
	// reply 必须是指针
	//创建一个replyv指针实例
	replyv :=reflect.New(m.ReplyType.Elem())
	//根据方法中ReplyType中元素的基本类型来决定replyv元素类型
	switch m.ReplyType.Elem().Kind(){
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(),0,0))
	}
	return replyv
}
//定义服务结构体
type service struct{
	name string
	typ reflect.Type
	rcvr reflect.Value
	method map[string]*methodType
}

//newService(rcvr interface{}):实例化服务,初始化结构体成员
func newService(rcvr interface{}) *service{
	//给service分配内存空间
	s:=new(service)
	//结构体成员赋值
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	s.registerMethods()
	return s
}

func (s *service) registerMethods()  {
	s.method = make(map[string]*methodType)
	for i:=0;i<s.typ.NumMethod();i++{
		method := s.typ.Method(i)
		mType := method.Type
		if mType.NumIn() !=3 || mType.NumOut()!=1{
			 continue
		}
		if mType.Out(0)!=reflect.TypeOf(nil).Elem(){
			 continue
		}
		argType,replyType := mType.In(1),mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType){
			continue
		}
		s.method[method.Name] = &methodType{
			method:method,
			ArgType: argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server:register %s.%s\n",s.name,method.Name)
	}
}

func isExportedOrBuiltinType(t reflect.Type) bool{
	return ast.IsExported(t.Name()) || t.PkgPath() =="d"
}

















