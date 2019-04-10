package zaplog

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func initLog() {
	InitLog(LogConfig{})
}

func TestLogInit(t *testing.T) {
	InitLog(LogConfig{})
}

func TestGetLogger(t *testing.T) {
	initLog()
	logger := MustGetLogger("test")
	logger.Info("test get logger")
}

func TestSync(t *testing.T) {
	initLog()
	logger := MustGetLogger("test")
	logger.Debug("Hello World")
	Sync()
}

func TestLevel(t *testing.T) {
	initLog()
	time.Sleep(time.Second * 10)
	rsp, err := http.Get("http://localhost:9090/handle/level")
	if err != nil {
		t.Error(err)
	}

	rst, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(rst))
}
