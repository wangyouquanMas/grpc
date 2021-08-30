package codec

import "io"

/*
  1 根据RPC调用，抽象出请求和响应中的参数和返回值为body,剩余信息放在header中

	err = client.Call("Arith.Multiply", args, &reply)


  2 定义编码解码接口
 其中包含 编码、解码、关闭、


参考 ： 接口中方法分为4个部分，分别是

   编码 (write(*header, interface{}) error) 同时现实header和body编码]
   解码
      ReadHeader(*header) 仅实现header解码
      ReadBody(interface()) 仅实现body解码
   关闭
      io.closer 用于关闭文件流


    接口的作用就是多态

Q1 为什么编码解码要这样划分？ 为什么不直接使用 encode，decode方法？

   1 提供的api调用顺序就是 newDecoder -> Decoder -> decode 。
   所以自然的就需要先写 newDecoder函数，Decoder 结构体，decode方法



*/

type Header struct {
	//服务名和方法名，通常与 Go 语言中的结构体和方法相映射。
	ServiceMethod string
	//请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求。
	Seq uint64
	//是错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中。
	Error string
}

//抽象出对消息体进行编解码的接口 Codec，
//抽象出接口是为了实现不同的 Codec 实例
type Codec interface {
	//encode
	//decode
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{})error
	Write(*Header,interface{})error
}


type NewCodeFunc func(io.ReadWriteCloser) Codec

type Type string

const(
	JsonType Type = "application/json"
	GobType Type = "application/gob"
)

var NewCodeFuncMap map[Type]NewCodeFunc

func init()  {
	NewCodeFuncMap = make(map[Type]NewCodeFunc)
	//NewGobCodec 返回的GobCodec 实现了Codec接口
	NewCodeFuncMap[GobType] = NewGobCodec
}













