package gofind

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRegex(t *testing.T, text string) *regexp.Regexp {
	re, err := regexp.Compile(text)
	if err != nil {
		t.Fatal(err)
	}

	return re
}
func TestSearchReplace_ReplaceNone(t *testing.T) {
	testData := []byte(`The Turing machine was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automatic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a machine exist that can determine whether any arbitrary machine on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a machine exist that can determine whether any arbitrary machine on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	patterns := []SearchReplacePattern{
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "Moore"),
			ReplacePattern: []byte(""),
			Occurrences:    -1,
		},
	}
	replaced, err := SearchReplace(testData, patterns)

	assert.NoError(t, err)
	assert.Equal(t, testData, replaced)
}

func TestSearchReplace_ReplaceOne(t *testing.T) {
	testData := []byte(`The Turing machine was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automatic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a machine exist that can determine whether any arbitrary machine on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a machine exist that can determine whether any arbitrary machine on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	patterns := []SearchReplacePattern{
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "Moore"),
			ReplacePattern: []byte(""),
			Occurrences:    -1,
		},
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "machine(\\s)"), // 'machine' followed by a whitespace
			ReplacePattern: []byte("m$1"),
			Occurrences:    -1,
		},
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "Mealy"),
			ReplacePattern: []byte(""),
			Occurrences:    -1,
		},
	}
	replaced, err := SearchReplace(testData, patterns)

	expectedOutput := []byte(`The Turing m was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automatic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a m exist that can determine whether any arbitrary m on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a m exist that can determine whether any arbitrary m on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, replaced)
}

func TestSearchReplace_ReplaceTwo(t *testing.T) {
	testData := []byte(`The Turing machine was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automatic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a machine exist that can determine whether any arbitrary machine on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a machine exist that can determine whether any arbitrary machine on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	patterns := []SearchReplacePattern{
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "Moore"),
			ReplacePattern: []byte(""),
			Occurrences:    -1,
		},
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "machine(\\s)"), // 'machine' followed by a whitespace
			ReplacePattern: []byte("m$1"),
			Occurrences:    -1,
		},
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "at"),
			ReplacePattern: []byte("AT"),
			Occurrences:    -1,
		},
	}
	replaced, err := SearchReplace(testData, patterns)

	expectedOutput := []byte(`The Turing m was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automATic machine).
		With this model, Turing was able to answer two questions in the negATive:
			(1) Does a m exist thAT can determine whether any arbitrary m on its tape is "circular" (e.g., freezes, or fails to continue its computATional task);
			(2) does a m exist thAT can determine whether any arbitrary m on its tape ever prints a given symbol.
		Thus by providing a mAThemATical description of a very simple device capable of
		arbitrary computATions, he was able to prove properties of computATion in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, replaced)
}

func TestSearchReplace_ReplaceOneOccurrence(t *testing.T) {
	testData := []byte(`The Turing machine was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automatic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a machine exist that can determine whether any arbitrary machine on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a machine exist that can determine whether any arbitrary machine on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	patterns := []SearchReplacePattern{
		SearchReplacePattern{
			SearchRegex:    makeRegex(t, "at"),
			ReplacePattern: []byte("AT"),
			Occurrences:    1,
		},
	}
	replaced, err := SearchReplace(testData, patterns)

	expectedOutput := []byte(`The Turing machine was invented in 1936 by Alan Turing,
		who called it an "a-machine" (automATic machine).
		With this model, Turing was able to answer two questions in the negative:
			(1) Does a machine exist that can determine whether any arbitrary machine on its tape is "circular" (e.g., freezes, or fails to continue its computational task);
			(2) does a machine exist that can determine whether any arbitrary machine on its tape ever prints a given symbol.
		Thus by providing a mathematical description of a very simple device capable of
		arbitrary computations, he was able to prove properties of computation in
		general and in particular, the uncomputability of the Entscheidungsproblem ('decision problem')`)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, replaced)
}
