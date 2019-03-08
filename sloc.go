package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/tabwriter"
)

type Matcher func(string) bool

type Language struct {
	Name              string
	ExtensionMatcher  func(string) bool
	LineComment       string
	BlockCommentStart string
	BlockCommentEnd   string
}

type LanguageStats struct {
	Lang         Language
	FileCount    int
	CodeLines    int
	CommentLines int
	EmptyLines   int
}

var info = map[string]*LanguageStats{}
var files []string

var languages = []Language{
	Language{"Bash", matchExt(".sh", ".bashrc"), "#", "", ""},
	Language{"Golang", matchExt(".go"), "//", "/*", "*/"},
	Language{"Html", matchExt(".html"), "", "/*", "*/"},
	Language{"JSON", matchExt(".json"), "", "", ""},
	Language{"SCSS", matchExt(".scss"), "//", "/*", "*/"},
	Language{"Typescript", matchExt(".ts"), "//", "/*", "*/"},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getLang(filename string) Language {
	for _, lang := range languages {
		if lang.ExtensionMatcher(filename) {
			return lang
		}
	}

	return Language{"none", matchExt(), "", "", ""}
}

func matchExt(extensions ...string) Matcher {
	return func(filename string) bool {
		ext := path.Ext(filename)
		for _, e := range extensions {
			if ext == e {
				return true
			}
		}

		return false
	}
}

func countFileLines(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	check(err)

	lines := strings.Split(string(data), "\n")

	codeLines := 0
	commentLines := 0
	emptyLines := 0
	inComment := false

	lang := getLang(filepath)

	if lang.Name == "none" {
		return
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			emptyLines++
			continue
		}

		if inComment && lang.BlockCommentEnd != "" && strings.Contains(line, lang.BlockCommentEnd) {
			inComment = false
			commentLines++
			continue
		}

		if inComment || lang.LineComment != "" && strings.HasPrefix(line, lang.LineComment) {
			commentLines++
			continue
		}

		if lang.BlockCommentStart != "" && strings.HasPrefix(line, lang.BlockCommentStart) {
			inComment = true
			commentLines++
			continue
		}

		codeLines++
	}

	if val, ok := info[lang.Name]; ok {
		val.CodeLines += codeLines
		val.CommentLines += commentLines
		val.EmptyLines += emptyLines
		val.FileCount++
	} else {
		info[lang.Name] = &LanguageStats{lang, 1, codeLines, commentLines, emptyLines}
	}
}

func handleFile(filepath string) {
	fileInfo, err := os.Stat(filepath)
	check(err)

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(filepath)
		check(err)

		for _, file := range files {
			name := file.Name()
			if name == "." || name == ".." {
				continue
			}

			handleFile(path.Join(filepath, name))
		}
	} else if fileInfo.Mode()&os.ModeType == 0 {
		countFileLines(filepath)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	for _, path := range args {
		handleFile(path)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Language\tFiles\tCode\tComment\tBlank\t")

	for _, langInfo := range info {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", langInfo.Lang.Name, langInfo.FileCount, langInfo.CodeLines, langInfo.CommentLines, langInfo.EmptyLines)
	}

	w.Flush()
}
