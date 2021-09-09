package internal

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/yugovtr/exec"
	"github.com/yugovtr/mock"
)

func TestSelectImage(t *testing.T) {
	const (
		Hour    = "15:04"
		Content = `{"setting": [{"hour": "8:00","image": "1.jpg"},{"hour": "18:00","image": "2.jpg"},{"hour": "22:00","image": "3.jpg"}]}`
	)

	tests := []struct {
		config    io.Reader
		name      string
		now       string
		wantImage string
		wantErr   bool
	}{
		{
			name:      "when is it dawn",
			now:       "3:10",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "1.jpg",
			wantErr:   false,
		},
		{
			name:      "when it is day",
			now:       "13:42",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "2.jpg",
			wantErr:   false,
		},
		{
			name:      "when it is nightfall",
			now:       "19:00",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "3.jpg",
			wantErr:   false,
		},
		{
			name:      "when it is night",
			now:       "23:58",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "3.jpg",
			wantErr:   false,
		},
		{
			name:      "when has error in hour",
			now:       "00:00",
			config:    bytes.NewBuffer([]byte(`{"setting": [{"hour": "12","image": "1.jpg"}]}`)),
			wantImage: "",
			wantErr:   true,
		},
		{
			name:      "when has error in file",
			now:       "00:00",
			config:    bytes.NewBuffer([]byte(`Lorem Ipsum is simply dummy text of the printing and typesetting industry.`)),
			wantImage: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := func() time.Time {
				now, err := time.Parse(Hour, tt.now)

				assert.NoError(t, err, "SelectImage() parse test case %s", tt.now)
				return now
			}

			gotI, err := SelectImage(tt.config, mock)

			assert.Equal(t, (err != nil), tt.wantErr, "SelectImage() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, gotI, tt.wantImage, "SelectImage() = %v, want %v", gotI, tt.wantImage)
		})
	}
}

func TestChangeWallpaper(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		image   string
		prepare func(*mock.MockCmd)
	}{
		{
			name:    "when has no errors",
			wantErr: false,
			image:   "1.jpg",
			prepare: func(m *mock.MockCmd) {
				m.EXPECT().Run().Return(nil)
			},
		},
		{
			name:    "when has errors",
			wantErr: true,
			image:   "2.jpg",
			prepare: func(m *mock.MockCmd) {
				m.EXPECT().Run().Return(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock.NewMockCmd(ctrl)
			if tt.prepare != nil {
				tt.prepare(m)
			}

			mock := func(command string, args ...string) exec.Cmd {
				wantCommand := "gsettings"
				wantArgs := []string{"set", "org.gnome.desktop.background", "picture-uri", tt.image}

				assert.Equal(t, wantCommand, command, "ChangeWallpaper() command = %v, want %v", command, wantCommand)
				assert.Equal(t, wantArgs, args, "ChangeWallpaper() args = %v, want %v", args, wantArgs)

				return m
			}

			err := ChangeWallpaper(tt.image, mock)
			assert.Equal(t, (err != nil), tt.wantErr, "ChangeWallpaper() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
