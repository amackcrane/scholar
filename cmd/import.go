// Copyright © 2018 Eiji Onchi <eiji@onchi.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/cgxeiji/bib"
	"github.com/cgxeiji/scholar/scholar"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a bibtex/biblatex file",
	Long: `Scholar: a CLI Reference Manager

Import a bibtex/biblatex file into a library.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			importParse(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importParse(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	entries, err := bib.Unmarshal(file)
	if err != nil {
		panic(err)
	}

	var es []*scholar.Entry

	for _, entry := range entries {
		e, err := scholar.NewEntry(entry["type"])
		if err != nil {
			if err == scholar.TypeNotFoundError {
				e, err = scholar.NewEntry("misc")
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
		delete(entry, "type")
		e.Key = entry["key"]
		delete(entry, "key")

		if file, ok := entry["file"]; ok {
			e.File = file
			delete(entry, "file")
		}

		for req := range e.Required {
			e.Required[req] = entry[req]
			delete(entry, req)
		}

		for opt := range entry {
			e.Optional[opt] = entry[opt]
			delete(entry, opt)
		}

		es = append(es, e)
	}

	// Only commit after all entries have been validated
	for _, e := range es {
		commit(e)
		if e.File != "" {
			attach(e, e.File)
		}
	}

	fmt.Println("Import from", filename, "successful!")
}
