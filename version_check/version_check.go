package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"github.com/fatih/color"
)

var (
	helpFlag      bool
	debugFlag     bool
	rulesFlag     bool
	typeFlag      bool
	buildTypeFlag bool
	noColor       bool
	versionFlag   bool
	VERSION       = "0.2"
	logger        *log.Logger
	debugLogger   *log.Logger
)

var (
	versionRelease      = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+$`)
	versionPrerelease   = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+(\-|\~)(alpha|beta|rc|pre)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
	versionPostrelease  = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+\.(fix|next|post)(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
	versionIntermediate = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+\_[a-zA-Z]+(\.[0-9]+|\_[a-zA-Z]+(\.[0-9]+)*)*$`)
)

func setupLoggers() {
	if noColor {
		logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime)
		debugLogger = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime)
	} else {
		infoPrefix := color.New(color.FgGreen).Sprint("INFO: ")
		debugPrefix := color.New(color.FgCyan).Sprint("DEBUG: ")

		logger = log.New(os.Stderr, infoPrefix, log.Ldate|log.Ltime)
		debugLogger = log.New(os.Stderr, debugPrefix, log.Ldate|log.Ltime)
	}
}

func logInfo(format string, v ...interface{}) {
	if debugFlag {
		logger.Printf(format, v...)
	}
}

func logDebug(format string, v ...interface{}) {
	if debugFlag {
		debugLogger.Printf(format, v...)
	}
}

func printHelp() {
	fmt.Println("version_check check that version is set correct according to project rules")
	fmt.Println("usage: version_check [-h|--help] [-v|--version] version")
	fmt.Println("options:")
	fmt.Println("    -h|--help        print this help and exit") // +
	fmt.Println("    -V|--version     print version number") // +
	fmt.Println("    -r|--rules       print regexp rules for check versions and exit") // +
	fmt.Println("    -d|--debug       debug output") // +
	fmt.Println("    -t|--type        output version type: release, prerelease, postrelease, intermediate") // +
	fmt.Println("    -b|--build-type  output build type for cmake: Release for release version, Debug for other") // +
	fmt.Println("    --no-color       no colored output") // +
	fmt.Println("arguments:")
	fmt.Println("    version          version string to check") // +
}

func printRules() {
	fmt.Println("version rules in precedence order:")
	fmt.Println("    release:     ^v?[0-9]+\\.[0-9]+\\.[0-9]+$")
	fmt.Println("    prerelease:  ^v?[0-9]+\\.[0-9]+\\.[0-9]+(\\-|\\~)(alpha|beta|rc|pre)(\\.[0-9]+|\\_[a-zA-Z]+(\\.[0-9]+)*)*$")
	fmt.Println("    postrelease: ^v?[0-9]+\\.[0-9]+\\.[0-9]+\\.(fix|next|post)(\\.[0-9]+|\\_[a-zA-Z]+(\\.[0-9]+)*)*$")
	fmt.Println("    intermediate:^v?[0-9]+\\.[0-9]+\\.[0-9]+\\_[a-zA-Z]+(\\.[0-9]+|\\_[a-zA-Z]+(\\.[0-9]+)*)*$")
}

func checkVersion(versionStr string) {
	logDebug("Checking version '%s'", versionStr)
	if versionRelease.MatchString(versionStr) {
		logDebug("Version '%s' is release", versionStr)
		if typeFlag {
			fmt.Println("Release")
		}
		if buildTypeFlag {
			fmt.Println("Release")
		}
		os.Exit(0)
	}
	if versionPrerelease.MatchString(versionStr) {
		logDebug("Version '%s' is pre release", versionStr)
		if typeFlag {
			fmt.Println("Pre release")
		}
		if buildTypeFlag {
			fmt.Println("Debug")
		}
		os.Exit(0)
	}
	if versionPostrelease.MatchString(versionStr) {
		logDebug("Version '%s' is post release", versionStr)
		if typeFlag {
			fmt.Println("Post release")
		}
		if buildTypeFlag {
			fmt.Println("Debug")
		}
		os.Exit(0)
	}
	if versionIntermediate.MatchString(versionStr) {
		logDebug("Version '%s' is intermediate", versionStr)
		if typeFlag {
			fmt.Println("Intermediate release")
		}
		if buildTypeFlag {
			fmt.Println("Debug")
		}
		os.Exit(0)
	}
	logInfo("Wrong version '%s'", versionStr)
	os.Exit(1)
}

func main() {
	setupLoggers()
	logDebug("Starting program")
	flag.BoolVar(&helpFlag, "h", false, "Print help and exit")
	flag.BoolVar(&helpFlag, "help", false, "Print help and exit")
	flag.BoolVar(&debugFlag, "d", false, "Debug output")
	flag.BoolVar(&debugFlag, "debug", false, "Debug output")
	flag.BoolVar(&rulesFlag, "r", false, "Print regexp rules for checking versions and exit")                                      
	flag.BoolVar(&rulesFlag, "rules", false, "Print regexp rules for checking versions and exit")                                   
	flag.BoolVar(&typeFlag, "t", false, "Output version type: release, prerelease, postrelease, intermediate")                     
	flag.BoolVar(&typeFlag, "type", false, "Output version type: release, prerelease, postrelease, intermediate")                   
	flag.BoolVar(&buildTypeFlag, "b", false, "Output build type for cmake: Release for release version, Debug for others")       
	flag.BoolVar(&buildTypeFlag, "build-type", false, "Output build type for cmake: Release for release version, Debug for others") 
	flag.BoolVar(&noColor, "no-color", false, "No color output")                                                                    
	flag.BoolVar(&versionFlag, "V", false, "Print version number")                                                                
	flag.BoolVar(&versionFlag, "version", false, "Print version number")                                                        
	flag.Parse()
	logDebug("Flags - help: %v, debug: %v, noColor: %v, version: %v, rules: %v, type: %v, buildType: %v",
		helpFlag, debugFlag, noColor, versionFlag, rulesFlag, typeFlag, buildTypeFlag)
	if helpFlag {
		logInfo("Help flag detected, printing help")
		printHelp()
		os.Exit(0)
	}
	if versionFlag {
		logInfo("Version flag detected, printing version")
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if rulesFlag {
		logInfo("Rules flag detected, printing rules")
		printRules()
		os.Exit(0)
	}
	if noColor {
		setupLoggers()
		logInfo("No-color flag set, disabling colors")
		
	}
	args := flag.Args()
	if len(args) == 0 {
		logInfo("No version argument provided")
		//fmt.Println("Error: No version argument provided")
		os.Exit(1)
	}
	logDebug("Processing version argument: %s", args[0])
	checkVersion(args[0])
}
