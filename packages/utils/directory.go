package utils

import (
	"fmt"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
	"unicode"
)

func GetDirectory(dirName string) (string, error) {
	suffix := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), rand.Intn(1000000000))
	dir := configs.GetRootPath() + "/" + dirName + "/" + suffix

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}

func GetSubDirs(dirPath string) ([]string, error) {
	var subDirs []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, filepath.Join(dirPath, entry.Name()))
		}
	}

	return subDirs, nil
}

func GetFilesInDir(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return naturalLess(filepath.Base(files[i]), filepath.Base(files[j]))
	})

	return files, nil
}

func naturalLess(a, b string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		aIsDigit := unicode.IsDigit(rune(a[i]))
		bIsDigit := unicode.IsDigit(rune(b[i]))

		if aIsDigit && bIsDigit {
			aNum, aLen := extractNumber(a[i:])
			bNum, _ := extractNumber(b[i:])

			if aNum != bNum {
				return aNum < bNum
			}
			i += aLen - 1
		} else if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func extractNumber(s string) (int, int) {
	i := 0
	for i < len(s) && unicode.IsDigit(rune(s[i])) {
		i++
	}
	num, _ := strconv.Atoi(s[:i])
	return num, i
}
