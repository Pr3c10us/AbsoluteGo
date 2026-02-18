package utils

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/disintegration/imaging"
)

type Effect string

const (
	EffectZoomIn   Effect = "zoomIn"
	EffectZoomOut  Effect = "zoomOut"
	EffectPanLeft  Effect = "panLeft"
	EffectPanRight Effect = "panRight"
	EffectPanUp    Effect = "panUp"
	EffectPanDown  Effect = "panDown"
	EffectNone     Effect = "none"
)

type VideoData struct {
	Panel    string
	Duration float64
	Effect   Effect
}

type TransitionEffect string

const (
	TransitionFade        TransitionEffect = "fade"
	TransitionWipeLeft    TransitionEffect = "wipeleft"
	TransitionWipeRight   TransitionEffect = "wiperight"
	TransitionWipeUp      TransitionEffect = "wipeup"
	TransitionWipeDown    TransitionEffect = "wipedown"
	TransitionSlideLeft   TransitionEffect = "slideleft"
	TransitionSlideRight  TransitionEffect = "slideright"
	TransitionSlideUp     TransitionEffect = "slideup"
	TransitionSlideDown   TransitionEffect = "slidedown"
	TransitionDissolve    TransitionEffect = "dissolve"
	TransitionSmoothLeft  TransitionEffect = "smoothleft"
	TransitionSmoothRight TransitionEffect = "smoothright"
	TransitionCircleOpen  TransitionEffect = "circleopen"
	TransitionCircleClose TransitionEffect = "circleclose"
)

type CreateVideoOptions struct {
	FPS                int
	Width              int
	Height             int
	BackgroundImage    string
	HWAccel            configs.HWAccel
	TransitionDuration float64
	TransitionEffect   TransitionEffect
}

func defaultCreateVideoOptions() CreateVideoOptions {
	return CreateVideoOptions{
		FPS:                30,
		Width:              1920,
		Height:             1080,
		HWAccel:            configs.HWAccelNone,
		TransitionDuration: 0.5,
		TransitionEffect:   TransitionFade,
	}
}

type MergeAudioOptions struct {
	AudioFade bool
	Loop      bool
	Volume    float64
}

func defaultMergeAudioOptions() MergeAudioOptions {
	return MergeAudioOptions{
		AudioFade: true,
		Loop:      false,
		Volume:    1.0,
	}
}

type MergeVideosOptions struct {
	TransitionDuration float64
	TransitionEffect   TransitionEffect
	HWAccel            configs.HWAccel
	GapDuration        float64
}

func defaultMergeVideosOptions() MergeVideosOptions {
	return MergeVideosOptions{
		TransitionDuration: 0.5,
		TransitionEffect:   TransitionFade,
		HWAccel:            configs.HWAccelNone,
		GapDuration:        1,
	}
}

type AddBackgroundMusicOptions struct {
	Volume  float64
	FadeIn  float64
	FadeOut float64
}

func defaultAddBackgroundMusicOptions() AddBackgroundMusicOptions {
	return AddBackgroundMusicOptions{
		Volume:  0.3,
		FadeIn:  1,
		FadeOut: 2,
	}
}

func easeExpr(totalFrames int) string {
	t := fmt.Sprintf("(on/%d)", totalFrames)
	return fmt.Sprintf("(%s*%s*(3-2*%s))", t, t, t)
}

func easeInvExpr(totalFrames int) string {
	t := fmt.Sprintf("(on/%d)", totalFrames)
	return fmt.Sprintf("(1-%s*%s*(3-2*%s))", t, t, t)
}

