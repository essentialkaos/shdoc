package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/shdoc/parser"
	"github.com/essentialkaos/shdoc/render/template"
	"github.com/essentialkaos/shdoc/render/terminal"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SHDoc"
	VER  = "0.9.0"
	DESC = "Tool for viewing and exporting docs for shell scripts"
)

const (
	OPT_OUTPUT   = "o:output"
	OPT_TEMPLATE = "t:template"
	OPT_NAME     = "n:name"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_TEMPLATE: {Value: "html"},
	OPT_NAME:     {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main function
func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	switch len(args) {
	case 1:
		process(args[0], "")
	case 2:
		process(args[0], args[1])
	default:
		showUsage()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// process starts source processing
func process(file string, pattern string) {
	err := fsutil.ValidatePerms("FRS", file)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	doc, errs := parser.Parse(file)

	if len(errs) != 0 {
		printErrorsAndExit(errs)
	}

	if !doc.IsValid() {
		printWarn("File %s doesn't contains any documentation", file)
		os.Exit(2)
	}

	if options.GetS(OPT_NAME) != "" {
		doc.Title = options.GetS(OPT_NAME)
	}

	if options.GetS(OPT_OUTPUT) == "" {
		err = terminal.Render(doc, pattern)
	} else {
		err = template.Render(
			doc,
			options.GetS(OPT_TEMPLATE),
			options.GetS(OPT_OUTPUT),
		)
	}

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorsAndExit prints errors and exit with exit code 1
func printErrorsAndExit(errs []error) {
	printError("Shell script docs parsing errors:")

	for _, err := range errs {
		printError("  %s", err.Error())
	}

	os.Exit(1)
}

// printErrorAndExit prints error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "file")

	info.AddOption(OPT_OUTPUT, "Path to output file", "file")
	info.AddOption(OPT_TEMPLATE, "Name of template", "name")
	info.AddOption(OPT_NAME, "Overwrite default name", "name")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"script.sh",
		"Parse shell script and show docs in console",
	)

	info.AddExample(
		"script.sh -t markdown -o my_script.md",
		"Parse shell script and save docs using given export template",
	)

	info.AddExample(
		"script.sh someEntity",
		"Parse shell script and show docs for some constant, variable or method",
	)

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "shdoc"))
	case "fish":
		fmt.Printf(fish.Generate(info, "shdoc"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "shdoc"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout shows info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/shdoc", update.GitHubChecker},
	}

	about.Render()
}
