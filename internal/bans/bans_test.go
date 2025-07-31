package bans

import (
	"testing"
)

func TestFetch(t *testing.T) {
	bans, err := Fetch()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("bans: %+v", bans)
}
