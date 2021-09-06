package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	exec "github.com/yugovtr/executor"
	"github.com/yugovtr/internal"
	"github.com/yugovtr/timer"
)

const (
	configFile = "settings.json"
)

func main() {
	hour := flag.String("hour", "", "specify hour in format hh:mm")
	config := flag.String("config", configFile, "specify settings.json path")

	flag.Parse()

	t := timer.New()
	if h, err := time.Parse("15:04", *hour); len(*hour) == 5 && err == nil {
		t.Travel(&h)
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

	if err := internal.ChangeWallpaper(img, exec.New()); err != nil {
		exit(err)
	}

	os.Exit(0)
}

func exit(e error) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}
