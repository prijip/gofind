# Introduction
A tool to search & replace using multiple patterns in a set of files
## Features
- Configurable input directory
- Configurable output directory
- Specify file-name inclusion pattern as regular expression
- Specify file-name exclusion patterns as regular expression
- Multiple patterns to be searched for and replaced

# Installation

# Usage

```sh
gofind -search <search-string> [-replace <replace-string> -files <file-name-pattern>] -in-dir <path> -out-dir <path>
gofind -config <path/to/config.json>

Where:
 -config string
        Configuration File Name (JSON)
  -files string
        Filename pattern
  -in-dir string
        Input Directory
  -out-dir string
        Output Directory
  -replace value
        String to replace with
  -search string
        Regular expression to search for
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
# If an 'include' pattern is provided, a file is selected only if
# any one of thepatterns match
#
# If an 'exclude' pattern is provided, a file is ignored if
# any one of the patterns patch
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
# If an 'include' pattern is provided, a file is selected only if
# any one of thepatterns match
#
# If an 'exclude' pattern is provided, a file is ignored if
# any one of the patterns patch
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
  - search: ^
    occurrences: 1
    replace: |+
      // Copyright (C) Foo
      // All rights reserved

    filter:
      include:
      exclude:
      - ^// Copyright.* # Skip adding header if file already has a header
```