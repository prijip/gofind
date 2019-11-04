package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/prijip/gofind"
)

// FilterOptions to define the search criteria
type FilterOptions struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// SearchReplaceOption holds the string to search for, and the string to be replaced with
type SearchReplaceOption struct {
	Search      string        `json:"search"`
	Replace     StringOption  `json:"replace"`
	Occurrences string        `json:"occurrences"`
	Filter      FilterOptions `json:"filter"`
}

// AppConfig stores the application configuration
type AppConfig struct {
	Patterns        []SearchReplaceOption `json:"patterns"`
	InputDirectory  string                `json:"inputDirectory"`
	OutputDirectory string                `json:"outputDirectory"`
	FileNames       FilterOptions         `json:"fileNamePatterns"`
	Filter          FilterOptions         `json:"filter"`
}

var (
	config         AppConfig
	searchPattern  string
	replacePattern StringOption
	occurrences    string
	showVersion    bool

	configFileName         string
	inputDirectory         string
	outputDirectory        string
	fileNameIncludePattern string
	generateConfigFileName string
)

// Version and Build Date set externally during linking
var (
	Version   = "undefined"
	BuildDate = "undefined"
)

func init() {
	flag.Usage = printUsage

	flag.StringVar(&configFileName, "config", "", "Configuration File Name (JSON/YAML)")
	flag.StringVar(&searchPattern, "search", "", "Regular expression to search for")
	flag.Var(&replacePattern, "replace", "String to replace with")
	flag.StringVar(&searchPattern, "occurrences", "", "Number of occurrences to be replaced. Default is all occurrences")
	flag.StringVar(&fileNameIncludePattern, "files", "", "Filename pattern")
	flag.StringVar(&inputDirectory, "in-dir", "", "Input Directory")
	flag.StringVar(&outputDirectory, "out-dir", "", "Output Directory")
	flag.StringVar(&generateConfigFileName, "generate-config", "", "Generate sample configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")
}

func printVersion() {
	fmt.Fprintln(flag.CommandLine.Output(),
		"gofind version", Version, "build date", BuildDate)
}

func printUsage() {
	printVersion()

	fmt.Fprintln(flag.CommandLine.Output(),
		"gofind -config <path/to/configfile>")
	fmt.Fprintln(flag.CommandLine.Output(), "configfile can be in JSON or YAML format")

	flag.PrintDefaults()
}

func parseFlags() error {
	flag.Parse()

	config = AppConfig{}
	// If a config file is specified, load it
	if len(configFileName) > 0 {
		configData, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Printf("Error reading config file %s. err=%v", configFileName, err)
			return err
		}

		switch configFileType := strings.ToLower(filepath.Ext(configFileName)); configFileType {
		case ".json":
			err = json.Unmarshal(configData, &config)

		case ".yaml":
			err = yaml.Unmarshal(configData, &config)

		default:
			err = fmt.Errorf("Unknown config file type '%s'", configFileType)
		}

		if err != nil {
			log.Printf("Error parsing config file %s. err=%v", configFileName, err)
			return err
		}
	}

	// Override config file setting, if any, with command line options
	if len(searchPattern) != 0 {
		pattern := SearchReplaceOption{
			Search:      searchPattern,
			Replace:     replacePattern,
			Occurrences: occurrences,
		}

		config.Patterns = append(config.Patterns, pattern)
	}

	if len(inputDirectory) > 0 {
		config.InputDirectory = inputDirectory
	}

	if len(fileNameIncludePattern) > 0 {
		config.FileNames.Include = append(config.FileNames.Include, fileNameIncludePattern)
	}

	if len(outputDirectory) > 0 {
		config.OutputDirectory = outputDirectory
	}

	if len(config.OutputDirectory) == 0 {
		config.OutputDirectory = config.InputDirectory
	}

	return nil
}

func validateFlags() error {
	if len(config.Patterns) == 0 || len(config.InputDirectory) == 0 {
		return fmt.Errorf("Incorrect Usage")
	}

	return nil
}

