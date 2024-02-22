package core

import (
	"reflect"
	"sync/atomic"
	"syscall"

	"github.com/iooikaak/frame/config"
	log "github.com/iooikaak/frame/log"
	"github.com/iooikaak/frame/protocol"
)

// RegisterAndServe 服务注册并开启监听
func RegisterAndServe(sd *ServiceDesc, ss interface{}, cfg *config.BaseCfg) {
	//返回sd的元素类型
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	//返回reflect.Type对象
	st := reflect.TypeOf(ss)
	//如果st没有继承ht
	if !st.Implements(ht) {
		log.Fatalf("warhorse.com: RegisterAndServe found the handler of type %v that does not implement %v", st, ht)
		return
	}
	//创建信号处理对象
	sh := NewSignalHandler()
	//优雅退出信号
	var h Singal
	sht := reflect.TypeOf((*Singal)(nil)).Elem()
	if !st.Implements(sht) {
		h = defaultServerSignal
	} else {
		h = ss.(Singal)
	}
	sh.Register(syscall.SIGTERM, h)
	sh.Register(syscall.SIGQUIT, h)
	sh.Register(syscall.SIGINT, h)
	sh.Start()
	//初始化frame.Discover对象
	s := Instance()
	s.startListen(sd, ss, cfg)
}

// GetInnerID 获取内部服务ID
func GetInnerID() int64 {
	return atomic.AddInt64(&Instance().innerid, 1)
}

// MeTables 获取方法 集合
func MeTables() map[string]Medesc {
	var mt = make(map[string]Medesc)
	Instance().mtLocker.RLock()
	for k, v := range Instance().mdtables {
		mt[k] = *v
	}
	Instance().mtLocker.RUnlock()
	return mt
}

// DeliverTo deliver request to anthor serve
func DeliverTo(task *protocol.Proto) (*protocol.Proto, error) {
	conn, err := Instance().Get(task.GetServeURI())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var b []byte
	//将结构体转化为字节数组
	if b, err = task.Serialize(); err != nil {
		return nil, err
	}
	resp, err := conn.RequestAndReponse(b, task.GetRequestID())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Nodes can get all instance of
func Nodes(uri string) []string {
	s := Instance()
	mgr, err := s.BalancerMgr(uri)
	if err != nil {
		return []string{}
	}
	return mgr.AllNodeAddr()
}

// Prepare 添加prepare middleware
func Prepare(mw ...Middleware) {
	Instance().prepare = append(Instance().prepare, mw...)
}

// After 添加after middleware
func After(mw ...Middleware) {
	Instance().after = append(Instance().after, mw...)
}
