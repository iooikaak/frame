package helper

import (
	"fmt"
	"testing"
	"time"
)

func TestRetry3(t *testing.T) {
	res := Retry3(func() error {
		return fmt.Errorf("nil")
	})
	fmt.Println(res)
}

func TestRetry5BackOff(t *testing.T) {
	res := Retry5BackOff(func() error {
		fmt.Println(time.Now())
		return fmt.Errorf("nil")
	})
	fmt.Println(res)
}
