package utils

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func AddPageNumberToOverlay(imagePath string, pageNumber int) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return fmt.Errorf("could not open image: %w", err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w == 0 || h == 0 {
		return fmt.Errorf("could not read image dimensions")
	}

	sampleSize := min(100, min(w, h)/4)
	sampleX := max(0, w-sampleSize-10)
	sampleY := max(0, h-sampleSize-10)

	var totalBrightness float64
	var pixelCount int
	for y := sampleY; y < sampleY+sampleSize && y < h; y++ {
		for x := sampleX; x < sampleX+sampleSize && x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			totalBrightness += gray
			pixelCount++
		}
	}

	avgBrightness := totalBrightness / float64(pixelCount)
	isLightBG := avgBrightness > 128

	var textColor, strokeColor color.Color
	var bgColor color.NRGBA
	if isLightBG {
		textColor = color.Black
		strokeColor = color.White
		bgColor = color.NRGBA{255, 255, 255, 216} // ~0.85 alpha
	} else {
		textColor = color.White
		strokeColor = color.Black
		bgColor = color.NRGBA{0, 0, 0, 216}
	}

	fontSize := math.Max(48, float64(min(w, h))/12.0)
	padding := fontSize * 0.15
	labelW := fontSize * 1.2
	labelH := fontSize * 1.1
	cornerOffset := fontSize * 0.3

	labelX := float64(w) - labelW - cornerOffset
	labelY := float64(h) - labelH - cornerOffset

	dc := gg.NewContext(w, h)
	dc.DrawImage(img, 0, 0)

	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(labelX, labelY, labelW, labelH, padding)
	dc.Fill()

	dc.SetColor(strokeColor)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(labelX, labelY, labelW, labelH, padding)
	dc.Stroke()

	dc.SetColor(textColor)
	if err := dc.LoadFontFace("arial.ttf", fontSize); err != nil {
		if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", fontSize); err != nil {
			_ = dc.LoadFontFace("/System/Library/Fonts/Helvetica.ttc", fontSize)
		}
	}
	label := fmt.Sprintf("%d", pageNumber)
	dc.DrawStringAnchored(label, labelX+labelW/2, labelY+labelH*0.78, 0.5, 0.5)

	return dc.SavePNG(imagePath)
}

func GenerateOverlaySvg(panels []OutputPanel, width, height int) []byte {
	var buf bytes.Buffer
	canvas := svg.New(&buf)
	canvas.Start(width, height)

	for i, p := range panels {
		c := p.Width / 6
		if h := p.Height / 6; h < c {
			c = h
		}
		if c > 20 {
			c = 20
		}

		r := p.Left + p.Width
		b := p.Top + p.Height

		canvas.Rect(p.Left, p.Top, p.Width, p.Height,
			"fill:none;stroke:#ff3366;stroke-width:3")

		canvas.Path(fmt.Sprintf("M%d,%dL%d,%dL%d,%d",
			p.Left, p.Top+c, p.Left, p.Top, p.Left+c, p.Top),
			"fill:none;stroke:#ffdd00;stroke-width:5")

		canvas.Path(fmt.Sprintf("M%d,%dL%d,%dL%d,%d",
			r-c, p.Top, r, p.Top, r, p.Top+c),
			"fill:none;stroke:#ffdd00;stroke-width:5")

		canvas.Path(fmt.Sprintf("M%d,%dL%d,%dL%d,%d",
			p.Left, b-c, p.Left, b, p.Left+c, b),
			"fill:none;stroke:#ffdd00;stroke-width:5")

		canvas.Path(fmt.Sprintf("M%d,%dL%d,%dL%d,%d",
			r-c, b, r, b, r, b-c),
			"fill:none;stroke:#ffdd00;stroke-width:5")

		labelW := 130
		if i >= 9 {
			labelW += 12
		}
		canvas.Roundrect(p.Left+8, b-62, labelW, 54, 4, 4,
			"fill:rgba(0,0,0,0.85);stroke:#ff3366;stroke-width:2")

		canvas.Text(p.Left+16, b-17, fmt.Sprintf("Panel %d", i+1),
			"font-family:Arial,sans-serif;font-size:32px;font-weight:bold;fill:white")
	}

	canvas.End()
	return buf.Bytes()
}
