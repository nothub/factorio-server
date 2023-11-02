package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

/*
0.307 Info ServerMultiplayerManager.cpp:814: updateTick(4294967295) changing state from(Ready) to(PreparedToHostGame)
0.307 Info ServerMultiplayerManager.cpp:814: updateTick(4294967295) changing state from(PreparedToHostGame) to(CreatingGame)
0.823 Info ServerMultiplayerManager.cpp:814: updateTick(5) changing state from(CreatingGame) to(InGame)
*/
var reChangeState = regexp.MustCompile("\\d+\\.\\d+ Info \\w+\\.cpp:\\d+: updateTick\\(\\d+\\) changing state from\\((\\w+)\\) to\\((\\w+)\\)")

/*
2023-11-02 17:58:12 [CHAT] <server>: asdf
2023-11-02 18:07:28 [CHAT] hub: asdfasdf
*/
var reChat = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[CHAT\\] ([\\w<>]+): (.+)")

// 2023-11-02 18:07:26 [JOIN] hub joined the game
var reJoined = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[JOIN\\] (\\w+) joined the game")

// 2023-11-02 18:07:32 [LEAVE] hub left the game
var reLeft = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2} \\[LEAVE\\] (\\w+) left the game")

func run(flags []string) {

	for _, flag := range flags {
		switch flag {
		case "--create", "--start-server", "--start-server-load-scenario", "--start-server-load-latest":
			log.Fatalf("The %s flag is reserved and will be provided by the process wrapper!\n", flag)
		}
	}

	// if a "saves" directory exists and contains any zip files:
	// --start-server-load-latest
	//     start a multiplayer server and load the
	//     latest available save

	// if there is no "saves" directory yet:
	// --start-server FILE
	//     start a multiplayer server
	// or
	// --start-server-load-scenario [MOD/]NAME
	//     start a multiplayer server and load the
	//     specified scenario. The scenario is looked for
	//     inside the given mod. If no mod is given, it is
	//     looked for in the top-level scenarios directory.
	// or
	// --create FILE
	//     create a new map

	flags = append(flags, "--start-server")
	flags = append(flags, "template.zip")

	cmd := exec.Command("./bin/x64/factorio", flags...)
	cmd.Dir = serverDir

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

		// run server process and wait for it to finish
		err = cmd.Run()
		if err != nil {
			log.Fatalln(err)
		}

		// wait for leftover operations and release resources
		err = cmd.Wait()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	sc := bufio.NewScanner(r)
	for {
		if sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if len(line) > 0 {
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
		}
		time.Sleep(10 * time.Millisecond)
	}
}
