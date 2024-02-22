package helper

import (
	"time"
)

const (
	BackOffDelayMin = 100 * time.Millisecond
	BackOffDelayMax = 10000 * time.Millisecond
)

func Sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}

func Retry(callback func() error, retry int, sleep time.Duration) (err error) {
	for i := 0; i < retry; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(sleep)
	}
	return
}

func Retry3(callback func() error) (err error) {
	for i := 0; i < 3; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}

func Retry3BackOff(callback func() error) (err error) {
	for i := 0; i < 3; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(backOff(i+1, BackOffDelayMin, BackOffDelayMax))
	}
	return
}

func Retry5BackOff(callback func() error) (err error) {
	for i := 0; i < 5; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(backOff(i+1, BackOffDelayMin, BackOffDelayMax))
	}
	return
}

func Retry10BackOff(callback func() error) (err error) {
	for i := 0; i < 10; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(backOff(i+1, BackOffDelayMin, BackOffDelayMax))
	}
	return
}

func backOff(attempt int, min time.Duration, max time.Duration) time.Duration {
	d := time.Duration(attempt*attempt) * min
	if d > max {
		d = max
	}
	return d
}
