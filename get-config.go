// get-config.go
package main

import (
	"os"
	"fmt"
	"bytes"
	"io/ioutil"
	"encoding/xml"
)
type Mark struct {
	XMLName xml.Name `xml:"mark"`
	ID string      `xml:"geodetic-code"`
	Dome string   `xml:"domes-number"`
}
type Site struct {
	XMLName xml.Name `xml:"site"`
	Marks []Mark `xml:"mark"`	
}

 func (m Mark) String() string {
         return fmt.Sprintf("\t Code : %s - Dome : %s \n", m.ID,  m.Dome)
 }
func f(x string) string { 
	xmlFile, err := os.Open(x)
	if err != nil {
		var buffer bytes.Buffer
		buffer.WriteString("Error opening file:")		
		buffer.WriteString(x)
		y := buffer.String()
		return y
	}
	y := xmlFile.Name()
	defer xmlFile.Close()	
	XMLdata, _ := ioutil.ReadAll(xmlFile)
	var s Site
	xml.Unmarshal(XMLdata, &s)
	fmt.Println(s.Marks)
	return y
}
func main() {
	fmt.Println("Hello World!")
	z := f("Site.xml")
	
	fmt.Println(z)
}