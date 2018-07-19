// Package opcodeml - builds utility opcode arrays out of DAG XML data
// Author: Paul Duncanson 07/15/2018
package opcodeml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// DAG - Directed Acyclic Graph
//    Shapes     : Array of shapes representing each GoRoutine
//       ShapeID : Unique shape identifier assigned by Libre Office Draw
//       Values  : Comments and Arguments placed in each shape
//    Connectors : Array of connectors representing a Go channel
//       Values  : Messages being passed between each GoRoutine
type DAG struct {
	XMLName xml.Name `xml:"document"`

	Shapes []struct {
		ShapeID string   `xml:"id,attr"`
		Values  []string `xml:"shptext"`
	} `xml:"body>drawing>page>custom-shape"`

	Connectors []struct {
		Value     string `xml:"shptext"`
		FromShape string `xml:"start-shape,attr"` // Shape ID that connector is connected from
		ToShape   string `xml:"end-shape,attr"`   // Shape ID that connector is connected to
	} `xml:"body>drawing>page>connector"`
}

// OpCodeML - Extracts OpCode from provided DAG file name.
//    Arguments:
//       dagFileName - Name of Directed Acyclic Graph
//                     file created by Libre Office Draw
//
//    Returns:
//       DAG:
//            Struct that contains shapes and connectors
//            where GoFunc names, GoFunc Arguments and
//            GoFunc Channel Messages are defined.
//       map[string]string:
//            OpCodeTable comprised of each
//            unique instance of functions found in the
//            provided DAG file.
func OpCodeML(dagFileName string) (DAG, map[string]string) {
	xmlFile, err := os.Open(dagFileName)
	check(err)

	defer xmlFile.Close()

	// read xmlFile as a byte array
	xmlData, err := ioutil.ReadAll(xmlFile)
	check(err)

	// Initialize Directed Acyclic Graph
	v := DAG{}

	// Populate struct from xml nodes of interest
	check(xml.Unmarshal(xmlData, &v))

	// Build OpCode table
	opCodeTable := make(map[string]string)

	for i := 0; i < len(v.Shapes); i++ {
		opCodeTable[v.Shapes[i].ShapeID] = v.Shapes[i].Values[len(v.Shapes[i].Values)-1]
	}

	fmt.Printf("Successfully extracted nodes of interest from %s\n%s\n%s\n", dagFileName, v, opCodeTable)
	return v, opCodeTable

}

// BuildNameAndShapeMaps - returns k,v pairs for optimal access
// to and from funcName and shapesID from provided OpCode Table
func BuildNameAndShapeMaps(opCodeTable map[string]string) (map[string]string, map[string]string) {

	// shapeToName - Key shapeID, Value funcName
	// nameToShape - Key funcName, Value shapeID
	nameToShape := make(map[string]string)
	shapeToName := make(map[string]string)
	for shapeID, shapeArgs := range opCodeTable {
		// First string in shapeArgs is the function name
		// (e.g. [doPDFToText, docName])
		arrayOfParams := strings.Split(shapeArgs, ",")
		var funcName = ""
		if len(arrayOfParams) >= 1 {
			funcName = strings.Trim(arrayOfParams[0], " []")
		}

		shapeToName[shapeID] = funcName
		nameToShape[funcName] = shapeID
		fmt.Printf("nameToShape[%s] is %s\n", shapeID, funcName)
	}
	return nameToShape, shapeToName
}

// BuildShapeToArgs - Creates mapping between shapeID and input argument string
func BuildShapeToArgs(dag DAG) map[string]string {
	shapeToArgs := make(map[string]string)
	for i := 0; i < len(dag.Shapes); i++ {
		shapeToArgs[dag.Shapes[i].ShapeID] = dag.Shapes[i].Values[len(dag.Shapes[i].Values)-1]
	}
	return shapeToArgs
}

// BuildShapeToToConnectors - Creates mapping of current shapeID + returnValue to ToShape
func BuildShapeToToConnectors(dag DAG) map[string]string {
	shapeToToConnectors := make(map[string]string)
	for i := 0; i < len(dag.Connectors); i++ {
		shapeToToConnectors[dag.Connectors[i].FromShape+dag.Connectors[i].Value] = dag.Connectors[i].ToShape
		fmt.Printf("shapeToToConnector %s, %s\n", dag.Connectors[i].FromShape+dag.Connectors[i].Value, dag.Connectors[i].ToShape)
	}
	return shapeToToConnectors
}

func check(e error) {
	if e != nil {
		fmt.Println("Panic level error: ", e)
		log.Fatal("Panic level error from opcodeml package: ", e)
		//panic(e) // for the stack as needed
	}
}

/*  Uncomment for Unit testing.  !!!!! Convert main to a test routine for use with >go test TestOpCodeML
func main() {

	v, opCodeTable := OpCodeML("EmailDocRouter.xml")

	fmt.Println("=======")
	for i := 0; i < len(v.Shapes); i++ {
		fmt.Printf("Shape Id: %s\n", v.Shapes[i].ShapeID)
		for j := 0; j < len(v.Shapes[i].Values); j++ {
			fmt.Printf("Shapes[%d].Values[%d].Value is %s\n", i, j, v.Shapes[i].Values[j])
		}
		fmt.Println("-------")
	}

	for i := 0; i < len(v.Connectors); i++ {
		fmt.Printf("Connector Message  : %s\n", v.Connectors[i].Value)
		fmt.Printf("Connector FromShape: %s\n", v.Connectors[i].FromShape)
		fmt.Printf("Connector ToShape  : %s\n", v.Connectors[i].ToShape)
		fmt.Println("-------")
	}
	fmt.Println("========")

	for i := 0; i < len(v.Shapes); i++ {
		opCodeTable[v.Shapes[i].ShapeID] = v.Shapes[i].Values[len(v.Shapes[i].Values)-1]
		fmt.Printf("opCodeTable[%s] is %s\n", v.Shapes[i].ShapeID, v.Shapes[i].Values[len(v.Shapes[i].Values)-1])
	}

	fmt.Printf("opCodeTable[%s] is %s\n", "id1", opCodeTable["id1"])
	fmt.Printf("opCodeTable[%s] is %s\n", "id2", opCodeTable["id2"])
	fmt.Printf("opCodeTable[%s] is %s\n", "id8", opCodeTable["id8"])
}
*/
