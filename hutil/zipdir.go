package hutil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
)

func archiveFiles(w *zip.Writer, folderName string, folderPath string, fileNameToExclude string) error {
	walker := func(path string, info os.FileInfo, err error) error {
		var pathInZip string

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.Contains(path, fileNameToExclude) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			re := regexp2.MustCompile(fmt.Sprintf(`.+?(?=%s)`, folderName), 0)
			if zipDir, _ := re.FindStringMatch(path); zipDir != nil {
				pathInZip = zipDir.String()
			}

			f, err := w.Create(strings.TrimPrefix(path, pathInZip))
			if err != nil {
				return err
			}
			if _, err = io.Copy(f, file); err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.Walk(folderPath, walker); err != nil {
		return err
	}

	return nil
}

// Archive folder and subfolders inside a zip file
func ZipFolder(folderName string, folderPath string, zipfilePath string, fileNameToExeclude string) error {

	zipFile, err := os.Create(zipfilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err = archiveFiles(zipWriter, folderName, folderPath, fileNameToExeclude); err != nil {
		return err
	}

	return nil
}
