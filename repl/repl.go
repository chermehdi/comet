package repl

import (
	"bufio"
	"fmt"
	"github.com/chermehdi/comet/eval"
	"github.com/chermehdi/comet/parser"
	"io"
)

type Repl struct{}

func (r *Repl) Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	evaluator := eval.Evaluator{}

	for {
		fmt.Fprint(writer, ">>")
		if !scanner.Scan() {
			fmt.Fprint(writer, "Good by !")
			break
		}

		line := scanner.Text()
		parser := parser.New(line)
		rootNode := parser.Parse()

		res := evaluator.Eval(rootNode)
		if res != nil {
			fmt.Fprintln(writer, res.ToString())
		}
	}
}
