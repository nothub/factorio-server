package bans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func FetchAndWrite(filePath string) ([]string, error) {

	bans, err := Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bans: %w", err)
	}

	err = Write(filePath, bans)
	if err != nil {
		return nil, fmt.Errorf("failed to write bans: %w", err)
	}

	return bans, nil
}

func Fetch() (bans []string, err error) {

	{
		url := "https://gist.githubusercontent.com/nothub/e0b6e1962aa551ce41580c34eb0ae5f6/raw/d7403fd1fbbfdbf06a2d1e29db243e868f5a47ba/bans.json"
		var data []struct {
			Username string `json:"username"`
			Reason   string `json:"reason"`
		}

		err := getJson(url, &data)
		if err != nil {
			return nil, err
		}

		for _, ban := range data {
			bans = append(bans, ban.Username)
		}
	}

	{
		url := "https://m45sci.xyz:8443/server-banlist.json"
		var data []struct {
			Username string `json:"username"`
			Reason   string `json:"reason"`
		}

		err := getJson(url, &data)
		if err != nil {
			return nil, err
		}

		for _, ban := range data {
			bans = append(bans, ban.Username)
		}
	}

	{
		url := "https://m45sci.xyz:8443/composite.json"
		var data []struct {
			Username string `json:"username"`
			Reason   string `json:"reason"`
		}

		err := getJson(url, &data)
		if err != nil {
			return nil, err
		}

		for _, ban := range data {
			bans = append(bans, ban.Username)
		}
	}

	{
		url := "https://getcomfy.eu/api/v1/bans"
		var data [][]struct {
			Username string `json:"username"`
			Reason   string `json:"reason"`
			DateTime string `json:"dateTime"`
		}

		err := getJson(url, &data)
		if err != nil {
			return nil, err
		}

		for _, dat := range data {
			for _, ban := range dat {
				bans = append(bans, ban.Username)
			}
		}
	}

	slices.Sort(bans)
	bans = slices.Compact(bans)
	{
		var filtered []string
		for _, ban := range bans {
			if len(strings.TrimSpace(ban)) > 0 {
				filtered = append(filtered, ban)
			}
		}
		bans = filtered
	}

	return bans, nil
}

func Write(filePath string, bans []string) error {

	if filepath.Base(filePath) != "server-banlist.json" {
		return fmt.Errorf("file name must be server-banlist.json")
	}

	data, err := json.MarshalIndent(bans, "", "")
	if err != nil {
		return fmt.Errorf("failed to marshal ban list: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write ban list to %s: %w", filePath, err)
	}

	return nil
}

func getJson(url string, data any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching bans failed with status %s for %s", res.Status, url)
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}
