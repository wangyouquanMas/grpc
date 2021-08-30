package day1_codec

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"qrpc/day1-codec/codec"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

//1 约定服务端与客户端协商内容： 消息的编码方式
type Option struct {
	MagicNumber int     //确保是qrpc请求
	CodecType  codec.Type  //可能选择了不同的编解码类型
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType: codec.GobType,
}


//2 rpc服务
type RpcServer struct {}

// 3 
func NewServer() *RpcServer {
	return &RpcServer{}
}

var DefaultServer = NewServer()

func (server *RpcServer) Accept(lis net.Listener){
	for	{
		conn,err:=lis.Accept()
		if err!=nil{
			log.Println("rpc  server:accept error:",err)
			return
		}
		go server.ServerConn(conn)
	}
}

func Accept(lis net.Listener){DefaultServer.Accept(lis)}


func (server *RpcServer) ServerConn(conn io.ReadWriteCloser){
	defer func() {_=conn.Close()}()

	var opt Option
	if err:=json.NewDecoder(conn).Decode(&opt); err!=nil{
		log.Println("rpc server:options error:",err)
		return
	}
	if opt.MagicNumber!=MagicNumber{
		log.Println("rpc server:invalid magic number %x",opt.MagicNumber)
		return
	}
	f:=codec.NewCodeFuncMap[opt.CodecType]
	if f==nil{
		log.Println("rpc server:invalid codec  type %s",opt.CodecType)
		return
	}
	//验证没有问题后，处理请求
	server.ServerCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

func (server *RpcServer) ServerCodec(cc codec.Codec){
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req,err:=server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			//Error客户端为空，服务端来赋值
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc,req,sending,wg)
	}
	wg.Wait()
	_=cc.Close()
}

// request stores all information of a call
type request struct {
	h	 *codec.Header
	argv,reply reflect.Value
}


// 解码header 信息，并返回
func (server *RpcServer) readRequestHeader(cc codec.Codec)(*codec.Header,error){
	var h codec.Header
	if err := cc.ReadHeader(&h);err!=nil{
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil,err
	}
	return &h,nil
}

//decode 请求
func (server *RpcServer) readRequest(cc codec.Codec)(*request,error){
	h,err:=server.readRequestHeader(cc)
	if err!=nil{
		return nil,err
	}

	req := &request{h:h}

	//请求参数这里被设置为 ""
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface());err!=nil{
		log.Println("rpc server: read argv err:",err)
	}
	return req,nil
}

func (server *RpcServer) sendResponse(cc codec.Codec,h *codec.Header,body interface{},sending *sync.Mutex){
	sending.Lock()
	defer sending.Unlock()
	if err:=cc.Write(h,body);err!=nil{
		log.Println("rpc server:write response error:",err)
	}
}

func (server *RpcServer) handleRequest(cc codec.Codec,req *request,sending *sync.Mutex, wg *sync.WaitGroup){
	defer wg.Done()
	log.Println(req.h,req.argv.Elem())
	req.reply =  reflect.ValueOf(fmt.Sprintf("qrpc resp %d",req.h.Seq))
	//server.sendResponse(cc,req.h,req.reply,sending)
	server.sendResponse(cc,req.h,req.reply.Interface(),sending)
}

