package main
// Generated content
var templateConfigData = []byte(`
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
  - search: (?m)^(.+)$
    occurrences: 1
    # Preserve the inline new-lines, strip the last one
    replace: |-
      // Copyright (C) Foo
      // All rights reserved
      $1
    filter: # Apply filter on the matched string
      include:
      exclude:
      - ^// Copyright.* # Skip adding header if file already has a header
`)
