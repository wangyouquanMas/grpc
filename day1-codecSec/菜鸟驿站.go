package day1_codecSec

import "io"

//说明
//菜鸟驿站作为接口，包含了中通，圆通，申通等各家快递寄送方法
//面向接口  抽象  封装变化

//ExpressInfo: Header 用与信息校验
type ExpressInfo struct {
	Name 	 	string
	PhoneNum  	int
	Adress 		string
	TrackNum	 string //对应seq；用于区分不同的运单号
	UserEvalution string  //对应err: 用户响应，客户和服务端都有。
}

//Goods: body 作为内容实体
type Goods struct {

}


type Cainiao_Station interface{
	//装车:编码
	Pack(*ExpressInfo,interface{})error
	//卸货：解码
	UnPackTrackNum(interface{}) error
	UnPackGoods(*ExpressInfo) error
	io.Closer
}

type New_Cainiao_Station func(io.ReadWriteCloser)Cainiao_Station

var Cainiao_StationMap map[string]New_Cainiao_Station

func init()  {
	Cainiao_StationMap = make(map[string]New_Cainiao_Station)
	Cainiao_StationMap["ZHONGTONG"] = New_CainiaoStation
}


func New_CainiaoStation(io.ReadWriteCloser)Cainiao_Station{

}


