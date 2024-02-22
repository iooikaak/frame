package mq

type Message struct {
	Body      []byte      //消息内容
	ID        string      //源ID
	Timestamp int64       //消息产生时间戮，纳秒
	Attempts  int         //推送次数
	Object    interface{} //原消息对象，视具体消息队列服务而定自行断言
}