func filterPatternsFromOptions(options FilterOptions) (patterns gofind.Filter, err error) {
	var includePatterns []*regexp.Regexp

	for i := range options.Include {
		var regExp *regexp.Regexp
		regExp, err = regexp.Compile(options.Include[i])
		if err != nil {
			log.Print("Failed to compile regex: ", options.Include[i])
			return
		}

		includePatterns = append(includePatterns, regExp)
	}

	var excludePatterns []*regexp.Regexp

	for i := range options.Exclude {
		var regExp *regexp.Regexp
		regExp, err = regexp.Compile(options.Exclude[i])
		if err != nil {
			log.Print("Failed to compile regex: ", options.Exclude[i])
			return
		}

		excludePatterns = append(excludePatterns, regExp)
	}

	patterns.Include = includePatterns
	patterns.Exclude = excludePatterns

	return
}

func searchReplacePatternsFromOptions(options []SearchReplaceOption) []gofind.SearchReplacePattern {
	var patterns []gofind.SearchReplacePattern

	// Compile the search text patterns
	for i := range options {
		searchRegex, err := regexp.Compile(options[i].Search)
		if err != nil {
			log.Print("Failed to compile regex: ", options[i].Search)
			return nil
		}
		var replacePattern []byte
		if options[i].Replace.IsValid() {
			replacePattern = []byte(config.Patterns[i].Replace.String())
		}

		occOpt := options[i].Occurrences
		occInt := -1
		if len(occOpt) == 0 || occOpt == "all" {
			occInt = -1
		} else {
			if v, err := strconv.ParseInt(occOpt, 10, 32); err != nil {
				log.Printf("Error parsing occurrences: err=%v", err)
			} else {
				occInt = int(v)
			}
		}

		filter, err := filterPatternsFromOptions(options[i].Filter)
		if err != nil {
			log.Print("Failed to compile regex: ", options[i].Search)
			return nil
		}

		pattern := gofind.SearchReplacePattern{
			SearchRegex:    searchRegex,
			ReplacePattern: replacePattern,
			Occurrences:    occInt,
			Filter:         &filter,
		}
		patterns = append(patterns, pattern)
	}

	return patterns
}

func fileHandler(patterns []gofind.SearchReplacePattern, fnFilter, filter *gofind.Filter, updatedFiles *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		bPass, _, excludeFile := fnFilter.TestFilters([]byte(path))
		if excludeFile && info.IsDir() {
			return filepath.SkipDir
		}

		if bPass && !info.IsDir() {
			fileName, err := filepath.Rel(config.InputDirectory, path)
			if err != nil {
				log.Print("Failed to find relative path for ", path, ", err=", err)
				return err
			}
			outputFilePath := filepath.Join(config.OutputDirectory, fileName)

			updated, _ := gofind.FileSearchReplace(path, outputFilePath, patterns, filter)
			if updated {
				*updatedFiles = append(*updatedFiles, path)
			}
		}

		return nil
	}
}

func doFind() {
	var err error

	// Compile the search text patterns
	patterns := searchReplacePatternsFromOptions(config.Patterns)
	filter, err := filterPatternsFromOptions(config.Filter)
	if err != nil {
		log.Print("Error compiling global filter patterns")
		return
	}

	fnFilter, err := filterPatternsFromOptions(config.FileNames)
	if err != nil {
		log.Print("Error compiling file name filter patterns")
		return
	}

	var updatedFiles []string
	err = filepath.Walk(config.InputDirectory, fileHandler(patterns, &fnFilter, &filter, &updatedFiles))

	updatedFileCount := len(updatedFiles)
	if updatedFileCount > 0 {
		log.Print("==== Summary ====")
		log.Print(updatedFileCount, " file(s) updated:")
		for _, file := range updatedFiles {
			log.Print(file)
		}
	}

	if err != nil {
		log.Print(err)
	}
}

func main() {
	if err := parseFlags(); err != nil {
		return
	}

	if showVersion {
		printVersion()
		return
	}

	if len(generateConfigFileName) > 0 {
		ioutil.WriteFile(generateConfigFileName, templateConfigData, 0777)
		return
	}

	if err := validateFlags(); err != nil {
		log.Print(err)
		printUsage()
	}

	log.Printf("Starting")
	defer log.Printf("Ending")

	doFind()
}
