package internal

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/yugovtr/executor"
	"github.com/yugovtr/setting"
	"github.com/yugovtr/timer"
)

var (
	errConfigFileFormat = errors.New("parse hours error")
)

func SelectImage(config io.Reader, t timer.Timer) (image string, err error) {
	const Hour = "15:04"

	var now string = t.Now().Format(Hour)

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
		t2, err2 := time.Parse(Hour, now)

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

func ChangeWallpaper(i string, e executor.Exec) error {
	args := []string{"set", "org.gnome.desktop.background", "picture-uri", i}

	if err := e.Run(nil, "gsettings", args...); err != nil {
		return err
	}
	return nil
}
