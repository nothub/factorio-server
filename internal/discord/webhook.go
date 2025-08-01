package discord

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func SendWebhook(message string, addr string) {

	req, err := http.NewRequest("POST", addr, strings.NewReader(fmt.Sprintf(`{"content": "%s"}`, message)))
	if err != nil {
		log.Fatalf("failed creating request: %v\n\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "factorio-server (+https://github.com/nothub/factorio-server)")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("failed sending to webhook: %v\n", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Printf("webhook returned status: %d\n", res.StatusCode)
	}
}
