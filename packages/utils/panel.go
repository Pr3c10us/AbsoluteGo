package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Pr3c10us/absolutego/packages/utils/manga"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type BoundingBox struct {
	MinRow, MinCol, MaxRow, MaxCol int
}

type OutputPanel struct {
	Left, Top, Width, Height int
}

type DetectResult struct {
	Panels  int
	Success bool
}

type floodFillPanel struct {
	MinX, MaxX, MinY, MaxY int
	PixelCount             int
}

func isContainedIn(inner, outer OutputPanel) bool {
	return inner.Left >= outer.Left &&
		inner.Top >= outer.Top &&
		inner.Left+inner.Width <= outer.Left+outer.Width &&
		inner.Top+inner.Height <= outer.Top+outer.Height
}

func removeNestedPanels(panels []OutputPanel) []OutputPanel {
	var result []OutputPanel
	for i, p := range panels {
		nested := false
		for j, o := range panels {
			if i != j && isContainedIn(p, o) {
				nested = true
				break
			}
		}
		if !nested {
			result = append(result, p)
		}
	}
	return result
}

func imageToGray(img image.Image) []uint8 {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	buf := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			buf[y*w+x] = uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(bl)) / 256)
		}
	}
	return buf
}

func imageToGrayFloat(img image.Image) ([]float64, int, int) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	buf := make([]float64, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			buf[y*w+x] = 0.299*float64(r)/256 + 0.587*float64(g)/256 + 0.114*float64(bl)/256
		}
	}
	return buf, w, h
}

func detectGutterColor(buf []uint8, w, h int) int {
	hist := make([]int, 256)
	for x := 0; x < w; x++ {
		hist[buf[x]]++
		hist[buf[(h-1)*w+x]]++
	}
	for y := 0; y < h; y++ {
		hist[buf[y*w]]++
		hist[buf[y*w+w-1]]++
	}
	midX, midY := w/2, h/2
	for y := 0; y < h; y++ {
		hist[buf[y*w+midX]]++
	}
	for x := 0; x < w; x++ {
		hist[buf[midY*w+x]]++
	}
	peak, peakIdx := 0, 0
	for i := 0; i < 256; i++ {
		if hist[i] > peak {
			peak = hist[i]
			peakIdx = i
		}
	}
	return peakIdx
}

func buildGutterMask(buf []uint8, w, h, gutterColor, tolerance int) []uint8 {
	mask := make([]uint8, w*h)
	for i := 0; i < w*h; i++ {
		diff := int(buf[i]) - gutterColor
		if diff < 0 {
			diff = -diff
		}
		if diff <= tolerance {
			mask[i] = 1
		}
	}
	return mask
}

func detectPanelsFloodFill(grayBuf []uint8, w, h int) []OutputPanel {
	gutterColor := detectGutterColor(grayBuf, w, h)
	gutterMask := buildGutterMask(grayBuf, w, h, gutterColor, 30)
	visited := make([]uint8, w*h)
	var panels []floodFillPanel

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := y*w + x
			if gutterMask[i] == 0 && visited[i] == 0 {
				p := floodFillPanel{MinX: x, MaxX: x, MinY: y, MaxY: y}
				queue := []int{x, y}
				visited[i] = 1
				for len(queue) > 0 {
					cy := queue[len(queue)-1]
					cx := queue[len(queue)-2]
					queue = queue[:len(queue)-2]
					p.PixelCount++
					if cx < p.MinX {
						p.MinX = cx
					}
					if cx > p.MaxX {
						p.MaxX = cx
					}
					if cy < p.MinY {
						p.MinY = cy
					}
					if cy > p.MaxY {
						p.MaxY = cy
					}
					for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
						nx, ny := cx+d[0], cy+d[1]
						if nx >= 0 && nx < w && ny >= 0 && ny < h {
							ni := ny*w + nx
							if gutterMask[ni] == 0 && visited[ni] == 0 {
								visited[ni] = 1
								queue = append(queue, nx, ny)
							}
						}
					}
				}
				panels = append(panels, p)
			}
		}
	}

	var out []OutputPanel
	for _, p := range panels {
		pw, ph := p.MaxX-p.MinX, p.MaxY-p.MinY
		if pw > 150 && ph > 250 {
			out = append(out, OutputPanel{Left: p.MinX, Top: p.MinY, Width: pw, Height: ph})
		}
	}
	return out
}

