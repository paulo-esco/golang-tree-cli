package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"

	// "path/filepath"
	// "strings"
	"sort"
)

const (
	SPACE_MODE int = iota
	TAB_MODE
)

const INDENT_MODE int = TAB_MODE

/*

sort
size
tabs
adaptive slash

*/

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	lines, err := makeBranch(path, printFiles)
	if err != nil {
		return err
	}

	indentLines(lines)

	output := ""

	for _, l := range lines {
		output += l + "\n"
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()
	out.Write([]byte(output))

	return nil
}

type ByName []os.DirEntry

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name() < n[j].Name() }

func makeBranch(path string, printFiles bool) ([]string, error) {
	entries, err := os.ReadDir(path)
	sort.Sort(ByName(entries))

	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		name := info.Name()
		size := info.Size()
		isDir := info.IsDir()

		if isDir {
			subBranchLines, err := makeBranch(
				path+string(os.PathSeparator)+name,
				printFiles,
			)
			if err != nil {
				return nil, err
			}

			indentLines(subBranchLines)

			line := name
			lines = append(lines, line)
			lines = append(lines, subBranchLines...)
		} else if printFiles {
			sizeLabel := strconv.FormatInt(size, 10) + "b"
			if size == 0 {
				sizeLabel = "empty"
			}
			line := name + " (" + sizeLabel + ")"
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func indentLines(lines []string) {
	var intermediaryIndent = []string{
		"│   ",
		"│\t",
	}

	var finalIndent = []string{
		"    ",
		"\t",
	}

	isIntermediaryIndent := false
	lastFirstLevel := false
	for i := len(lines) - 1; i >= 0; i-- {
		subLine := lines[i]
		isBrRoot := isBranchRoot(subLine)

		if lastFirstLevel {
			isIntermediaryIndent = true
		}

		if isBrRoot {
			lastFirstLevel = true
		}

		isLast := len(lines)-1 == i

		lines[i] = finalIndent[INDENT_MODE] + subLine

		if isBrRoot {
			lines[i] = "├───" + subLine
		} else if isIntermediaryIndent {
			lines[i] = intermediaryIndent[INDENT_MODE] + subLine
		}

		if isBrRoot && (isLast || !isIntermediaryIndent) {
			lines[i] = "└───" + subLine
		}
	}
}

func isBranchRoot(line string) bool {
	return unicode.IsLetter([]rune(line)[0])
}
