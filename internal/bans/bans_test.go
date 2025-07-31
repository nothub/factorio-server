package bans

import (
	"fmt"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	filePath := fmt.Sprintf("%s/%s", os.TempDir(), "server-banlist.json")
	bans, err := FetchAndWrite(filePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v bans at %s", len(bans), filePath)
}
