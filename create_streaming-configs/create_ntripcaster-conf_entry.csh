#! /bin/csh -f

if ($#argv == 0) then
  cat << END
    create primary ntripcaster.conf entry reading ip from network repo

    usage: create_ntripcasterconf-entry.csh   site

END
exit
endif


set site = `echo $1 | awk '{print toupper($1)}'`

set ipfile = /home/elidana/git/geonet/network/data/devices.csv

set ip = `grep $site $ipfile | grep "^gps" | awk 'BEGIN {FS=","} {print $3}' | awk 'BEGIN {FS="/"} {print $1}'`

echo "relay pull -m /"$site"_RTCM" $ip":8855 -2"
echo "relay pull -m /"$site"_RTCM-MSM" $ip":8857 -2"
#echo "relay pull -i rtgps:ejlhO8rU -m /"$site"_RTCM ntrip-px-wa.geonet.org.nz:8888/"$site"0 -2"
