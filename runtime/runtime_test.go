package runtime

import (
	"testing"
)

func TestGetServerRun(t *testing.T) {
	SetServerRun(FakeServer{})
	if GetServerRun() == nil {
		t.Error("Runtime server assigned but return nil")
	}
}

func TestSetServerRun(t *testing.T) {
	SetServerRun(FakeServer{})
	if GetServerRun() == nil {
		t.Error("Runtime server assigned but return nil")
	}
}

func TestFakeServer_Workers(t *testing.T) {
	var f FakeServer
	if len(f.Workers()) != 0 {
		t.Error("FakeServer should return empty slice")
	}
}
