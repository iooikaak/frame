package core

import (
	"testing"

	"github.com/iooikaak/frame/protocol"
)

func TestCtxLog(t *testing.T) {
	c := &icecontext{
		req: &protocol.Proto{
			Bizid: "ZZZZZZZZZZZZ",
		},
	}
	c.Debug("a", "b", "c")
	c.Debugf("A=%d B=%d C=%d", 48, 49, 50)

	c.Warn("a", "b", "c")
	c.Warnf("A=%d B=%d C=%d", 48, 49, 50)

	c.Info("a", "b", "c")
	c.Infof("A=%d B=%d C=%d", 48, 49, 50)

	c.Error("a", "b", "c")
	c.Errorf("A=%d B=%d C=%d", 48, 49, 50)

	c.Fatal("a", "b", "c")
	c.Fatalf("A=%d B=%d C=%d", 48, 49, 50)

	//c.Monitor("a", "b", "c")
	//c.Monitorf("A=%d B=%d C=%d", 48, 49, 50)
}
