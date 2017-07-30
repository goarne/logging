package logging

import (
	"io/ioutil"
	"os"
	"testing"
)

//The class performs vaious tests on the exported functions for the RotatingFileWriter
var (
	size                int64 = 20
	filename                  = "test.log"
	maxNumberOfFiles          = 2
	logFile                   = os.File{}
	logEntry1                 = "This is loggentry: 1"
	logEntry2                 = "This is loggentry: 2"
	logEntry3                 = "This is loggentry: 3"
	logEntry4                 = "This is loggentry: 4"
	expectedLogile1           = "test.log.1"
	expectedLogile2           = "test.log.2"
	nonExistingLogfile3       = "test.log.3"
)

func TestWriteAndRotate(t *testing.T) {
	defer func() {
		os.Remove(filename)
		os.Remove(expectedLogile1)
		os.Remove(expectedLogile2)
	}()

	fr := RotatingFileWriter{FileName: filename, Size: size, MaxNumberOfFiles: maxNumberOfFiles, File: &logFile}
	fr.OpenFile()

	//Writing to logfiles which shall cause 3 rotations.
	//The congtents of the first logfile, will be removed.
	if _, err := fr.Write([]byte(logEntry1)); err != nil {
		t.Error(err.Error())
	}

	if _, err := fr.Write([]byte(logEntry2)); err != nil {
		t.Error(err.Error())
	}

	if _, err := fr.Write([]byte(logEntry3)); err != nil {
		t.Error(err.Error())
	}

	if _, err := fr.Write([]byte(logEntry4)); err != nil {
		t.Error(err.Error())
	}

	if err := fr.CloseFile(); err != nil {
		t.Error(err.Error())
	}

	//Checking that whats expecting to exist actually exists.
	verifyFileExist(t, filename)
	verifyFileExist(t, expectedLogile1)
	verifyFileExist(t, expectedLogile2)

	verifyFileContent(t, expectedLogile1, logEntry4)
	verifyFileContent(t, expectedLogile2, logEntry3)

	//Checking that rotation does not create more files than expected.
	verifyFileNotExist(t, nonExistingLogfile3)

}

//Functions used for verifying the core functionallity of the RotatingFileWriter.
func verifyFileNotExist(t *testing.T, fileName string) {
	_, err := os.Stat(fileName)

	if err == nil {
		t.Error("Expected error because file should not exist.")
	}
}

func verifyFileContent(t *testing.T, fileName string, content string) {
	fileContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		t.Error(err.Error())
	}

	if content != string(fileContent) {
		t.Errorf("Filecontent '%s' should equal source content '%s'", string(fileContent), content)
	}
}

func verifyFileExist(t *testing.T, fileName string) {
	_, err := os.Stat(fileName)

	if err != nil {
		t.Error(err.Error())
	}
}
