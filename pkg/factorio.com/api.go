package factorio_com

import (
	"encoding/json"
	"github.com/nothub/semver"
	"log"
	"net/http"
)

type VersionModel struct {
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

func LatestRelease() semver.Version {
	if cachedLatestRelease != nil {
		return *cachedLatestRelease
	}

	res, err := http.Get("https://factorio.com/api/latest-releases")
	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		// TODO: return error value for api failure response status
		log.Fatalf("factorio api response status: %s\n", res.Status)
	}

	var infos VersionModel

	err = json.NewDecoder(res.Body).Decode(&infos)
	if err != nil {
		log.Fatalln(err)
	}

	ver, err := semver.Parse(infos.Stable.Headless)
	if err != nil {
		log.Fatalln(err)
	}

	cachedLatestRelease = &ver

	return ver
}