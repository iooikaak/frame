package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/iooikaak/frame/balance"
	_ "github.com/iooikaak/frame/balance/roundrobin"
	"github.com/iooikaak/frame/config"
	log "github.com/iooikaak/frame/log"

	"google.golang.org/grpc"
)

const root = "/services"

// TopoChange 拓扑变化通过的数据结构
// URI 拓扑在服务体系树中的位置;
// Conn 到新实例的连接;
// NodeHashKey 实例节点的HashKey;
// 标识当前变更是否是新增一个实点;
type TopoChange struct {
	URI         string
	Conn        *ConnActor
	NodeHashKey string
	NewNode     bool
}

var (
	errForgetSelfURI        = errors.New("Forget self uri ?")
	errConnectURIIsNil      = errors.New("Connect uri not found")
	errNotFoundConnect      = errors.New("Not found connection")
	errConnectRemoteAddrnil = errors.New("Connect remoteAddr is null")
	errNotSupportMultiURI   = errors.New("Not support multi uris")
)

// Discover 服务发现的类结构
type Discover struct {

	// 全局服务拓扑结构
	// 每个warhorse.com体系下的服务都会有一份同样的拓扑结构表
	// 并通过监听ETCD时间变化，跟新自己的拓扑结构
	// key 是接口的URI Value 是服务多个实例
	topology   map[string]*balance.BalanceMgr
	topoLocker sync.RWMutex

	// etcd client api
	kapi *clientv3.Client

	// 服务间相互调用时，首次调用会与对端建立TCP连接，并维护此连接
	// 下次再次调用时，直接则可直接冲连接池中拿到
	connholder map[string]*ConnActor
	connLocker sync.RWMutex

	// 服务唯一标识，注册到ETCD
	selfURI []string

	// 服务名称
	name string

	// your server
	service interface{}

	// middleware
	// 调用具体的方法前/后会执行
	// 如果middleware执行失败则会返回失败
	// 即使方法执行成功也会返回错误
	prepare []Middleware
	after   []Middleware

	// 当前服务所有的方法
	// 因为一个服务所有的方法在服务启动时是固定的
	// 在服务启动后不会变更，故此结构不会变更
	// 不需要加锁
	md map[string]*MethodDesc

	// 其他服务方法映射
	mtLocker sync.RWMutex
	mdtables map[string]*Medesc

	// 服务监TCP听地址
	localListenAddr string

	// 内部请求ID，服务内部唯一
	innerid int64

	ctx    context.Context
	cancel context.CancelFunc
}

var discoverOnce sync.Once
var instance *Discover

// Instance 返回GateSvr的单例对象
func Instance() *Discover {
	discoverOnce.Do(func() {
		instance = new(Discover)
		instance.md = make(map[string]*MethodDesc)
		instance.mdtables = make(map[string]*Medesc)
		instance.ctx, instance.cancel = context.WithCancel(context.TODO())
		instance.topology = make(map[string]*balance.BalanceMgr)
		instance.connholder = make(map[string]*ConnActor)
	})
	return instance
}

// BalancerMgr 获取某个实例集群管理者
func (discover *Discover) BalancerMgr(uri string) (*balance.BalanceMgr, error) {
	discover.topoLocker.RLock()
	mgr, ok := discover.topology[uri]
	discover.topoLocker.RUnlock()
	if !ok {
		return nil, errConnectURIIsNil
	}
	return mgr, nil
}

// Start 开启服务发现机制
func (discover *Discover) Start(srvName string, cfg *config.BaseCfg, selfURI []string, address string) *Discover {
	discover.selfURI = selfURI
	discover.name = srvName
	discover.localListenAddr = address
	//初始化etcd kapi对象
	if err := discover.readyEtcd(&cfg.Etcd); err != nil {
		panic(err.Error())
	}
	if len(selfURI) == 0 {
		panic(errForgetSelfURI.Error())
	}
	//注册自己
	if err := discover.selfRegist(); err != nil {
		panic(err.Error())
	}
	//持续监听新加入的或者删除的实例,维护在内存map类型里面
	go discover.discover()
	// 程序启动告警
	log.Infof("%s start up,local listen addr:%s,serve uri:%v", discover.name, discover.localListenAddr, discover.selfURI)
	return discover
}

