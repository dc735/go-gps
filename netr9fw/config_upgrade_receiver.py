#! /usr/bin/env python
"""
Python source code - upgrade firmware and/or send GNSS receiver settings to 
                     a specific site equipped with TrimbleNetR9
                     read ip from network repo locally copied
                     read metadata from delta repo locally copied
"""

# Script configuration

import os
from lxml import etree
import ConfigParser
import pycurl
from io import BytesIO
import re
import argparse
import time
import getpass
import subprocess

## Setup argument options

parser = argparse.ArgumentParser()
parser.add_argument("-U","--user", type=str, help="Enter computer username")
parser.add_argument("-u","--credentials", type=str, help="Enter receiver username:password")
parser.add_argument("-t","--test", help="Test a site", action="store_true")
parser.add_argument("-s","--sitecode", type=str,  help="Enter site Code")
parser.add_argument("-c","--clonename", type=str,  help="Enter clone file name to use [gns/linz]")
parser.add_argument("-f","--firmware", type=float,  help="Use this option to update the firmware")
args = parser.parse_args()

username = args.credentials[:args.credentials.index(':')]
password = args.credentials[args.credentials.index(':'):]

## General configurations

config = ConfigParser.ConfigParser()
config.readfp(open(r'/home/'+args.user+'/git/gps-scripts/netr9/config_files/def_nertr9.config'))

clone_file_prefix = config.get('Station', 'clone_file_prefix').replace('$USER', args.user)
ftppushserver = config.get('Station', 'ftppushserver').replace('$USER', args.user)
confdir = config.get('Application', 'confdir').replace('$USER', args.user)
clonefiledir = config.get('Application', 'clonedir').replace('$USER', args.user)
fwfiledir = config.get('Application', 'firmwaredir').replace('$USER', args.user)
networkfile = config.get('Application', 'networkfile').replace('$USER', args.user)
crondir = config.get('Application', 'crontabsdir').replace('$USER', args.user)
antennacombinations = (confdir+'antenna_radome_pairs.txt').replace('$USER', args.user)
fwreleasedate = (confdir+'firmware_release_date.txt').replace('$USER', args.user)

## Functions

def get_ip(site):
    sitegrep = "station="+site
    delay_offset = int('30')
    try:
        return_string=os.popen("grep " + site + " " + networkfile).read() ##TODO : add grep "^gps"
        ip = re.findall( r'[0-9]+(?:\.[0-9]+){3}', return_string )
        return_string=os.popen("grep " + "-h " + sitegrep + " " + crondir + "gpsin_*_crontab").read()
        ftppushdelay = int(re.match(r'\d+', return_string).group())
        ftppushdelay = abs(ftppushdelay - delay_offset)
        ftppushdelay = str(ftppushdelay)
    except:
        ip = "0.0.0.0"
        ftppushdelay = int('10')
        pass
    return ip[0],ftppushdelay

def get_rec_fw(ipaddress):
    try:
       data = BytesIO()
       curl = pycurl.Curl()
       curl.setopt(curl.URL,"http://localhost:8080/prog/Show?FirmwareVersion")
       curl.setopt(curl.USERPWD, username+':'+password)
       curl.setopt(curl.WRITEFUNCTION, data.write)
       curl.perform()
       rec_fw_txt = data.getvalue()
       rec_fw = rec_fw_txt.split("version=")
       rec_fw = rec_fw[1]
       rec_fw = rec_fw.split(" ")
       rec_fw = rec_fw[0]
    except:
        rec_fw = "0.00"
    return rec_fw

def get_rec_fw_warranty(ipaddress):
    try:
       data = BytesIO()
       curl = pycurl.Curl()
       curl.setopt(curl.URL,"http://localhost:8080/prog/Show?FirmwareWarranty")
       curl.setopt(curl.USERPWD, username+':'+password)
       curl.setopt(curl.WRITEFUNCTION, data.write)
       curl.perform()
       rec_fw_warranty_txt = data.getvalue()
       rec_fw_warranty = rec_fw_warranty_txt.split("date=")
       rec_fw_warranty = rec_fw_warranty[1]
       rec_fw_warranty = rec_fw_warranty.replace("\n", "")
#       rec_fw = rec_fw.split(" ")
#       rec_fw = rec_fw[0]
    except:
        rec_fw_warranty = ""
    pass
    return rec_fw_warranty

