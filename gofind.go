package gofind

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// Filter stores the inclusion and exclusion patterns
//
type Filter struct {
	Include []*regexp.Regexp
	Exclude []*regexp.Regexp
}

// TestFilters applies the inclusion and exclusion tests on the given data
// Returns true if:
//		there are no exclusion patterns, or none of the exclusion tests pass
//  	And
//  	there are no inclusion patterns, or one of the inclusion patterns pass
// Returns false if:
//		At least one of the exclusion patterns pass or
//		All of the inclusion patterns fails
func (f *Filter) TestFilters(data []byte) (canSelect, include, exclude bool) {
	include = true
	exclude = false

	// If an inclusion pattern is not specified inclusion test pass
	// If an inclusion pattern is specified and any one of the patterns match, pass
	if len(f.Include) > 0 {
		include = false

		for _, regEx := range f.Include {
			if regEx.Match(data) {
				include = true
				break
			}
		}
	}

	// If one of the exclusion patterns match, fail
	if len(f.Exclude) > 0 {
		for _, regEx := range f.Exclude {
			if regEx.Match(data) {
				exclude = true
				break
			}
		}
	}

	return !exclude && include, include, exclude
}

// SearchReplacePattern stores the pattern to be searched for and replaced
type SearchReplacePattern struct {
	SearchRegex    *regexp.Regexp
	ReplacePattern []byte
	Occurrences    int
	Filter         *Filter
}

// SearchReplace searches the inData for the given patterns
func SearchReplace(inData []byte, patterns []SearchReplacePattern) ([]byte, error) {
	replaced := inData
	for i := range patterns {
		if patterns[i].ReplacePattern == nil {
			continue
		}

		// If all occurrences need to be replaced, with no filters to be applied
		// for each replacement, replace everything in one go
		if patterns[i].Occurrences < 0 && (patterns[i].Filter == nil || len(patterns[i].Filter.Include)+len(patterns[i].Filter.Exclude) == 0) {
			replaced = patterns[i].SearchRegex.ReplaceAll(replaced, patterns[i].ReplacePattern)
			continue
		}

		// TODO: Optimize
		count := patterns[i].Occurrences
		var segs [][]byte // Segments of data
		searchBuf := replaced
		for {
			// Loop exit condition on count
			// If occurrences is provided, use that in the loop exit condition
			// Otherwise loop until all matches are processed
			if count == 0 {
				break
			}
			if count > 0 {
				count--
			}

			loc := patterns[i].SearchRegex.FindIndex(searchBuf)
			if loc == nil { // No match
				break
			}

			if loc[0] != 0 {
				seg := searchBuf[0:loc[0]]
				segs = append(segs, seg)
			}

			s := searchBuf[loc[0]:loc[1]]
			if len(s) == 0 { // Some regex trouble, break out anyway
				log.Print("Warning: RegExp: '", patterns[i].SearchRegex.String(), "' causing ZERO length match; please verify RegExp")
				break
			}

			shouldReplace := true
			if patterns[i].Filter != nil {
				shouldReplace, _, _ = patterns[i].Filter.TestFilters(s)
			}
			r := s
			if shouldReplace {
				r = patterns[i].SearchRegex.ReplaceAll(s, patterns[i].ReplacePattern)
			}
			segs = append(segs, r)
			searchBuf = searchBuf[loc[1]:]
		}
		if len(searchBuf) > 0 {
			segs = append(segs, searchBuf)
		}

		replaced = bytes.Join(segs, []byte{})
	}

	return replaced, nil
}

// FileSearchReplace searches the input file for the patters and updates
// the output file with the updated content
// Returns true if there is any content update (output file is written)
// Returns false if there is no content update
func FileSearchReplace(inFilePath, outFilePath string, patterns []SearchReplacePattern, filter *Filter) (bool, error) {
	fileContent, err := ioutil.ReadFile(inFilePath)
	if err != nil {
		log.Printf("Error processing file %s. err=%v", inFilePath, err)
		return false, err
	}

	// If file level filters are specified, test them and skip file accordingly
	if filter != nil {
		if bPass, _, _ := filter.TestFilters(fileContent); !bPass {
			return false, nil
		}
	}

	replaced, err := SearchReplace(fileContent, patterns)
	if err != nil {
		log.Printf("SearchReplace failed for file %s. err=%v", inFilePath, err)
		return false, err
	}

	if bytes.Equal(replaced, fileContent) {
		log.Printf("%s [No Change]", inFilePath)
		return false, nil
	}

	if inFilePath != outFilePath {
		outputDir := filepath.Dir(outFilePath)
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating path %s. err=%v", outputDir, err)
			return false, err
		}
	}
	err = ioutil.WriteFile(outFilePath, replaced, 0777)
	if err != nil {
		log.Printf("%s - Failed to write to %s, err=%v", inFilePath, outFilePath, err)
	} else {
		log.Printf("%s - Updated - %s", inFilePath, outFilePath)
	}

	return true, nil
}
