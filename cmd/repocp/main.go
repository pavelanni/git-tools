package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pavelanni/git-tools/gitrepo"
	flag "github.com/spf13/pflag"
)

func main() {

	var srcRepo *string = flag.StringP("repo", "r", "", "source repo, required")
	var dstDir *string = flag.StringP("dst", "d", "", "destination directory, required")
	var branch *string = flag.StringP("branch", "b", "", "branch to copy from; default is HEAD")
	var help *bool = flag.BoolP("help", "h", false, "help")

	flags := flag.CommandLine
	flags.SortFlags = false

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *srcRepo == "" {
		fmt.Println("source repo is required")
		os.Exit(1)
	}
	if *dstDir == "" {
		fmt.Println("destination directory is required")
		os.Exit(1)
	}
	if *branch == "" {
		fmt.Println("branch not specified, defaulting to HEAD")
	}

	err := gitrepo.Copy(*srcRepo, *branch, *dstDir)
	if err != nil {
		log.Fatal(err)
	}
}
