package icelog

import (
	"fmt"
	"testing"
	"time"
)

func TestConsoleLog(t *testing.T) {
	Debug("a=1", " b=2")
	Debugf("a=%d b=%d", 1, 2)
	Info("a=1", " b=2")
	Infof("a=%d b=%d", 1, 2)
	Warn("a=1", " b=2")
	Warnf("a=%d b=%d", 1, 2)

	Error("a=1", " b=2")
	Errorf("a=%d b=%d", 1, 2)

	Fatal("a=1", " b=2")
	Fatalf("a=%d b=%d", 1, 2)
}
func TestError(t *testing.T) {
	SetLevel("Warn")
	Error(fmt.Sprintf("%s-%v", hostname, time.Now().UnixNano()/1e6))
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("Hello Test Zap logger.")
		Info("a", "b", "c")
		Infof("A=%d B=%d C=%d", 48, 49, 50)
	}
}

// func TestArgs(t *testing.T) {
// 	t.Log(os.Args[0])
// }

// func TestFileWriter(t *testing.T) {
// 	fw := NewFileWriter("/Users/iooikaak/goproj/warhorse.com/src/warhorse.com/frame/log/horse.log")
// 	for i := 0; i < math.MaxInt64; i++ {
// 		fw.Write(fmt.Sprint(i) + "\n")
// 		time.Sleep(time.Microsecond * 10)
// 	}
// 	fw.Close()
// }

// func TestFileWriter(t *testing.T) {
// 	SetLog("/Users/iooikaak/goproj/warhorse.com/src/warhorse.com/frame/log/horse.log", "debug")
// 	Debug("c==1", " d=2")
// 	Debugf("c==%d d=%d", 1, 2)
// 	Info("c==1", " d=2")
// 	Infof("c==%d d=%d", 1, 2)
// 	Warn("c==1", " d=2")
// 	Warnf("c==%d d=%d", 1, 2)

// 	Error("c==1", " d=2")
// 	Errorf("c==%d d=%d", 1, 2)

// 	Fatal("c==1", " d=2")
// 	Fatalf("c==%d d=%d", 1, 2)
// 	Close()
// }

// func TestFileWalk(t *testing.T) {
// 	filepath.Walk("/Users/iooikaak/goproj/warhorse.com/src/warhorse.com/frame/log/horse.go", func(path string, info os.FileInfo, err error) error {
// 		if info == nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			return nil
// 		}
// 		t.Log(info.Name())
// 		return nil
// 	})
// }
