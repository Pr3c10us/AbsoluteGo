package manga

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

type ProcessOptions struct {
	WhiteThreshold      int
	MinArea             int
	MaxAreaRatio        float64
	EdgeMargin          int
	OCRConfidenceThresh float64 // Minimum ratio of recognized characters to total dark pixels area to confirm text (0-1)
	OCRMinConfidence    int     // Minimum average Tesseract confidence for recognized text (0-100)
	OCRMinChars         int     // Minimum number of recognized characters to consider it as text
	OCRMinCharDensity   float64 // Minimum character density (recognized chars per 1000px of region area)
	OCRLangs            string  // Tesseract language(s)
}

func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		WhiteThreshold:      230,
		MinArea:             0,
		MaxAreaRatio:        0,
		EdgeMargin:          0,
		OCRConfidenceThresh: 0,
		OCRMinConfidence:    70,
		OCRMinChars:         10,
		OCRMinCharDensity:   0.3,
		OCRLangs:            "eng",
	}
}

type Component struct {
	Pixels map[int]bool
	MinX   int
	MinY   int
	MaxX   int
	MaxY   int
	Area   int
}

type BubbleCandidate struct {
	Component *Component
	HasText   bool
}

func getPixelIndex(x, y, width int) int {
	return y*width + x
}

func findConnectedComponents(isWhite []bool, width, height int) []*Component {
	visited := make([]bool, width*height)
	var components []*Component

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := getPixelIndex(x, y, width)
			if visited[idx] || !isWhite[idx] {
				continue
			}

			component := &Component{
				Pixels: make(map[int]bool),
				MinX:   x,
				MinY:   y,
				MaxX:   x,
				MaxY:   y,
				Area:   0,
			}

			queue := [][2]int{{x, y}}
			visited[idx] = true

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				cx, cy := current[0], current[1]
				cIdx := getPixelIndex(cx, cy, width)

				component.Pixels[cIdx] = true
				component.Area++

				if cx < component.MinX {
					component.MinX = cx
				}
				if cy < component.MinY {
					component.MinY = cy
				}
				if cx > component.MaxX {
					component.MaxX = cx
				}
				if cy > component.MaxY {
					component.MaxY = cy
				}

				neighbors := [][2]int{
					{cx - 1, cy}, {cx + 1, cy},
					{cx, cy - 1}, {cx, cy + 1},
				}

				for _, neighbor := range neighbors {
					nx, ny := neighbor[0], neighbor[1]
					if nx < 0 || nx >= width || ny < 0 || ny >= height {
						continue
					}
					nIdx := getPixelIndex(nx, ny, width)
					if visited[nIdx] || !isWhite[nIdx] {
						continue
					}
					visited[nIdx] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}

			components = append(components, component)
		}
	}

	return components
}

func calculateSolidity(comp *Component) float64 {
	boxWidth := comp.MaxX - comp.MinX + 1
	boxHeight := comp.MaxY - comp.MinY + 1
	boxArea := boxWidth * boxHeight
	if boxArea > 0 {
		return float64(comp.Area) / float64(boxArea)
	}
	return 0
}

func countWordChars(text string) int {
	count := 0
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
		if r >= 0x3000 && r <= 0x9FFF || r >= 0xAC00 && r <= 0xD7AF {
			count++
		}
	}
	return count
}

func countMeaningfulChars(text string) int {
	count := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}

func verifyTextWithOCR(
	img image.Image,
	comp *Component,
	client *gosseract.Client,
	options ProcessOptions,
) (bool, error) {
	padding := 4
	imgBounds := img.Bounds()
	imgWidth := imgBounds.Dx()
	imgHeight := imgBounds.Dy()

	x0 := max(0, comp.MinX-padding)
	y0 := max(0, comp.MinY-padding)
	x1 := min(imgWidth-1, comp.MaxX+padding)
	y1 := min(imgHeight-1, comp.MaxY+padding)

	cropW := x1 - x0 + 1
	cropH := y1 - y0 + 1

	if cropW < 20 || cropH < 20 {
		return false, nil
	}

	croppedImg := imaging.Crop(img, image.Rect(x0, y0, x1+1, y1+1))

	tempFile, err := os.CreateTemp("", "ocr_*.png")
	if err != nil {
		return false, err
	}
	defer os.Remove(tempFile.Name())

	err = png.Encode(tempFile, croppedImg)
	tempFile.Close()
	if err != nil {
		return false, err
	}

	client.SetImage(tempFile.Name())
	recognizedText, err := client.Text()
	if err != nil {
		return false, err
	}
	recognizedText = strings.TrimSpace(recognizedText)

	avgConfidence := 70

	_ = countMeaningfulChars(recognizedText)

	wordChars := countWordChars(recognizedText)

	regionArea := cropW * cropH
	charDensity := (float64(wordChars) / float64(regionArea)) * 1000

	hasEnoughChars := wordChars >= options.OCRMinChars
	hasEnoughConfidence := avgConfidence >= options.OCRMinConfidence
	hasEnoughDensity := charDensity >= options.OCRMinCharDensity
	highConfidence := avgConfidence > 60

	confirmed := hasEnoughChars && hasEnoughConfidence && (highConfidence || hasEnoughDensity)

	return confirmed, nil
}

