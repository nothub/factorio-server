package factorio_com

import (
	"encoding/json"
	"fmt"
	"github.com/nothub/factorio-server/internal/config"
	"github.com/nothub/factorio-server/internal/utils/files"
	"github.com/nothub/semver"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
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

type Mod struct {
	Category          string `json:"category"`
	DownloadsCount    int    `json:"downloads_count"`
	LastHighlightedAt string `json:"last_highlighted_at"`
	Name              string `json:"name"`
	Owner             string `json:"owner"`
	Releases          []struct {
		DownloadUrl string `json:"download_url"`
		FileName    string `json:"file_name"`
		InfoJson    struct {
			FactorioVersion string `json:"factorio_version"`
		} `json:"info_json"`
		ReleasedAt time.Time `json:"released_at"`
		Sha1       string    `json:"sha1"`
		Version    string    `json:"version"`
	} `json:"releases"`
	Score     float64 `json:"score"`
	Summary   string  `json:"summary"`
	Thumbnail string  `json:"thumbnail"`
	Title     string  `json:"title"`
}

var cachedLatestRelease *semver.Version = nil

func GetModInfo(name string) (Mod, error) {
	escapedName := html.EscapeString(name)
	modUrl := fmt.Sprintf("https://mods.factorio.com/api/mods/%s", escapedName)

	res, err := http.Get(modUrl)
	if err != nil {
		return Mod{}, fmt.Errorf("failed to fetch mod info: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Mod{}, fmt.Errorf("factorio.com error status: %s", res.Status)
	}

	var modResponse Mod

	err = json.NewDecoder(res.Body).Decode(&modResponse)
	if err != nil {
		return Mod{}, fmt.Errorf("failed to decode mod response: %w", err)
	}

	return modResponse, nil
}

func DownloadMod(mod Mod) (string, error) {

	if len(mod.Releases) == 0 {
		return "", fmt.Errorf("no releases found for mod %s", mod.Name)
	}

	latest := mod.Releases[len(mod.Releases)-1]
	filePath := filepath.Join("mods", latest.FileName)

	// TODO: check game version compatibility

	if files.IsFile(filePath) {
		log.Printf("Skipping %s (exists)", latest.FileName)
		return latest.FileName, nil
	}

	fullUrl := fmt.Sprintf("https://mods.factorio.com%s?username=%s&token=%s",
		latest.DownloadUrl,
		url.QueryEscape(config.Loaded.FactorioUser),
		url.QueryEscape(config.Loaded.FactorioToken))

	log.Printf("Downloading %s ...\n", latest.FileName)

	resp, err := http.Get(fullUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download mod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download mod, status: %s", resp.Status)
	}

	// create file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create mod file: %w", err)
	}
	defer out.Close()

	// write file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write mod file: %w", err)
	}

	return latest.FileName, nil
}

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
