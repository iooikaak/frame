package unique_id

import (
	"errors"
	"sync"
	"time"

	"github.com/iooikaak/frame/util"
)

const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 8                                // bit length of sequence number
	BitLenMachineID = 53 - BitLenTime - BitLenSequence // bit length of machine id
	TimeUnit        = 1e8                              // 100微秒
)

type UniqueID struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

func NewUniqueID() *UniqueID {
	u := new(UniqueID)
	u.mutex = new(sync.Mutex)
	u.startTime = time.Date(2016, 3, 4, 0, 0, 0, 0, time.Local).UnixNano() / TimeUnit
	u.sequence = uint16(1<<BitLenSequence - 1)
	u.machineID = uint16(util.GetWorkerID())
	return u
}

func (u *UniqueID) Next() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)
	u.mutex.Lock()
	defer u.mutex.Unlock()
	current := time.Now().UnixNano()/TimeUnit - u.startTime
	if u.elapsedTime < current {
		u.elapsedTime = current
		u.sequence = 0
	} else {
		u.sequence = (u.sequence + 1) & maskSequence
		if u.sequence == 0 {
			u.elapsedTime++
			overTime := u.elapsedTime - current
			time.Sleep(time.Duration(overTime)*10*time.Millisecond - time.Duration(time.Now().Local().UnixNano()%TimeUnit)*time.Nanosecond)
		}
	}
	return u.toID()
}

func (u *UniqueID) toID() (uint64, error) {
	if u.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}
	return uint64(u.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(u.sequence)<<BitLenMachineID |
		uint64(u.machineID), nil
}
