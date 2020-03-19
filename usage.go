package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	pf := func(format string, a ...interface{}) {
		fmt.Fprintf(flag.CommandLine.Output(), format, a)
	}
	pln := func(s string) {
		fmt.Fprintln(flag.CommandLine.Output(), s)
	}

	pf("Usage of %s [options] <filename>\n\n", os.Args[0])
	pln("Options:\n")
	flag.PrintDefaults()
	pln("")
	pln("  Read more about date and time layouts here: https://golang.org/pkg/time/#Parse")
	pln("")
	pln("File format:\n")

	pln("  Each non-empty line of the input file specified by <filename> represents an entry.")
	pln("  An entry consists of a date and a duration, separated by a space.")
	pln("  The duration is either an absolute value or a time range.")
	pln("    Absolute value: 6h12m0s")
	pln("    Time range: 0810-1422")
	pln("")
	pln("  Special case:")
	pln("    The last entry of a file can be an open range (e.g. `0930-`) and trk")
	pln("    will assume you are still working, using the current time to calculate")
	pln("    the duration.")
	pln("")
	pln("  Comments:")
	pln("    Everything following a number sign (#) on the same line is ignored.")
	pln("")
	pln("  Examples:")
	pln("    19-9-25 0950-1830")
	pln("    19-9-26 8h # took a day off")
	pln("    19-9-27 1015-")
	pln("")
	pln("Output:\n")
	pln("  Output consists of seven columns:\n")
	pln("  Date\tDate of the entry")
	pln("  From\tStart time of the entry")
	pln("  To\tEnd time of the entry")
	pln("  Dur.\tDuration of the entry")
	pln("  Day\tHow much you worked that day (see special case above)")
	pln("  Week\tHow much you worked that week")
	pln("  Total\tIn combination with -weekly option, keeps a running")
	pln("       \ttotal of your overtime")
	pln("")
}
