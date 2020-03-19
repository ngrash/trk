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
	pln("  Output consists of four columns separated by a tab:\n")
	pln("  date:\t\tDate of the entry")
	pln("  duration:\tHow much you worked that day (see special case above)")
	pln("  week total:\tHow much you worked this week")
	pln("  overtime:\tIn combination with -weekly option, keeps a running")
	pln("           \ttotal of your overtime")

	pln("")
	pln("  Example (with -weekly 24h):")
	pln("    Wed 04.03.  6h20m0s  22h40m0s  -1h20m0s")
	pln("    Fri 06.03.  1h45m0s  24h25m0s  25m0s")
	pln("    Mon 09.03.  6h0m0s   6h0m0s    -17h35m0s")
	pln("")
	pln("    On Wednesday I worked for 6 hours and 20 minutes which makes for a week total of")
	pln("     22 hours and 40 minutes. I still had to work for 1 hour and 20 minutes that week.")
	pln("    On Friday I worked for 1 hour and 45 minutes. Having worked 24 hours and 25 minutes")
	pln("     in total that week I left with 25 minutes of overtime.")
	pln("    On Monday the next week I worked 6 hours and 10 minutes. Since I had 25 minutes")
	pln("     of overtime I had only 17 hours and 35 minutes left to work that week.")
	pln("")
}
