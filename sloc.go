package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/tabwriter"
)

// Matcher : A function that attempts to match a given string.
type Matcher func(string) bool

// Language : A programming language and it's details.
type Language struct {
	Name              string
	ExtensionMatcher  func(string) bool
	LineComment       string
	BlockCommentStart string
	BlockCommentEnd   string
}

// LanguageStats : The statistics gatherd about a source file.
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
	Language{"Batch", matchExt(".bat", ".cmd"), "REM", "", ""},
	Language{"C", matchExt(".c", ".cc", ".h"), "//", "/*", "*/"},
	Language{"C++", matchExt(".cpp", ".cxx"), "//", "/*", "*/"},
	Language{"C#", matchExt(".cs"), "//", "/*", "*/"},
	Language{"Clojure", matchExt(".clj", ".cljs", ".cljc", ".cljx", ".clojure", ".edn"), ";;", "", ""},
	Language{"Coffeescript", matchExt(".coffee", ".cson"), "#", "###", "###"},
	Language{"CSS", matchExt(".css"), "", "/*", "*/"},
	Language{"Bash", matchExt(".sh", ".bashrc", ".bash_profile"), "#", "", ""},
	Language{"Golang", matchExt(".go"), "//", "/*", "*/"},
	Language{"Html", matchExt(".html"), "", "<!--", "-->"},
	Language{"JavaScript", matchExt(".js"), "//", "/*", "*/"},
	Language{"JSON", matchExt(".json"), "", "", ""},
	Language{"LESS", matchExt(".less"), "//", "/*", "*/"},
	Language{"Python", matchExt(".py"), "#", `"""`, `"""`},
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

func handleFile(filepath string, exclude string) {
	matched, _ := regexp.MatchString(exclude, filepath)
	if exclude != "" && matched {
		return
	}

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

			handleFile(path.Join(filepath, name), exclude)
		}
	} else if fileInfo.Mode()&os.ModeType == 0 {
		countFileLines(filepath)
	}
}

func main() {
	var exclude string

	flag.StringVar(&exclude, "exclude", "", "A regular expression for directories/files to exclude.")
	flag.Parse()

	args := flag.Args()

	for _, path := range args {
		handleFile(path, exclude)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Language\tFiles\tCode\tComment\tBlank\t")
	total := LanguageStats{Language{}, 0, 0, 0, 0}
	total.Lang.Name = "Total"

	for _, langInfo := range info {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", langInfo.Lang.Name, langInfo.FileCount, langInfo.CodeLines, langInfo.CommentLines, langInfo.EmptyLines)
		total.FileCount += langInfo.FileCount
		total.CodeLines += langInfo.CodeLines
		total.CommentLines += langInfo.CommentLines
		total.EmptyLines += langInfo.EmptyLines
	}

	fmt.Fprintln(w, "\t\t\t\t\t")
	fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n", total.Lang.Name, total.FileCount, total.CodeLines, total.CommentLines, total.EmptyLines)

	w.Flush()
}
