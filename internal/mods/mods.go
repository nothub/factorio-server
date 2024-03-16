package mods

import (
	"fmt"
	"html"
	"log"
	"time"
)

type mod struct {
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

// mod-list.json

// TODO: factorio api client

// TODO: download / update mods

var modList = []string{"Nanobots"}

func downloadMods() {
	for _, name := range modList {
		if name == "base" {
			continue
		}
		u := fmt.Sprintf("https://mods.factorio.com/api/mods/%s", html.EscapeString(name))
		log.Printf("downloading: %s\n", u)
		// u returns mod
		//
		// url_segm: '.releases | last | .download_url'
		// filename: '.releases | last | .file_name'
		//
		// skip existing
		//
		// delete old
		//   find "${dir}" -type f -name "${mod_name}_*.zip" -exec rm -f {} \;
		//
		//  download new
		//    echo "downloading ${file_name}" >&2
		//    curl -sSL -o "${dir}/${file_name}" "https://mods.factorio.com${url_segm}?username=${username}&token=${token}"
	}
}
