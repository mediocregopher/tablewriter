// Package tablewriter implements a Writer middleman which takes in table data
// where rows are delimited by newlines and columns by tabs. It outputs this
// same data but formatted into an actual ascii table, where each column has a
// specific width that the text is wrapped into.
//
// Each cell has both a width and a padding, with the padding applied to the
// right edge of the cell and contributing to the width of the cell:
//
//		|----------------- WIDTH -----------------|
//		|------ CONTENT ------|----- PADDING -----|
//
// Check out the example/ folder of the package for a more in-depth usage
// example
package tablewriter

import (
	"fmt"
	"io"
	"bytes"
	"strings"
	"errors"
)

type Writer struct {
	bottomPad int
	width int
	cols []colSpec
	output io.Writer
	buf *bytes.Buffer
	filledcols bool
}

type colSpec struct {
	rightPad int
	width int
}

// Constants defining the defaults of a number of values
const (
	// The default padding for the bottom edge of each row in the table
	DEFAULT_BOTTOM_PAD = 1

	// The default padding for the right edge of each column in the table
	DEFAULT_RIGHT_PAD = 2

	// The default full width of the table
	DEFAULT_TABLE_WIDTH = 120
)

// Initiates a new tablewriter, taking in the io.Writer it will finally be
// written to (e.g. os.Stdout, a file, or a network connection)
func New(out io.Writer) *Writer {
	return &Writer{
		bottomPad: DEFAULT_BOTTOM_PAD,
		width: DEFAULT_TABLE_WIDTH,
		cols: make([]colSpec, 0, 8),
		output: out,
		buf: bytes.NewBuffer(make([]byte,0,64)),
		filledcols: false,
	}
}

// Sets the padding to be applied to the bottom edge of each row in the table to
// something other than the default
func (w *Writer) SetBottomPadding(p int) {
	w.bottomPad = p
}

// Sets the total table width to something other than the default
func (w *Writer) SetTableWidth(wid int) {
	w.width = wid
}

// Adds a column to the table, given the column's padding on the right edge and
// total width (including padding). If rightPad is -1 then the default is used, 
// if width is -1 then the column's size is calculated from the unaccounted for
// width in the table divided by the number of columns who don't specify their
// width
func (w *Writer) AddColumn(rightPad, width int) {
	w.cols = append(w.cols, colSpec{rightPad, width})
}

func (w *Writer) fillColSpecs() error {
	leftoverWidth := w.width
	numWidthUnspec := 0
	for i := range w.cols {
		col := &w.cols[i]
		if col.rightPad < 0 {
			col.rightPad = DEFAULT_RIGHT_PAD
		}
		if col.width < 0 {
			numWidthUnspec++
		} else {
			leftoverWidth -= col.width
		}
	}

	if numWidthUnspec > 0 {
		if leftoverWidth < 1 {
			return errors.New("Table not wide enough for given columns")
		}

		leftoverPer := leftoverWidth / numWidthUnspec
		leftoverExtra := leftoverWidth % numWidthUnspec

		// Go through and set the widths for all the columns that need it. We
		// only want to give the leftoverExtra (the remainder from taking the
		// average) to the final unspecified column in the table.
		for i := len(w.cols)-1; i >= 0; i-- {
			if w.cols[i].width < 0 {
				w.cols[i].width = leftoverPer + leftoverExtra
				leftoverExtra = 0
			}
		}
	}
	
	return nil
}

var space = []byte{' '}
var newline = []byte{'\n'}

func wrapFill(buf *bytes.Buffer, l, rp int) {
	fill := bytes.Repeat(space, (l - buf.Len()) + rp)
	buf.Write(fill)
}

func wrappedString(b string, l, rp int) ([]string, error) {
	fields := strings.Fields(b)
	buf := bytes.NewBuffer(make([]byte,0,(l+rp)*4))
	ret := make([]string, 0, 3)

	for f := 0;f < len(fields); {
		if len(fields[f]) > l {
			return ret, fmt.Errorf("Word '%s' too big for its cell", fields[f])
		} else if len(fields[f]) + buf.Len() <= l {
			buf.Write([]byte(fields[f]))
			if 1 + buf.Len() <= l {
				buf.Write(space)
			}
			f++
		} else {
			wrapFill(buf, l, rp)
			ret = append(ret, buf.String())
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		wrapFill(buf, l, rp)
		ret = append(ret, buf.String())
	}
	return ret, nil
}

func (w *Writer) writeRow(row string) (written int, err error) {
	cells := strings.Split(row, "\t")
	if len(cells) != len(w.cols) {
		return 0, errors.New("incorrect number of columns in given row")
	}
	cellswrapped := make([][]string, len(cells))
	for i, cellb := range cells {
		col := w.cols[i]
		w := col.width - col.rightPad
		rp := col.rightPad
		cellswrapped[i], err = wrappedString(cellb, w, rp)
		if err != nil {
			return
		}
	}
	lbuf := new(bytes.Buffer)
	for linecnt := 0;;linecnt++ {
		any := false
		var n int
		for i, cellwrapped := range cellswrapped {
			if linecnt >= len(cellwrapped) {
				lbuf.Write(bytes.Repeat(space, w.cols[i].width))
			} else {
				lbuf.Write([]byte(cellwrapped[linecnt]))
				any = true
			}
			written += n
			if err != nil {
				return
			}
		}
		if any {
			n, err = w.output.Write(lbuf.Bytes())
			written += n
			if err != nil {
				return
			}
			n, err = w.output.Write(newline)
			written += n
			if err != nil {
				return
			}
			lbuf.Reset()
		} else {
			break
		}
	}
	bottomPadding := bytes.Repeat(newline, w.bottomPad)
	n, err := w.output.Write(bottomPadding)
	written += n
	return written + n, err
}

// Implements the Writer interface. This will write the output stream as it gets
// enough information to write it, instead of buffering until the end.
func (w *Writer) Write(b []byte) (written int, err error) {
	if !w.filledcols {
		if err = w.fillColSpecs(); err != nil {
			return
		}
		w.filledcols = true
	}

	if _, err = w.buf.Write(b); err != nil {
		return
	}

	var row string
	var n int
	for {
		row, err = w.buf.ReadString('\n')
		if err == io.EOF {
			err = nil
			return
		} else if err != nil {
			return
		}
		if len(row) == 0 {
			return
		}
		n, err = w.writeRow(strings.TrimRight(row,"\n"))
		written += n
		if err != nil {
			return
		}
	}
}