func getEffectFilter(effect Effect, duration float64, fps int) string {
	totalFrames := int(math.Ceil(duration * float64(fps)))
	eased := easeExpr(totalFrames)
	easedInv := easeInvExpr(totalFrames)

	centerX := func(z string) string { return fmt.Sprintf("(iw-iw/(%s))/2", z) }
	centerY := func(z string) string { return fmt.Sprintf("(ih-ih/(%s))/2", z) }

	switch effect {
	case EffectZoomIn:
		// 120% → 100%: z eases from 1.2 down to 1.0
		z := fmt.Sprintf("1.2-0.2*%s", eased)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, centerX(z), centerY(z))

	case EffectZoomOut:
		// 100% → 120%: z eases from 1.0 up to 1.2
		z := fmt.Sprintf("1.0+0.2*%s", eased)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, centerX(z), centerY(z))

	case EffectPanRight:
		// 150% width, left-edge aligned → right-edge aligned
		// x: 0 → iw - iw/1.5   y: vertically centred
		z := "1.5"
		x := fmt.Sprintf("(iw-iw/%s)*%s", z, eased)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, x, centerY(z))

	case EffectPanLeft:
		// 150% width, right-edge aligned → left-edge aligned
		// x: iw - iw/1.5 → 0   y: vertically centred
		z := "1.5"
		x := fmt.Sprintf("(iw-iw/%s)*%s", z, easedInv)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, x, centerY(z))

	case EffectPanDown:
		// 150% height, top-edge aligned → bottom-edge aligned
		// y: 0 → ih - ih/1.5   x: horizontally centred
		z := "1.5"
		y := fmt.Sprintf("(ih-ih/%s)*%s", z, eased)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, centerX(z), y)

	case EffectPanUp:
		// 150% height, bottom-edge aligned → top-edge aligned
		// y: ih - ih/1.5 → 0   x: horizontally centred
		z := "1.5"
		y := fmt.Sprintf("(ih-ih/%s)*%s", z, easedInv)
		return fmt.Sprintf("z='%s':x='%s':y='%s'", z, centerX(z), y)

	default:
		// No effect – show image at exactly 100%, centred
		return "z='1.0':x='0':y='0'"
	}
}

func ffprobeGetDuration(filePath string) (float64, error) {
	out, err := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed for %s: %w", filePath, err)
	}
	d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse duration for %s: %w", filePath, err)
	}
	return d, nil
}

func ffprobeHasAudio(filePath string) (bool, error) {
	out, err := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "a",
		"-show_entries", "stream=codec_type",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()
	if err != nil {
		return false, fmt.Errorf("ffprobe failed for %s: %w", filePath, err)
	}
	return strings.TrimSpace(string(out)) != "", nil
}

type processedImage struct {
	path   string
	isTemp bool
}

func preprocessImages(videoData []VideoData, width, height int, padColor color.Color) ([]processedImage, string, error) {
	tempDir, err := os.MkdirTemp("", "vidgen-")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	results := make([]processedImage, len(videoData))
	for i, data := range videoData {
		outPath := filepath.Join(tempDir, fmt.Sprintf("panel_%d.png", i))

		src, err := imaging.Open(data.Panel, imaging.AutoOrientation(true))
		if err != nil {
			os.RemoveAll(tempDir)
			return nil, "", fmt.Errorf("failed to open image %s: %w", data.Panel, err)
		}

		// Fit the image within the target dimensions, preserving aspect ratio
		resized := imaging.Fit(src, width, height, imaging.Lanczos)

		// Create a canvas of the exact target size filled with the pad color
		canvas := imaging.New(width, height, padColor)

		// Paste the resized image centered on the canvas
		offsetX := (width - resized.Bounds().Dx()) / 2
		offsetY := (height - resized.Bounds().Dy()) / 2
		canvas = imaging.Paste(canvas, resized, image.Pt(offsetX, offsetY))

		if err := imaging.Save(canvas, outPath); err != nil {
			os.RemoveAll(tempDir)
			return nil, "", fmt.Errorf("failed to save image %s: %w", outPath, err)
		}

		results[i] = processedImage{path: outPath, isTemp: true}
	}
	return results, tempDir, nil
}

func cleanupTempDir(tempDir string) {
	if tempDir != "" {
		os.RemoveAll(tempDir)
	}
}

