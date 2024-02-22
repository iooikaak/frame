package encode

import (
	"bytes"
)

// Pack  按消息类型打包消息
func Pack(dataType byte, data []byte) []byte {
	b := bytes.NewBuffer(nil)
	b.WriteByte(dataType)
	b.Write(data)
	return b.Bytes()
}

// Unpack 解包消息
func Unpack(data []byte) (byte, []byte) {
	return data[0], data[1:]
}