func gaussianBlur(img []float64, w, h int, sigma float64) []float64 {
	size := int(math.Ceil(sigma*6)) | 1
	kernel := make([]float64, size)
	center := size >> 1
	var sum float64
	for i := 0; i < size; i++ {
		d := float64(i - center)
		kernel[i] = math.Exp(-(d * d) / (2 * sigma * sigma))
		sum += kernel[i]
	}
	for i := range kernel {
		kernel[i] /= sum
	}

	half := center
	temp := make([]float64, w*h)
	result := make([]float64, w*h)

	for y := 0; y < h; y++ {
		off := y * w
		for x := 0; x < w; x++ {
			var s float64
			for k := 0; k < size; k++ {
				ix := x + k - half
				if ix < 0 {
					ix = 0
				} else if ix >= w {
					ix = w - 1
				}
				s += img[off+ix] * kernel[k]
			}
			temp[off+x] = s
		}
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s float64
			for k := 0; k < size; k++ {
				iy := y + k - half
				if iy < 0 {
					iy = 0
				} else if iy >= h {
					iy = h - 1
				}
				s += temp[iy*w+x] * kernel[k]
			}
			result[y*w+x] = s
		}
	}
	return result
}

func cannyEdgeDetection(gray []float64, w, h int) []uint8 {
	blurred := gaussianBlur(gray, w, h, 1.4)
	length := w * h
	magnitude := make([]float64, length)
	direction := make([]uint8, length)

	for y := 1; y < h-1; y++ {
		yOff := y * w
		yOffM := (y - 1) * w
		yOffP := (y + 1) * w
		for x := 1; x < w-1; x++ {
			tl := blurred[yOffM+x-1]
			tc := blurred[yOffM+x]
			tr := blurred[yOffM+x+1]
			ml := blurred[yOff+x-1]
			mr := blurred[yOff+x+1]
			bl := blurred[yOffP+x-1]
			bc := blurred[yOffP+x]
			br := blurred[yOffP+x+1]
			gx := -tl + tr - 2*ml + 2*mr - bl + br
			gy := -tl - 2*tc - tr + bl + 2*bc + br
			magnitude[yOff+x] = math.Sqrt(gx*gx + gy*gy)
			angle := math.Mod(math.Atan2(gy, gx)*180/math.Pi+180, 180)
			if angle < 22.5 || angle >= 157.5 {
				direction[yOff+x] = 0
			} else if angle < 67.5 {
				direction[yOff+x] = 1
			} else if angle < 112.5 {
				direction[yOff+x] = 2
			} else {
				direction[yOff+x] = 3
			}
		}
	}

	suppressed := make([]float64, length)
	for y := 1; y < h-1; y++ {
		yOff := y * w
		for x := 1; x < w-1; x++ {
			idx := yOff + x
			mag := magnitude[idx]
			dir := direction[idx]
			var n1, n2 float64
			switch dir {
			case 0:
				n1, n2 = magnitude[idx-1], magnitude[idx+1]
			case 1:
				n1, n2 = magnitude[idx-w+1], magnitude[idx+w-1]
			case 2:
				n1, n2 = magnitude[idx-w], magnitude[idx+w]
			default:
				n1, n2 = magnitude[idx-w-1], magnitude[idx+w+1]
			}
			if mag >= n1 && mag >= n2 {
				suppressed[idx] = mag
			}
		}
	}

	var maxVal float64
	for _, v := range suppressed {
		if v > maxVal {
			maxVal = v
		}
	}
	high := maxVal * 0.15
	low := maxVal * 0.05
	edges := make([]uint8, length)
	for i, v := range suppressed {
		if v >= high {
			edges[i] = 255
		} else if v >= low {
			edges[i] = 128
		}
	}

	for changed := true; changed; {
		changed = false
		for y := 1; y < h-1; y++ {
			yOff := y * w
			for x := 1; x < w-1; x++ {
				idx := yOff + x
				if edges[idx] != 128 {
					continue
				}
				if edges[idx-w-1] == 255 || edges[idx-w] == 255 || edges[idx-w+1] == 255 ||
					edges[idx-1] == 255 || edges[idx+1] == 255 ||
					edges[idx+w-1] == 255 || edges[idx+w] == 255 || edges[idx+w+1] == 255 {
					edges[idx] = 255
					changed = true
				}
			}
		}
	}
	for i := range edges {
		if edges[i] == 128 {
			edges[i] = 0
		}
	}
	return edges
}

