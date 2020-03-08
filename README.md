# trk - plain text time tracking tool

## Installation
`go get github.com/ngrash/trk`

## Usage
```
$ trk -help
Usage of [./trk] [options] <filename>

Options:

  -date-in string
    	Layout of date input (default "06-1-2")
  -date-out string
    	Layout of date output (default "Mon 02.01.")
  -time-in string
    	Layout of time input (default "1504")
  -time-out string
    	Layout of time output (default "15:04")
  -weekly duration
    	Weekly working time (default 24h0m0s)

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

  Output consists of four columns separated by a tab:

  date:		Date of the entry
  duration:	How much you worked that day (see special case above)
  week total:	How much you worked this week
  overtime:	In combination with -weekly option, keeps a running
           	total of your overtime

  Example (with -weekly 24h):
    Wed 04.03.  6h20m0s  22h40m0s  -1h20m0s
    Fri 06.03.  1h45m0s  24h25m0s  25m0s
    Mon 09.03.  6h0m0s   6h0m0s    -17h35m0s

    On Wednesday I worked for 6 hours and 20 minutes which makes for a week total of
     22 hours and 40 minutes. I still had to work for 1 hour and 20 minutes that week.
    On Friday I worked for 1 hour and 45 minutes. Having worked 24 hours and 25 minutes
     in total that week I left with 25 minutes of overtime.
    On Monday the next week I worked 6 hours and 10 minutes. Since I had 25 minutes
     of overtime I had only 17 hours and 35 minutes left to work that week.
```
