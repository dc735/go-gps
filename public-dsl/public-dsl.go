package main

import (
	"fmt"
	"log"
	"bufio"
    "os"
	"strings"
)

import (
        "github.com/zone"
)
func check(e error) {
    if e != nil {
        panic(e)
    }
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
        // all Hongdian Cellular Modem's
        list, err = z.MatchByModel("^Spark")
        if err != nil {
                log.Fatal(err)
        }
    	fn, err := os.Create("spark-public-noconn.cfg")
    	check(err)
		defer fn.Close()
    	f, err := os.Create("spark-public-conn.cfg")
    	check(err)
		defer f.Close()
        for _, e := range list {
				s := (strings.Replace(e.Place, " ", "-", -1))
				fmt.Printf("%s %s #\n", e.IP, s)
				w := bufio.NewWriter(f)
				o := fmt.Sprintf("%s %s #\n", e.IP, s)
				_, err := w.WriteString(o)
				check(err)
				w.Flush()
				w = bufio.NewWriter(fn)
				o = fmt.Sprintf("%s %s # noconn\n", e.IP, s)
				_, err = w.WriteString(o)
				check(err)
				w.Flush()
        }
}
