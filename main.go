package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
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
