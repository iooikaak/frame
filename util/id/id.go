package id

import (
	"crypto/md5"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var staticMachine = getMachineHash()
var staticIncrement int64
var staticPid = int32(os.Getpid())

func getMachineHash() int32 {
	machineName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	buf := md5.Sum([]byte(machineName))
	return (int32(buf[0])<<0x10 + int32(buf[1])<<8) + int32(buf[2])
}

// GenerateID 生产陪玩订单ID
func GenerateID() string {
	timeStr := time.Now().Format("20060102150405")
	return timeStr + fmt.Sprint(staticMachine) +
		fmt.Sprint(staticPid) + fmt.Sprint(atomic.AddInt64(&staticIncrement, 1))
}

// GenerateCouponID 生成优惠券唯一码
func GenerateCouponID(biz string) string {
	return biz + "-" + time.Now().Format("20060102150405") +
		fmt.Sprint(atomic.AddInt64(&staticIncrement, 1))
}
