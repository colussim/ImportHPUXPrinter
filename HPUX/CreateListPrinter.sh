#!/usr/bin/ksh

#################################################################
#								#
#	@project : MILAN					#
#	@package : CreateListPrinter.sh				#
#	@argument :						#
#	@description : create a list of printers		# 
#		       PhyPrinterList.csv: Physical printer	#
#			Physicalname,IP_ADDRESS			#
#	               LogPrinterList.csv: Logical printer	#
#			Loginame,Physicalname			# 
#								#
#################################################################

HOMEREP="/etc/lp/interface/"
CMLISTINT=`ll /var/spool/lp/interface/*[[:upper:]]*|awk '{print $9}'|awk '{FS="/"}{print $6}'`
CMLISTLOG=`grep "PRINTER=" $HOMEREP*`

PHYPRINTER="PhyPrinterList.csv"
LOGPRINTER="LogPrinterList.csv"

cat /dev/null > $PHYPRINTER
cat /dev/null > $LOGPRINTER


echo "⇨ Generation of the printers files ...." 

# Generate Physical Printer File

for PERIPH in $CMLISTINT
do

INTERFACEIP=`grep -E '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)' $HOMEREP$PERIPH|grep PERIPH|awk '{FS="="}{ print $2 }'`

GETLOGICIP=`grep "PRINTER=$PERIPH" $HOMEREP* |awk '{FS=":"}{print $1}'|awk '{FS="/"}{print $5}'`

if [ -z "$INTERFACEIP" ]
then echo ""
else echo "$PERIPH,$INTERFACEIP" >> $PHYPRINTER
fi
done

# Generate Logical Printer File

for LOG in $CMLISTLOG
do
IFS="/"
set -A array $LOG
IFS=":"
set -A array2 ${array[4]}
IFS="="
set -A array3 ${array[4]}

echo "${array2[0]},${array3[1]}" >> $LOGPRINTER

done
