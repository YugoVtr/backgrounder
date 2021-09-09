package internal

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/yugovtr/exec"
	"github.com/yugovtr/setting"
)

var (
	errConfigFileFormat = errors.New("parse hours error")
)

type Now func() time.Time

func SelectImage(config io.Reader, now Now) (image string, err error) {
	const Hour = "15:04"

	var currentHour string = now().Format(Hour)

	content, err := io.ReadAll(config)
	if err != nil {
		return
	}

	var s = setting.Config{}
	if err = json.Unmarshal(content, &s); err != nil {
		return
	}

	for _, s := range s.Setting {
		t1, err1 := time.Parse(Hour, s.Hour)
		t2, err2 := time.Parse(Hour, currentHour)

		if err1 != nil || err2 != nil {
			err = errConfigFileFormat
			return
		}

		if t2.Before(t1) {
			return s.Image, nil
		}
	}

	return s.Setting[len(s.Setting)-1].Image, nil
}

func ChangeWallpaper(i string, exec exec.Commander) error {
	args := []string{"set", "org.gnome.desktop.background", "picture-uri", i}

	return exec("gsettings", args...).Run()
}
