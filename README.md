[![Build Status](https://travis-ci.org/prijip/gofind.svg?branch=master)](https://travis-ci.org/prijip/gofind)
# Introduction
A tool to search & replace using regular expression in a set of files.
## Features
- Configurable output directory
- Select/Filter files by name
- Select/Filter files by content
- Conditional replacement - In addition to the search regular expression, additional conditions/filters can be checked on the selected text before replacing it
- Multiple search replace on a file in one go
- All filters (file name / content / conditional replacement) support inclusion and exclusion conditions to be specified

# Installation

# Usage

```
gofind -config <path/to/configfile>
configfile can be in JSON or YAML format
  -config string
        Configuration File Name (JSON/YAML)
  -files string
        Filename pattern
  -generate-config string
        Generate sample configuration file
  -in-dir string
        Input Directory
  -occurrences string
        Number of occurrences to be replaced. Default is all occurrences
  -out-dir string
        Output Directory
  -replace value
        String to replace with
  -search string
        Regular expression to search for
  -version
        Show version and exit
```
# Sample Configuration
A sample YAML configuration file:

```yaml
# Name of the directory to search for files
inputDirectory: ./testdata/input

# Name of the directory to place the updated files
# If not provided, the original file will be replaced
outputDirectory: ./testdata/output

# Regular expressions to select the files based on their name
# Default is to select all files
#
# If an 'include' pattern is provided, a file is selected only if one of the patterns match
#
# If an 'exclude' pattern is provided, a file is ignored if one of the patterns match
#
# If a file matches both 'exclude' and 'include' patterns,
# 'exclude' takes priority and the file will be ignored
fileNamePatterns:
  include:
  - .*\.in$
  exclude:
  - .*\.ex$

# Regular expressions to select the files based on their content
# Default is to select all files
#
# If an 'include' pattern is provided, a file is selected only if one of the patterns match
#
# If an 'exclude' pattern is provided, a file is ignored if one of the patterns match
#
# If a file matches both 'exclude' and 'include' patterns,
# 'exclude' takes priority and the file will be ignored
filter:
  include:
  exclude:
  - ^// Skip This.*

# Search Replace patterns
patterns:
  # Search and replace all 'one's with 'ONE'
  - search: one
    replace: ONE

  # Insert a file header in files that does not already have one
  - search: (?m)^(.+)$
    occurrences: 1
    # "|-" is used below to preserve the inline new-lines, strip the last one
    replace: |-
      // Copyright (C) Foo
      // All rights reserved
      $1
    # Apply filter on the matched string
    filter:
      include:
      exclude:
      - ^// Copyright.* # Skip adding header if file already has a header
```
