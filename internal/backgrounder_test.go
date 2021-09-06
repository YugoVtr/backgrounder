package internal

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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
		prepare   func(*mock.MockTimer, string)
	}{
		{
			name:      "when is it dawn",
			now:       "3:10",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "1.jpg",
			wantErr:   false,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
		{
			name:      "when it is day",
			now:       "13:42",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "2.jpg",
			wantErr:   false,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
		{
			name:      "when it is nightfall",
			now:       "19:00",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "3.jpg",
			wantErr:   false,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
		{
			name:      "when it is night",
			now:       "23:58",
			config:    bytes.NewBuffer([]byte(Content)),
			wantImage: "3.jpg",
			wantErr:   false,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
		{
			name:      "when has error in hour",
			now:       "00:00",
			config:    bytes.NewBuffer([]byte(`{"setting": [{"hour": "12","image": "1.jpg"}]}`)),
			wantImage: "",
			wantErr:   true,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
		{
			name:      "when has error in file",
			now:       "00:00",
			config:    bytes.NewBuffer([]byte(`Lorem Ipsum is simply dummy text of the printing and typesetting industry.`)),
			wantImage: "",
			wantErr:   true,
			prepare: func(m *mock.MockTimer, now string) {
				m.EXPECT().Now().DoAndReturn(func() time.Time {
					t, _ := time.Parse(Hour, now)
					return t
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock.NewMockTimer(ctrl)
			if tt.prepare != nil {
				tt.prepare(m, tt.now)
			}

			gotI, err := SelectImage(tt.config, m)

			if (err != nil) != tt.wantErr {
				t.Errorf("SelectImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotI != tt.wantImage {
				t.Errorf("SelectImage() = %v, want %v", gotI, tt.wantImage)
			}
		})
	}
}

func TestChangeWallpaper(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		image   string
		prepare func(*mock.MockExec)
	}{
		{
			name:    "",
			wantErr: false,
			image:   "1.jpg",
			prepare: func(m *mock.MockExec) {
				m.EXPECT().Run(nil, "gsettings", "set", "org.gnome.desktop.background", "picture-uri", "1.jpg").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock.NewMockExec(ctrl)
			if tt.prepare != nil {
				tt.prepare(m)
			}

			if err := ChangeWallpaper(tt.image, m); (err != nil) != tt.wantErr {
				t.Errorf("ChangeWallpaper() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
