package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vicjohnson1213/sloc"
)

func main() {
	options := sloc.Options{}

	flag.StringVar(&options.Include, "include", "", "A regular expression for directories/files to include.")
	flag.StringVar(&options.Exclude, "exclude", "", "A regular expression for directories/files to exclude.")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	info := sloc.CountLines(args[0], options)

	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Language\tFiles\tCode\tComment\tBlank\t")
	total := sloc.LanguageStats{
		Lang:         sloc.Language{},
		FileCount:    0,
		CodeLines:    0,
		CommentLines: 0,
		EmptyLines:   0,
	}

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