package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/mediocregopher/tablewriter"

)

func main() {
	t := tablewriter.New(os.Stdout)
	t.SetTableWidth(200)
	t.AddColumn(-1, 75)
	t.AddColumn(3,  25)
	t.AddColumn(10, -1)
	t.AddColumn(0,  -1)

	n := 91

	_, err := fmt.Fprintf(
		t,
		"%s\t%s\t%s\t%s\n",
		strings.Repeat("a ",n),
		strings.Repeat("b ",n),
		strings.Repeat("c ",n),
		strings.Repeat("d ",n),
	)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(
		t,
		"%s\t%s\t%s\t%s\n",
		strings.Repeat("e ",n),
		strings.Repeat("f ",n),
		strings.Repeat("g ",n),
		strings.Repeat("h ",n),
	)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(
		t,
		"%s\t%s\t%s\t%s\n",
		strings.Repeat("i ",n),
		strings.Repeat("j ",n),
		strings.Repeat("k ",n),
		strings.Repeat("l ",n),
	)
	if err != nil {
		panic(err)
	}
}
