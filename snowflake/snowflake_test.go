package snowflake

import (
	"strconv"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	id, err := NewID()

	if err != nil {
		t.Error("生成ID失败：", err)
	}

	if id == 0 {
		t.Error("生成ID失败：", err)
	}
	t.Log(id, len(strconv.FormatInt(int64(id), 10)))

	id2, err := NewId()

	if err != nil {
		t.Error("生成ID失败：", err)
	}

	if id2 == 0 {
		t.Error("生成ID失败：", err)
	}
	t.Log(id2, len(strconv.FormatInt(int64(id2), 10)))
}
