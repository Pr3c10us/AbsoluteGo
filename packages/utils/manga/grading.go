package manga

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

type GradientMapOptions struct {
	Shadows    RGB
	Highlights RGB
}

var Presets = map[string]GradientMapOptions{
	"manga": {
		Shadows:    RGB{R: 30, G: 25, B: 60},
		Highlights: RGB{R: 245, G: 238, B: 210},
	},
	"horror": {
		Shadows:    RGB{R: 60, G: 5, B: 10},
		Highlights: RGB{R: 180, G: 210, B: 45},
	},
	"action": {
		Shadows:    RGB{R: 50, G: 5, B: 5},
		Highlights: RGB{R: 220, G: 50, B: 40},
	},
}

type PresetName string

const (
	PresetManga  PresetName = "manga"
	PresetHorror PresetName = "horror"
	PresetAction PresetName = "action"
)

func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a) + (float64(b)-float64(a))*t + 0.5)
}

func buildLUT(shadows, highlights RGB) [256]RGB {
	var lut [256]RGB

	for i := 0; i < 256; i++ {
		t := float64(i) / 255.0
		lut[i] = RGB{
			R: lerp(shadows.R, highlights.R, t),
			G: lerp(shadows.G, highlights.G, t),
			B: lerp(shadows.B, highlights.B, t),
		}
	}

	return lut
}

func getLuminance(r, g, b uint8) uint8 {
	return uint8((int(r)*299 + int(g)*587 + int(b)*114) / 1000)
}

func ApplyGradientMapWithPreset(imagePath, outputDir string, presetName PresetName) (string, error) {
	preset, ok := Presets[string(presetName)]
	if !ok {
		return "", fmt.Errorf("unknown preset: %s", presetName)
	}
	return ApplyGradientMap(imagePath, outputDir, &preset)
}

func ApplyGradientMap(imagePath, outputDir string, options *GradientMapOptions) (string, error) {
	shadows := RGB{R: 30, G: 25, B: 60}
	highlights := RGB{R: 245, G: 238, B: 210}

	if options != nil {
		shadows = options.Shadows
		highlights = options.Highlights
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	lut := buildLUT(shadows, highlights)

	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return "", fmt.Errorf("unable to read dimensions of %q", imagePath)
	}

	output := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origColor := img.At(x, y)
			r, g, b, a := origColor.RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			luminance := getLuminance(r8, g8, b8)

			gradedColor := lut[luminance]

			output.Set(x, y, color.RGBA{
				R: gradedColor.R,
				G: gradedColor.G,
				B: gradedColor.B,
				A: a8,
			})
		}
	}

	inputBasename := filepath.Base(imagePath)
	ext := filepath.Ext(inputBasename)
	baseName := strings.TrimSuffix(inputBasename, ext)
	outputPath := filepath.Join(outputDir, baseName+"_graded.png")

	if err := imaging.Save(output, outputPath); err != nil {
		return "", fmt.Errorf("failed to save graded image: %w", err)
	}

	return outputPath, nil
}
