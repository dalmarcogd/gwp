package internal

import (
	"testing"
)

func TestParseServerInfos(t *testing.T) {
	infos := ParseServerInfos(FakeServer{})
	if _, ok := infos["cpus"]; !ok {
		t.Error("Expected key cpus on result from ParseServerInfos")
	}
	if _, ok := infos["goroutines"]; !ok {
		t.Error("Expected key goroutines on result from ParseServerInfos")
	}
	if _, ok := infos["workers"]; !ok {
		t.Error("Expected key workers on result from ParseServerInfos")
	}
}