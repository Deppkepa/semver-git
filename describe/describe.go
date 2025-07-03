package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"github.com/fatih/color"
)

var (
	helpFlag    bool
	debugFlag   bool
	noColor     bool
	versionFlag bool
	releaseFlag bool
	strategy    string
	RELEASE     = 0
	VERSIONING  = "tag"
	logger      *log.Logger
	debugLogger *log.Logger
)

const (
	VERSION = "0.4"
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

func Help() {
	fmt.Println("describe project version and release from git describe")
	fmt.Println("usage: describe [-h|--help] [-v|--version] [-r|--release] project|version|release|full")
	fmt.Println("options:")
	fmt.Println("    -h|--help      print this help and exit") // +
	fmt.Println("    -V|--version   print version number")                                                // +
	fmt.Println("    -d|--debug     debug output")                                                        // +
	fmt.Println("    --no-color     no color output")                                                     // +
	fmt.Println("    -r|--release   use commit number as release number, default is no and release is 1") // +
	fmt.Println("    -s|--strategy  versioning strategy type: tag, abbrev, rank; default is tag")         // +
	fmt.Println("commands:")
	fmt.Println("    project        print project name")                      // +
	fmt.Println("    module         print module name")                       // +
	fmt.Println("    version        print project version")                   // +
	fmt.Println("    release        print project release")                   // +
	fmt.Println("    full           print full project name-version-release") // +
}

func Project() string {
	projectName := os.Getenv("PROJECT_NAME")
	if projectName != "" {
		logInfo("Using PROJECT_NAME from environment: %s", projectName)
		return projectName
	}
	logInfo("PROJECT_NAME not set, extracting from git remote")
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		logInfo("Failed to get git remote: %v", err)
		return ""
	}
	logDebug("Git remote output: %s", string(output))
	lines := strings.Split(string(output), "n")
	for _, line := range lines {
		if strings.Contains(line, "fetch") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			url := parts[0]
			logDebug("Processing git remote URL: %s", url)
			var repoPath string
			if strings.Contains(url, "://") {
				// HTTPS
				pathParts := strings.Split(url, "/")
				if len(pathParts) < 2 {
					continue
				}
				repoPath = strings.Join(pathParts[len(pathParts)-2:], "/")
			} else {
				// SSH
				sshParts := strings.SplitN(url, ":", 2)
				if len(sshParts) != 2 {
					continue
				}
				repoPath = sshParts[1]
			}
			repoPath = strings.TrimSuffix(repoPath, ".git")
			result := strings.ReplaceAll(repoPath, "/", "-")
			logInfo("Extracted project name: %s", result)
			return result
		}
	}
	logInfo("Could not extract project name from git remote")
	return ""
}

func Module() string {
	projectName := os.Getenv("PROJECT_NAME")
	if projectName != "" {
		logInfo("Using PROJECT_NAME from environment for module: %s", projectName)
		return projectName
	}
	logInfo("PROJECT_NAME not set, extracting module from git remote")
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		logInfo("Failed to get git remote for module: %v", err)
		return ""
	}
	logDebug("Git remote output for module: %s", string(output))
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "fetch") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			url := parts[1]
			afterColon := strings.SplitN(url, ":", 2)
			if len(afterColon) < 2 {
				continue
			}
			path := strings.TrimSuffix(afterColon[1], ".git")
			result := filepath.Base(path)
			logInfo("Extracted module name: %s", result)
			return result
		}
	}
	logInfo("Could not extract module name from git remote")
	return ""
}

func Release() string {
	if RELEASE == 1 {
		logInfo("Using commit number as release number")
		cmd := exec.Command("git", "describe", "--match", "v[0-9]*", "--abbrev=2", "--tags", "HEAD")
		output, err := cmd.Output()
		if err != nil {
			logInfo("Failed to get git describe for release: %v", err)
			return ""
		}
		logDebug("Git describe output for release: %s", string(output))
		re := regexp.MustCompile(`-g[a-f0-9]+$`)
		result := re.ReplaceAllString(string(output), "")
		parts := strings.Split(result, "-")
		if len(parts) > 1 {
			releaseNum := parts[len(parts)-2]
			logInfo("Extracted release number: %s", releaseNum)
			return releaseNum
		}
		logInfo("No release number found, returning default: 0")
		return "0"
	}
	logInfo("Release number not using commit, defaulting to 1")
	return "1"
}

func Version() string {
	logInfo("Determining version with strategy: %s", VERSIONING)
	switch VERSIONING {
	case "tag":
		return versionTag()
	case "abbrev":
		return versionAbbrev()
	case "rank":
		return versionRank()
	default:
		logInfo("Unknown versioning strategy: %s", VERSIONING)
		os.Exit(1)
		return ""
	}
}

