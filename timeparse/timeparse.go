package timeparse

import (
	"strings"
	"sync"
)

/*
Transform

	%%
		a literal %
	%a
		locale's abbreviated weekday name (e.g., Sun)
	%A
		locale's full weekday name (e.g., Sunday)
	%b
		locale's abbreviated month name (e.g., Jan)
	%B
		locale's full month name (e.g., January)
	%c
		locale's date and time (e.g., Thu Mar 3 23:05:25 2005)
	%C
		century; like %Y, except omit last two digits (e.g., 20)
	%d
		day of month (e.g, 01)
	%D
		date; same as %m/%d/%y
	%e
		day of month, space padded; same as %_d
	%F
		full date; same as %Y-%m-%d
	%g not implemented
		last two digits of year of ISO week number (see %G)
	%G not implemented
		year of ISO week number (see %V); normally useful only with %V
	%h
		same as %b
	%H
		hour (00..23)
	%I
		hour (01..12)
	%j
		day of year (001..366)
	%k not implemented
		hour ( 0..23)
	%l not implemented
		hour ( 1..12)
	%m
		month (01..12)
	%M
		minute (00..59)
	%n
		a newline
	%N
		nanoseconds (000000000..999999999)
	%p
		locale's equivalent of either AM or PM; blank if not known
	%P
		like %p, but lower case
	%r
		locale's 12-hour clock time (e.g., 11:11:04 PM)
	%R
		24-hour hour and minute; same as %H:%M
	%s not implemented
		seconds since 1970-01-01 00:00:00 UTC
	%S
		second (00..60)
	%t
		a tab
	%T
		time; same as %H:%M:%S
	%u not implemented
		day of week (1..7); 1 is Monday
	%U not implemented
		week number of year, with Sunday as first day of week (00..53)
	%V not implemented
		ISO week number, with Monday as first day of week (01..53)
	%w not implemented
		day of week (0..6); 0 is Sunday
	%W not implemented
		week number of year, with Monday as first day of week (00..53)
	%x not implemented
		locale's date representation (e.g., 12/31/99)
	%X not implemented
		locale's time representation (e.g., 23:13:48)
	%y
		last two digits of year (00..99)
	%Y
		year
	%z
		+hhmm numeric timezone (e.g., -0400)
	%:z
		+hh:mm numeric timezone (e.g., -04:00)
	%::z
		+hh:mm:ss numeric time zone (e.g., -04:00:00)
	%:::z
		numeric time zone with : to necessary precision (e.g., -04, +05:30)
	%Z
		alphabetic time zone abbreviation (e.g., EDT)
*/

var replacer *strings.Replacer
var replacerInit sync.Once

func Transform(format string) (goFormat string) {

	replacerInit.Do(func() {
		replacer = strings.NewReplacer(
			`%%`, "%",
			`%a`, "Mon",
			`%A`, "Monday",
			`%b`, "Jan",
			`%B`, "January",
			`%c`, "Mon Jan 2 15:04:05 2006",
			`%C`, "20",
			`%d`, "02",
			`%D`, "01/02/06",
			`%e`, "_2",
			`%F`, "2006-01-02",
			`%h`, "Jan",
			`%H`, "15",
			`%I`, "03",
			`%j`, "002",
			`%m`, "01",
			`%M`, "04",
			`%n`, "\n",
			`%N`, "000000000",
			`%p`, "PM",
			`%P`, "pm",
			`%r`, "03:04:05 PM",
			`%R`, "15:04",
			`%S`, "05",
			`%t`, "\t",
			`%T`, "15:04:05",
			`%y`, "06",
			`%Y`, "2006",
			`%z`, "-0700",
			`%:z`, "-07:00",
			`%::z`, "-07:00:00",
			`%:::z`, "-07",
			`%Z`, "MST",
		)
	})

	return replacer.Replace(format)
}