func (discover *Discover) startListen(sd *ServiceDesc, ss interface{}, cfg *config.BaseCfg) {
	//组装注册服务对象
	discover.service = ss
	for i := range sd.Methods {
		d := &sd.Methods[i]
		discover.md[d.MethodName] = d
	}
	//生成一个tcp协议内网ip的listener
	var err error
	var listener net.Listener
	for {
		// 随机端口
		netIp := Netip()
		listener, err = net.Listen("tcp", netIp+":0")
		if err != nil {
			log.Warnf("warhorse.com:net.Listen(%s,%s) error, error is %s", "tcp", netIp+":0", err.Error())
			time.Sleep(time.Millisecond * 200)
			continue
		}
		break
	}
	//向ETCD注册信息(microservice1,配置文件,"/services/v1/microservice1",监听器生成的内网IP端口)
	discover.Start(sd.ServiceName, cfg, sd.ServiceURI, listener.Addr().String())
	//持续监听并处理到来的连接
	for {
		//接收进来的连接返回conn对象
		c, err := listener.Accept()
		if err != nil {
			log.Error("warhorse.com:", err.Error())
			continue
		}
		//初始化连接对象
		ca := ConnActor{c: c, reconn: false}
		ca.ctx, ca.cancel = context.WithCancel(context.TODO())
		ca.id = atomic.AddUint32(&connActorID, CA_BROKEN)
		ca.connType = passiveConnActor
		ca.p = &sync.Pool{
			New: func() interface{} {
				return new(icecontext)
			}}
		ca.initConnActor(c)
	}
}

// Get 获取URI对应的一个可用连接
func (discover *Discover) Get(URI string) (*ConnActor, error) {
	discover.topoLocker.RLock()
	if balancer, ok := discover.topology[URI]; ok {
		discover.topoLocker.RUnlock()
		tcpAddr, err := balancer.Pick(balance.RoundRobin())
		if err != nil {
			log.Error(err.Error())
			return nil, fmt.Errorf("load balance:%s", err.Error())
		}
		return discover.getConnActor(tcpAddr.String(), URI)
	}
	discover.topoLocker.RUnlock()
	return nil, fmt.Errorf("%s not found in topology", URI)
}

// Authorization 返回RPC方法的调用认证方法
func (discover *Discover) Authorization(path string) Authorization {
	ps := strings.Split(strings.ToLower(path), "/")
	if len(ps) <= 4 {
		return ""
	}
	mk := strings.Join(intercept(ps), "@")
	discover.mtLocker.RLock()
	if md := discover.mdtables[mk]; md == nil {
		discover.mtLocker.RUnlock()
		return ""
	} else {
		discover.mtLocker.RUnlock()
		return md.A
	}
}

func intercept(ss []string) []string {
	return ss[2:5]
}

// 从连接池中拿到远端连接句柄
func (discover *Discover) getConnActor(remoteAddr, uri string) (*ConnActor, error) {
	if len(remoteAddr) == 0 {
		return nil, errConnectRemoteAddrnil
	}
	if len(uri) == 0 {
		return nil, errConnectURIIsNil
	}
	// 找到了节点。取出/新建连接
	var connactor *ConnActor

	createConn := func() error {
		log.Debug("try ot connect:", remoteAddr)
		c, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		log.Debugf("connect backend serve %s success[%s]", remoteAddr, uri)
		connactor = NewActiveConnActor(c)
		discover.connLocker.Lock()
		discover.connholder[remoteAddr] = connactor
		discover.connLocker.Unlock()
		return nil
	}

	discover.connLocker.RLock()
	connactor, found := discover.connholder[remoteAddr]
	discover.connLocker.RUnlock()
	if !found {
		if err := createConn(); err != nil {
			return nil, err
		}
	} else if connactor.Status() == CA_ABANDON {
		if err := createConn(); err != nil {
			return nil, err
		}
	}
	return connactor, nil
}

