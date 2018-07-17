// Package opcodefuncs - corresponding functions for each DAG shape
// Author: Paul Duncanson 07/15/2018
package opcodefuncs

import (
	"fmt"
	"math/rand"
	"time"
)

// SaveFile - Saves the specified file and returns success in channel
func SaveFile(fileName string, count chan<- string) {
	time.Sleep(6 * time.Second)
	fmt.Println("SaveFile", fileName)
	count <- "1"
}

// DoOCRToText - Converts OCR to Text using Tesseract
func DoOCRToText(ocrName string, message chan<- string) {
	time.Sleep(20 * time.Second)
	fmt.Println("DoOCRToText", ocrName)
	message <- "id1[docStr]"
}

// DoPDFToText - Converts PDF to Text using Tesseract
func DoPDFToText(pdfName string, message chan<- string) {
	time.Sleep(10 * time.Second)
	fmt.Println("DoPDFToText", pdfName)
	message <- "id4[docStr]"
}

// SaveStringToFile - Saves provided string to a flat file
func SaveStringToFile(stringToSave string, count chan<- string) {
	time.Sleep(3 * time.Second)
	fmt.Println("SaveStringToFile", stringToSave)
	count <- "1"
}

// GetNextEmail - Retrieves next email from email servers' specified mail folder
func GetNextEmail(notUsed string, message chan<- string) {
	time.Sleep(4 * time.Second)
	fmt.Println("GetNextEmail", notUsed)
	fmt.Println("id5[emailXML]")
	message <- "id5[emailXML]"
}

// GetFileType - Sends File Type of specified file into a channel for caller access
func GetFileType(fileName string, message chan<- string) {
	time.Sleep(2 * time.Second)
	fmt.Println("GetFileType", fileName)

	var a [4]string
	a[0] = "id3[\"pdf\"]"
	a[1] = "id3[(ELSE)]"
	a[2] = "id3[\"tiff\", \"jpg\", \"png\", \"jfif\", \"bmp\"]"
	a[3] = "id3[\"doc\",\"docx\"]"

	// Get random number from the channel instead of re-generating the seeded value
	seed := rand.NewSource(time.Now().UnixNano())
	rndm := rand.New(seed)

	fmt.Println(a[rndm.Intn(3)])
	message <- a[rndm.Intn(3)]
}

// GetNextAttachment - Retrieves next attachment from specified email
func GetNextAttachment(notUsed string, message chan<- string) {
	time.Sleep(2 * time.Second)
	fmt.Println("GetNextAttachment", notUsed)
	fmt.Println("id6[docName]")
	message <- "id6[docName]"
}

// DoDocToText - Converts Word Doc format to Text
func DoDocToText(docName string, message chan<- string) {
	time.Sleep(4 * time.Second)
	fmt.Println("DoDocToText", docName)
	fmt.Println("id7[docStr]")
	message <- "id7[docStr]"
}
