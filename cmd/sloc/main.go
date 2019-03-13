package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vicjohnson1213/sloc"
)

func main() {
	options := sloc.Options{}
	var format string

	flag.StringVar(&options.Include, "include", "", "A regular expression for directories/files to include.")
	flag.StringVar(&options.Include, "i", "", "A regular expression for directories/files to include.")

	flag.StringVar(&options.Exclude, "exclude", "", "A regular expression for directories/files to exclude.")
	flag.StringVar(&options.Exclude, "e", "", "A regular expression for directories/files to exclude.")

	flag.StringVar(&format, "format", "table", "The desired output format (table, JSON).")
	flag.StringVar(&format, "f", "table", "The desired output format (table, json, csv).")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	info := sloc.CountLines(args[0], options)

	switch format {
	case "json":
		outputJSON(info)
	case "csv":
		outputCSV(info)
	default:
		outputTable(info)
	}
}

func outputTable(info sloc.Info) {
	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Language\tFiles\tCode\tComment\tBlank\t")
	total := sloc.LanguageStats{
		Lang:         sloc.Language{Name: "Total"},
		FileCount:    0,
		CodeLines:    0,
		CommentLines: 0,
		EmptyLines:   0,
	}

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

func outputCSV(info sloc.Info) {
	println("Language, Files, Code, Comment, Blank")

	for _, langInfo := range info {
		fmt.Printf("%s, %d, %d, %d, %d\n", langInfo.Lang.Name, langInfo.FileCount, langInfo.CodeLines, langInfo.CommentLines, langInfo.EmptyLines)
	}
}

func outputJSON(info sloc.Info) {
	result, _ := json.MarshalIndent(info, "", "  ")
	println(string(result))
}
