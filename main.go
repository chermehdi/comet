package main

import (
	"fmt"
	"github.com/chermehdi/comet/repl"
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

func main() {
	fmt.Print(BANNER)
	repler := repl.Repl{}
	repler.Start(os.Stdin, os.Stdout)
}
