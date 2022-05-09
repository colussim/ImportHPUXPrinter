#! /bin/bash
PRINTER=10.96.249.59 
SRV=`uname -n`
# Have debug info in /var/log/cups/error_log:
set -x
# Output "device discovery" information on stdout: 
if test "$#" = "0"
then echo 'network socket://""$PRINTER"":9100 ""MAGELLAN-'$SRV"" ""$PRINTER':9100"'
exit 0 
fi
# Set INPUTFILE to where the input comes from: 
INPUTFILE="-"
if test -n "$6"
then INPUTFILE="$6"
fi
# Send the data to the remote port:

SRV=`uname -n`
echo "\033%-12345X@PJL">ENT$$
echo "@PJL COMMENT CANPJL SET DOCNAME=\"MAGELLAN-$SRV\"">>ENT$$ 
echo "@PJL COMMENT CANPJL SET USERNAME=\"$2\"">>ENT$$
echo "@PJL ENTER LANGUAGE=PCL">>ENT$$
type=pcl
copies=$4

echo "INFO: sending data to $PRINTER:9100" 1>&2

perl -pi.bak -e 's/\#\*\#/\*/g' $file 
grep -l ZEBRA $file>/dev/null retour=$?
if [ $retour = 0 ]
then
/appli/printer/program/SPOOLETIQ "$file"
cat ENT$$ $fic | /appli/printer/program/sicFilter -m1 -t/appli/printer/param/TRAYMAP7080 -p/appli/printer/param/PAPERMAP7080| netcat -w 1 $PRINTER 9100 1>&2 
else
cat ENT$$ $file | /appli/printer/program/sicFilter -m1 -t/appli/printer/param/TRAYMAP7080 - p/appli/printer/param/PAPERMAP7080 | netcat -w 1 $PRINTER 9100 1>&2
fi