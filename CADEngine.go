// main - Engine that drives invocation of DAG processes
// Author: Paul Duncanson 07/15/2018
package main

import (
	forceExport "DocRouter/forceexport"
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
func BuildShapeToFuncMap(nameToShape map[string]string) map[string]func(interface{}, string) {
	shapeToFunc := make(map[string]func(interface{}, string))

	message := make(chan string)

	// Need to explicitly call each of the operands to force inclusion of routines in symbolic link table.
	// CAD Engine does not, by design, statically invoke opCode routines as invocation is determined at run time.
	// Locate Linker flag to override optimization at specific package level when there's time!!!!!
	opCodeFuncs.SaveFile("", message)
	opCodeFuncs.DoOCRToText("", message)
	opCodeFuncs.DoPDFToText("", message)
	opCodeFuncs.SaveStringToFile("", message)
	opCodeFuncs.GetNextEmail("", message)
	opCodeFuncs.GetFileType("", message)
	opCodeFuncs.GetNextAttachment("", message)
	opCodeFuncs.DoDocToText("", message)

	for funcName, shapeID := range nameToShape {
		funcNameWPkg := "DocRouter/opcodefuncs." + funcName
		fmt.Printf("Going to search for %s\n", funcNameWPkg)

		var funcPtr func(interface{}, string)

		err := forceExport.GetFunc(&funcPtr, funcNameWPkg)
		if err == nil {
			shapeToFunc[shapeID] = funcPtr
			fmt.Println("Found:", funcNameWPkg, funcName)
		} else {
			fmt.Println("Missing function specified in DAG:", funcNameWPkg, funcName)
			//log.Fatal("Missing function specified in DAG as:", funcName, " for ", funcNameWPkg) // !!!!! for production
		}
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

	fmt.Printf("%s\n%s\n", nameToShape, shapeToName)

	shapeToArgs := opCodeML.BuildShapeToArgs(dag)
	shapeToToConnectors := opCodeML.BuildShapeToToConnectors(dag)

	fmt.Println(shapeToToConnectors)
	//toConnectorToMsgs := opCodeML.BuildToConnectorToArgs(dag.Connectors)

	quit := make(chan bool)
	message := make(chan string)
	count := make(chan string)

	maxCount := 100 // number of email attachements to process
	currentShapeID := startShape

	type OpCodeFuncType func(interface{}, string)
	var opCodeFunc OpCodeFuncType

	go func() {

		// While count is less then maxCount trips through DAG
		for cnt := 0; cnt < maxCount; {

			select {
			case msg := <-message:
				currentShapeID = shapeToToConnectors[msg]
				opCodeFunc = shapeToFunc[currentShapeID]
				go opCodeFunc(nil, shapeToArgs[currentShapeID])
			case <-count:
				cnt++
				currentShapeID = startShape
				opCodeFunc = shapeToFunc[currentShapeID]
				go opCodeFunc(nil, shapeToArgs[currentShapeID])
			}
		}
		quit <- true
	}()

	<-quit // Hold up main thread until CADEngine loop has processed maxCount attachments
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
