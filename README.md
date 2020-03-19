# trk - plain text time tracking tool

## Installation
`go get github.com/ngrash/trk`

## Example
```
$ trk -this-week ~/companytime
Date       From  To    Dur.   Day    Week   Total
Tue 17.03. 10:30-12:30  2h00m  2h00m  2h00m -21h30m
Tue 17.03. 13:00-16:30  3h30m  5h30m  5h30m -18h00m
Tue 17.03. 17:00-18:00  1h00m  6h30m  6h30m -17h00m
Wed 18.03. 08:00-09:00  1h00m  1h00m  7h30m -16h00m
Wed 18.03. 10:00-12:45  2h45m  3h45m 10h15m -13h15m
Wed 18.03. 13:50-15:00  1h10m  4h55m 11h25m -12h05m
Wed 18.03. 16:00-17:15  1h15m  6h10m 12h40m -10h50m
Thu 19.03. 13:20-15:00  1h40m  1h40m 14h20m  -9h10m
Thu 19.03. 15:34-18:30  2h56m  4h36m 17h16m  -6h14m
```

## Usage
```
$ trk -help
Usage of [trk] [options] <filename>

Options:

  -date-layout string
    	Layout of date input (default "06-1-2")
  -quiet
    	Do not print column names
  -this-week
    	Print only items from the current week
  -time-layout string
    	Layout of time input (default "1504")
  -weekly duration
    	Weekly working time (default 24h0m0s)

  Read more about date and time layouts here: https://golang.org/pkg/time/#Parse

File format:

  Each non-empty line of the input file specified by <filename> represents an entry.
  An entry consists of a date and a duration, separated by a space.
  The duration is either an absolute value or a time range.
    Absolute value: 6h12m0s
    Time range: 0810-1422

  Special case:
    The last entry of a file can be an open range (e.g. `0930-`) and trk
    will assume you are still working, using the current time to calculate
    the duration.

  Comments:
    Everything following a number sign (#) on the same line is ignored.

  Examples:
    19-9-25 0950-1830
    19-9-26 8h # took a day off
    19-9-27 1015-

Output:

  Output consists of seven columns:

  Date	Date of the entry
  From	Start time of the entry
  To	End time of the entry
  Dur.	Duration of the entry
  Day	How much you worked that day (see special case above)
  Week	How much you worked that week
  Total	In combination with -weekly option, keeps a running
       	total of your overtime
```
