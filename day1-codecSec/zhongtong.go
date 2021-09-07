package day1_codecSec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

//中通要实现菜鸟驿站提供的服务

//ZhongTong: 实际干活人
type ZhongTong struct {
	//装车
	UnpackGoods *gob.Decoder
	//卸货
	PackGoods *gob.Encoder
	//实际承载物品
	Car  io.ReadWriteCloser
	//
	Buf *bufio.Writer
}


//zhogntong实现菜鸟提供的服务
func (ZhongTong *ZhongTong) Pack( expressInfo *ExpressInfo, goods interface{})(err error){
	if err:=ZhongTong.PackGoods.Encode(expressInfo);err!=nil{
		log.Println("rpc : gob error encoding header:", err)
		return
	}
	if err:=ZhongTong.PackGoods.Encode(goods);err!=nil{
		log.Println("rpc : gob error encoding header:", err)
		return
	}
	return
}

func (ZhongTong *ZhongTong)	UnPackTrackNum()(){

}


