//usr/bin/env -S go run "$0" "$@" ; exit
//go:build exclude

package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os/exec"
	"text/template"
)

//go:embed readme.tmpl
var fs embed.FS

func init() {
	log.SetFlags(0)
}

func main() {

	err := exec.Command("go", "build", "-o", "/tmp/factorio-server", "./cmd/").Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	data, err := exec.Command("/tmp/factorio-server", "--help").CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
	}

	tmpl, err := template.ParseFS(fs, "readme.tmpl")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var buf = bytes.Buffer{}
	err = tmpl.Execute(&buf, "```\n"+string(data)+"```")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Print(buf.String())

}