func floodFillComponent(img *image.RGBA, comp *Component, width int) {
	for idx := range comp.Pixels {
		x := idx % width
		y := idx / width
		img.Set(x, y, color.RGBA{255, 255, 255, 255})
	}

	visited := make(map[int]bool)

	for y := comp.MinY; y <= comp.MaxY; y++ {
		for x := comp.MinX; x <= comp.MaxX; x++ {
			idx := getPixelIndex(x, y, width)
			if comp.Pixels[idx] || visited[idx] {
				continue
			}

			region := make(map[int]bool)
			queue := [][2]int{{x, y}}
			touchesBoundary := false
			visited[idx] = true

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				cx, cy := current[0], current[1]
				cIdx := getPixelIndex(cx, cy, width)
				region[cIdx] = true

				if cx <= comp.MinX || cx >= comp.MaxX || cy <= comp.MinY || cy >= comp.MaxY {
					touchesBoundary = true
				}

				neighbors := [][2]int{
					{cx - 1, cy}, {cx + 1, cy},
					{cx, cy - 1}, {cx, cy + 1},
				}

				for _, neighbor := range neighbors {
					nx, ny := neighbor[0], neighbor[1]
					if nx < comp.MinX || nx > comp.MaxX || ny < comp.MinY || ny > comp.MaxY {
						touchesBoundary = true
						continue
					}
					nIdx := getPixelIndex(nx, ny, width)
					if visited[nIdx] || comp.Pixels[nIdx] {
						continue
					}
					visited[nIdx] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}

			if !touchesBoundary {
				for fillIdx := range region {
					fx := fillIdx % width
					fy := fillIdx / width
					img.Set(fx, fy, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
}

func RemoveSpeechBubbleText(imagePath, outputDir string, options *ProcessOptions) (string, error) {
	opts := DefaultProcessOptions()
	if options != nil {
		if options.WhiteThreshold != 0 {
			opts.WhiteThreshold = options.WhiteThreshold
		}
		if options.EdgeMargin != 0 {
			opts.EdgeMargin = options.EdgeMargin
		}
		if options.OCRMinConfidence != 0 {
			opts.OCRMinConfidence = options.OCRMinConfidence
		}
		if options.OCRMinChars != 0 {
			opts.OCRMinChars = options.OCRMinChars
		}
		if options.OCRMinCharDensity != 0 {
			opts.OCRMinCharDensity = options.OCRMinCharDensity
		}
		if options.OCRLangs != "" {
			opts.OCRLangs = options.OCRLangs
		}
	}

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
	imgArea := width * height

	rgbaImg := imaging.Clone(img)

	isWhite := make([]bool, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := rgbaImg.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			gray := (int(r8) + int(g8) + int(b8)) / 3
			isWhite[getPixelIndex(x, y, width)] = gray >= opts.WhiteThreshold
		}
	}

	components := findConnectedComponents(isWhite, width, height)

	var candidates []BubbleCandidate

	for _, comp := range components {
		onEdge := comp.MinX < opts.EdgeMargin ||
			comp.MinY < opts.EdgeMargin ||
			comp.MaxX > width-opts.EdgeMargin ||
			comp.MaxY > height-opts.EdgeMargin

		if onEdge && comp.Area > int(float64(imgArea)*0.05) {
			continue
		}

		solidity := calculateSolidity(comp)

		if solidity > 0.5 {
			hasText := false

			for y := comp.MinY + 2; y < comp.MaxY-2 && !hasText; y++ {
				for x := comp.MinX + 2; x < comp.MaxX-2 && !hasText; x++ {
					idx := getPixelIndex(x, y, width)
					if !comp.Pixels[idx] && !isWhite[idx] {
						hasText = true
					}
				}
			}

			if hasText && solidity > 0.4 {
				candidates = append(candidates, BubbleCandidate{Component: comp, HasText: true})
			}
		}
	}

	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetLanguage(opts.OCRLangs); err != nil {
		return "", fmt.Errorf("failed to set Tesseract language: %w", err)
	}

	outputImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			outputImg.Set(x, y, rgbaImg.At(x, y))
		}
	}

	clearedCount := 0

	for _, candidate := range candidates {
		comp := candidate.Component

		isConfirmedText, err := verifyTextWithOCR(rgbaImg, comp, client, opts)
		if err != nil {
			continue
		}

		if isConfirmedText {
			floodFillComponent(outputImg, comp, width)
			clearedCount++
		}
	}

	inputBasename := filepath.Base(imagePath)
	ext := filepath.Ext(inputBasename)
	baseName := strings.TrimSuffix(inputBasename, ext)
	outputPath := filepath.Join(outputDir, baseName+"_cleaned.png")

	if err := imaging.Save(outputImg, outputPath); err != nil {
		return "", fmt.Errorf("failed to save output image: %w", err)
	}

	return outputPath, nil
}
