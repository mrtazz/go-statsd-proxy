package statsdproxy

import (
	"testing"
)

func TestGetHashRingPosition(t *testing.T) {
	id, err := GetHashRingPosition("foo")

	if err != nil {
		t.Errorf("GetHashRingPosition() failed with %v", err)
		t.FailNow()
	}
	if id != 3675831724 {
		t.Errorf("wrong id returned, expected 3675831724 and got %d", id)
	}
}
