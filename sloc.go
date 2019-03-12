package sloc

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// Language is a single programming language and it's details.
type Language struct {
	Name       string      `json:"Name"`
	Extensions []string    `json:"Extensions"`
	Comments   CommentData `json:"Comments"`
}

// Info are the statistics gathered from all files counted.
type Info map[string]*LanguageStats

// LanguageStats are the statistics gatherd about a source file type.
type LanguageStats struct {
	Lang         Language `json:"Language"`
	FileCount    int      `json:"FileCount"`
	CodeLines    int      `json:"CodeLines"`
	CommentLines int      `json:"CommentLines"`
	EmptyLines   int      `json:"EmptyLines"`
}

// CommentData describes comment lines/blocks for a longuage.
type CommentData struct {
	LineCommentPrefixes []string `json:"LineCommentPrefixes"`
	BlockCommentPrefix  string   `json:"BlockCommentPrefix"`
	BlockCommentSuffix  string   `json:"BlockCommentSuffix"`
}

// Options are the options available to the line counter.
type Options struct {
	Include string
	Exclude string
	Format  string
}

// CountLines counts the lines of a file or recursively counts the lines of all files in a directory.
func CountLines(filepath string, options Options) Info {
	countLines(filepath, options)
	return info
}

var info = Info{}
var files []string

func getLang(filename string) Language {
	ext := path.Ext(filename)
	for _, lang := range languages {
		if contains(lang.Extensions, ext) {
			return lang
		}
	}

	return Language{"none", []string{}, CommentData{}}
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

		lineComment := startsWithOneOf(line, lang.Comments.LineCommentPrefixes)
		startsComment := lang.Comments.BlockCommentPrefix != "" && strings.HasPrefix(line, lang.Comments.BlockCommentPrefix)
		endsComment := strings.Contains(line, lang.Comments.BlockCommentSuffix)

		if lineComment || startsComment && endsComment {
			commentLines++
			continue
		}

		if inComment && endsComment {
			inComment = false
			commentLines++
			continue
		}

		if inComment || startsComment {
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
		info[lang.Name] = &LanguageStats{
			lang,
			1,
			codeLines,
			commentLines,
			emptyLines,
		}
	}
}

func countLines(filepath string, options Options) {
	shouldInclude, _ := regexp.MatchString(options.Include, filepath)
	shouldExclude, _ := regexp.MatchString(options.Exclude, filepath)

	if options.Exclude != "" && shouldExclude {
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

			CountLines(path.Join(filepath, name), options)
		}
	} else if shouldInclude && fileInfo.Mode()&os.ModeType == 0 {
		countFileLines(filepath)
	}
}

var languages = []Language{
	Language{"Assembly", []string{".asm", ".s"}, semiComments},
	Language{"Bash", []string{".bash", ".bashrc", ".bash_profile"}, bashComments},
	Language{"Batch", []string{".bat", ".cmd"}, batchComments},
	Language{"C", []string{".c", ".cc", ".h"}, cComments},
	Language{"C++", []string{".cpp", ".cxx"}, cComments},
	Language{"C#", []string{".cs"}, cComments},
	Language{"Clojure", []string{".clj", ".cljs", ".cljc", ".cljx", ".clojure", ".edn"}, clojureComments},
	Language{"Coffeescript", []string{".coffee", ".cson"}, coffeeComments},
	Language{"CSS", []string{".css"}, cssComments},
	Language{"Erlang", []string{".erl"}, erlangComments},
	Language{"Golang", []string{".go"}, cComments},
	Language{"Groovy", []string{".groovy"}, cComments},
	Language{"Haskell", []string{".hs", ".lhs"}, haskellComments},
	Language{"Html", []string{".html"}, xmlComments},
	Language{"Java", []string{".java"}, cComments},
	Language{"JavaScript", []string{".js"}, cComments},
	Language{"JSON", []string{".json"}, noComments},
	Language{"Kotlin", []string{".kt"}, cComments},
	Language{"LESS", []string{".less"}, cComments},
	Language{"Lisp", []string{".lsp", ".lisp"}, semiComments},
	Language{"Lua", []string{".lua"}, luaComments},
	Language{"Make", []string{"makefile", "Makefile", "MAKEFILE"}, bashComments},
	Language{"Markdown", []string{".md"}, noComments},
	Language{"Objective-C", []string{".m", ".mm", ".M"}, cComments},
	Language{"Perl", []string{".pl", ".pm"}, perlComments},
	Language{"PHP", []string{".php", ".php3", ".php4", ".php5", ".phtml"}, phpComments},
	Language{"Python", []string{".py"}, pythonComments},
	Language{"R", []string{".r", ".R"}, bashComments},
	Language{"Ruby", []string{".rb"}, rubyComments},
	Language{"Rust", []string{".rs", ".rc"}, cComments},
	Language{"Scala", []string{".scala"}, cComments},
	Language{"Scheme", []string{".scm", ".scheme"}, semiComments},
	Language{"Swift", []string{".swift"}, cComments},
	Language{"SASS", []string{".sass"}, cComments},
	Language{"SCSS", []string{".scss"}, cComments},
	Language{"Shell", []string{".sh"}, bashComments},
	Language{"SQL", []string{".sql"}, sqlComments},
	Language{"Typescript", []string{".ts"}, cComments},
	Language{"VimL", []string{".vim"}, vimComments},
	Language{"Visual Basic", []string{".vb"}, vbComments},
	Language{"XML", []string{".xml"}, xmlComments},
}

var (
	noComments      = CommentData{[]string{}, "", ""}
	bashComments    = CommentData{[]string{"#"}, "", ""}
	batchComments   = CommentData{[]string{"REM"}, "", ""}
	cComments       = CommentData{[]string{"//"}, "/*", "*/"}
	clojureComments = CommentData{[]string{";;"}, "", ""}
	coffeeComments  = CommentData{[]string{"#"}, "###", "###"}
	cssComments     = CommentData{[]string{}, "/*", "*/"}
	erlangComments  = CommentData{[]string{"%"}, "", ""}
	haskellComments = CommentData{[]string{"--"}, "{-", "-}"}
	luaComments     = CommentData{[]string{"--"}, "--[[", "]]"}
	perlComments    = CommentData{[]string{"#"}, "###", "###"}
	phpComments     = CommentData{[]string{"//", "#"}, "/*", "*/"}
	pythonComments  = CommentData{[]string{"#"}, `"""`, `"""`}
	rubyComments    = CommentData{[]string{"#"}, "=begin", "=end"}
	semiComments    = CommentData{[]string{";"}, "", ""}
	sqlComments     = CommentData{[]string{"--"}, "/*", "*/"}
	vbComments      = CommentData{[]string{"'"}, "", ""}
	vimComments     = CommentData{[]string{`"`}, "", ""}
	xmlComments     = CommentData{[]string{}, "<!--", "-->"}
)

func contains(array []string, element string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}

	return false
}

func startsWithOneOf(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}

	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