func (discover *Discover) put(uri string) (<-chan *clientv3.LeaseKeepAliveResponse,
	*clientv3.LeaseGrantResponse, error) {

	resp, err := discover.kapi.Grant(context.TODO(), 21)
	if err != nil {
		return nil, nil, err
	}
	leaseResp, err := discover.kapi.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return nil, nil, err
	}

	svrURI := uri + "/provider/name"
	log.Debugf("set %s=%s", svrURI, discover.name)
	_, err = discover.kapi.Put(context.TODO(), svrURI, discover.name, clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, nil, err
	}

	// 先KeepAlive 在Put临时节点
	svrURI = uri + "/provider/instances/" + discover.localListenAddr
	log.Debugf("set %s=%s with leaseid=%x", svrURI, discover.localListenAddr, resp.ID)
	_, err = discover.kapi.Put(context.TODO(), svrURI, discover.localListenAddr, clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, nil, err
	}

	// 注册方法表
	for k, v := range discover.md {
		mdname := uri + "/" + strings.ToLower(v.MethodName) + "/provider/authorization/" + string(v.A)
		_, err := discover.kapi.Put(context.TODO(), mdname, k, clientv3.WithLease(resp.ID))
		if err != nil {
			return nil, nil, err
		}
	}
	return leaseResp, resp, nil
}

func (discover *Discover) selfRegist() error {
	if len(discover.selfURI) == 0 {
		return errForgetSelfURI
	}
	// 暂时不支持多个URI
	if len(discover.selfURI) > 1 {
		return errNotSupportMultiURI
	}
	for _, uri := range discover.selfURI {
		leaseResp, grantResp, err := discover.put(uri)
		if err != nil {
			return err
		}
		//维持心跳
		go func(uri string, leaseid clientv3.LeaseID) {
			t := time.NewTicker(time.Second * 10)
			svrURI := uri + "/provider/instances/" + discover.localListenAddr
			for {
				select {
				//如果etcd返回心跳
				case <-leaseResp:

				//每隔10s执行一次
				case <-t.C:
					//获取key为/provider/instances/内网IP的值
					gResp, err := discover.kapi.Get(context.TODO(), svrURI)
					//err!=nil或者获取对象的键值对长度为0
					if err != nil || len(gResp.Kvs) == 0 {
						//计入日志
						log.Fatalf("warhorse.com:%s svr uri %s get fail,detail=%v", discover.name, uri, err)
						//重新放置该服务的uri /services/v1/microservice1
						leaseResp, grantResp, err = discover.put(uri)
						if err != nil {
							log.Error(err.Error())
						}
					}
				//上下文对象死掉了，退出goroutine
				case <-discover.ctx.Done():
					return
				}
			}
		}(uri, grantResp.ID)
	}
	return nil
}

