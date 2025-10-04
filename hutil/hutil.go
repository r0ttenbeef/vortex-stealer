package hutil

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"vortex/encrypt"
)

// Trace logs and write it into a log file
func LogTrace(mainFolder string, errlog error) {
	logFile, _ := os.OpenFile(filepath.Join(mainFolder, encrypt.B64Util("CrashHandler.log", 0)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println(errlog)
}

// Copy files directly without needing to open the file
func CopyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

// Copy folders from source to destination
func CopyFolder(src string, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyFolder(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Find files recursivly by file extension
func FindFilesExt(path string, fileExt string) ([]string, error) {
	var fileList []string
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	filepath.WalkDir(path, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == fileExt {
			fileList = append(fileList, s)
		}
		return nil
	})
	return fileList, nil
}