func dilatePanel(img []uint8, w, h, iterations int) []uint8 {
	src := make([]uint8, len(img))
	copy(src, img)
	dst := make([]uint8, len(img))

	for iter := 0; iter < iterations; iter++ {
		for i := range dst {
			dst[i] = 0
		}
		for y := 1; y < h-1; y++ {
			yOff := y * w
			for x := 1; x < w-1; x++ {
				idx := yOff + x
				m := src[idx]
				for _, d := range []int{idx - w - 1, idx - w, idx - w + 1, idx - 1, idx + 1, idx + w - 1, idx + w, idx + w + 1} {
					if src[d] > m {
						m = src[d]
					}
				}
				dst[idx] = m
			}
		}
		src, dst = dst, src
	}
	return src
}

func binaryFillHoles(edges []uint8, w, h int) []uint8 {
	length := w * h
	result := make([]uint8, length)
	for i := range result {
		result[i] = 255
	}
	visited := make([]uint8, length)
	stack := make([]int, 0, length/4)

	for x := 0; x < w; x++ {
		if edges[x] == 0 {
			stack = append(stack, x)
		}
		bot := (h-1)*w + x
		if edges[bot] == 0 {
			stack = append(stack, bot)
		}
	}
	for y := 1; y < h-1; y++ {
		left := y * w
		right := left + w - 1
		if edges[left] == 0 {
			stack = append(stack, left)
		}
		if edges[right] == 0 {
			stack = append(stack, right)
		}
	}

	for len(stack) > 0 {
		idx := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[idx] != 0 {
			continue
		}
		visited[idx] = 1
		if edges[idx] != 0 {
			continue
		}
		result[idx] = 0
		x := idx % w
		y := idx / w
		if x > 0 && visited[idx-1] == 0 {
			stack = append(stack, idx-1)
		}
		if x < w-1 && visited[idx+1] == 0 {
			stack = append(stack, idx+1)
		}
		if y > 0 && visited[idx-w] == 0 {
			stack = append(stack, idx-w)
		}
		if y < h-1 && visited[idx+w] == 0 {
			stack = append(stack, idx+w)
		}
	}
	return result
}

func labelComponentsWithBoxes(binary []uint8, w, h int) []BoundingBox {
	length := w * h
	labels := make([]int, length)
	var boxes []BoundingBox
	stack := make([]int, 0, 1024)

	for y := 0; y < h; y++ {
		yOff := y * w
		for x := 0; x < w; x++ {
			idx := yOff + x
			if binary[idx] != 255 || labels[idx] != 0 {
				continue
			}
			label := len(boxes) + 1
			minRow, maxRow, minCol, maxCol := y, y, x, x
			stack = append(stack[:0], idx)
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if labels[cur] != 0 {
					continue
				}
				labels[cur] = label
				cx := cur % w
				cy := cur / w
				if cy < minRow {
					minRow = cy
				}
				if cy > maxRow {
					maxRow = cy
				}
				if cx < minCol {
					minCol = cx
				}
				if cx > maxCol {
					maxCol = cx
				}
				if cx > 0 && binary[cur-1] == 255 && labels[cur-1] == 0 {
					stack = append(stack, cur-1)
				}
				if cx < w-1 && binary[cur+1] == 255 && labels[cur+1] == 0 {
					stack = append(stack, cur+1)
				}
				if cy > 0 && binary[cur-w] == 255 && labels[cur-w] == 0 {
					stack = append(stack, cur-w)
				}
				if cy < h-1 && binary[cur+w] == 255 && labels[cur+w] == 0 {
					stack = append(stack, cur+w)
				}
			}
			boxes = append(boxes, BoundingBox{minRow, minCol, maxRow, maxCol})
		}
	}
	return boxes
}