// encoderOptions returns ffmpeg encoder flags for the given HWAccel.
func encoderOptions(accel configs.HWAccel) []string {
	switch accel {
	case configs.HWAccelNvidia:
		return []string{"-c:v", "h264_nvenc", "-preset", "p1", "-rc", "vbr", "-cq", "26"}
	case configs.HWAccelApple:
		return []string{"-c:v", "h264_videotoolbox", "-q:v", "65"}
	default:
		return []string{"-c:v", "libx264", "-preset", "ultrafast", "-tune", "fastdecode", "-crf", "26"}
	}
}

func encoderOptionsMerge(accel configs.HWAccel) []string {
	switch accel {
	case configs.HWAccelNvidia:
		return []string{"-c:v", "h264_nvenc", "-preset", "p4", "-rc", "vbr", "-cq", "23"}
	case configs.HWAccelApple:
		return []string{"-c:v", "h264_videotoolbox", "-q:v", "65"}
	default:
		return []string{"-c:v", "libx264", "-preset", "fast", "-crf", "23"}
	}
}

func CreateVideoFromImages(videoData []VideoData, outputPath string, opts *CreateVideoOptions) error {
	o := defaultCreateVideoOptions()
	if opts != nil {
		if opts.FPS > 0 {
			o.FPS = opts.FPS
		}
		if opts.Width > 0 {
			o.Width = opts.Width
		}
		if opts.Height > 0 {
			o.Height = opts.Height
		}
		if opts.BackgroundImage != "" {
			o.BackgroundImage = opts.BackgroundImage
		}
		if opts.HWAccel != "" {
			o.HWAccel = opts.HWAccel
		}
		if opts.TransitionDuration > 0 {
			o.TransitionDuration = opts.TransitionDuration
		}
		if opts.TransitionEffect != "" {
			o.TransitionEffect = opts.TransitionEffect
		}
	}

	if len(videoData) == 0 {
		return fmt.Errorf("videoData slice cannot be empty")
	}

	minDur := videoData[0].Duration
	for _, d := range videoData[1:] {
		if d.Duration < minDur {
			minDur = d.Duration
		}
	}
	if o.TransitionDuration >= minDur {
		return fmt.Errorf("transition duration (%.2fs) must be less than shortest clip (%.2fs)", o.TransitionDuration, minDur)
	}

	processed, tempDir, err := preprocessImages(videoData, o.Width, o.Height, color.NRGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF})
	if err != nil {
		return err
	}
	defer cleanupTempDir(tempDir)

	return createVideo(videoData, processed, outputPath, o)
}

