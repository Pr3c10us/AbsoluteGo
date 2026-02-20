package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
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
		newName := fmt.Sprintf("%d%s", globalIndex, ext)
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

// Orientation represents the orientation category of an image.
type Orientation string

const (
	OrientationVertical   Orientation = "vertical"
	OrientationHorizontal Orientation = "horizontal"
	OrientationSquare     Orientation = "square"
)

func DetectOrientation(filePath string) (Orientation, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening image file: %w", err)
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return "", fmt.Errorf("decoding image config: %w", err)
	}

	w := cfg.Width
	h := cfg.Height

	if h == 0 {
		return "", fmt.Errorf("image has zero height")
	}

	switch {
	case 4*w < 3*h:
		return OrientationVertical, nil
	case 3*w > 4*h:
		return OrientationHorizontal, nil
	default:
		return OrientationSquare, nil
	}
}
