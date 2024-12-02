package template

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"text/template"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/path"

	"github.com/essentialkaos/shdoc/script"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Render prints script info into terminal
func Render(doc *script.Document, tmpl, output string) error {
	templateFile := getPathToTemplate(tmpl)

	if templateFile == "" {
		return fmt.Errorf("Can't find template %q", tmpl)
	}

	t, err := readTemplate(templateFile)

	if err != nil {
		return err
	}

	fd, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer fd.Close()

	err = t.Execute(fd, doc)

	if err != nil {
		return err
	}

	printDocumentStats(doc, output)

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getPathToTemplate returns path to template file
func getPathToTemplate(tmpl string) string {
	if fsutil.IsExist(tmpl) {
		return tmpl
	}

	templateFile := path.Join(
		os.Getenv("GOPATH"),
		"src/github.com/essentialkaos/shdoc/templates",
		tmpl+".tpl",
	)

	if fsutil.IsExist(templateFile) {
		return templateFile
	}

	return ""
}

// readTemplate reads template
func readTemplate(templateFile string) (*template.Template, error) {
	err := fsutil.ValidatePerms("FRS", templateFile)

	if err != nil {
		return nil, err
	}

	templateData, err := os.ReadFile(templateFile)

	if err != nil {
		return nil, err
	}

	tmpl := template.New("Template")

	return tmpl.Parse(string(templateData))
}

// printDocumentStats prints information about document
func printDocumentStats(doc *script.Document, output string) {
	fmtutil.Separator(false, doc.Title)

	fmtc.Printfn("  {*}Constants:{!} %d", len(doc.Constants))
	fmtc.Printfn("  {*}Variables:{!} %d", len(doc.Variables))
	fmtc.Printfn("  {*}Methods:{!}   %d", len(doc.Methods))

	fmtc.NewLine()

	fmtc.Printfn(
		"  {*}Output:{!} %s {s-}(%s){!}", output,
		fmtutil.PrettySize(fsutil.GetSize(output)),
	)

	fmtutil.Separator(false)
}