def get_rectype(site):
    deltasiteinfo = etree.parse("http://magma.geonet.org.nz/services/gps/reception?markCode="+site)
    rectype = deltasiteinfo.xpath('/SITE/cgps-session/receiver/igs-designation')
    deltafirmware = deltasiteinfo.xpath('/SITE/cgps-session/receiver/firmware-history/version')
    lat = deltasiteinfo.xpath('/SITE/location/latitude')
    long = deltasiteinfo.xpath('/SITE/location/longitude')
    height = deltasiteinfo.xpath('/SITE/location/height')
    rec = deltasiteinfo.xpath('/SITE/cgps-session/receiver/igs-designation')
    anttype = deltasiteinfo.xpath('/SITE/cgps-session/installed-cgps-antenna/cgps-antenna/igs-designation')
    antradome = deltasiteinfo.xpath('/SITE/cgps-session/installed-cgps-antenna/radome')
    antoffset = deltasiteinfo.xpath('/SITE/cgps-session/installed-cgps-antenna/height')
    antserialno = deltasiteinfo.xpath('/SITE/cgps-session/installed-cgps-antenna/cgps-antenna/serial-number')
    domenumber = deltasiteinfo.xpath('/SITE/mark/domes-number')
    return (rectype[0].text, deltafirmware[0].text, lat[0].text, long[0].text,
            height[0].text, rec[0].text, anttype[0].text, antradome[0].text, antoffset[0].text,
            antserialno[0].text, domenumber[0].text)

def get_antnum(anttype,antradome,antennacombinations):
    antennacombinations = open(antennacombinations,'r')
    antradpairs = antennacombinations.readlines()
    antennacombinations.close()
    for antrad in antradpairs:
        a = antrad.split()
        if a[0] in anttype and a[1] in antradome:
            trimblenum = a[2]
    return trimblenum

def get_fw_release_date(new_fw):
    f = open(fwreleasedate,'r')
    dates = f.readlines()
    f.close()
    try:
        for fw in dates:
            a = fw.split()
            if a[0] == new_fw:
                release_date = a[1]
                release_date = release_date.replace("\n", "")
        return release_date
    except:
        print "firmware release config file needs updating"
        exit()


def url_systemname(ipaddress,site):
    a = "http://"+ipaddress
    b = "/cgi-bin/sys_name.xml?sysName="+site
    string = (a,b)
    url = "".join(string)
    return url

def url_setantenna(ipaddress,antnum,antoffset,antserialno):
    a = "http://"+ipaddress
    b = "/cgi-bin/antenna.xml?antennaID="+antnum
    c = "&rinexName="+antnum
    d = "&antennaSerial="+antserialno
    e = "&radomeSerial=&antennaMeasMethod=0&antennaHeight="+antoffset
    string = (a,b,c,d,e)
    url = "".join(string)
    return url

def url_refstation(site,ipaddress,lat,long,height,domenumber):
    a = "http://"+ipaddress
    b = "/prog/set?RefStation&lat="+lat
    c = "&lon="+long
    d = "&height="+height
    e = "&Name="+site
    if domenumber == None:
        f = "&Code="+site
    else:
        f = "&Code="+domenumber
    string = (a,b,c,d,e,f)
    url = "".join(string)
    return url

## hard-coded country need to be automated
def url_rnxmetadata(site,ipaddress):
    a = "http://"+ipaddress
    b = "/xml/dynamic/dataLogger.xml"
    c = "?RinexHeaderInfo=set&observerName=Geonet&agencyName=GNS"
    d = "&igsStationInfo="+site+"00NZL"
    string = (a,b,c,d)
    url = "".join(string)
    return url

def url_reboot(ipaddress):
    a = "http://"+ipaddress
    b = "/prog/reset?system"
    string = (a,b)
    url = "".join(string)
    return url

def get_system_name(ipaddress):
    try:
        data = BytesIO()
        curl = pycurl.Curl()
        curl.setopt(curl.URL,"http://localhost:8080/xml/dynamic/sysData.xml")
        curl.setopt(curl.USERPWD, username+':'+password)
        curl.setopt(curl.WRITEFUNCTION, data.write)
        curl.perform()
        sysdata = data.getvalue()
        systemdata = etree.fromstring(sysdata)
        systemdata = systemdata.xpath('/sysData/ownerString1')
        systemname = systemdata[0].text
    except:
        systemname = "XXXX"
    pass
    return systemname

def url_clonefile(ipaddress,dest_fw):
    a = "http://"+ipaddress
    if float(dest_fw) <= 4.70 :
        b = "/cgi-bin/clone_fileUpload.xml"
    else:
        b = "/cgi-bin/clone_fileUpload.html"
    c = "?installCloneFile=true&installStaticIpAddr=false&cloneUploadName="
    string = (a,b,c)
    url = "".join(string)
    return url

