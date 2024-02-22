package utils

import (
	"time"

	"github.com/pkg/errors"
)

// Retry 重试方法, attempts 重试次数, sleep 失败等待时间, fn 要执行的方法
func Retry(attempts int, sleep time.Duration, fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			err = errors.Wrap(err, "retry method error")
			return
		}
	}()

	if err = fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, sleep, fn)
		}
		return err
	}
	return nil
}