func mergeOverlapping(regions []BoundingBox) []BoundingBox {
	var panels []BoundingBox
	for _, r := range regions {
		merged := false
		for i := range panels {
			p := &panels[i]
			if r.MinRow < p.MaxRow && r.MaxRow > p.MinRow && r.MinCol < p.MaxCol && r.MaxCol > p.MinCol {
				if r.MinRow < p.MinRow {
					p.MinRow = r.MinRow
				}
				if r.MinCol < p.MinCol {
					p.MinCol = r.MinCol
				}
				if r.MaxRow > p.MaxRow {
					p.MaxRow = r.MaxRow
				}
				if r.MaxCol > p.MaxCol {
					p.MaxCol = r.MaxCol
				}
				merged = true
				break
			}
		}
		if !merged {
			panels = append(panels, r)
		}
	}
	return panels
}

func clusterPanels(bboxes []BoundingBox, axis string, depth int) []interface{} {
	if depth > 10 || len(bboxes) <= 1 {
		sorted := make([]BoundingBox, len(bboxes))
		copy(sorted, bboxes)
		if axis == "row" {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].MinRow < sorted[j].MinRow })
		} else {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].MinCol < sorted[j].MinCol })
		}
		result := make([]interface{}, len(sorted))
		for i, b := range sorted {
			result[i] = b
		}
		return result
	}

	var clusters [][]BoundingBox
	for _, bbox := range bboxes {
		added := false
		for ci := range clusters {
			for _, b := range clusters[ci] {
				aligned := false
				if axis == "row" {
					aligned = b.MinRow < bbox.MaxRow && bbox.MinRow < b.MaxRow
				} else {
					aligned = b.MinCol < bbox.MaxCol && bbox.MinCol < b.MaxCol
				}
				if aligned {
					clusters[ci] = append(clusters[ci], bbox)
					added = true
					break
				}
			}
			if added {
				break
			}
		}
		if !added {
			clusters = append(clusters, []BoundingBox{bbox})
		}
	}

	if len(clusters) == 1 && len(clusters[0]) == len(bboxes) {
		sorted := make([]BoundingBox, len(bboxes))
		copy(sorted, bboxes)
		if axis == "row" {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].MinRow < sorted[j].MinRow })
		} else {
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].MinCol < sorted[j].MinCol })
		}
		result := make([]interface{}, len(sorted))
		for i, b := range sorted {
			result[i] = b
		}
		return result
	}

	if axis == "row" {
		sort.Slice(clusters, func(i, j int) bool { return clusters[i][0].MinRow < clusters[j][0].MinRow })
	} else {
		sort.Slice(clusters, func(i, j int) bool { return clusters[i][0].MinCol < clusters[j][0].MinCol })
	}

	var result []interface{}
	for _, c := range clusters {
		nextAxis := "col"
		if axis == "col" {
			nextAxis = "row"
		}
		if len(c) > 1 {
			result = append(result, clusterPanels(c, nextAxis, depth+1))
		} else {
			result = append(result, c[0])
		}
	}
	return result
}

func flattenPanels(nested []interface{}) []BoundingBox {
	var result []BoundingBox
	for _, item := range nested {
		switch v := item.(type) {
		case BoundingBox:
			result = append(result, v)
		case []interface{}:
			result = append(result, flattenPanels(v)...)
		}
	}
	return result
}

func detectPanelsCanny(img image.Image, w, h int) []OutputPanel {
	gray, _, _ := imageToGrayFloat(img)
	edges := cannyEdgeDetection(gray, w, h)

	// --- FIX START ---
	// Reduced kernel size from 21 to 9.
	// 21 was too large and merged distinct panels together (filling the gutters).
	// 9 is safer: it connects broken lines without jumping across panel gaps.
	const kSize = 9

	// 1. Close Horizontal gaps
	closedH := morphologicalClose(edges, w, h, kSize, 1)

	// 2. Close Vertical gaps
	closedV := morphologicalClose(edges, w, h, 1, kSize)

	// 3. Combine
	combined := make([]uint8, w*h)
	for i := range combined {
		if closedH[i] > 128 || closedV[i] > 128 || edges[i] > 128 {
			combined[i] = 255
		} else {
			combined[i] = 0
		}
	}

	// 4. Minor smoothing (keep this small)
	finalEdges := dilateDirectional(combined, w, h, 3, 3)
	// --- FIX END ---

	filled := binaryFillHoles(finalEdges, w, h)
	regions := labelComponentsWithBoxes(filled, w, h)
	merged := mergeOverlapping(regions)

	imageArea := w * h
	var filtered []BoundingBox
	for _, b := range merged {
		area := (b.MaxRow - b.MinRow) * (b.MaxCol - b.MinCol)

		// Filter noise (too small)
		if area < imageArea/150 {
			continue
		}

		// --- NEW SAFETY CHECK ---
		// If a single detected panel covers more than 90% of the image,
		// it is likely a "merge error" where gutters were bridged.
		// We discard it so we don't lose the Flood Fill results.
		if area > int(float64(imageArea)*0.90) {
			continue
		}

		filtered = append(filtered, b)
	}

	clustered := clusterPanels(filtered, "row", 0)
	ordered := flattenPanels(clustered)

	out := make([]OutputPanel, len(ordered))
	for i, p := range ordered {
		out[i] = OutputPanel{Left: p.MinCol, Top: p.MinRow, Width: p.MaxCol - p.MinCol, Height: p.MaxRow - p.MinRow}
	}
	return out
}

