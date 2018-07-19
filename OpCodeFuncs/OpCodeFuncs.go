// Package opcodefuncs - corresponding functions for each DAG shape
// Author: Paul Duncanson 07/15/2018
package opcodefuncs

import (
	"fmt"
	"math/rand"
	"time"
)

// GetNameToFunc - Returns association between a function name to the actual function
func GetNameToFunc() map[string]func(string, chan<- string) {
	seed := rand.NewSource(time.Now().UnixNano())
	rndm := rand.New(seed)

	nameToFunc := make(map[string]func(string, chan<- string))

	// SaveFile - Saves the specified file and returns success in channel
	saveFile := func(fileName string, message chan<- string) {
		time.Sleep(4 * time.Second)
		message <- "LastShape"
	}
	nameToFunc["SaveFile"] = saveFile

	// DoOCRToText - Converts OCR to Text using Tesseract
	doOCRToText := func(ocrName string, message chan<- string) {
		time.Sleep(4 * time.Second)
		message <- "id1[docStr]"
	}
	nameToFunc["DoOCRToText"] = doOCRToText

	// DoPDFToText - Converts PDF to Text using Tesseract
	doPDFToText := func(pdfName string, message chan<- string) {
		time.Sleep(3 * time.Second)
		message <- "id4[docStr]"
	}
	nameToFunc["DoPDFToText"] = doPDFToText

	// SaveStringToFile - Saves provided string to a flat file
	saveStringToFile := func(stringToSave string, message chan<- string) {
		time.Sleep(3 * time.Second)
		message <- "LastShape"
	}
	nameToFunc["SaveStringToFile"] = saveStringToFile

	// GetNextEmail - Retrieves next email from email servers' specified mail folder
	getNextEmail := func(notUsed string, message chan<- string) {
		time.Sleep(4 * time.Second)
		message <- "id5[emailXML]"
	}
	nameToFunc["GetNextEmail"] = getNextEmail

	// GetFileType - Sends File Type of specified file into a channel for caller access
	getFileType := func(fileName string, message chan<- string) {
		time.Sleep(2 * time.Second)

		var a [4]string
		a[0] = "id3[“pdf”]"
		a[1] = "id3[(ELSE)]"
		a[2] = "id3[“tiff”, ”jpg”, ”png”, ”jfif”, ”bmp”]"
		a[3] = "id3[“doc”,”docx”]"

		fmt.Printf("Random return value from getFileType: '%s'\n", a[rndm.Intn(3)])
		message <- a[rndm.Intn(3)]
	}
	nameToFunc["GetFileType"] = getFileType

	// GetNextAttachment - Retrieves next attachment from specified email
	getNextAttachment := func(notUsed string, message chan<- string) {
		time.Sleep(2 * time.Second)
		message <- "id6[docName]"
	}
	nameToFunc["GetNextAttachment"] = getNextAttachment

	// DoDocToText - Converts Word Doc format to Text
	doDocToText := func(docName string, message chan<- string) {
		time.Sleep(4 * time.Second)
		message <- "id7[docStr]"
	}
	nameToFunc["DoDocToText"] = doDocToText

	return nameToFunc
}
