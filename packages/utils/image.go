package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GenerateBlurCover(outputDir string, image string) (string, error) {
	coverBlurPath := filepath.Join(outputDir, "cover_blur.jpg")

	img, err := imaging.Open(image)
	if err != nil {
		return "", err
	}

	blurred := imaging.Blur(img, 20)

	err = imaging.Save(blurred, coverBlurPath, imaging.JPEGQuality(80))
	if err != nil {
		return "", err
	}
	return coverBlurPath, nil
}

func FindImages(dir string) ([]string, error) {
	var images []string
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".bmp": true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if imageExts[ext] {
				images = append(images, path)
			}
		}
		return nil
	})

	return images, err
}

func SortImages(tempDir string, outputDir string) ([]string, error) {
	globalIndex := 0
	var allImageFiles []string

	comicImages, err := FindImages(tempDir)
	if err != nil {
		return nil, err
	}
	sort.Strings(comicImages)

	for _, imgPath := range comicImages {
		globalIndex++
		ext := filepath.Ext(imgPath)
		newName := fmt.Sprintf("%05d%s", globalIndex, ext)
		destPath := filepath.Join(outputDir, newName)

		err = os.Rename(imgPath, destPath)
		if err != nil {
			return nil, err
		}

		allImageFiles = append(allImageFiles, destPath)
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		return nil, err
	}

	return allImageFiles, nil
}
