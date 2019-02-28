package log

import (
	"io"
	"testing"
	"time"
)

func TestSizeSplit(t *testing.T) {
	writer, err := NewSizeSplitWriter(10000)
	if err != nil {
		t.Error(err)
	}
	writer.SetCheckInterval(time.Second * 30)
	writer.StartCheck()
	startWrite(writer)
}

func startWrite(w io.Writer) {
	for {
		w.Write([]byte("hello world\n"))
		time.Sleep(time.Millisecond * 10)
	}
}

func TestTimeSplit(t *testing.T) {
	writer, err := NewTimeSplitWriter(time.Second * 30)
	if err != nil {
		t.Error(err)
	}
	writer.StartCheck()
	startWrite(writer)
}

func TestDateSplit(t *testing.T) {
	writer, err := NewDateSplitWriter()
	if err != nil {
		t.Error(err)
	}
	writer.SetCheckInterval(time.Second * 60)
	writer.StartCheck()
	startWrite(writer)
}
