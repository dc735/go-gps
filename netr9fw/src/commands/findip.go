package commands

import (
	"net"
	"github.com/zone"
	"log"
//    "io/ioutil"
)



func FindIP(site string) (ip net.IP) {
	var SiteIP  net.IP
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
        // And match the Trimble NetR9's with a Site code
        list, err = z.MatchByModelAndCode("^Trimble",site)
        if err != nil {
                log.Fatal(err)
        }
		for _, e := range list {
			SiteIP = e.IP
		}
return SiteIP
}