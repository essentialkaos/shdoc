package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/pager"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/apps"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	term "github.com/essentialkaos/ek/v13/terminal"

	"github.com/essentialkaos/shdoc/parser"
	"github.com/essentialkaos/shdoc/render/template"
	"github.com/essentialkaos/shdoc/render/terminal"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SHDoc"
	VER  = "0.10.2"
	DESC = "Tool for viewing and exporting docs for shell scripts"
)

const (
	OPT_OUTPUT   = "o:output"
	OPT_TEMPLATE = "t:template"
	OPT_NAME     = "n:name"
	OPT_NO_PAGER = "np:no-pager"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {},
	OPT_TEMPLATE: {Value: "html"},
	OPT_NAME:     {},
	OPT_NO_PAGER: {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main application function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		term.Error("Options parsing errors:")
		term.Error(errs.Error(" - "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithApps(apps.Bash()).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP), len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err := readDocs(
		args.Get(0).Clean().String(),
		args.Get(1).String(),
	)

	if err != nil {
		term.Error(err)
		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// readDocs reads the file and prints documentation from it
func readDocs(file string, pattern string) error {
	err := fsutil.ValidatePerms("FRS", file)

	if err != nil {
		return err
	}

	doc, errs := parser.Parse(file)

	if !errs.IsEmpty() {
		term.Error("Shell script documentation parsing errors:")
		term.Error(errs.Error(" - "))
		fmtc.NewLine()
		return fmt.Errorf("Can't parse script documentation")
	}

	if !doc.IsValid() {
		return fmt.Errorf("File %s doesn't contains any documentation", file)
	}

	if options.GetS(OPT_NAME) != "" {
		doc.Title = options.GetS(OPT_NAME)
	}

	if !options.Has(OPT_OUTPUT) {
		if !options.GetB(OPT_NO_PAGER) {
			if tty.IsTTY() {
				if pager.Setup() == nil {
					defer pager.Complete()
				}
			}
		}

		err = terminal.Render(doc, pattern)
	} else {

		err = template.Render(
			doc,
			options.GetS(OPT_TEMPLATE),
			options.GetS(OPT_OUTPUT),
		)
	}

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, "shdoc"))
	case "fish":
		fmt.Print(fish.Generate(info, "shdoc"))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, "shdoc"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "script")

	info.AddOption(OPT_OUTPUT, "Path to output file", "file")
	info.AddOption(OPT_TEMPLATE, "Name of template", "name")
	info.AddOption(OPT_NAME, "Overwrite default name", "name")
	info.AddOption(OPT_NO_PAGER, "Disable pager for long output")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"script.sh",
		"Parse shell script and show documentation in console",
	)

	info.AddExample(
		"script.sh -t markdown -o my_script.md",
		"Parse shell script and render documentation to markdown file",
	)

	info.AddExample(
		"script.sh -t /path/to/template.tpl -o my_script.ext",
		"Parse shell script and render documentation with given template",
	)

	info.AddExample(
		"script.sh myFunction",
		"Parse shell script and show documentation for some constant, variable or method",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{
			"essentialkaos/shdoc",
			update.GitHubChecker,
		}
	}

	return about
}
