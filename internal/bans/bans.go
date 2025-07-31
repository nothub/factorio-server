package bans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

// server-banlist.json

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
	slices.Compact(bans)

	return bans, nil
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
