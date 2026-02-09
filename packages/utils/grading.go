package utils

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func ApplyCinematicTealOrange(imagePath, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	result := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origColor := img.At(x, y)
			rr, gg, bb, aa := origColor.RGBA()

			r := float64(rr >> 8)
			g := float64(gg >> 8)
			b := float64(bb >> 8)
			a := uint8(aa >> 8)

			lum := (0.299*r + 0.587*g + 0.114*b) / 255.0


			r = sCurve(r, 0.6)
			g = sCurve(g, 0.6)
			b = sCurve(b, 0.6)

			shadowWeight := 1.0 - lum
			highlightWeight := lum

			tealStrength := 0.3
			orangeStrength := 0.3

			r = mixChannel(r, r*0.85, shadowWeight*tealStrength)
			g = mixChannel(g, g*1.05, shadowWeight*tealStrength)
			b = mixChannel(b, b*1.15, shadowWeight*tealStrength)

			r = mixChannel(r, r*1.12, highlightWeight*orangeStrength)
			g = mixChannel(g, g*1.03, highlightWeight*orangeStrength)
			b = mixChannel(b, b*0.88, highlightWeight*orangeStrength)

			gray := 0.299*r + 0.587*g + 0.114*b
			satBoost := 1.08
			r = clampFloat(gray + (r-gray)*satBoost)
			g = clampFloat(gray + (g-gray)*satBoost)
			b = clampFloat(gray + (b-gray)*satBoost)

			result.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: a,
			})
		}
	}

	inputBasename := filepath.Base(imagePath)
	ext := filepath.Ext(inputBasename)
	baseName := strings.TrimSuffix(inputBasename, ext)
	outputPath := filepath.Join(outputDir, baseName+"_graded.png")

	if err := imaging.Save(result, outputPath); err != nil {
		return "", fmt.Errorf("failed to save graded image: %w", err)
	}

	return outputPath, nil
}

func sCurve(value, strength float64) float64 {
	normalized := value / 255.0
	curved := 1.0 / (1.0 + math.Exp(-strength*(normalized*12.0-6.0)))
	return clampFloat(curved * 255.0)
}

func mixChannel(original, target, factor float64) float64 {
	return clampFloat(original + (target-original)*factor)
}

func clampFloat(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return math.Round(v)
}