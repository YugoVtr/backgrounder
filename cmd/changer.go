package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yugovtr/exec"
	"github.com/yugovtr/internal"
)

const (
	configFile = "settings.json"
)

func main() {
	hour := flag.String("hour", "", "specify hour in format hh:mm")
	config := flag.String("config", configFile, "specify settings.json path")

	flag.Parse()

	t := time.Now
	if h, err := time.Parse("15:04", *hour); len(*hour) == 5 && err == nil {
		t = func() time.Time {
			return h
		}
	}

	content, err := os.ReadFile(*config)
	if err != nil {
		exit(err)
	}

	reader := bytes.NewBuffer(content)

	img, err := internal.SelectImage(reader, t)
	if err != nil {
		exit(err)
	}

	fmt.Printf("Set image \"%s\"\n", img)

	if err := internal.ChangeWallpaper(img, exec.Command); err != nil {
		exit(err)
	}

	os.Exit(0)
}

func exit(e error) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}
