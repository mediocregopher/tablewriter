package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/mediocregopher/tablewriter"

)

func main() {
	big := strings.Repeat("o ", 52)

	// Create  a table writer which will write the processed data to Stdout
	t := tablewriter.New(os.Stdout)

	// Set the width to not be the default but to instead be 200 chars
	t.SetTableWidth(200)

	// Add a column with the default padding and a width of 75 chars
	t.AddColumn(-1, 75)

	// Add a column with padding of 3 and a width of 25 chars
	t.AddColumn(3,  25)

	// Add a column with padding of 10 and an automatically specified width
	t.AddColumn(10, -1)

	// Add a column with no padding on the right (since its the rightmost column
	// the padding isn't necessary) and an automatically specified width
	t.AddColumn(0,  -1)

	// The automatically specified widths are calculated by taking the leftover
	// width (200 - (75 + 25) = 100) and dividing it by the number of columns
	// that need their widths filled in (2). So in this case the third and
	// fourth columns will have widths of 50.

	// Write 3 rows to the tablewriter. Once we start writing data we can't
	// change any of the options or add any more columns. If we do, santa won't
	// give us any presents on christmas.
	for i := 0; i < 3; i++ {
		// Whenever we write a full row (ending in a newline) tablewriter will
		// immediately process it and output it to the Writer given in the call
		// to New. So you can effectively stream as much data as you like.
		_, err := fmt.Fprintf(t, "%s\t%s\t%s\t%s\n", big, big, big, big)
		if err != nil {
			panic(err)
		}
	}
}
