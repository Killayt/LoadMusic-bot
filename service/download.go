package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var ErrDurationLimitExceeded = errors.New("max duration limit")

const expression = "^(http(s)?:\\/\\/)?((w){3}.)?(music\\.)?youtu(be|.be)?(\\.com)?\\/.+"

type Service interface {
	IsValidURL(url string) bool
	Download(ctx context.Context, url string) (string, error)
}

type Downloader struct {
	maxVideoDuration time.Duration

	r *regexp.Regexp
}

func AcceptDownloader(maxVideoDuration int32) (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		maxVideoDuration: time.Minute * time.Duration(maxVideoDuration),
		r:                r,
	}, nil
}

func (d *Downloader) Download(ctx context.Context, url string) (string, error) {
	// this command downloads video and extracts mp3
	cmd := exec.CommandContext(ctx, "youtube-dl", "-x", "--audio-format", "mp3", url)
	data, err := cmd.CombinedOutput()
	if err != nil {
		os.Remove("*")
		return "", err
	}

	if strings.Contains(string(data), "ERROR") {
		os.Remove("*")
		return "", errors.New(fmt.Sprintf("error downloading video with youtube-dl, output: %s", string(data)))
	}

	return string(data) + ".mp3", nil
}

func (d *Downloader) IsValidURL(url string) bool {
	return d.r.MatchString(url)
}
