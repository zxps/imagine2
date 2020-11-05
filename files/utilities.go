package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// RemoveEmptyDirectories ...
func RemoveEmptyDirectories(targetPath string) {
	filepath.Walk(targetPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			return nil
		}

		currentTime := time.Now()
		fileTime := f.ModTime()
		todayCreated := false
		if currentTime.Year() == fileTime.Year() {
			if currentTime.Month() == fileTime.Month() {
				if currentTime.Day() == fileTime.Day() {
					todayCreated = true
				}
			}
		}

		if !todayCreated {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return nil
			}

			if len(files) < 1 {
				os.RemoveAll(path)
			}
		}

		return nil
	})

}
