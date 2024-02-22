package util

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
)

const (
	MESSAGE_TYPE_ONLINE_BROADCAST = int(1)
	MESSAGE_TYPE_GROUP            = int(2)
	MESSAGE_TYPE_PUBLIC           = int(3)
	MESSAGE_TYPE_PRIVATE          = int(4)
	MESSAGE_TYPE_CHATROOM         = int(5)
	MESSAGE_TYPE_NOTIFICATION     = int(9)
)

var MESSAGE_TYPES = map[int]int{
	MESSAGE_TYPE_ONLINE_BROADCAST: MESSAGE_TYPE_ONLINE_BROADCAST,
	MESSAGE_TYPE_GROUP:            MESSAGE_TYPE_GROUP,
	MESSAGE_TYPE_PUBLIC:           MESSAGE_TYPE_PUBLIC,
	MESSAGE_TYPE_PRIVATE:          MESSAGE_TYPE_PRIVATE,
	MESSAGE_TYPE_CHATROOM:         MESSAGE_TYPE_CHATROOM,
	MESSAGE_TYPE_NOTIFICATION:     MESSAGE_TYPE_NOTIFICATION,
}

type MessageThread struct {
	MessageType int
	MasterId    int64
	SlaveId     int64
}

func CreateMessageThread(messageType int, masterId, slaveId int64) *MessageThread {
	return &MessageThread{
		MessageType: messageType,
		MasterId:    masterId,
		SlaveId:     slaveId,
	}
}

func CreateMessageThreadFromString(thread string) *MessageThread {
	thread = strings.TrimSpace(thread)
	sections := strings.Split(thread, ":")
	if len(sections) != 3 {
		return CreateMessageThread(0, 0, 0)
	}
	messageType, _ := strconv.Atoi(sections[0])
	masterId, _ := strconv.ParseInt(sections[1], 10, 64)
	slaveId, _ := strconv.ParseInt(sections[2], 10, 64)
	return CreateMessageThread(messageType, masterId, slaveId)
}

func CreatePrivateMessageThread(from, to int64) *MessageThread {
	if from > to {
		return CreateMessageThread(MESSAGE_TYPE_PRIVATE, to, from)
	}
	return CreateMessageThread(MESSAGE_TYPE_PRIVATE, from, to)
}

func CreatePubMessageThread(pubId, from, to int64) *MessageThread {
	var masterId, slaveId int64
	if from == 0 {
		from = pubId
	}
	masterId = pubId
	if to > 0 {
		if from != masterId {
			slaveId = from
		} else {
			slaveId = to
		}
		if masterId == slaveId {
			slaveId = 0
		}
	} else {
		slaveId = 0
	}
	return CreateMessageThread(MESSAGE_TYPE_PUBLIC, masterId, slaveId)
}

func CreateGroupMessageThread(groupId, to int64) *MessageThread {
	return CreateMessageThread(MESSAGE_TYPE_GROUP, groupId, to)
}

func CreateSysMessageThread(sysId, to int64) *MessageThread {
	return CreateMessageThread(MESSAGE_TYPE_ONLINE_BROADCAST, sysId, to)
}

func CreateNotificationMessageThread(subtype int64) *MessageThread {
	return CreateMessageThread(MESSAGE_TYPE_NOTIFICATION, subtype, 0)
}

func CreateChatroomMessageThread(roomID int64) *MessageThread {
	return CreateMessageThread(MESSAGE_TYPE_CHATROOM, roomID, 0)
}

//master thread string, 2:1:0
func (thread *MessageThread) ThreadString() string {
	master := fmt.Sprintf("%d:%d", thread.MessageType, thread.MasterId)
	if thread.MessageType == MESSAGE_TYPE_PRIVATE {
		return fmt.Sprintf("%s:%d", master, thread.SlaveId)
	}
	return fmt.Sprintf("%s:0", master)
}

// sub thread string, 2:1:4
func (thread *MessageThread) String() string {
	return fmt.Sprintf("%d:%d:%d", thread.MessageType, thread.MasterId, thread.SlaveId)
}

func (thread *MessageThread) IsValid() bool {
	if _, ok := MESSAGE_TYPES[thread.MessageType]; ok && thread.MasterId > 0 && thread.SlaveId >= 0 {

		if thread.MessageType == MESSAGE_TYPE_PRIVATE && thread.MasterId > thread.SlaveId {
			return false
		}

		return true
	}
	return false
}

func (thread *MessageThread) Subscribe(redisConn redis.Conn, userId int64) {
	redisConn.Do("SADD", fmt.Sprintf("userthreads:%d", userId), thread.ThreadString())
	if thread.MessageType != MESSAGE_TYPE_PRIVATE {
		redisConn.Do("SADD", fmt.Sprintf("thread:%s", thread.ThreadString()), userId)
	}
}

func (thread *MessageThread) UnSubscribe(redisConn redis.Conn, userId int64) {
	redisConn.Do("SREM", fmt.Sprintf("userthreads:%d", userId), thread.ThreadString())
	if thread.MessageType != MESSAGE_TYPE_PRIVATE {
		redisConn.Do("SREM", fmt.Sprintf("thread:%s", thread.ThreadString()), userId)
	}
}

func (thread *MessageThread) IsSubscribed(redisConn redis.Conn, userId int64) (bool, error) {
	return redis.Bool(redisConn.Do("SISMEMBER", fmt.Sprintf("userthreads:%d", userId), thread.ThreadString()))
}
