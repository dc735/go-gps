package main

import (
	"fmt"
	"log"
	"flag"
    "github.com/zone"
	"strings"
//		"bufio" 
//  	"net/http"
//		"net"
//  	"io/ioutil"
//	"github.com/delta/meta"

)
//var stn string
//func init(){
//}

func file_extension(fx *string){
	*fx = ".xml"
}
//func mark(m meta.Marks, s *meta.Stations) bool {
//	err := store(m, s)
//	if err != nil {
//		log.Fatal(err)
//	return ok
//	}	
//}
func main() {
	stnCode := flag.String("stn", "TEST", "Enter a stataion code")
	flag.Parse()
	//	fmt.Println("Station input is:", *stnCode)
	_, stn := fmt.Printf(strings.ToUpper(*stnCode)) 
//	fp := "/home/gpsin/sites/"
//	fn := "akto.xml"
//	_, stnx := fmt.Printf(strings.ToLower(sN))
//	_ = f(fp+fn)
	fmt.Println(stn)
//	fmt.Println(Mark.ID)

        // default geonet connection ....
        z := zone.Equipment{
                Zone:   "wan.geonet.org.nz.",
                Server: "rhubarb.geonet.org.nz",
                Port:   "53",
        }

        // get all equipment ...
        list, err := z.List()
        if err != nil {
                log.Fatal(err)
        }
        // all Trimble NetR9's
        list, err = z.MatchByModelAndCode("^Trimble NetR9",*stnCode)
        if err != nil {
                log.Fatal(err)
        }
		fmt.Println(list)	

}