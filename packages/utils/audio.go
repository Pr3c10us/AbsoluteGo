package utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
)

func WriteWAV(filePath string, base64Audio string, sampleRate uint32, channels uint16, bitsPerSample uint16, extraSilenceSecs float64) error {
	pcm, err := base64.StdEncoding.DecodeString(base64Audio)
	if err != nil {
		return fmt.Errorf("decode base64: %w", err)
	}

	return WriteWAVFromPCM(filePath, pcm, sampleRate, channels, bitsPerSample, extraSilenceSecs)
}

func WriteWAVFromPCM(filePath string, pcm []byte, sampleRate uint32, channels uint16, bitsPerSample uint16, extraSilenceSecs float64) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	byteRate := sampleRate * uint32(channels) * uint32(bitsPerSample) / 8
	blockAlign := channels * bitsPerSample / 8

	silenceBytes := uint32(extraSilenceSecs * float64(byteRate))
	dataSize := uint32(len(pcm)) + silenceBytes

	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, uint32(36+dataSize))
	f.Write([]byte("WAVE"))

	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, channels)
	binary.Write(f, binary.LittleEndian, sampleRate)
	binary.Write(f, binary.LittleEndian, byteRate)
	binary.Write(f, binary.LittleEndian, blockAlign)
	binary.Write(f, binary.LittleEndian, bitsPerSample)

	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, dataSize)
	_, err = f.Write(pcm)
	if err != nil {
		return fmt.Errorf("write pcm data: %w", err)
	}

	if silenceBytes > 0 {
		silence := make([]byte, silenceBytes)
		_, err = f.Write(silence)
		if err != nil {
			return fmt.Errorf("write silence: %w", err)
		}
	}

	return nil
}

func BufDuration(buf []byte, rate, ch, bps int) float64 {
	if rate == 0 {
		rate = 24000
	}
	if ch == 0 {
		ch = 1
	}
	if bps == 0 {
		bps = 2
	}
	return float64(len(buf)) / float64(rate*ch*bps)
}
