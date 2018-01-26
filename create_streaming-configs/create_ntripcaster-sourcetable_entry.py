#! /usr/bin/python
# Filename: delta2siteinfo.py

############################################################################
#### create BKG ntripcaster sourcetable entry based on Delta info

#### TODO: 
####    - modify hardcoded ISO country code (string)
####    - modify hardcoded receiver type (string)
############################################################################

# vim: tabstop=8 expandtab shiftwidth=4 softtabstop=4 
import glob
import os
import pycurl
import cStringIO
import ConfigParser
import time
from lxml import etree
import matplotlib.path as mplPath
import numpy as np


config = ConfigParser.ConfigParser()
config.readfp(open(r'/home/elidana/git/geonet/gps-scripts/sitelogs/delta2siteinfo.config'))
delta = config.get('General','deltadatabase')
urlnetwork = config.get('General','deltanetwork')
staticname = config.get('LocalDirs','deltastatic')


def check_date(date):
    return date != 'open'

def check_dateforxml(date):
    if date == 'CCYY-MM-DDThh:mm':
        outdate = ''
    else:
        outdate = "%sZ" % (date)
    return outdate

def get_site_list(urlnetwork):
    xmllist = etree.parse(urlnetwork)
    root  = xmllist.getroot()
    sitelist = [] 
    for child in root:
        sitecode = child.get("code")
        sitename = child.get("name")
        a = sitecode,sitename
        sitelist.append(a)
    return sitelist

def get_site_network(urlnetwork):
        xmllist = etree.parse(urlnetwork)
        root  = xmllist.getroot()
        sitelist = []
        for child in root:
            sitecode = child.get("code")
            sitenetw = child.get("network")
            if (sitenetw in 'LI'):
                netw = 'LINZ'
            else:
                netw = 'GeoNet'
            a = sitecode,netw
            sitelist.append(a)
        return sitelist

def get_delta_siteinfo(siteinfo):
    domenumber = siteinfo.xpath('/SITE/mark/domes-number')
    lat = siteinfo.xpath('/SITE/location/latitude')
    lon = siteinfo.xpath('/SITE/location/longitude')
    hgt = siteinfo.xpath('/SITE/location/height')
    datum = siteinfo.xpath('/SITE/location/datum')
    operator = siteinfo.xpath('/SITE/cgps-session/rinex/header-comment-name')
    lat = float(lat[0].text)
    lon = float(lon[0].text)
    hgt = float(hgt[0].text)
    domenumber = domenumber[0].text
    if not domenumber:
        domenumber = 'none'
    return domenumber, lat, lon, hgt, operator[0].text, datum[0].text


sitein = raw_input("   Enter site code (or comma separated list of sites): ")
siteslist = sitein.rsplit(',')

type = raw_input("   Enter streaming type [r = RTCM / m = RTCM_MSM : ")
stringrtcm = ("_RTCM", "RTCM 3;1004(1),1012(1),1006(10),1008(10),1013(10),1033(10);2;GPS+GLONASS")
stringrtcmmsm = ("_RTCM-MSM", "RTCM 3.2;1004(1),1012(1),1006(10),1008(10),1013(10),1033(10),1074(1),1084(1),1094(1),1114(1),1124(1);2;GPS+GLO+GAL+BDS+QZS")


deltasitelist = get_site_list(urlnetwork)
deltasitelist = dict(deltasitelist)

networklist = get_site_network(urlnetwork)
networklist = dict(networklist)

for site in siteslist:
    site = str.upper(site)
    siteplace = deltasitelist.get(site)
    sitenetw = networklist.get(site)
    urldelta = ("%s%s" % (delta,site))
    siteinfo = etree.parse(urldelta)
    (domenumber,lat,lon,hgt,operator,datum) = get_delta_siteinfo(siteinfo)
    if type in 'r':
        string = ("STR;%s%s;%s;%s;%s;NZL;%6.2f;%6.2f;0;0;Trimble NetR9;none;B;N;9600;" % (site,stringrtcm[0],siteplace,stringrtcm[1],sitenetw,lat,lon))
        print string
    elif type in 'm':
        string = ("STR;%s%s;%s;%s;%s;NZL;%6.2f;%6.2f;0;0;Trimble NetR9;none;B;N;9600;" % (site,stringrtcmmsm[0],siteplace,stringrtcmmsm[1],sitenetw,lat,lon))
        print string
    else:
        print "type must be r or m"
    