def url_ftp_push(ipaddress,ftppushserver,ftppushdelay):
    a = "http://"+ipaddress
    b = "/cgi-bin/ftpPush.xml"
    c = "?request=0&fakePassword=0&serverNumber=0"
    d = "&ftpServer="+ftppushserver
    e = "&ftpPort=21&ftpUser=netr9"
    f = "&ftpPassword=drift"
    g = "&ftpPassword2=drift"
    h = "&delay="+ftppushdelay
    i = "&ftpPath=FTPIN"
    l = "&dirStyle=Flat&rename=No"
    m = "&transferMode=PassiveThenActive&request=0"
    string = (a,b,c,d,e,f,g,h,i,l,m)
    ftp_url = "".join(string)
    return ftp_url

def run_curl(url):
    curl = pycurl.Curl()
    curl.setopt(curl.URL, '%s' % (url))
    curl.setopt(curl.USERPWD, username+':'+password)
    #curl.setopt(curl.VERBOSE, True)
    curl.perform()

def run_curl_upload_clone(url,clonefile):
    curl = pycurl.Curl()
    curl.setopt(curl.POST, 1)
    curl.setopt(curl.URL, '%s' % (url))
    curl.setopt(curl.USERPWD, username+':'+password)
    curl.setopt(curl.HTTPPOST, [("cloneUploadName", (curl.FORM_FILE, clonefile.name))])
    #curl.setopt(curl.VERBOSE, True)
    curl.perform()
    curl.close()

def run_fw_update(ipaddress,fwfile):
    curl = pycurl.Curl()
    curl.setopt(curl.POST, 1)
    curl.setopt(curl.URL,"http://localhost:8080/prog/Upload?FirmwareFile")
    curl.setopt(curl.USERPWD, username+':'+password)
    curl.setopt(curl.HTTPPOST, [("firmwareFile", (curl.FORM_FILE, fwfile.name))])
    print fwfile.name
    print "Starting Firmware file upload to"
    curl.perform()
    curl.close()
    print "Finished file upload"

def show_install_progress(ipaddress):
    data = BytesIO()
    curl = pycurl.Curl()
    curl.setopt(curl.URL,"http://localhost:8080/prog/Show?InstallFirmwareStatus")
    curl.setopt(curl.USERPWD, username+':'+password)
    curl.setopt(curl.WRITEFUNCTION, data.write)
    count = 0
    while count < 40: # Lets give it 10 minutes
        count = count + 1
        time.sleep(15)
        try:
          curl.perform()
          progress = data.getvalue()
          print progress
          if 'Idle' in progress :
            break
          else:
            if 'Done' in progress :
                break
        except:
          print "Rebooting Receiver??"
          pass
        if count == 39:
          print "We Gave up"

# Program execution

if args.sitecode and args.clonename:
    
    site = str.upper(args.sitecode)
    clonefilename = str.upper(args.clonename)
    (rec_type,deltafirmware,lat,long,height,rec,anttype,antradome,antoffset,antserialno,domenumber) = get_rectype(site);
    (ipaddress,ftppushdelay) = get_ip(site)
    sitedata = os.popen("grep " + site + " " + networkfile).read()
    sitename = sitedata[:sitedata.index(',')]
    
    # Establish port forwarding to the receiver using tunneling through avjmp01
    # this will close automatically shortly after the script finishes running
    
    subprocess.Popen('ssh -f avjmp01.geonet.org.nz -L 8080:'+sitename+'.wan.geonet.org.nz:80 sleep 10', shell=True) 
    # use this to check open connections:
#    subprocess.Popen('ps ax | grep ssh | grep -v grep', shell=True) 
    time.sleep(1) # wait for the subshell to close so only one tunnel instance exists
    
    system_name = get_system_name(ipaddress)
    rec_fw = get_rec_fw(ipaddress)
    url_ftp_push_str = url_ftp_push(ipaddress,ftppushserver,ftppushdelay)
    antnum = get_antnum(anttype,antradome,antennacombinations)
    url_setantenna_str = url_setantenna(ipaddress,antnum,antoffset,antserialno)
    url_refstation_str = url_refstation(site,ipaddress,lat,long,height,domenumber)
    url_rnxmetadata_str = url_rnxmetadata(site,ipaddress)
    url_systemname_str = url_systemname(ipaddress,site)
    url_reboot_str = url_reboot(ipaddress)
else:
    print ""
    print "   Please enter a site code with a \" -s XXXX \" "
    print "      and a clone file name with a \" -c [gns/linz] \" "
    print ""
    print "       use -h for options"
    print ""
    exit()
