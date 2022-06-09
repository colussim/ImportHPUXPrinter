/*******************************************************************/
/*  															   */
/*  @project     : Milan							               */
/*  @package     : main   										   */
/*  @subpackage  :												   */
/*  @access      :												   */
/*  @paramtype   : 												   */
/*  @argument    :												   */
/*  @description : Read files :									   */
/*					 PhyPrinterList.csv	- LogPrinterList.csv	   */
/*                 Create printer in CUPS Spooler                  */
/*				                                                   */
/*																   */
/*  @author Emmanuel COLUSSI									   */
/*  @version 1.00												   */
/******************************************************************/

package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Declare a struct for Config fields

type Config struct {
	FilePhysical string
	FileLogical  string
	Template     string
	Fprinter     string
	Description  string
}

// Declare a Printer Struct
type Printer struct {
	PhyName   string
	IPprinter string
	LogName   []string
}

// Declare a Physical Printer Struct
type HWPrinter struct {
	PhyName   string
	IPprinter string
}

// Declare a Logical Printer Struct

type LogPrinter struct {
	LogName string
	PhyName string
}

var CMD = ""

func GetConfig(config Config) Config {

	fconfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Problem with the configuration file : config.json")
		os.Exit(1)
	}
	json.Unmarshal(fconfig, &config)
	return config
}

// Func read CSV File
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

// Func exec sys cmd
func sys_cmd(cmd1 string, Printer string) {

	cmd := exec.Command("/bin/bash", "-c", cmd1)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	_, err := cmd.Output()
	if err != nil {
		log.Printf("⇨  error occured : %s", err)
		os.Exit(1)
	} else {
		log.Printf("⇨ Printer Add in CUPS Spooler : %s\n", Printer)
	}
}
// Func  Delete Space in string Soket
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), "")
}
/* ------------------------------ Main ------------------------------*/

func main() {

	// Read Config File

	var configapp Config
	var AppConfig = GetConfig(configapp)

	// Read PHYSICAL Printer File

	lines, err := ReadCsv(AppConfig.FilePhysical)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	sizehw := len(lines)

	if sizehw > 0 {

		/*----- Create Physical Printer struct --------*/
		HWprinter := make([]HWPrinter, sizehw)
		var i = 0

		for _, line := range lines {
			HWprinter[i] = HWPrinter{
				PhyName:   line[0],
				IPprinter: line[1],
			}
			i++
		}
		log.Println("⇨ Physical Printer struct created")

		/*--------------------------------------------*/

		// Read LOGICAL Printer File

		lines1, err := ReadCsv(AppConfig.FileLogical)
		if err != nil {
			panic(err)
			os.Exit(1)
		}

		/*----- Create Logical Printer struct --------*/
		sizelog := len(lines1)
		LOGprinter := make([]LogPrinter, sizelog)
		i = 0

		// Loop through lines & turn into object
		for _, lines1 := range lines1 {
			LOGprinter[i] = LogPrinter{
				LogName: lines1[0],
				PhyName: lines1[1],
			}
			i++
		}
		log.Println("⇨ Logical Printer struct created")

		/*--------------------------------------------*/

		/*----- Create Printer struct --------*/

		ImpPrinter := make([]Printer, sizehw)

		for i := 0; i < sizehw; i++ {
			ImpPrinter[i].PhyName = HWprinter[i].PhyName
			ImpPrinter[i].IPprinter = HWprinter[i].IPprinter
			for j := 0; j < sizelog; j++ {
				if LOGprinter[j].PhyName == HWprinter[i].PhyName {
					ImpPrinter[i].LogName = append(ImpPrinter[i].LogName, LOGprinter[j].LogName)

				}

			}
		}
		log.Println("⇨ Printer struct created")

		/*--------------------------------------------*/

		/*------ Create Printer CUPS Spool -----------*/

		// Get Hostname
		Hname, err := os.Hostname()
		if err != nil {
			panic(err)
			os.Exit(1)
		}
		Description := AppConfig.Description + Hname

		// Read Template
		data, err := ioutil.ReadFile(AppConfig.Template)
		if err != nil {
			log.Panicf("⇨ failed reading data from file: %s", err)
			os.Exit(1)
		}
		data1 := string(data)
		searchRegexp := regexp.MustCompile(`PRINTER=((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3})\d{1}`)

		for i := 0; i < sizehw; i++ {

			Newvalue := "PRINTER=" + ImpPrinter[i].IPprinter
			Socket1 := ImpPrinter[i].IPprinter + ":9100"
			Socket := standardizeSpaces(Socket1)
			result := searchRegexp.ReplaceAllString(data1, Newvalue)

			Fprinter1 := AppConfig.Fprinter + ImpPrinter[i].PhyName
			if err = ioutil.WriteFile(Fprinter1, []byte(result), 0755); err != nil {
				log.Println(err)
				os.Exit(1)
			}
			log.Println("⇨ Printer backend created : ", ImpPrinter[i].PhyName)

			NbrLogic := len(ImpPrinter[i].LogName)
			if NbrLogic > 0 {
				for nbr := 0; nbr < NbrLogic; nbr++ {
					CMD = CMD + "sudo /usr/sbin/lpadmin -p " + ImpPrinter[i].LogName[nbr] + " -E -v socket://" + Socket + " -D MAGELLAN-" + Hname + ";"
					sys_cmd(CMD, ImpPrinter[i].LogName[nbr])

				}
			} else {
				CMD = CMD + "sudo /usr/sbin/lpadmin -p " + ImpPrinter[i].PhyName + " -E -v socket://" + Socket + " -D " + Description + ";"
				sys_cmd(CMD, ImpPrinter[i].PhyName)
			}

		}
	} else {
		log.Println("⇨ No declared printer")
	}

}
