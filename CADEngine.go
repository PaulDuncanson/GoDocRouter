// main - Engine that concurrently invokes DAG processes
// Author: Paul Duncanson 07/15/2018
package main

import (
	opCodeFuncs "DocRouter/opcodefuncs"
	opCodeML "DocRouter/opcodeml"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// BuildShapeToFuncMap - binds corresponding function with each shape id.
//   input:   nameToShape - Key funcName with shapeID Value
//   returns: Key shapeID with function pointer Value
func BuildShapeToFuncMap(nameToShape map[string]string) map[string]func(string, chan<- string) {

	nameToFunc := opCodeFuncs.GetNameToFunc()
	shapeToFunc := make(map[string]func(string, chan<- string))

	for eachNameToShapeKey, eachNameToShapeValue := range nameToShape {
		shapeToFunc[eachNameToShapeValue] = nameToFunc[eachNameToShapeKey]
	}

	return shapeToFunc
}

// CADEngine - GoRoutine scheduler
func CADEngine(dag opCodeML.DAG, opCodeTable map[string]string, startShape string) int {

	// !!!!!Run tests to verify sweet spot for number of channels in relation to CPUs
	//cpuCount := runtime.NumCPU()
	//fmt.Println("CPU Count is ", cpuCount)

	nameToShape, shapeToName := opCodeML.BuildNameAndShapeMaps(opCodeTable)
	shapeToFunc := BuildShapeToFuncMap(nameToShape)

	shapeToArgs := opCodeML.BuildShapeToArgs(dag)
	shapeToToConnectors := opCodeML.BuildShapeToToConnectors(dag)

	for shapeToToConnectorsKey, shapeToToConnectorsValue := range shapeToToConnectors {
		fmt.Printf("shapeToToConnectors[%s] is %s\n", shapeToToConnectorsKey, shapeToToConnectorsValue)
	}

	quit := make(chan bool)
	message := make(chan string)

	maxCount := 20 // number of email attachements to process
	currentShapeID := startShape

	fmt.Printf("======================== Setup Complete ========================\n")

	go func() {

		fmt.Printf("'%s', '%s'(%s)\n", currentShapeID, shapeToName[currentShapeID], shapeToArgs[currentShapeID])
		go shapeToFunc[currentShapeID](shapeToArgs[currentShapeID], message)

		// While count is less then maxCount trips through DAG
		for cnt := 0; cnt < maxCount; {

			select {
			case msg := <-message:
				fmt.Printf("returned - '%s'\n", msg)
				if msg != "LastShape" {
					currentShapeID = shapeToToConnectors[msg]
				} else {
					cnt++
					fmt.Printf("Processed %d attachments\n", cnt)
					currentShapeID = startShape
				}
				fmt.Printf("'%s', '%s'(%s)\n", currentShapeID, shapeToName[currentShapeID], shapeToArgs[currentShapeID])

				go shapeToFunc[currentShapeID](shapeToArgs[currentShapeID], message)
			}
		}
		fmt.Printf("============= Successfully performed %d DAG cycles =============\n", maxCount)
		quit <- true
	}()

	<-quit // Block main thread until CADEngine has processed maxCount attachments

	return 0
}

// Command-line flags.
var (
	xmlFileName = flag.String("xmlname", "./EmailDocRouter.xml", "Filename and path of Libre Office DAG file.")
	//emailbox    = flag.String("emailbox", "Inbox", "eMail Box to perform DocRouting of attached emails")
	//emailserver = flag.String("emailserver", "imap.gmail.com", "Email Server")
	//version     = flag.String("version", "1.10", "Go version")
	//email       = flag.String("email", "pbduncanson@gmail.com", "email")
	//password    = flag.String("password", "<PASSWORD>", "password")
)

func main() {
	fmt.Printf("========================= Setup Start ==========================\n")
	flag.Parse()

	// For random channel communications from the OpCode routine, getFileType, that triggers the fan-out.
	rand.Seed(time.Now().UTC().UnixNano())

	// if flag.NFlag() < 4 { !!!!! Implement and add command line support for:
	//						 !!!!! start invocation at a given shapeID,
	//						 !!!!! shape highlight feature for run-time demonstration,
	//						 !!!!! debug mode,
	//						 !!!!! trace and performance measurements, etc.
	/*
		if flag.NFlag() < 0 {
			fmt.Printf("Usage: %s [options]\n", os.Args[0])
			fmt.Println("Options:")
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/
	// Retrieve DAG and OpCodeTable from Libre Office Draw
	dag, opCodeTable := opCodeML.OpCodeML(*xmlFileName)

	fmt.Println(dag, opCodeTable)

	// Transfer control to CADEngine that will distribute work
	// to corresponding GoRoutines as illustrated in the DAG
	log.Fatal(CADEngine(dag, opCodeTable, "id5"))

}
