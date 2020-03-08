package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	weekly        = flag.Duration("weekly", 24*time.Hour, "Weekly working time")
	dateInLayout  = flag.String("date-in", "06-1-2", "Layout of date input")
	timeInLayout  = flag.String("time-in", "1504", "Layout of time input")
	dateOutLayout = flag.String("date-out", "Mon 02.01.", "Layout of date output")
	timeOutLayout = flag.String("time-out", "15:04", "Layout of time output")
)

type entry struct {
	date     time.Time
	from     *time.Time
	to       *time.Time
	duration *time.Duration
}

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

func main() {
	flag.Parse()
	flag.Usage = usage
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	entries, err := readFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// sort by date
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].date.Before(entries[j].date)
	})

	year, week := 0, 0
	weekTotal := time.Duration(0)
	carry := time.Duration(0)

	for i, e := range entries {
		y, w := e.date.ISOWeek()
		if year != y || week != w {
			year, week = y, w
			weekTotal = time.Duration(0)
			carry -= *weekly
		}

		// assume last entry without duration ends now, if it is from today
		if lastEntry := i == len(entries)-1; lastEntry && e.to == nil && today(e.date) {
			d := time.Since(*e.from).Truncate(time.Minute)
			e.duration = &d
		}

		if e.duration != nil {
			weekTotal += *e.duration
		} else {
			fmt.Printf("Missing duration on %v\n", e.date)
			continue
		}

		carry += *e.duration
		date := e.date.Format(*dateOutLayout)
		fmt.Printf("%v\t%v\t%v\t%v\n", date, e.duration, weekTotal, carry)
	}
}

func today(t time.Time) bool {
	o := time.Now()
	return t.Day() == o.Day() && t.Month() == o.Month() && t.Year() == o.Year()
}

func readFile(filename string) ([]*entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := make([]*entry, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		e, err := parseEntry(line)
		if err != nil {
			log.Printf("Failed to parse entry: %v, err: %v", line, err)
			continue
		}

		if e == nil {
			continue
		}

		entries = append(entries, e)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func parseEntry(line string) (*entry, error) {
	// remove comments (starting with #)
	if i := strings.Index(line, "#"); i > -1 {
		line = line[:i]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	components := strings.SplitN(line, " ", 2)
	dateStr := components[0]
	durationStr := components[1]

	date, err := time.Parse(*dateInLayout, dateStr)
	if err != nil {
		return nil, err
	}

	var from, to *time.Time
	var duration *time.Duration

	// does not start with but contains a dash
	if strings.Index(durationStr, "-") > 0 {
		timeComponents := strings.SplitN(durationStr, "-", 2)
		fromStr := timeComponents[0]
		toStr := timeComponents[1]

		from, to, err = maybeParseSpecifiedDuration(dateStr, fromStr, toStr)
		if err != nil {
			return nil, err
		}

		if from != nil && to != nil {
			d := to.Sub(*from)
			duration = &d
		}

	} else {
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			return nil, err
		}

		duration = &d
	}

	return &entry{date, from, to, duration}, nil
}

func maybeParseSpecifiedDuration(dateStr, fromStr, toStr string) (from, to *time.Time, err error) {
	if fromStr == "" {
		err = errors.New("missing `from` date")
		return
	}

	from, err = parseWithDate(dateStr, fromStr)
	if err != nil {
		return
	}

	if toStr == "" {
		return
	}

	to, err = parseWithDate(dateStr, toStr)
	return
}

func parseWithDate(dateStr, timeStr string) (*time.Time, error) {
	layout := *dateInLayout + "/" + *timeInLayout
	str := dateStr + "/" + timeStr
	t, err := time.Parse(layout, str)
	if err != nil {
		return nil, err
	} else {
		return &t, err
	}
}