func panelsOverlap(a, b OutputPanel) bool {
	return a.Left < b.Left+b.Width &&
		a.Left+a.Width > b.Left &&
		a.Top < b.Top+b.Height &&
		a.Top+a.Height > b.Top
}

func mergeTouchingPanels(panels []OutputPanel) []OutputPanel {
	if len(panels) == 0 {
		return nil
	}
	parent := make([]int, len(panels))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]]
			i = parent[i]
		}
		return i
	}
	union := func(a, b int) {
		parent[find(a)] = find(b)
	}

	for i := 0; i < len(panels); i++ {
		for j := i + 1; j < len(panels); j++ {
			if panelsOverlap(panels[i], panels[j]) {
				union(i, j)
			}
		}
	}

	groups := map[int][]OutputPanel{}
	for i := range panels {
		root := find(i)
		groups[root] = append(groups[root], panels[i])
	}

	var result []OutputPanel
	for _, group := range groups {
		left, top := math.MaxInt32, math.MaxInt32
		right, bottom := math.MinInt32, math.MinInt32
		for _, p := range group {
			if p.Left < left {
				left = p.Left
			}
			if p.Top < top {
				top = p.Top
			}
			if p.Left+p.Width > right {
				right = p.Left + p.Width
			}
			if p.Top+p.Height > bottom {
				bottom = p.Top + p.Height
			}
		}
		result = append(result, OutputPanel{Left: left, Top: top, Width: right - left, Height: bottom - top})
	}
	return result
}

func isColorful(img image.Image) bool {
	small := imaging.Fit(img, 100, 100, imaging.Lanczos)
	b := small.Bounds()
	w, h := b.Dx(), b.Dy()
	pixelCount := w * h
	colorfulPixels := 0
	threshold := uint32(15 * 257)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := small.At(x, y).RGBA()
			maxC := r
			if g > maxC {
				maxC = g
			}
			if bl > maxC {
				maxC = bl
			}
			minC := r
			if g < minC {
				minC = g
			}
			if bl < minC {
				minC = bl
			}
			if maxC-minC > threshold {
				colorfulPixels++
			}
		}
	}
	return float64(colorfulPixels)/float64(pixelCount) > 0.1
}

