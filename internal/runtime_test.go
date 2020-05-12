package internal

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

func TestFakeServer_Healthy(t *testing.T) {
	var f FakeServer
	if !f.Healthy() {
		t.Error("FakeServer should return true healthy, but returned false")
	}
}

func TestFakeServer_Infos(t *testing.T) {
	var f FakeServer
	if infos := f.Infos(); infos == nil {
		t.Error("FakeServer should return non nil infos, but returned nil")
	}
}