package utils

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type BubbleComponent struct {
	Label  int
	Pixels []image.Point
	Area   int
	MinX   int
	MaxX   int
	MinY   int
	MaxY   int
}

type BubbleRemovalOptions struct {
	WhiteThreshold              int // Threshold for considering a pixel as white (default: 245)
	MaxTextArea                 int // Maximum area of text components to remove (default: 2000)
	NeighborhoodPad             int // Padding around component for neighborhood check (default: 4)
	NeighborhoodWhiteThreshold  int // Threshold for neighborhood whiteness (default: 200)
}

func DefaultBubbleRemovalOptions() BubbleRemovalOptions {
	return BubbleRemovalOptions{
		WhiteThreshold:             245,
		MaxTextArea:                2000,
		NeighborhoodPad:            4,
		NeighborhoodWhiteThreshold: 200,
	}
}

func RemoveSpeechBubbleText(inputPath, outputDir string, options *BubbleRemovalOptions) (string, error) {
	opts := DefaultBubbleRemovalOptions()
	if options != nil {
		if options.WhiteThreshold != 0 {
			opts.WhiteThreshold = options.WhiteThreshold
		}
		if options.MaxTextArea != 0 {
			opts.MaxTextArea = options.MaxTextArea
		}
		if options.NeighborhoodPad != 0 {
			opts.NeighborhoodPad = options.NeighborhoodPad
		}
		if options.NeighborhoodWhiteThreshold != 0 {
			opts.NeighborhoodWhiteThreshold = options.NeighborhoodWhiteThreshold
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	img, err := imaging.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return "", fmt.Errorf("could not read image dimensions")
	}

	grayBuffer := make([]uint8, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (uint8(r>>8) + uint8(g>>8) + uint8(b>>8)) / 3
			grayBuffer[y*width+x] = gray
		}
	}

	result := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			if a8 < 255 {
				alpha := float64(a8) / 255.0
				r8 = uint8(float64(r8)*alpha + 255*(1-alpha))
				g8 = uint8(float64(g8)*alpha + 255*(1-alpha))
				b8 = uint8(float64(b8)*alpha + 255*(1-alpha))
			}
			
			result.Set(x, y, color.RGBA{r8, g8, b8, 255})
		}
	}

	nonWhiteMask := make([]uint8, width*height)
	for i := 0; i < len(grayBuffer); i++ {
		if grayBuffer[i] <= uint8(opts.WhiteThreshold) {
			nonWhiteMask[i] = 1
		}
	}

	labels := make([]int32, width*height)
	for i := range labels {
		labels[i] = -1
	}
	
	var components []BubbleComponent
	currentLabel := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			if nonWhiteMask[idx] == 1 && labels[idx] == -1 {
				component := bubbleFloodFill(nonWhiteMask, labels, width, height, x, y, currentLabel)
				components = append(components, component)
				currentLabel++
			}
		}
	}

	removedCount := 0
	for _, comp := range components {
		if comp.Area >= opts.MaxTextArea {
			continue
		}

		x1 := max(0, comp.MinX-opts.NeighborhoodPad)
		y1 := max(0, comp.MinY-opts.NeighborhoodPad)
		x2 := min(width, comp.MaxX+opts.NeighborhoodPad+1)
		y2 := min(height, comp.MaxY+opts.NeighborhoodPad+1)

		var neighborSum int64
		var neighborCount int

		for ny := y1; ny < y2; ny++ {
			for nx := x1; nx < x2; nx++ {
				nidx := ny*width + nx
				if labels[nidx] != int32(comp.Label) {
					neighborSum += int64(grayBuffer[nidx])
					neighborCount++
				}
			}
		}

		neighborMean := 0.0
		if neighborCount > 0 {
			neighborMean = float64(neighborSum) / float64(neighborCount)
		}

		if neighborMean > float64(opts.NeighborhoodWhiteThreshold) {
			for _, pixel := range comp.Pixels {
				result.Set(pixel.X, pixel.Y, color.RGBA{255, 255, 255, 255})
			}
			removedCount++
		}
	}

	finalGray := make([]uint8, width*height)
	for i := 0; i < width*height; i++ {
		r, g, b, _ := result.At(i%width, i/width).RGBA()
		finalGray[i] = uint8((int(r>>8) + int(g>>8) + int(b>>8)) / 3)
	}

	whiteMask := make([]uint8, width*height)
	for i := 0; i < len(finalGray); i++ {
		if finalGray[i] > 240 {
			whiteMask[i] = 1
		}
	}

	closedMask := morphClose(whiteMask, width, height, 2)

	for i := 0; i < len(closedMask); i++ {
		if closedMask[i] == 1 {
			x := i % width
			y := i / width
			result.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	inputFilename := filepath.Base(inputPath)
	ext := filepath.Ext(inputFilename)
	baseName := strings.TrimSuffix(inputFilename, ext)
	outputPath := filepath.Join(outputDir, baseName+"_cleaned.png")

	if err := imaging.Save(result, outputPath); err != nil {
		return "", fmt.Errorf("failed to save output image: %w", err)
	}

	return outputPath, nil
}

func bubbleFloodFill(mask []uint8, labels []int32, width, height, startX, startY, label int) BubbleComponent {
	var pixels []image.Point
	stack := []image.Point{{X: startX, Y: startY}}

	minX, maxX, minY, maxY := startX, startX, startY, startY

	for len(stack) > 0 {
		pos := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		x, y := pos.X, pos.Y
		idx := y*width + x

		if x < 0 || x >= width || y < 0 || y >= height {
			continue
		}
		if mask[idx] != 1 || labels[idx] != -1 {
			continue
		}

		labels[idx] = int32(label)
		pixels = append(pixels, image.Point{X: x, Y: y})

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}

		stack = append(stack,
			image.Point{X: x + 1, Y: y},
			image.Point{X: x - 1, Y: y},
			image.Point{X: x, Y: y + 1},
			image.Point{X: x, Y: y - 1},
			image.Point{X: x + 1, Y: y + 1},
			image.Point{X: x - 1, Y: y - 1},
			image.Point{X: x + 1, Y: y - 1},
			image.Point{X: x - 1, Y: y + 1},
		)
	}

	return BubbleComponent{
		Label:  label,
		Pixels: pixels,
		Area:   len(pixels),
		MinX:   minX,
		MaxX:   maxX,
		MinY:   minY,
		MaxY:   maxY,
	}
}

func morphClose(mask []uint8, width, height, radius int) []uint8 {
	dilated := dilate(mask, width, height, radius)
	return erode(dilated, width, height, radius)
}

func dilate(mask []uint8, width, height, radius int) []uint8 {
	result := make([]uint8, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var maxVal uint8 = 0
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					ny := y + dy
					nx := x + dx
					if ny >= 0 && ny < height && nx >= 0 && nx < width {
						if mask[ny*width+nx] > maxVal {
							maxVal = mask[ny*width+nx]
						}
					}
				}
			}
			result[y*width+x] = maxVal
		}
	}

	return result
}

func erode(mask []uint8, width, height, radius int) []uint8 {
	result := make([]uint8, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var minVal uint8 = 1
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					ny := y + dy
					nx := x + dx
					if ny >= 0 && ny < height && nx >= 0 && nx < width {
						if mask[ny*width+nx] < minVal {
							minVal = mask[ny*width+nx]
						}
					}
				}
			}
			result[y*width+x] = minVal
		}
	}

	return result
}