func versionTag() string {
	logInfo("Using 'tag' versioning strategy")
	cmd := exec.Command("git", "describe", "--match", "v[0-9]*", "--abbrev=0", "--tags", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		logInfo("Failed to get git describe for tag version: %v", err)
		return ""
	}
	versionStr := strings.TrimSpace(string(output))
	logDebug("Raw version output (tag): %s", versionStr)
	versionStr = strings.TrimPrefix(versionStr, "v")
	versionStr = strings.Replace(versionStr, "-", "~", 1)
	logInfo("Extracted version (tag): %s", versionStr)
	return versionStr
}

func versionAbbrev() string {
	logInfo("Using 'abbrev' versioning strategy")
	cmd := exec.Command("git", "describe", "--match", "v[0-9]*", "--abbrev=2", "--always", "--tags")
	output, err := cmd.Output()
	if err != nil {
		logInfo("Failed to get git describe for abbrev version: %v", err)
		return ""
	}
	logDebug("Raw version output (abbrev): %s", string(output))
	lines := strings.Split(strings.TrimSpace(string(output)), "n")
	versions := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			line = strings.Replace(line, "-g", "", 1)
			line = strings.Replace(line, "-", "~", 1)
			versions = append(versions, line)
		}
	}
	sort.Strings(versions)
	if len(versions) > 0 {
		result := strings.TrimPrefix(versions[len(versions)-1], "v")
		logInfo("Extracted version (abbrev): %s", result)
		return result
	}
	logInfo("No version found for abbrev strategy")
	return ""
}

func versionRank() string {
	logInfo("Using 'rank' versioning strategy")
	cmd := exec.Command("git", "describe", "--match", "v[0-9]*", "--abbrev=0", "--always", "--tags")
	output, err := cmd.Output()
	if err != nil {
		logInfo("Failed to get git describe for rank version: %v", err)
		return ""
	}
	logDebug("Raw version output (rank): %s", string(output))
	lines := strings.Split(strings.TrimSpace(string(output)), "n")
	versions := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			versions = append(versions, line)
		}
	}
	sort.Strings(versions)
	if len(versions) > 0 {
		result := strings.TrimPrefix(versions[len(versions)-1], "v")
		logInfo("Extracted version (rank): %s", result)
		return result
	}
	logInfo("No version found for rank strategy")
	return ""
}

func main() {
	setupLoggers()
	flag.BoolVar(&helpFlag, "h", false, "Print help and exit")
	flag.BoolVar(&helpFlag, "help", false, "Print help and exit")
	flag.BoolVar(&debugFlag, "d", false, "Debug output")
	flag.BoolVar(&debugFlag, "debug", false, "Debug output")
	flag.BoolVar(&noColor, "no-color", false, "No color output")
	flag.BoolVar(&versionFlag, "V", false, "Print version number")
	flag.BoolVar(&versionFlag, "version", false, "Print version number")
	flag.BoolVar(&releaseFlag, "r", false, "Use commit number as release number")
	flag.BoolVar(&releaseFlag, "release", false, "Use commit number as release number")
	flag.StringVar(&strategy, "s", "tag", "Versioning strategy: tag, abbrev, rank")
	flag.StringVar(&strategy, "strategy", "tag", "Versioning strategy: tag, abbrev, rank")
	flag.Parse()

	args := flag.Args()
	logInfo("Starting describe tool with args: %v", args)
	logDebug("Flags - help: %v, debug: %v, noColor: %v, version: %v, release: %v, strategy: %s",
		helpFlag, debugFlag, noColor, versionFlag, releaseFlag, strategy)

	if helpFlag {
		logInfo("Help flag set, printing help and exiting")
		Help()
		os.Exit(0)
	}
	if versionFlag {
		logInfo("Version flag set, printing version and exiting")
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if releaseFlag {
		RELEASE = 1
		logInfo("Release flag set, using commit number as release")
	}
	if noColor {
		setupLoggers()
		logInfo("No-color flag set, disabling colors")

	}
	if len(args) == 0 {
		logInfo("No command provided, exiting with error")
		fmt.Println("Error: No command provided. Usage: describe module|project|version|release|full")
		os.Exit(1)
	}
	command := args[0]
	logInfo("Executing command: %s", command)

	switch strategy {
	case "tag":
		VERSIONING = "tag"
	case "abbrev":
		VERSIONING = "abbrev"
	case "rank":
		VERSIONING = "rank"
	default:
		logInfo("Unknown versioning strategy: '%s', expected tag, abbrev or rank", strategy)
		os.Exit(1)
	}
	logInfo("Versioning strategy set to: %s", VERSIONING)

	switch command {
	case "project":
		logInfo("Executing 'project' command")
		fmt.Println(Project())
	case "module":
		logInfo("Executing 'module' command")
		fmt.Println(Module())
	case "version":
		logInfo("Executing 'version' command")
		fmt.Println(Version())
	case "release":
		logInfo("Executing 'release' command")
		fmt.Println(Release())
	case "full":
		logInfo("Executing 'full' command")
		fmt.Printf("%s-%s-%s", Project(), Version(), Release())
	default:
		logInfo("Unknown command: '%s'. Usage: describe module|project|version|release|fulln", command)
		os.Exit(1)
	}
	logInfo("Command execution completed")
}
