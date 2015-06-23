// get-config.go
package main

import (
	"os"
	"fmt"
	"bytes"
)

func f(x string) string { 
	xmlFile, err := os.Open(x)
	if err != nil {
		var buffer bytes.Buffer
		buffer.WriteString("Error opening file:")		
		buffer.WriteString(x)
//		y := "Error opening file:"
//		y += x		
		y := buffer.String()
		return y
	}
	y := xmlFile.Name()
	defer xmlFile.Close()	
	return y
}
func main() {
	fmt.Println("Hello World!")
	z := f("Sites.xml")
	fmt.Println(z)
}