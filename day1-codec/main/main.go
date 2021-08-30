package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	day1_codec "qrpc/day1-codec"
	"qrpc/day1-codec/codec"
	"time"
)


func startServer(addr chan string){
	l,err:=net.Listen("tcp",":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	day1_codec.Accept(l)
}


func main()  {

	//1 创建连接

	addr :=make(chan string)
	go startServer(addr)

	conn,_:=net.Dial("tcp",<-addr)
	defer func() {_=conn.Close()}()

	time.Sleep(time.Second)
	//2 对发送内容编码

	json.NewEncoder(conn).Encode(day1_codec.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i:=0; i<5 ; i++{
		h:=&codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq: uint64(i),
		}
		_= cc.Write(h,fmt.Sprintf("grpc req %d",h.Seq))
		_= cc.ReadHeader(h)
		var reply string
		_= cc.ReadBody(&reply)
		log.Println("reply:",reply)
	}
	//3 对返回内容解码


}