func createVideo(videoData []VideoData, images []processedImage, outputPath string, o CreateVideoOptions) error {
	const ss = 1.5
	ssW := int(float64(o.Width) * ss)
	ssH := int(float64(o.Height) * ss)
	padColor := "0xFF00FF"

	var args []string

	// Inputs
	if o.BackgroundImage != "" {
		args = append(args, "-loop", "1", "-i", o.BackgroundImage)
	}
	for _, img := range images {
		args = append(args, "-i", img.path)
	}

	// Build filter_complex
	var filters []string

	if o.BackgroundImage != "" {
		splitOuts := ""
		for i := range videoData {
			splitOuts += fmt.Sprintf("[bg%d]", i)
		}
		filters = append(filters,
			fmt.Sprintf("[0:v]scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d,setsar=1,split=%d%s",
				o.Width, o.Height, o.Width, o.Height, len(videoData), splitOuts))
	}

	for i, data := range videoData {
		idx := i
		if o.BackgroundImage != "" {
			idx = i + 1
		}
		isLast := i == len(videoData)-1
		effectiveDuration := data.Duration
		if !isLast {
			effectiveDuration += o.TransitionDuration
		}

		effectParams := getEffectFilter(data.Effect, effectiveDuration, o.FPS)
		frames := int(math.Ceil(effectiveDuration * float64(o.FPS)))

		filters = append(filters,
			fmt.Sprintf("[%d:v]scale=%d:%d:flags=lanczos,zoompan=%s:d=%d:s=%dx%d:fps=%d,scale=%d:%d:flags=lanczos,colorkey=%s:0.3:0.1[effect%d]",
				idx, ssW, ssH, effectParams, frames, ssW, ssH, o.FPS, o.Width, o.Height, padColor, i))

		if o.BackgroundImage != "" {
			filters = append(filters,
				fmt.Sprintf("[bg%d][effect%d]overlay=(W-w)/2:(H-h)/2:shortest=1[v%d]", i, i, i))
		} else {
			filters = append(filters,
				fmt.Sprintf("[effect%d]format=yuv420p[v%d]", i, i))
		}
	}

	// xfade chain
	if len(videoData) == 1 {
		filters = append(filters, "[v0]copy[outv]")
	} else {
		offset := 0.0
		for i := 0; i < len(videoData)-1; i++ {
			offset += videoData[i].Duration
			inA := fmt.Sprintf("[v0]")
			if i > 0 {
				inA = fmt.Sprintf("[xf%d]", i-1)
			}
			inB := fmt.Sprintf("[v%d]", i+1)
			out := fmt.Sprintf("[xf%d]", i)
			if i == len(videoData)-2 {
				out = "[outv]"
			}
			filters = append(filters,
				fmt.Sprintf("%s%sxfade=transition=%s:duration=%g:offset=%g%s",
					inA, inB, o.TransitionEffect, o.TransitionDuration, offset-o.TransitionDuration, out))
		}
	}

	filterComplex := strings.Join(filters, "; ")
	args = append(args, "-filter_complex", filterComplex)
	args = append(args, "-map", "[outv]", "-threads", "0",
		"-r", strconv.Itoa(o.FPS), "-pix_fmt", "yuv420p", "-movflags", "+faststart")
	args = append(args, encoderOptions(o.HWAccel)...)
	args = append(args, "-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}

// MergeAudioToVideo merges an audio track onto a video file.
func MergeAudioToVideo(videoPath, audioPath, outputPath string, opts *MergeAudioOptions) error {
	o := defaultMergeAudioOptions()
	if opts != nil {
		o.AudioFade = opts.AudioFade
		o.Loop = opts.Loop
		if opts.Volume > 0 {
			o.Volume = opts.Volume
		}
	}

	var args []string
	args = append(args, "-i", videoPath)

	if o.Loop {
		args = append(args, "-stream_loop", "-1")
	}
	args = append(args, "-i", audioPath)

	var audioFilters []string
	if o.Volume != 1.0 {
		audioFilters = append(audioFilters, fmt.Sprintf("volume=%g", o.Volume))
	}

	if o.AudioFade {
		duration, err := ffprobeGetDuration(videoPath)
		if err != nil {
			return err
		}
		fadeStart := math.Max(0, duration-1)
		audioFilters = append(audioFilters,
			fmt.Sprintf("afade=t=in:st=0:d=0.5,afade=t=out:st=%g:d=1", fadeStart))
	}

	outOpts := []string{"-c:v", "copy", "-c:a", "aac", "-b:a", "192k", "-shortest"}
	if len(audioFilters) > 0 {
		outOpts = append(outOpts, "-af", strings.Join(audioFilters, ","))
	}

	args = append(args, outOpts...)
	args = append(args, "-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}

// MergeVideos concatenates multiple videos with xfade transitions.
func MergeVideos(videoPaths []string, outputPath string, opts *MergeVideosOptions) error {
	o := defaultMergeVideosOptions()
	if opts != nil {
		if opts.TransitionDuration > 0 {
			o.TransitionDuration = opts.TransitionDuration
		}
		if opts.TransitionEffect != "" {
			o.TransitionEffect = opts.TransitionEffect
		}
		if opts.HWAccel != "" {
			o.HWAccel = opts.HWAccel
		}
		if opts.GapDuration >= 0 {
			o.GapDuration = opts.GapDuration
		}
	}

	if len(videoPaths) == 0 {
		return fmt.Errorf("videoPaths slice cannot be empty")
	}

	if len(videoPaths) == 1 {
		data, err := os.ReadFile(videoPaths[0])
		if err != nil {
			return err
		}
		return os.WriteFile(outputPath, data, 0644)
	}

	durations := make([]float64, len(videoPaths))
	audioChecks := make([]bool, len(videoPaths))
	for i, vp := range videoPaths {
		d, err := ffprobeGetDuration(vp)
		if err != nil {
			return err
		}
		durations[i] = d

		ha, err := ffprobeHasAudio(vp)
		if err != nil {
			return err
		}
		audioChecks[i] = ha
	}

	hasAudio := true
	for _, ha := range audioChecks {
		if !ha {
			hasAudio = false
			break
		}
	}

	minDur := durations[0]
	for _, d := range durations[1:] {
		if d < minDur {
			minDur = d
		}
	}
	if o.TransitionDuration >= minDur {
		return fmt.Errorf("transition duration (%.2fs) must be less than shortest video (%.2fs)", o.TransitionDuration, minDur)
	}

	numVideos := len(videoPaths)

	var args []string
	for _, vp := range videoPaths {
		args = append(args, "-i", vp)
	}

	var filters []string

	// Pad videos with gap
	if o.GapDuration > 0 {
		for i := 0; i < numVideos-1; i++ {
			filters = append(filters,
				fmt.Sprintf("[%d:v]tpad=stop_duration=%g[vpad%d]", i, o.GapDuration, i))
		}
	}

	effectiveDurations := make([]float64, numVideos)
	for i, d := range durations {
		if i < numVideos-1 {
			effectiveDurations[i] = d + o.GapDuration
		} else {
			effectiveDurations[i] = d
		}
	}

	cumDur := 0.0
	for i := 0; i < numVideos-1; i++ {
		cumDur += effectiveDurations[i]

		inputA := fmt.Sprintf("[%d:v]", 0)
		if o.GapDuration > 0 && i == 0 {
			inputA = "[vpad0]"
		}
		if i > 0 {
			inputA = fmt.Sprintf("[vfade%d]", i-1)
		}

		inputB := fmt.Sprintf("[%d:v]", i+1)
		if i+1 < numVideos-1 && o.GapDuration > 0 {
			inputB = fmt.Sprintf("[vpad%d]", i+1)
		}

		outLabel := fmt.Sprintf("[vfade%d]", i)
		if i == numVideos-2 {
			outLabel = "[outv]"
		}

		offset := cumDur - o.TransitionDuration*float64(i+1)
		filters = append(filters,
			fmt.Sprintf("%s%sxfade=transition=%s:duration=%g:offset=%.3f%s",
				inputA, inputB, o.TransitionEffect, o.TransitionDuration, offset, outLabel))
	}

	if hasAudio {
		numTransitions := numVideos - 1
		totalTrimNeeded := o.TransitionDuration * float64(numTransitions)
		trimPerClip := totalTrimNeeded / float64(numVideos)
		const audioFade = 0.03

		for i := 0; i < numVideos; i++ {
			audioDur := durations[i] - trimPerClip
			isLast := i == numVideos-1
			padFilter := ""
			if !isLast && o.GapDuration > 0 {
				padFilter = fmt.Sprintf(",apad=pad_dur=%g", o.GapDuration)
			}
			filters = append(filters,
				fmt.Sprintf("[%d:a]atrim=0:%.3f,afade=t=out:st=%.3f:d=%g%s[a%d]",
					i, audioDur, audioDur-audioFade, audioFade, padFilter, i))
		}

		audioInputs := ""
		for i := 0; i < numVideos; i++ {
			audioInputs += fmt.Sprintf("[a%d]", i)
		}
		filters = append(filters,
			fmt.Sprintf("%sconcat=n=%d:v=0:a=1[outa]", audioInputs, numVideos))
	}

	filterComplex := strings.Join(filters, "; ")
	args = append(args, "-filter_complex", filterComplex)
	args = append(args, "-map", "[outv]")
	if hasAudio {
		args = append(args, "-map", "[outa]")
	}
	args = append(args, "-movflags", "+faststart")
	args = append(args, encoderOptionsMerge(o.HWAccel)...)
	if hasAudio {
		args = append(args, "-c:a", "aac", "-b:a", "192k")
	}
	args = append(args, "-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg merge error: %w", err)
	}
	return nil
}

// AddBackgroundMusic mixes background music tracks onto a video.
func AddBackgroundMusic(videoPath string, audioPaths []string, outputPath string, opts *AddBackgroundMusicOptions) error {
	o := defaultAddBackgroundMusicOptions()
	if opts != nil {
		if opts.Volume > 0 {
			o.Volume = opts.Volume
		}
		if opts.FadeIn > 0 {
			o.FadeIn = opts.FadeIn
		}
		if opts.FadeOut > 0 {
			o.FadeOut = opts.FadeOut
		}
	}

	if len(audioPaths) == 0 {
		return fmt.Errorf("audioPaths slice cannot be empty")
	}

	videoDuration, err := ffprobeGetDuration(videoPath)
	if err != nil {
		return err
	}

	audioDurations := make([]float64, len(audioPaths))
	totalAudioDur := 0.0
	for i, ap := range audioPaths {
		d, err := ffprobeGetDuration(ap)
		if err != nil {
			return err
		}
		audioDurations[i] = d
		totalAudioDur += d
	}

	loopsNeeded := int(math.Ceil(videoDuration / totalAudioDur))

	var args []string
	args = append(args, "-i", videoPath)
	for _, ap := range audioPaths {
		args = append(args, "-i", ap)
	}

	numAudios := len(audioPaths)
	var filters []string

	// Concatenate audio inputs
	if numAudios == 1 {
		filters = append(filters, "[1:a]acopy[concataudio]")
	} else {
		audioInputs := ""
		for i := range audioPaths {
			audioInputs += fmt.Sprintf("[%d:a]", i+1)
		}
		filters = append(filters,
			fmt.Sprintf("%sconcat=n=%d:v=0:a=1[concataudio]", audioInputs, numAudios))
	}

	// Loop if needed
	if loopsNeeded > 1 {
		loopSize := int(math.Ceil(totalAudioDur * 48000))
		filters = append(filters,
			fmt.Sprintf("[concataudio]aloop=loop=%d:size=%d[loopedaudio]", loopsNeeded-1, loopSize))
		filters = append(filters,
			fmt.Sprintf("[loopedaudio]atrim=0:%g[trimmedaudio]", videoDuration))
	} else {
		filters = append(filters,
			fmt.Sprintf("[concataudio]atrim=0:%g[trimmedaudio]", videoDuration))
	}

	// Volume and fades
	fadeOutStart := math.Max(0, videoDuration-o.FadeOut)
	filters = append(filters,
		fmt.Sprintf("[trimmedaudio]volume=%g,afade=t=in:st=0:d=%g,afade=t=out:st=%g:d=%g[bgmusic]",
			o.Volume, o.FadeIn, fadeOutStart, o.FadeOut))

	// Mix with existing audio or use alone
	hasExistingAudio, err := ffprobeHasAudio(videoPath)
	if err != nil {
		return err
	}

	if hasExistingAudio {
		filters = append(filters, "[0:a][bgmusic]amix=inputs=2:duration=first:dropout_transition=2[outa]")
	} else {
		filters = append(filters, "[bgmusic]acopy[outa]")
	}

	filterComplex := strings.Join(filters, "; ")
	args = append(args, "-filter_complex", filterComplex)
	args = append(args, "-map", "0:v", "-map", "[outa]",
		"-c:v", "copy", "-c:a", "aac", "-b:a", "192k",
		"-movflags", "+faststart",
		"-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}
