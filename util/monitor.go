package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/iooikaak/frame/cache"
)

type Monitor struct {
	redis_addr   string
	redis_pwd    string
	redis_key    string
	redispool    *cache.RedisPool
	log_module   string
	log_workerid string
	msg          chan *Msg
	maxRoutine   chan int
}

type Msg struct {
	msgType string
	data    string
}

var m *Monitor = nil

const MAX_BUF = 4096
const ROUTINE_NUM = 1
const WRITEBUF_TIMEOUT = 50

func NewMonitor(module string, workerId int64, redis_addr, redis_pwd, redis_key string) *Monitor {

	if m == nil {

		pool := cache.NewRedisPool(redis_addr, redis_pwd, 3, 300*time.Second)

		id := strconv.FormatInt(workerId, 10)
		m = &Monitor{redis_addr, redis_pwd, redis_key, pool, module, id, nil, nil}
		m.msg = make(chan *Msg, MAX_BUF)
		m.maxRoutine = make(chan int, ROUTINE_NUM)

		go func() {

			for {
				m.maxRoutine <- 1
				go uploadMonitor()
			}
		}()
	}

	return m
}

func (m *Monitor) Log2Monitor(level, data string) {

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		_, file = filepath.Split(file)
	}

	msgType := "log"
	curTime := time.Now().Format("2006/01/02 15:04:05.000000")
	//msg := fmt.Sprintf("{\"msgtype\":\"%s\",\"msgbody\":\"%s\"}", msgType, data)
	msg := data
	val := fmt.Sprintf("[%s:%s] %s %s:%d: [%s]:%s", m.log_module, m.log_workerid, curTime, file, line, level, msg)

	// m.msg <- &Msg{msgType, val}
	select {
	case m.msg <- &Msg{msgType, val}:

	case <-time.After(time.Millisecond * WRITEBUF_TIMEOUT):
		LogWarning("Log2Monitor channel timeout")
	}
}

func (m *Monitor) SendMonitor(msgType, data string) {

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		_, file = filepath.Split(file)
	}

	curTime := time.Now().Format("2006/01/02 15:04:05.000000")
	msg := fmt.Sprintf("{\"msgtype\":\"%s\",\"msgbody\":%s}", msgType, data)
	val := fmt.Sprintf("[%s:%s] %s %s:%d: [MONITOR]:%s", m.log_module, m.log_workerid, curTime, file, line, msg)

	// m.msg <- &Msg{msgType, val}
	select {
	case m.msg <- &Msg{msgType, val}:

	case <-time.After(time.Millisecond * WRITEBUF_TIMEOUT):
		LogWarning("SendMonitor channel timeout")
	}
}

func uploadMonitor() error {

	defer func() {
		if err := recover(); err != nil {
			LogError("uploadMonitor panic %v", err)
			LogError("crash stack: %v", string(debug.Stack()))
		}

		<-m.maxRoutine
	}()

	for {

		newMsg := <-m.msg
		msgType := newMsg.msgType
		data := newMsg.data

		if msgType == "close" && data == "quit" {
			break
		}

		c := m.redispool.Get()
		for i := 0; i < 3; i++ {
			_, err := c.Do("LPUSH", m.redis_key, data)
			if err == nil {
				break
			} else {
				c.Close()
				c = m.redispool.Get()
			}
		}

		c.Close()

	}

	return nil
}

func (m *Monitor) Close() {

	if m != nil {
		for i := 0; i < ROUTINE_NUM; i++ {
			m.SendMonitor("close", "quit")
		}
		m.redispool.Close()
	}
}
