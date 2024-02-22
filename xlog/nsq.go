package xlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	MaxBlockLogNum   = 2000
	MaxLogBufNum     = 600
	MaxLogCommitNum  = 1300
	MaxLogByte       = 1 << 20 // 1024k
	MaxReportTimeOut = time.Second * 10
	MaxGoRoutineNum  = 60
)

type Nsq struct {
	maxGoroutineNum int64
	blockLogMsg     int64
	Topic           string
	Service         string
	producer        *nsq.Producer
	logMsg          chan []byte
	signalChan      chan struct{}
}

type NsqConfig struct {
	Addr    string `yaml:"addr"`
	Topic   string `yaml:"topic"`
	Service string `yaml:"service"`
	Conf    *nsq.Config
}

type nsqLookupNodes struct {
	Producers []nsqLookupProducer `json:"producers"`
}

type nsqLookupProducer struct {
	RemoteAddress    string   `json:"remote_address"`
	BroadcastAddress string   `json:"broadcast_address"`
	Hostname         string   `json:"hostname"`
	TCPPort          int      `json:"tcp_port"`
	HTTPPort         int      `json:"http_port"`
	Version          string   `json:"version"`
	Topics           []string `json:"topics"`
}

//NewNsq .
func NewNsq(conf *NsqConfig) (n *Nsq, err error) {
	addr, err := getNsqdAddr(conf.Addr)
	if err != nil {
		return nil, err
	}
	product, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	n = &Nsq{
		producer:        product,
		Topic:           conf.Topic,
		Service:         conf.Service,
		logMsg:          make(chan []byte, MaxLogBufNum),
		signalChan:      make(chan struct{}),
		maxGoroutineNum: 1,
	}

	//log producer
	go n.worker()

	//reg signal
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
		<-signalChan
		n.logMsg = nil      //close send
		close(n.signalChan) //batch close worker
	}()

	return
}

//Nsq 获取 nsqd tcp 连接地址
func getNsqdAddr(addr string) (tcpAddr string, err error) {
	lookupNodeAddr := fmt.Sprintf("%s/nodes", addr)
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(lookupNodeAddr)
	if err != nil {
		err = fmt.Errorf("NSQ：lookup request error：%s", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("NSQ：lookup response error：%d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("NSQ：lookup read body error：%s", err)
		return
	}
	nodes := &nsqLookupNodes{}
	err = json.Unmarshal(body, nodes)
	if err != nil {
		err = fmt.Errorf("NSQ：lookup unmarshal body error：%s", err)
		return
	}
	n := len(nodes.Producers)
	if n > 0 {
		t := time.Now().Unix()
		i := t % int64(n)
		tcpAddr = fmt.Sprintf("%s:%d", nodes.Producers[i].BroadcastAddress, nodes.Producers[i].TCPPort)
		return
	}
	err = errors.New("NSQ：nsqd not found")
	return
}

//Write implement write interface
/* old 方便对比压测
func (nsq *Nsq) Write(p []byte) (n int, err error) {
	err = nsq.producer.Publish(nsq.Topic, p)
	if err != nil {
		fmt.Println("nsq error ", err)
	}
	return len(p), err
}
*/
func (nsq *Nsq) Write(p []byte) (int, error) {

	select {
	case nsq.logMsg <- p:
	case <-nsq.signalChan: //signal close
		err := nsq.producer.Publish(nsq.Topic, p)
		if err != nil {
			fmt.Fprint(os.Stderr, string(p))
		}
	default:
		if n := atomic.LoadInt64(&nsq.maxGoroutineNum); n < MaxGoRoutineNum {
			atomic.AddInt64(&nsq.maxGoroutineNum, 1)
			go nsq.workerTemporary()
		} else {
			//Log blocking Accumulated {MaxBlockLogNum} times, report
			if n := atomic.LoadInt64(&nsq.blockLogMsg); n%MaxBlockLogNum == 0 {
				//stdout console
				fmt.Fprintf(os.Stderr, "[_maxLogNum_] log concurrency has reached the upper limit. current log info [%s]\n", string(p))
			}
			atomic.AddInt64(&nsq.blockLogMsg, 1)
		}

	}

	return len(p), nil
}

//worker
func (nsq *Nsq) worker() {
	var (
		ok          bool
		m           []byte
		msgBufBytes int64
		err         error
		timeOut     = time.NewTimer(MaxReportTimeOut)
		msgBuf      = make([][]byte, 0, MaxLogCommitNum)
	)
	defer timeOut.Stop()
	for {
		select {
		case m, ok = <-nsq.logMsg:
			if !ok {
				return
			}

			msgBuf = append(msgBuf, m)
			if len(msgBuf) >= MaxLogCommitNum || atomic.AddInt64(&msgBufBytes, int64(len(m))) >= MaxLogByte {
				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}

				msgBuf = make([][]byte, 0, MaxLogCommitNum)
				atomic.StoreInt64(&nsq.blockLogMsg, 0)
				atomic.StoreInt64(&msgBufBytes, 0)
				if !timeOut.Stop() {
					<-timeOut.C
				}
				timeOut.Reset(MaxReportTimeOut)

			}
		case <-timeOut.C:
			if len(msgBuf) > 0 {
				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}

				msgBuf = make([][]byte, 0, MaxLogCommitNum)
				atomic.StoreInt64(&nsq.blockLogMsg, 0)
				atomic.StoreInt64(&msgBufBytes, 0)
			}

			timeOut.Reset(MaxReportTimeOut)
		case <-nsq.signalChan: //close worker from signal
			if len(msgBuf) > 0 {
				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}
			}
			return
		}

	}

}

func (nsq *Nsq) workerTemporary() {
	var (
		ok          bool
		m           []byte
		msgBufBytes int64
		err         error
		msgBuf      = make([][]byte, 0, MaxLogCommitNum)
		timeOut     = time.NewTimer(MaxReportTimeOut + 5)
	)

	defer func() {
		timeOut.Stop()
		//归还maxGoroutineNum
		atomic.AddInt64(&nsq.maxGoroutineNum, -1)
	}()
	for {
		select {
		case m, ok = <-nsq.logMsg:
			if !ok {
				return
			}

			msgBuf = append(msgBuf, m)
			if len(msgBuf) >= MaxLogCommitNum || atomic.AddInt64(&msgBufBytes, int64(len(m))) >= MaxLogByte {

				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}
				msgBuf = make([][]byte, 0, MaxLogCommitNum)
				atomic.StoreInt64(&nsq.blockLogMsg, 0)
				atomic.StoreInt64(&msgBufBytes, 0)
				if !timeOut.Stop() {
					<-timeOut.C
				}
				timeOut.Reset(MaxReportTimeOut)
			}
		case <-nsq.signalChan: //close worker from signal
			if len(msgBuf) > 0 {
				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}
			}
			return
		case <-timeOut.C:
			if len(msgBuf) > 0 {
				err = nsq.producer.MultiPublishAsync(nsq.Topic, msgBuf, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "batch report log error, %v\n", err)
				}

				atomic.StoreInt64(&nsq.blockLogMsg, 0)
				atomic.StoreInt64(&msgBufBytes, 0)
			}
			return
		}

	}

}
