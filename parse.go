package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"time"
)

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
