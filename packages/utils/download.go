package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func DownloadPage(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpDir, err := GetDirectory("tmp")
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(parsedURL.Path)

	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%d-%d%s", time.Now().UnixMilli(), rand.Intn(1000000000), ext))

	out, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return tmpFile, nil
}
