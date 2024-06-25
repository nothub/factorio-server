package factorio_com

import (
	"encoding/json"
	"github.com/nothub/semver"
	"log"
	"net/http"
)

type versionModel struct {
	Experimental struct {
		Alpha    string `json:"alpha"`
		Demo     string `json:"demo"`
		Headless string `json:"headless"`
	} `json:"experimental"`
	Stable struct {
		Alpha    string `json:"alpha"`
		Demo     string `json:"demo"`
		Headless string `json:"headless"`
	} `json:"stable"`
}

var cachedLatestRelease *semver.Version = nil

func LatestRelease() string {
	if cachedLatestRelease != nil {
		return cachedLatestRelease.String()
	}

	res, err := http.Get("https://factorio.com/api/latest-releases")
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		// TODO: return error value for api failure response status
		log.Fatalf("factorio api response status: %s\n", res.Status)
	}

	var infos versionModel

	err = json.NewDecoder(res.Body).Decode(&infos)
	if err != nil {
		log.Fatalln(err)
	}

	ver, err := semver.Parse(infos.Stable.Headless)
	if err != nil {
		log.Fatalln(err)
	}

	cachedLatestRelease = &ver

	return ver.String()
}
