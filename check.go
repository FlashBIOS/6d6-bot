package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
)

func check(err error) {
	if err != nil {
		exitWithMessage(err.Error())
	}
}

func checkWithMessage(err error, message string) {
	if err != nil {
		exitWithMessage(message)
	}
}

func exitWithMessage(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s %v\n", aurora.Bold(aurora.Red("Error: ")), aurora.Red(message))
	os.Exit(1)
}
