//The package contains logic for implementing a filewriter which rotates the files
//at a given state.
package logging

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//The RotatingFileWriter implements the Writer interface with filerotation.
type RotatingFileWriter struct {
	FileName         string
	Size             int64
	MaxNumberOfFiles int
	File             *os.File
}

//The exported Write method. For each Write, it checks the state to decide if its time to rotate the files.
//This is done after the write.
func (w *RotatingFileWriter) Write(buf []byte) (n int, err error) {
	defer w.startClose()
	return executeWrite(buf, w.File)
}

func (w *RotatingFileWriter) OpenFile() (err error) {
	logFilepath := filepath.Dir(w.FileName)

	if _, err := os.Stat(logFilepath); os.IsNotExist(err) {
		os.MkdirAll(logFilepath, 0777)
	}

	w.File, err = os.OpenFile(w.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return err
}

func (w *RotatingFileWriter) CloseFile() (err error) {
	if w.File == nil {
		return nil
	}

	return w.File.Close()
}

//The function implements the statecheck on the logfile. Depending on the state, it can execute rotation of files. This is a defered called after the Write is completed.
func (w *RotatingFileWriter) startClose() (err error) {
	logFile, err := os.Stat(w.FileName)

	if err != nil {
		return errors.New("Failed with file:" + w.FileName + " " + err.Error())
	}

	if logFile.Size() >= w.Size {
		_, err = w.rotateAll()

		if err != nil {
			return errors.New("Failed rotating file: " + logFile.Name() + " " + err.Error())
		}
	}

	return nil
}

//The function implements the rotation of files. Files matching the name of the logfile, are renamed with an index postfixing the filename.
//Active logfile is renamed to <filename>.1 logfile indexed with .2 is renamed to .3 and on until maxNumberOfFiles is reached.
func (w *RotatingFileWriter) rotateAll() (rf []os.FileInfo, err error) {
	err = w.CloseFile()

	if err != nil {
		return nil, errors.New("Failed rotating file" + w.FileName + " " + err.Error())
	}

	dir := filepath.Dir(w.FileName)
	rotatorFileName := filepath.Base(w.FileName)

	sortedFiles := createSortedFileList(dir)
	var rotatedFiles []os.FileInfo

	for _, file := range sortedFiles {
		targetFileName := file.Name()

		if match, _ := filepath.Match(rotatorFileName+"*", targetFileName); match == true {
			if err = w.rotateFile(dir + "/" + targetFileName); err != nil {
				fmt.Println(err.Error())
				return nil, err
			}

			rotatedFiles = append(rotatedFiles, file)
		}
	}

	err = w.OpenFile()

	return rotatedFiles, err
}

//The function handles a single file. It either removes the file, because the number of files are equal maxNumberOfFiles,
//or renames the file accoring to index.
func (w *RotatingFileWriter) rotateFile(fileName string) (err error) {

	fileIndex := extractLogNumber(fileName)

	if fileIndex >= w.MaxNumberOfFiles {
		err = os.Remove(fileName)

		if err != nil {
			return err
		}
		return nil
	}

	fileIndex++
	newFileName := w.FileName + "." + strconv.Itoa(fileIndex)

	return os.Rename(fileName, newFileName)
}

//Using a functionpointer to simplify testing.
var executeWrite = func(buf []byte, file *os.File) (n int, err error) {
	return file.Write(buf)
}

func extractLogNumber(fileName string) int {
	fileExt := filepath.Ext(fileName)

	if logFileNum, e := strconv.Atoi(strings.Replace(fileExt, ".", "", 1)); e == nil {
		return logFileNum
	}

	return 0
}

//The function sorts a list based on the fileindex and returns the sorted array.
func createSortedFileList(dir string) FileInfoArr {
	files, _ := ioutil.ReadDir(dir)
	fileInfoArr := FileInfoArr{}

	for _, f := range files {
		fileInfoArr = append(fileInfoArr, f)
	}

	sort.Sort(sort.Reverse(fileInfoArr))

	return fileInfoArr
}

type FileInfoArr []os.FileInfo

func (d FileInfoArr) Len() int {
	return len(d)
}

func (d FileInfoArr) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d FileInfoArr) Less(i, j int) bool {
	fileNumA := extractLogNumber(d[i].Name())
	fileNumB := extractLogNumber(d[j].Name())

	return fileNumA < fileNumB
}
