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
	evaluator := eval.New()

	for {
		fmt.Fprint(writer, ">> ")
		if !scanner.Scan() {
			fmt.Fprintln(writer, "Good by !")
			break
		}
		line := scanner.Text()
		if strings.Trim(line, " \n\t\r") == "" {
			continue
		}
		p := parser.New(line)
		rootNode := p.Parse()

		res := evaluator.Eval(rootNode)
		if res != nil {
			fmt.Fprintln(writer, res.ToString())
		}
	}
}
