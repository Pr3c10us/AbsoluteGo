package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Pr3c10us/absolutego/packages/configs"
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
