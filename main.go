package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/xi2/xz"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	log.Println("The factory must grow!")
	prepare()
	run(flag.Args())
}

var archiveCachePath = filepath.Join("factorio-server", fmt.Sprintf("factorio_headless_x64_%s.tar.xz", factorioVersion))

func prepare() {
	binPath := filepath.Join(serverDir, "bin", "x64", "factorio")
	if _, err := os.Stat(binPath); err == nil {
		// server dir is already prepared
		return
	}

	archivePath, err := xdg.SearchCacheFile(archiveCachePath)
	if err != nil {
		// archive is not cached
		log.Println(err)

		archivePath, err = xdg.CacheFile(archiveCachePath)
		if err != nil {
			log.Fatalln(err)
		}

		u := fmt.Sprintf("https://www.factorio.com/get-download/%s/headless/linux64", factorioVersion)
		log.Printf("downloading server archive from: %s\n", u)

		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", "factorio-server (+https://github.com/nothub/factorio-server)")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		if res.StatusCode >= http.StatusBadRequest {
			log.Fatalf("http status %v\n", res.StatusCode)
		}

		f, err := os.Create(archivePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		_, err = io.Copy(f, res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// unpack archive

	log.Printf("extracting %s to: %s\n", filepath.Base(archivePath), serverDir)

	f, err := os.Open(archivePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	xzR, err := xz.NewReader(f, 0)
	if err != nil {
		log.Fatalln(err)
	}
	tarR := tar.NewReader(xzR)

	for {
		hdr, err := tarR.Next()
		if err == io.EOF {
			// end of archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		path, err := filepath.Abs(filepath.Join(serverDir, strings.TrimPrefix(hdr.Name, "factorio/")))
		if err != nil {
			log.Fatalln(err)
		}

		err = os.MkdirAll(filepath.Dir(path), 0750)
		if err != nil {
			log.Fatalln(err)
		}

		w, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
		if err != nil {
			log.Fatalln(err)
		}

		_, err = io.Copy(w, tarR)
		if err != nil {
			log.Fatalln(err)
		}

		err = w.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
