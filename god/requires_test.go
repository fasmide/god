package god

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestRequiresTimeout(t *testing.T) {
	rq := Requires{Exists: "/tmp/blarh", Timeout: time.Duration(time.Second)}
	err := rq.Wait()
	if err == nil {
		t.Fail()
	}
}

func TestRequires(t *testing.T) {
	rq := Requires{Exists: "/tmp/blarhhhhh", Timeout: time.Duration(time.Second)}
	go func() {
		ioutil.WriteFile("/tmp/blarhhhhh", nil, 0600)
	}()
	err := rq.Wait()
	if err != nil {
		t.Fail()
	}
}
