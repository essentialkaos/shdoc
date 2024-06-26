package template

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/path"

	"github.com/essentialkaos/shdoc/script"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Render prints script info into terminal
func Render(doc *script.Document, tmpl, output string) error {
	templateFile := getPathToTemplate(tmpl)

	if templateFile == "" {
		return fmt.Errorf("Can't find template \"%s\"", tmpl)
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

	templateData, err := ioutil.ReadFile(templateFile)

	if err != nil {
		return nil, err
	}

	tmpl := template.New("Template")

	return tmpl.Parse(string(templateData))
}

// printDocumentStats prints information about document
func printDocumentStats(doc *script.Document, output string) {
	fmtutil.Separator(false, doc.Title)

	fmtc.Printf("  {*}Constants:{!} %d\n", len(doc.Constants))
	fmtc.Printf("  {*}Variables:{!} %d\n", len(doc.Variables))
	fmtc.Printf("  {*}Methods:{!}   %d\n", len(doc.Methods))

	fmtc.NewLine()

	fmtc.Printf(
		"  {*}Output:{!} %s {s-}(%s){!}\n", output,
		fmtutil.PrettySize(fsutil.GetSize(output)),
	)

	fmtutil.Separator(false)
}
