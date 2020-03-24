package filefactory

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// Copy file from src to dest
func CopyFile(srcFile, destFile string) {
	file, err := os.Open(srcFile)
	if err != nil {
		log.Warn().Timestamp().Msg("Could not copy file")
		panic(err)
	}
	defer file.Close()
	dest, err := os.Create(destFile)
	if err != nil {
		log.Warn().Timestamp().Msg("Could not copy file")
		panic(err)
	}
	defer dest.Close()
	_, err = io.Copy(dest, file)
	if err != nil {
		log.Warn().Timestamp().Msg("Could not copy file")
		panic(err)
	}
}

// create the directory of filename.
func createDirectory(filename string) {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Warn().Timestamp().
			Msg("Could not create log directory.")
		panic(err)
	}
}

// create logfile and open
func NewFile(filename string) *os.File {
	createDirectory(filename)

	newFile, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Warn().Timestamp().
			Msg("Could not create log file.")
		panic(err)
	}
	return newFile
}