if args.test:
  if rec_type in ('TRIMBLE NETR9'):
    if site in system_name:
      if args.firmware:
	#new_fw =  "%s.2f" % (args.firmware)
        new_fw =  "{0:4.2f}".format(args.firmware)
        rec_fw_warranty = get_rec_fw_warranty(ipaddress)
        release_date = get_fw_release_date(new_fw)
        w_date = time.strptime(rec_fw_warranty,  "%Y-%m-%d")
        r_date = time.strptime(release_date,  "%Y-%m-%d")
	print "TESTING - Update FW to %s" % new_fw
        if w_date < r_date:
	  print "Warranty has Expired. Update not possible without new warranty"
	  exit()
        else:
	  print "Warranty OK"
          print "Receiver firmware is %s" % (rec_fw)
          print "Receiver firmware Warranty is %s" % (rec_fw_warranty)
          print "Firmware release date is %s" % (release_date)
        url_clonefile_str = url_clonefile(ipaddress,new_fw)
        clonefile = open(("%s%s_NETR9_%s.xml" % (clonefiledir, clonefilename, new_fw)), 'r')
	fwfile = open(("%s/NetR9_V%s.timg" % (fwfiledir, new_fw)), 'r')
	print fwfile.name
      else:
        url_clonefile_str = url_clonefile(ipaddress,rec_fw)
        clonefile = open(("%s%s_NETR9_%s.xml" % (clonefiledir, clonefilename, rec_fw)), 'r')
        print "TESTING - Update Reciever config only"
      print "\n--- Uploading clonefile ---"
      print ("%s" % clonefile.name)
      print url_clonefile_str
      print "\n--- Setting ftp-push ---"
      print ("ftp-push server %s - ftp-push delay %s min" % (ftppushserver,ftppushdelay))
      print url_ftp_push_str
      print "\n--- Setting antenna ---"
      print url_setantenna_str
      print "\n--- Setting reference station ---"
      print url_refstation_str
      print "\n--- Setting Rinex Metadata ---"
      print url_rnxmetadata_str
      print "\n--- Setting System Name ---"
      print url_systemname_str
      print "\n--- Rebooting ---"
      print url_reboot_str
    else:
      print "System Name does not match, run update_cronlist.sh or check the site"
  else:
    print "  %s receiver is a Trimble NetRS, must be configured manually" % (site)
else:
  if rec_type in ('TRIMBLE NETR9'):
    if site in system_name:
      print "Receiver firmware is %s" % (rec_fw)
      if args.firmware:
	new_fw =  "{0:4.2f}".format(args.firmware)
        rec_fw_warranty = get_rec_fw_warranty(ipaddress)
        release_date = get_fw_release_date(new_fw)
        w_date = time.strptime(rec_fw_warranty,  "%Y-%m-%d")
        r_date = time.strptime(release_date,  "%Y-%m-%d")
        print "Update FW to %s" % (new_fw)
        if float(rec_fw) <= new_fw:
	    print "Starting FW upgrade"
            if w_date < r_date:
		print "Warranty has Expired. Update not possible without new warranty"
		exit()
	    else:
		print "Warranty OK"
	    fwfile = open(("%s/NetR9_V%s.timg" % (fwfiledir, new_fw)), 'r')
	    run_fw_update(ipaddress,fwfile)
	    show_install_progress(ipaddress)
	    print "--- Firmware upgrade done ---"
            now = time.strftime("%Y-%m-%dT%H:%M:%S",time.gmtime())
            print now
	else:
	    print "New Firmware less then current Firmware"
      	print "\n--- Pause 90 seconds ---"
	time.sleep(90)
      else:
        print "Update Receiver config only"
        #print "------------ Starting site configuration --------"
      print "\n--- Uploading clonefile ---"
      rec_fw = get_rec_fw(ipaddress)
      print "current rec_fw is %s" % (rec_fw)
      url_clonefile_str = url_clonefile(ipaddress,rec_fw)
      print url_clonefile_str
      clonefile = open(("%s/%s_NETR9_%s.xml" % (clonefiledir, clonefilename, rec_fw)), 'r')
      print ("%s" % clonefile.name)
      run_curl_upload_clone(url_clonefile_str,clonefile)
      print "\n--- Setting ftp-push ---"
      print ("ftp-push server %s - ftp-push delay %s min" % (ftppushserver,ftppushdelay))
      run_curl(url_ftp_push_str)
      print "\n--- Setting antenna ---"
      antnum = get_antnum(anttype,antradome,antennacombinations)
      run_curl(url_setantenna_str)
      print "\n--- Setting reference station ---"
      run_curl(url_refstation_str)
      print "\n--- Setting Rinex Metadata ---"
      run_curl(url_rnxmetadata_str)
      print "\n--- Setting System Name ---"
      run_curl(url_systemname_str)
      print "\n--- Rebooting ---"
      run_curl(url_reboot_str)
    else:
      print "System Name does not match, run update_cronlist.sh or check the site"
  else:
    print "  %s receiver is a Trimble NetRS, must be configured manually" % (site)
