package utils

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	i := 3
	if err := Retry(3, time.Second*1, func() error {
		i--
		return errors.New("")
	}); err != nil {
		t.Fail()
	}

	if i != 0 {
		t.Logf("current i is %d", i)
		t.Error("retry function error")
	}

	if err := Retry(3, time.Second*1, func() error {
		i--
		return nil
	}); err != nil {
		t.Fail()
	}
}
