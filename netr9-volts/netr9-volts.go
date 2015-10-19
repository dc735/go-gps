package main

import (
	"fmt"
	"log"
//	"bufio"
//	"os"
	"strings"
    "net/http"
	"net"
    "io/ioutil"
)

import (
        "github.com/zone"
)
func check(e error) {
    if e != nil {
        panic(e)
    }
}
func gv(ip net.IP) string {
		response, err := http.Get("http://admin:cuseeme@"+ip.String()+"/prog/show?Voltages")
    if err != nil {
		return "volts=0.00"
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
			return "volts=0.00"
        }
		result := strings.Split(string(contents),"\n")
	    for i := range result {
			str := "port=2"
			if strings.Contains(result[i],str) {
				x := strings.Split(result[i], " ")
				for j := range x {
					str := "volts="
					if strings.Contains(x[j],str) {
					return x[j]
					}
				}
					
			}
		}
    }	
return "volts=0.00"
}
func main() {

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
        list, err = z.MatchByModel("^Trimble NetR9")
        if err != nil {
                log.Fatal(err)
        }
/*
    	fn, err := os.Create("netr9-noconn.cfg")
    	check(err)
		defer fn.Close()
    	f, err := os.Create("netr9-conn.cfg")
    	check(err)
		defer f.Close()
*/
        for _, e := range list {
				s := (strings.Replace(e.Place, " ", "-", -1))
				x := gv(e.IP)				
				fmt.Printf("%s %s %s %s #\n", e.IP, s, x,e.Code)
/*
				w := bufio.NewWriter(f)
				o := fmt.Sprintf("%s %s %s #\n", e.IP, s, x)
				_, err := w.WriteString(o)
				check(err)
				w.Flush()
				w = bufio.NewWriter(fn)
				o = fmt.Sprintf("%s %s # noconn\n", e.IP, s)
				_, err = w.WriteString(o)
				check(err)
				w.Flush()
*/
        }
}
