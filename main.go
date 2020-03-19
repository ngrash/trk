package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"
)

var (
	weekly       = flag.Duration("weekly", 24*time.Hour, "Weekly working time")
	dateInLayout = flag.String("date-in", "06-1-2", "Layout of date input")
	timeInLayout = flag.String("time-in", "1504", "Layout of time input")
)

const (
	dateOutLayout = "Mon 02.01."
	timeOutLayout = "15:04"
)

type entry struct {
	date     time.Time
	from     *time.Time
	to       *time.Time
	duration *time.Duration
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
		return entries[i].date.Before(entries[j].date) && entries[i].from.Before(*entries[j].from)
	})

	year, week, day := 0, 0, 0
	weekTotal := time.Duration(0)
	dayTotal := time.Duration(0)
	carry := time.Duration(0)

	fmt.Println("Date       From  To    Time   Day    Week   Total")

	for i, e := range entries {
		y, w := e.date.ISOWeek()
		if year != y || week != w {
			year, week = y, w
			weekTotal = time.Duration(0)
			carry -= *weekly
		}

		if d := e.date.YearDay(); d != day || y != year {
			// edge case: if you only logged once a year but at the
			// same day, we would not know it were different days
			// if we did not check the year.
			year = y
			day = d
			dayTotal = time.Duration(0)
		}

		// assume last entry without duration ends now, if it is from today
		if lastEntry := i == len(entries)-1; lastEntry && e.to == nil && today(e.date) {
			d := time.Since(*e.from).Truncate(time.Minute)
			e.duration = &d
		}

		if e.duration != nil {
			weekTotal += *e.duration
			dayTotal += *e.duration
		} else {
			fmt.Printf("Missing duration on %v\n", e.date)
			continue
		}

		carry += *e.duration
		date := e.date.Format(dateOutLayout)
		fmt.Printf("%v %v-%v %v %v %v %v\n",
			date,
			e.from.Format(timeOutLayout),
			e.to.Format(timeOutLayout),
			d2s(*e.duration, false),
			d2s(dayTotal, false),
			d2s(weekTotal, false),
			d2s(carry, true))
	}
}

func d2s(d time.Duration, negPossible bool) string {
	n := int64(d) / int64(time.Minute)
	min := int64(math.Abs(float64(n % 60)))
	hours := int64(n / 60)
	if negPossible {
		return fmt.Sprintf("%3vh%02dm", hours, min)
	} else {
		return fmt.Sprintf("%2vh%02dm", hours, min)
	}
}

func today(t time.Time) bool {
	o := time.Now()
	return t.Day() == o.Day() && t.Month() == o.Month() && t.Year() == o.Year()
}
