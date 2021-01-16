package repl

import (
	"bufio"
	"fmt"
	"github.com/chermehdi/comet/eval"
	"github.com/chermehdi/comet/parser"
	"io"
	"strings"
)

type Repl struct{}

func (r *Repl) Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	evaluator := eval.NewEvaluator()

	for {
		fmt.Fprint(writer, ">> ")
		if !scanner.Scan() {
			fmt.Fprintln(writer, "Good by !")
			break
		}
		line := scanner.Text()
		if line = strings.Trim(line, " \n\t\r"); line == "" {
			continue
		}
		if line[0] == '/' {
			if line == "/exit" {
				break
			}
			if line == "/scope" {
				printScope(evaluator)
				continue
			}
		}
		p := parser.New(line)
		rootNode := p.Parse()
		if p.Errors.HasAny() {
			fmt.Fprintln(writer, p.Errors)
			continue
		}
		res := evaluator.Eval(rootNode)
		if res != nil {
			fmt.Fprintln(writer, res.ToString())
		}
	}
}

func printScope(eval *eval.Evaluator) {
	fmt.Println("==== Variables ====")
	scope := eval.Scope
	for cur := scope; cur != nil; cur = cur.Parent {
		for k, v := range cur.Variables {
			fmt.Println(fmt.Sprintf("%s = %v", k, v.Type()))
		}
	}
	fmt.Println("==== Types ====")
	for _, t := range eval.Types {
		fmt.Println(t.Name)
	}
}
