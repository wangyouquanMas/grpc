package day3

import (
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














