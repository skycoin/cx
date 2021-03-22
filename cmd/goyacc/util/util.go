package util

import (
	"fmt"
	"io"
	"sort"
)

const (
	MaxInt  = 1<<(IntBits-1) - 1
	MinInt  = -MaxInt - 1
	MaxUint = 1<<IntBits - 1
	IntBits = 1 << (^uint(0)>>32&1 + ^uint(0)>>16&1 + ^uint(0)>>8&1 + 3)
)

// IsEOF reports whether err is an EOF condition.
func IsEOF(err error) bool { return err == io.EOF }

//Max returns the bigger of a and b
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min returns the smaller of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Formatter is an io.Writer extended by a fmt.Printf like function Format
type Formatter interface {
	io.Writer
	Format(format string, args ...interface{}) (n int, errno error)
}

type indentFormatter struct {
	io.Writer
	indent      []byte
	indentLevel int
	state       int
}

const (
	st0 = iota
	stBOL
	stPERC
	stBOLPERC
)

// IndentFormatter returns a new Formatter which interprets %i and %u in the
// Format() format string as indent and undent commands. The commands can
// nest. The Formatter writes to io.Writer 'w' and inserts one 'indent'
// string per current indent level value.
// Behaviour of commands reaching negative indent levels is undefined.
//	IndentFormatter(os.Stdout, "\t").Format("abc%d%%e%i\nx\ny\n%uz\n", 3)
// output:
//	abc3%e
//		x
//		y
//	z
// The Go quoted string literal form of the above is:
//	"abc%%e\n\tx\n\tx\nz\n"
// The commands can be scattered between separate invocations of Format(),
// i.e. the formatter keeps track of the indent level and knows if it is
// positioned on start of a line and should emit indentation(s).
// The same output as above can be produced by e.g.:
//	f := IndentFormatter(os.Stdout, " ")
//	f.Format("abc%d%%e%i\nx\n", 3)
//	f.Format("y\n%uz\n")
func IndentFormatter(w io.Writer, indent string) Formatter {
	return &indentFormatter{w, []byte(indent), 0, stBOL}
}

func (f *indentFormatter) format(flat bool, format string, args ...interface{}) (n int, errno error) {
	buf := []byte{}
	for i := 0; i < len(format); i++ {
		c := format[i]
		switch f.state {
		case st0:
			switch c {
			case '\n':
				cc := c
				if flat && f.indentLevel != 0 {
					cc = ' '
				}
				buf = append(buf, cc)
				f.state = stBOL
			case '%':
				f.state = stPERC
			default:
				buf = append(buf, c)
			}
		case stBOL:
			switch c {
			case '\n':
				cc := c
				if flat && f.indentLevel != 0 {
					cc = ' '
				}
				buf = append(buf, cc)
			case '%':
				f.state = stBOLPERC
			default:
				if !flat {
					for i := 0; i < f.indentLevel; i++ {
						buf = append(buf, f.indent...)
					}
				}
				buf = append(buf, c)
				f.state = st0
			}
		case stBOLPERC:
			switch c {
			case 'i':
				f.indentLevel++
				f.state = stBOL
			case 'u':
				f.indentLevel--
				f.state = stBOL
			default:
				if !flat {
					for i := 0; i < f.indentLevel; i++ {
						buf = append(buf, f.indent...)
					}
				}
				buf = append(buf, '%', c)
				f.state = st0
			}
		case stPERC:
			switch c {
			case 'i':
				f.indentLevel++
				f.state = st0
			case 'u':
				f.indentLevel--
				f.state = st0
			default:
				buf = append(buf, '%', c)
				f.state = st0
			}
		default:
			panic("unexpected state")
		}
	}
	switch f.state {
	case stPERC, stBOLPERC:
		buf = append(buf, '%')
	}
	return f.Write([]byte(fmt.Sprintf(string(buf), args...)))
}

func (f *indentFormatter) Format(format string, args ...interface{}) (n int, errno error) {
	return f.format(false, format, args...)
}

// Dedupe returns n, the number of distinct elements in data. The resulting
// elements are sorted in elements [0, n) or data[:n] for a slice.
func Dedupe(data sort.Interface) (n int) {
	if n = data.Len(); n < 2 {
		return n
	}

	sort.Sort(data)
	a, b := 0, 1
	for b < n {
		if data.Less(a, b) {
			a++
			if a != b {
				data.Swap(a, b)
			}
		}
		b++
	}
	return a + 1
}

// Uint64Slice attaches the methods of sort.Interface to []uint64, sorting in increasing order.
type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Sort is a convenience method.
func (s Uint64Slice) Sort() {
	sort.Sort(s)
}
