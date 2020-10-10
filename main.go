package main

import (
	"flag"
	"fmt"
	"github.com/chermehdi/comet/eval"
	"github.com/chermehdi/comet/parser"
	"github.com/chermehdi/comet/repl"
	"io/ioutil"
	"os"
)

const VERSION = "1.0.0"

var BANNER = fmt.Sprintf(`
   ______                     __ 
  / ____/___  ____ ___  ___  / /_
 / /   / __ \/ __/ __ \/ _ \/ __/
/ /___/ /_/ / / / / / /  __/ /_
\____/\____/_/ /_/ /_/\___/\__/
				
author : @chermehdi
version: %s 
`, VERSION)

var filePath = flag.String("file", "", "File to run")

func main() {
	flag.Parse()
	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Println("Could not read passed file")
			return
		}
		source, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Could not read passed file")
			return
		}
		p := parser.New(string(source))
		rootNode := p.Parse()
		if p.Errors.HasAny() {
			fmt.Println(p.Errors)
			return
		}
		evaluator := eval.NewEvaluator()
		evaluator.Eval(rootNode)
	} else {
		// REPL MODE
		fmt.Print(BANNER)
		repler := repl.Repl{}
		repler.Start(os.Stdin, os.Stdout)
	}
}