func drawRect(img *image.RGBA, p OutputPanel, c color.RGBA, thickness int) {
	for t := 0; t < thickness; t++ {
		for x := p.Left; x < p.Left+p.Width; x++ {
			if x >= 0 && x < img.Bounds().Dx() {
				if p.Top+t >= 0 && p.Top+t < img.Bounds().Dy() {
					img.SetRGBA(x, p.Top+t, c)
				}
				if p.Top+p.Height-1-t >= 0 && p.Top+p.Height-1-t < img.Bounds().Dy() {
					img.SetRGBA(x, p.Top+p.Height-1-t, c)
				}
			}
		}
		for y := p.Top; y < p.Top+p.Height; y++ {
			if y >= 0 && y < img.Bounds().Dy() {
				if p.Left+t >= 0 && p.Left+t < img.Bounds().Dx() {
					img.SetRGBA(p.Left+t, y, c)
				}
				if p.Left+p.Width-1-t >= 0 && p.Left+p.Width-1-t < img.Bounds().Dx() {
					img.SetRGBA(p.Left+p.Width-1-t, y, c)
				}
			}
		}
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA, thickness int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy
	b := img.Bounds()
	half := thickness / 2

	for {
		for ty := -half; ty <= half; ty++ {
			for tx := -half; tx <= half; tx++ {
				px, py := x0+tx, y0+ty
				if px >= b.Min.X && px < b.Max.X && py >= b.Min.Y && py < b.Max.Y {
					img.SetRGBA(px, py, c)
				}
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawOverlay(img *image.RGBA, panels []OutputPanel) {
	pink := color.RGBA{255, 51, 102, 255}
	yellow := color.RGBA{255, 221, 0, 255}
	bgColor := color.RGBA{0, 0, 0, 216} // ~0.85 opacity
	white := color.RGBA{255, 255, 255, 255}

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

		drawRect(img, p, pink, 3)

		drawLine(img, p.Left, p.Top+c, p.Left, p.Top, yellow, 5)
		drawLine(img, p.Left, p.Top, p.Left+c, p.Top, yellow, 5)
		drawLine(img, r-c, p.Top, r, p.Top, yellow, 5)
		drawLine(img, r, p.Top, r, p.Top+c, yellow, 5)
		drawLine(img, p.Left, b-c, p.Left, b, yellow, 5)
		drawLine(img, p.Left, b, p.Left+c, b, yellow, 5)
		drawLine(img, r-c, b, r, b, yellow, 5)
		drawLine(img, r, b, r, b-c, yellow, 5)

		labelW := 130
		if i >= 9 {
			labelW += 12
		}
		labelRect := image.Rect(p.Left+8, b-62, p.Left+8+labelW, b-62+54)
		draw.Draw(img, labelRect, &image.Uniform{bgColor}, image.Point{}, draw.Over)
		drawRect(img, OutputPanel{Left: p.Left + 8, Top: b - 62, Width: labelW, Height: 54}, pink, 2)

		dc := gg.NewContextForRGBA(img)
		if err := dc.LoadFontFace("arial.ttf", 32); err == nil {
			dc.SetColor(white)
			label := fmt.Sprintf("Panel %d", i+1)
			textX := float64(p.Left + 8 + labelW/2)
			textY := float64(b - 35)
			dc.DrawStringAnchored(label, textX, textY, 0.5, 0.5)
		}
	}
}

func DetectAndExtractPanels(imagePath string) DetectResult {
	fullPath, _ := filepath.Abs(imagePath)
	dir := filepath.Dir(fullPath)
	ext := filepath.Ext(fullPath)
	name := strings.TrimSuffix(filepath.Base(fullPath), ext)

	srcImg, err := imaging.Open(fullPath)
	if err != nil {
		return DetectResult{0, false}
	}

	bounds := srcImg.Bounds()
	origW, origH := bounds.Dx(), bounds.Dy()

	processingWidth := 1500
	if origW < processingWidth {
		processingWidth = origW
	}
	scale := float64(origW) / float64(processingWidth)

	resized := imaging.Resize(srcImg, processingWidth, 0, imaging.Lanczos)
	rb := resized.Bounds()
	w, h := rb.Dx(), rb.Dy()

	grayBuf := imageToGray(resized)

	floodPanels := detectPanelsFloodFill(grayBuf, w, h)
	cannyPanels := detectPanelsCanny(resized, w, h)

	allPanels := append(floodPanels, cannyPanels...)
	merged := mergeTouchingPanels(allPanels)

	outputPanels := make([]OutputPanel, len(merged))
	for i, p := range merged {
		outputPanels[i] = OutputPanel{
			Left:   int(float64(p.Left) * scale),
			Top:    int(float64(p.Top) * scale),
			Width:  int(math.Ceil(float64(p.Width) * scale)),
			Height: int(math.Ceil(float64(p.Height) * scale)),
		}
	}

	outputPanels = removeNestedPanels(outputPanels)
	sort.Slice(outputPanels, func(i, j int) bool {
		if outputPanels[i].Top != outputPanels[j].Top {
			return outputPanels[i].Top < outputPanels[j].Top
		}
		return outputPanels[i].Left < outputPanels[j].Left
	})

	if len(outputPanels) == 0 {
		return DetectResult{0, true}
	}

	panelsDir := filepath.Join(dir, name)
	os.MkdirAll(panelsDir, 0o755)

	for i, p := range outputPanels {
		region := image.Rect(
			max(0, p.Left),
			max(0, p.Top),
			min(origW, p.Left+p.Width),
			min(origH, p.Top+p.Height),
		)
		if region.Dx() <= 0 || region.Dy() <= 0 {
			continue
		}

		cropped := imaging.Crop(srcImg, region)
		tempPath := filepath.Join(panelsDir, fmt.Sprintf("temp_%d%s", i+1, ext))
		imaging.Save(cropped, tempPath)

		colorful := isColorful(cropped)

		if colorful {
			cleanedPath, err := RemoveSpeechBubbleText(tempPath, panelsDir, nil)
			if err != nil {
				cleanedPath = tempPath
			}
			finalPath := filepath.Join(panelsDir, fmt.Sprintf("%d.png", i+1))
			if cleanedPath != tempPath {
				os.Rename(cleanedPath, finalPath)
				os.Remove(tempPath)
			} else {
				os.Rename(tempPath, finalPath)
			}
		} else {
			cleanedPath, err := manga.RemoveSpeechBubbleText(tempPath, panelsDir, nil)
			if err != nil {
				cleanedPath = tempPath
			}
			finalPath := filepath.Join(panelsDir, fmt.Sprintf("%d.png", i+1))
			if cleanedPath != tempPath {
				os.Rename(cleanedPath, finalPath)
				os.Remove(tempPath)
			} else {
				os.Rename(tempPath, finalPath)
			}
		}
	}

	overlayImg := image.NewRGBA(image.Rect(0, 0, origW, origH))
	draw.Draw(overlayImg, overlayImg.Bounds(), srcImg, bounds.Min, draw.Src)
	drawOverlay(overlayImg, outputPanels)
	overlayPngPath := filepath.Join(dir, name+".png")
	f, err := os.Create(overlayPngPath)
	if err == nil {
		png.Encode(f, overlayImg)
		f.Close()
	}

	return DetectResult{Panels: len(outputPanels), Success: true}
}

// dilateDirectional performs dilation with a rectangular kernel (kW x kH)
func dilateDirectional(src []uint8, w, h, kW, kH int) []uint8 {
	dst := make([]uint8, len(src))
	// Radius
	rW := kW / 2
	rH := kH / 2

	for y := 0; y < h; y++ {
		yOff := y * w
		for x := 0; x < w; x++ {
			maxVal := src[yOff+x]

			// Check neighbors within the rectangle
			for ky := -rH; ky <= rH; ky++ {
				py := y + ky
				if py < 0 || py >= h {
					continue
				}
				pOff := py * w
				for kx := -rW; kx <= rW; kx++ {
					px := x + kx
					if px < 0 || px >= w {
						continue
					}
					if src[pOff+px] > maxVal {
						maxVal = src[pOff+px]
					}
				}
			}
			dst[yOff+x] = maxVal
		}
	}
	return dst
}

// erodeDirectional performs erosion with a rectangular kernel (kW x kH)
func erodeDirectional(src []uint8, w, h, kW, kH int) []uint8 {
	dst := make([]uint8, len(src))
	rW := kW / 2
	rH := kH / 2

	for y := 0; y < h; y++ {
		yOff := y * w
		for x := 0; x < w; x++ {
			minVal := uint8(255)

			// To erode, we find the minimum value in the kernel
			for ky := -rH; ky <= rH; ky++ {
				py := y + ky
				if py < 0 || py >= h {
					// Out of bounds is treated as "white" (255) for erosion logic
					// so it doesn't shrink edges near image borders artificially
					continue
				}
				pOff := py * w
				for kx := -rW; kx <= rW; kx++ {
					px := x + kx
					if px < 0 || px >= w {
						continue
					}
					if src[pOff+px] < minVal {
						minVal = src[pOff+px]
					}
				}
			}
			dst[yOff+x] = minVal
		}
	}
	return dst
}

// morphologicalClose bridges gaps by Dilating then Eroding
func morphologicalClose(src []uint8, w, h, kW, kH int) []uint8 {
	dilated := dilateDirectional(src, w, h, kW, kH)
	eroded := erodeDirectional(dilated, w, h, kW, kH)
	return eroded
}
