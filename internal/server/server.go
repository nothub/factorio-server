package server

import (
	"archive/tar"
	"bufio"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/nothub/factorio-server/internal/config"
	factorio_com "github.com/nothub/factorio-server/internal/factorio.com"
	"github.com/nothub/factorio-server/internal/files"
	"github.com/ulikunitz/xz"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

/*
0.307 Info ServerMultiplayerManager.cpp:814: updateTick(4294967295) changing state from(Ready) to(PreparedToHostGame)
0.307 Info ServerMultiplayerManager.cpp:814: updateTick(4294967295) changing state from(PreparedToHostGame) to(CreatingGame)
0.823 Info ServerMultiplayerManager.cpp:814: updateTick(5) changing state from(CreatingGame) to(InGame)
*/
var reChangeState = regexp.MustCompile("^\\d+\\.\\d+ Info \\w+\\.cpp:\\d+: updateTick\\(\\d+\\) changing state from\\((\\w+)\\) to\\((\\w+)\\)$")

/*
2023-11-02 17:58:12 [CHAT] <server>: asdf
2023-11-02 18:07:28 [CHAT] hub: asdfasdf
*/
var reChat = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[CHAT\\] ([\\w<>]+): (.+)$")

// 2023-11-02 18:07:26 [JOIN] hub joined the game
var reJoined = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[JOIN\\] (\\w+) joined the game$")

// 2023-11-02 18:07:32 [LEAVE] hub left the game
var reLeft = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[LEAVE\\] (\\w+) left the game$")

func Run(args []string) (shutdown func()) {
	setup()

	for _, arg := range args {
		switch arg {
		case "--create", "--start-server", "--start-server-load-scenario", "--start-server-load-latest":
			log.Fatalf("The %s flag is reserved and will be provided by the process wrapper!\n", arg)
		}
	}

	if savesExist() {
		// load the latest savegames
		args = append(args, "--start-server-load-latest")

	} else {
		// make sure map.zip exists
		if !files.IsFile("map.zip") {
			createMap()
		}

		args = append(args, "--start-server")
		args = append(args, "map.zip")
	}

	cmd := exec.Command("./bin/x64/factorio", args...)
	cmd.Dir = config.ServerDir

	// these pipes will read the process stdout and stderr
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalln(err)
	}

	// this pipe will write to the process stdin
	var in *io.WriteCloser

	go func() {
		// set custom pipe for process stdout and stderr
		cmd.Stdout = w
		cmd.Stderr = w

		// grab process stdin pipe
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatalln(err)
		}
		in = &stdin

		// Run server process and wait for it to finish
		log.Println("Starting the server")
		if err = cmd.Run(); err != nil {
			log.Fatalln(err)
		}
	}()

	// scan server console output periodically
	scan := bufio.NewScanner(r)
	tick := time.NewTicker(10 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				if scan.Scan() {
					line := strings.TrimSpace(scan.Text())
					if len(line) > 0 {
						handle(line, in)
					}
				}
			}
		}
	}()

	return func() {
		log.Println("Stopping server!")
		tick.Stop()
		done <- true
	}
}

func handle(line string, in *io.WriteCloser) {
	log.Println(line)
	if m := reChangeState.FindStringSubmatch(line); m != nil {
		log.Printf("[EVENT] state changed: %s -> %s\n", m[1], m[2])
		if m[2] == "InGame" {
			(*in).Write([]byte("lets go\n"))
		}
	} else if m := reChat.FindStringSubmatch(line); m != nil {
		log.Printf("[EVENT] player %s says: %s\n", m[1], m[2])
	} else if m := reJoined.FindStringSubmatch(line); m != nil {
		log.Printf("[EVENT] player joined: %s\n", m[1])
	} else if m := reLeft.FindStringSubmatch(line); m != nil {
		log.Printf("[EVENT] player left: %s\n", m[1])
	}
}

func createMap() {
	cmd := exec.Command("./bin/x64/factorio", "--create", "map.zip")
	cmd.Dir = config.ServerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Creating fresh map")
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func savesExist() (ok bool) {
	if !files.IsDir("saves") {
		return false
	}
	err := filepath.WalkDir("saves", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".zip") {
			ok = true
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return ok
}

var archiveCachePath = filepath.Join("factorio-server", fmt.Sprintf("factorio_headless_x64_%s.tar.xz", factorio_com.LatestRelease()))

func setup() {
	binPath := filepath.Join(config.ServerDir, "bin", "x64", "factorio")
	if _, err := os.Stat(binPath); err == nil {
		// server dir is already prepared
		return
	}

	archivePath, err := xdg.SearchCacheFile(archiveCachePath)
	if err != nil {
		// this error indicates: file does not exist
		log.Println(err)

		archivePath, err = xdg.CacheFile(archiveCachePath)
		if err != nil {
			log.Fatalln(err)
		}

		u := fmt.Sprintf("https://www.factorio.com/get-download/%s/headless/linux64", factorio_com.LatestRelease())
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

	log.Printf("extracting %s to: %s\n", filepath.Base(archivePath), config.ServerDir)

	f, err := os.Open(archivePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	xzR, err := xz.NewReader(f)
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

		path, err := filepath.Abs(filepath.Join(config.ServerDir, strings.TrimPrefix(hdr.Name, "factorio/")))
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