func (discover *Discover) readyEtcd(cfg *config.EtcdCfg) error {
	api, err := clientv3.New(clientv3.Config{
		Endpoints: cfg.EndPoints,
		Username:  cfg.User,
		Password:  cfg.Psw,

		DialOptions: []grpc.DialOption{
			grpc.WithTimeout(time.Second * 3),
			grpc.WithInsecure(),
		},

		DialTimeout: time.Second * cfg.Timeout,
	})
	if err != nil {
		return err
	}
	discover.kapi = api
	resp, err := discover.kapi.Get(context.TODO(), root, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, subNode := range resp.Kvs {
		if len(subNode.Key) == 0 || len(subNode.Value) == 0 {
			log.Warnf("warhorse.com:ready etcd key=%s value=%s", string(subNode.Key), string(subNode.Value))
		} else {
			discover.setTopo(string(subNode.Key), string(subNode.Value))
		}
	}
	return nil
}

// discover()持续监听集群中加入或者删除的实例，并维护内存map的结构
func (discover *Discover) discover() {
	//持续监听/services开头的key,过有这些key的值发生改变就发送一个channel
	ch := discover.kapi.Watch(context.TODO(), root, clientv3.WithPrefix())
	//持续监听channel信号
	for {
		select {
		// /services开头key的值发生了改变
		case notify := <-ch:
			if notify.Err() != nil {
				log.Warn("warhorse.com:", notify.Err())
				continue
			}
			//遍历[]clientV3.Event数组
			for _, event := range notify.Events {
				//获取key
				key := string(event.Kv.Key)
				//获取value
				value := string(event.Kv.Value)
				log.Debugf("warhorse.com:watch event:%s key:%s value:%s leasid:%x",
					event.Type.String(), key, value, event.Kv.Lease)
				//处理ectd监听到的事件类型
				switch event.Type {
				//有新的/services开头的key插入
				case clientv3.EventTypePut:
					//把键值对加入到etcd集群中
					discover.setTopo(key, value)
				// /services开头的key被删除
				case clientv3.EventTypeDelete:
					//集群中删除该键值对
					discover.rmTopo(key, value)
				}
			}
		//discover的上下文对象死掉了
		case <-discover.ctx.Done():
			//优雅关闭
			log.Infof("warhorse.com:dicover watch graceful exit.")
			return
		}
	}
}

func (discover *Discover) setTopo(key, value string) {
	segment := strings.Split(string(key), "/")
	var segl int
	if segl = len(segment); segl < 3 {
		return
	}
	if leafname := segment[segl-1]; leafname == "config" {
	} else if leafname == "name" {
	} else if segment[segl-2] == "instances" {
		interfaceURI := strings.Join(segment[:segl-3], "/")
		discover.regist(interfaceURI, value)
	} else if segment[segl-2] == "authorization" {
		discover.registMethod(key, value)
	} else if segment[segl-2] == "pprof" && discover.localListenAddr == value &&
		strings.ToLower(discover.name) == strings.ToLower(segment[3]) {
		go profile(segment[segl-1])
	}
}

func (discover *Discover) registMethod(mdkey, mdValue string) {
	if len(mdkey) < len(root) {
		return
	}

	ns := strings.Split(mdkey, "/")
	nsl := len(ns)
	if nsl < 4 {
		return
	}
	var md Medesc
	md.A = Authorization(ns[nsl-1])
	md.MdName = mdValue
	mk := strings.Join(intercept(ns), "@")
	discover.mtLocker.Lock()
	discover.mdtables[mk] = &md
	discover.mtLocker.Unlock()
}

func (discover *Discover) rmTopo(key, value string) {
	segment := strings.Split(string(key), "/")
	var l int
	if l = len(segment); l < 3 {
		return
	}

	if leafname := segment[l-1]; leafname == "config" {
		// TO DO
	} else if leafname == "name" {
		// TO DO
	} else if segment[l-2] == "instances" {
		interfaceURI := strings.Join(segment[:l-3], "/")
		log.Debug("rmTopo:", interfaceURI, " ", segment[l-1])
		discover.unRegist(interfaceURI, segment[l-1])
	} else if segment[l-2] == "scope" {
	} else if segment[l-2] == "authorization" {
		//服务都下线后移除
		discover.RemoveMdTables(key, 4)
		log.Debugf("rmTopo Check MdTables %s", key)
	}
}

// 注册一个后台服务接口
func (discover *Discover) regist(URI string, svrAddr string) {
	if len(URI) == 0 {
		return
	}

	// 过滤掉监听到自己的状态变化产生的通知
	if discover.localListenAddr == svrAddr {
		log.Debugf("warhorse.com:discover self node changed %s", svrAddr)
		return
	}

	discover.topoLocker.Lock()
	defer discover.topoLocker.Unlock()
	var blr *balance.BalanceMgr
	var found bool

	if blr, found = discover.topology[URI]; !found {
		blr = balance.Manager(discover.ctx)
		discover.topology[URI] = blr
		log.Debugf("Regist a new service at direction %s, the addr is %s", URI, svrAddr)
	}

	// 用后台服务的地址作为key来生成hash节点
	log.Debugf("AddNode: %s svrAddr:%s", URI, svrAddr)
	blr.Add(svrAddr)
}

// unRegist 注销一个后台服务接口 指定的实例hash节点，如果不指定则清空所有节点
func (discover *Discover) unRegist(URI string, svrAddr string) {
	if len(URI) == 0 {
		return
	}
	if len(svrAddr) > 0 {
		discover.topoLocker.Lock()
		defer discover.topoLocker.Unlock()
		if nodeList, found := discover.topology[URI]; found {
			remoteAddr := nodeList.Remove(svrAddr)
			log.Debugf("Remove backend serve %s, svrAddr %s remoteAddr %s.",
				URI, svrAddr, remoteAddr)
			// 清掉已经建立的连接
			if remoteAddr != "" {
				if connactor, found := discover.connholder[remoteAddr]; found {
					if connactor != nil {
						connactor.Close()
					}
					delete(discover.connholder, remoteAddr)
				}
			}
			if nodeList.Len() == 0 {
				log.Debugf("Remove backend topology %s", URI)
				delete(discover.topology, URI)
			}
		}
	} else {
		discover.topoLocker.Lock()
		defer discover.topoLocker.Unlock()
		if nodeList, found := discover.topology[URI]; found {
			log.Infof("Remove all backend serve %s.", URI)
			nodeList.Clear()
			delete(discover.topology, URI)

			// 清掉已经建立的连接
			for _, remoteAddr := range nodeList.AllNodeAddr() {
				if connactor, found := discover.connholder[remoteAddr]; found {
					if connactor != nil {
						connactor.Close()
					}
					delete(discover.connholder, remoteAddr)
				}
			}
		}
	}
}

func (discover *Discover) RemoveNode(url string, svrAddr string) {
	if url == "" || svrAddr == "" {
		return
	}

	discover.topoLocker.Lock()
	defer discover.topoLocker.Unlock()
	if nodeList, found := discover.topology[url]; found {
		remoteAddr := nodeList.Remove(svrAddr)
		log.Debugf("RemoveNode Remove backend serve %s, svrAddr %s remoteAddr %s.",
			url, svrAddr, remoteAddr)
		// 清掉已经建立的连接
		if remoteAddr != "" {
			if connactor, found := discover.connholder[remoteAddr]; found {
				if connactor != nil {
					connactor.Close()
				}
				delete(discover.connholder, remoteAddr)
			}
		}
		if nodeList.Len() == 0 {
			log.Debugf("RemoveNode Remove backend topology %s", url)
			delete(discover.topology, url)
		}
	}
}

func (discover *Discover) RemoveMdTables(url string, offset int) {
	ns := strings.Split(url, "/")
	nsl := len(ns)
	if nsl < 4 {
		return
	}
	interfaceURI := strings.Join(ns[:nsl-offset], "/")
	mk := strings.Join(intercept(ns), "@")

	if _, found := discover.topology[interfaceURI]; !found {
		log.Debugf("Remove mdtables key:%s, service:%s, path:%s", mk, interfaceURI, url)
		delete(discover.mdtables, mk)
	}
}

// Quit Quit
func (discover *Discover) quit() {
	// 停止Etcd Watch
	discover.cancel()
	// 先删除ETCD节点，再关闭连接，不然会出现ETCD节点丢失的情况
	for _, v := range discover.selfURI {
		uri := v + "/provider/instances/" + discover.localListenAddr
		_, err := discover.kapi.Delete(context.TODO(), uri)
		if err != nil {
			log.Errorf("discover uri %s Delete fail %s", uri, err.Error())
		}
		log.Debugf("warhorse.com:%s quit delete etcd key:%s", discover.name, uri)
	}
	discover.kapi.Close()
	discover.connLocker.RLock()
	defer discover.connLocker.RUnlock()
	for k, c := range discover.connholder {
		if c != nil {
			delete(discover.connholder, k)
			c.Close()
		}
	}
}
