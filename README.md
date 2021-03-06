[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)![Shell Script](https://img.shields.io/badge/shell_script-%23121011.svg?style=for-the-badge&logo=gnu-bash&logoColor=white)


# ImportHPUXPrinter
This is a simple application to import HP-UX Printer in CUPS Linux Spooler.

## Prerequisites

Before you get started, you’ll need to have these things:

* [Go language installed](https://go.dev/) on Linux server
* Git installed on Linux server

The repository contains the following files:
* **HPUX/CreateListPrinter.sh** : script to capture physical and logical printers under HP-UX 
* **src/SetPrinter.go** : source of the go script
* **SetPrinter** : binary of the go script
* **config.json** : configuration file of the script
* **template.sh**: backend template file

## Initial setup

Clone the repository :

```
[archi@mercure] git clone https://github.com/colussim/ImportHPUXPrinter.git
```

Copy the HPUX directory to the HP-UX server

## Setup

You will find 1 configuration files :

**config.json** : configuration file with the following parameters:
```json
{

    "FilePhysical" : "PhyPrinterList.csv",  Name of the file containing the names and ip addresses of the physical printers
    "FileLogical": "LogPrinterList.csv",    Name of the file containing the names of the logical printers
    "Template":"template.sh",               Name of the template file for the definition of the backends
    "Fprinter":"/usr/lib/cups/backend/",    Directory where the backends are stored
    "Description" :" MAGELLAN-"             Description added when creating the printer 

}
```

## Usage

On HP-UX server run the script :
```bash
[root@bandol HPUX]./CreateListPrinter.sh
```
This script generates two files in csv format:
* PhyPrinterList.csv
* LogPrinterList.csv
    
Copy these two files into the directory of the Linux server where the repository was cloned.

On Linux server run the script :
```bash
[archi@mercure ImportHPUXPrinter]sudo ./SetPrinter
```
This script reads: 
- the **PhyPrinterList.csv** and **LogPrinterList.csv** files 
- creates the printers backend files 
- adds the logical printers if they exist as well as the physical printers
