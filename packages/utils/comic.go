package utils

import (
	"fmt"
	_ "image/gif"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"golift.io/xtractr"
)

type ExtractEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ComicFormat string

const (
	FormatCBR ComicFormat = "cbr"
	FormatCBZ ComicFormat = "cbz"
	FormatCB7 ComicFormat = "cb7"
	FormatCBT ComicFormat = "cbt"
	FormatPDF ComicFormat = "pdf"
)

var supportedFormats = map[ComicFormat]string{
	FormatCBR: ".cbr",
	FormatCBZ: ".cbz",
	FormatCB7: ".cb7",
	FormatCBT: ".cbt",
	FormatPDF: ".pdf",
}

func GetComicFormat(filePath string) (ComicFormat, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	for format, extension := range supportedFormats {
		if ext == extension {
			return format, nil
		}
	}
	return "", fmt.Errorf("unsupported file format: %s", ext)
}

func extractCBR(filePath, outputDir string) (int, error) {
	xFile := &xtractr.XFile{
		FilePath:  filePath,
		OutputDir: outputDir,
		FileMode:  0644,
		DirMode:   0755,
	}

	_, files, _, err := xtractr.ExtractRAR(xFile)
	if err != nil {
		return 0, fmt.Errorf("failed to extract RAR: %w", err)
	}

	return len(files), nil
}

func extractCBZ(filePath, outputDir string) (int, error) {
	xFile := &xtractr.XFile{
		FilePath:  filePath,
		OutputDir: outputDir,
		FileMode:  0644,
		DirMode:   0755,
	}

	_, files, err := xtractr.ExtractZIP(xFile)
	if err != nil {
		return 0, fmt.Errorf("failed to extract ZIP: %w", err)
	}

	return len(files), nil
}

func extractCB7(filePath, outputDir string) (int, error) {
	xFile := &xtractr.XFile{
		FilePath:  filePath,
		OutputDir: outputDir,
		FileMode:  0644,
		DirMode:   0755,
	}

	_, files, _, err := xtractr.Extract7z(xFile)
	if err != nil {
		return 0, fmt.Errorf("failed to extract 7z: %w", err)
	}

	return len(files), nil
}

func extractCBT(filePath, outputDir string) (int, error) {
	xFile := &xtractr.XFile{
		FilePath:  filePath,
		OutputDir: outputDir,
		FileMode:  0644,
		DirMode:   0755,
	}

	_, files, err := xtractr.ExtractTar(xFile)
	if err != nil {
		return 0, fmt.Errorf("failed to extract TAR: %w", err)
	}

	return len(files), nil
}

func extractPDF(filePath, outputDir string) (int, error) {
	cmd := exec.Command("pdftoppm", "-jpeg", filePath, filepath.Join(outputDir, "page"))
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to convert PDF (ensure poppler-utils is installed): %w", err)
	}

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "page") && strings.HasSuffix(name, ".jpg") {
			count++
		}
	}
	return count, nil
}

func ExtractComicToDir(fullPath string, format ComicFormat, outputDir string) error {
	var err error
	switch format {
	case FormatCBR:
		_, err = extractCBR(fullPath, outputDir)
	case FormatCBZ:
		_, err = extractCBZ(fullPath, outputDir)
	case FormatCB7:
		_, err = extractCB7(fullPath, outputDir)
	case FormatCBT:
		_, err = extractCBT(fullPath, outputDir)
	case FormatPDF:
		_, err = extractPDF(fullPath, outputDir)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return err
}
