package pkg

import "time"

// these notations are supported by Go, the rule reports them anyway
var (
	// The month argument is 0 considered as January by Go, but the rule reports it anyway
	_ = time.Date(2023, 0, 2, 3, 4, 5, 6, time.UTC) // MATCH /time.Date month argument should not be zero/

	// The day argument is 0 considered as the first day of the month by Go, but the rule reports it anyway
	_ = time.Date(2023, 1, 0, 3, 4, 5, 6, time.UTC) // MATCH /time.Date day argument should not be zero/

	// there is no need to use a plus sign in front of numbers
	_ = time.Date(
		+2023, // MATCH /time.Date year argument contains a useless plus sign: +2023/
		+0b1,  // MATCH /time.Date month argument contains a useless plus sign: +0b1/
		+02,   // MATCH /time.Date day argument contains a useless plus sign: +02/
		+0o3,  // MATCH /time.Date hour argument contains a useless plus sign: +0o3/
		+4,    // MATCH /time.Date minute argument contains a useless plus sign: +4/
		+5,    // MATCH /time.Date second argument contains a useless plus sign: +5/
		+6,    // MATCH /time.Date nanosecond argument contains a useless plus sign: +6/
		time.UTC)

	// time.Date supports negative values, but it's uncommon to use them.
	// the rule reports them
	_ = time.Date(
		2023,
		-0b1, // MATCH /time.Date month argument is negative: -0b1/
		-02,  // MATCH /time.Date day argument is negative: -02/
		-0o3, // MATCH /time.Date hour argument is negative: -0o3/
		-0x4, // MATCH /time.Date minute argument is negative: -0x4/
		-5,   // MATCH /time.Date second argument is negative: -5/
		-6,   // MATCH /time.Date nanosecond argument is negative: -6/
		time.UTC)

	// using gigantic nanoseconds to move into the future
	// Go will compute 2054-09-10 04:50:45 +0000 UTC
	_ = time.Date(2023, 1, 2, 3, 4, 5, 1000000000000000000, time.UTC) // MATCH /time.Date nanosecond argument should be between 0 and 999999999: 1000000000000000000/

)

// Common user errors with dates
var (
	// June has only 30 days, so there is no June 31st
	_ = time.Date(2023, 6, 31, 3, 4, 5, 0, time.UTC) // MATCH /time.Date day argument is 31, but June has only 30 days/

	// there is no January 32nd
	_ = time.Date(2024, 1, 32, 3, 4, 5, 0, time.UTC) // MATCH /time.Date day argument is 32, but January has only 31 days/

	// there is no February 29th in 2023 (not a leap year)
	_ = time.Date(2023, 2, 29, 3, 4, 5, 0, time.UTC) // MATCH /time.Date day argument is 29, but February 2023 has only 28 days/

	// there is no February 30th
	_ = time.Date(2024, 2, 30, 3, 4, 5, 0, time.UTC) // MATCH /time.Date day argument is 30, but February 2024 has only 29 days/
)

// day and month arguments are swapped
var (
	_ = time.Date(2023, 30, 6, 3, 4, 5, 0, time.UTC)
	// MATCH:59 /time.Date month argument should be between 1 and 12: 30/
	// MATCH:59 /time.Date month and day arguments appear to be swapped: 2023-06-30 vs 2023-30-06/

	// month and day arguments are swapped, but there is no June 31st
	_ = time.Date(2023, 31, 6, 3, 4, 5, 0, time.UTC) // MATCH /time.Date month argument should be between 1 and 12: 31/
)

// edge cases to validate the order of the checks
var (

	// here the month argument could have been swapped with the day argument
	// but the day argument is zero
	_ = time.Date(2023, 31, 0, 3, 4, 5, 0, time.UTC)
	// MATCH:72 /time.Date day argument should not be zero/
	// MATCH:72 /time.Date month argument should be between 1 and 12: 31/

	// here the month argument could have been swapped with the day argument
	// but the day argument is negative
	_ = time.Date(2023, 31, -1, 3, 4, 5, 0, time.UTC)
	// MATCH:78 /time.Date day argument is negative: -1/
	// MATCH:78 /time.Date month argument should be between 1 and 12: 31/
)

// arguments are totally out of bounds
var (
	_ = time.Date(
		2023,
		13, // MATCH /time.Date month argument should be between 1 and 12: 13/
		// here we cannot detect the number of days in the month, so it falls back to the default of 31
		32,         // MATCH /time.Date day argument should be between 1 and 31: 32/
		25,         // MATCH /time.Date hour argument should be between 0 and 23: 25/
		60,         // MATCH /time.Date minute argument should be between 0 and 59: 60/
		61,         // MATCH /time.Date second argument should be between 0 and 60: 61/
		1000000000, // MATCH /time.Date nanosecond argument should be between 0 and 999999999: 1000000000/
		time.UTC)
)

// valid time.Date calls that should not be reported by the rule
// there are provided to prevent regression.
var (
	// a date before the start of era is valid
	_ = time.Date(-500, 1, 2, 3, 4, 5, 6, time.UTC)

	// a date in the future is valid
	_ = time.Date(3000, 1, 2, 3, 4, 5, 6, time.UTC)

	// 0, 1 and, -1 are valid and distinct years in Go
	_ = time.Date(0, 1, 2, 3, 4, 5, 6, time.UTC)
	_ = time.Date(-1, 1, 2, 3, 4, 5, 6, time.UTC)
	_ = time.Date(1, 1, 2, 3, 4, 5, 6, time.UTC)

	// midnight is a valid time
	_ = time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	// 2023 is not a leap year, so February has only 28 days
	_ = time.Date(2023, 2, 28, 3, 4, 5, 0, time.UTC)

	// 2024 is a leap year, so February has 29 days
	_ = time.Date(2024, 2, 29, 3, 4, 5, 0, time.UTC)

	// 2020 is not a leap year, so February has only 28 days
	_ = time.Date(2020, 2, 28, 3, 4, 5, 0, time.UTC)

	// a leap second is valid
	_ = time.Date(2016, 12, 31, 23, 59, 60, 0, time.UTC)
)

// these should not be reported by the rule
var (
	a int
	_ = time.Date(2023, 1, a, 3, 4, 5, 6, time.UTC)
	_ = time.Date(2023, 1, -a, 3, 4, 5, 6, time.UTC)
	_ = time.Date(2023, 1, +a, 3, 4, 5, 6, time.UTC)
	_ = time.Date(2023, 1, ^1, 3, 4, 5, 6, time.UTC)
	_ = time.Date(2023, 1, ^a, 3, 4, 5, 6, time.UTC)

	// obscure, but valid, notations that are ignored
	_ = time.Date(
		+-2023,
		-+1,
		+-+2,
		-+-3,
		+-+-4,
		-+-+5,
		+-+-+6,
		time.UTC)
)

// values that do not fit in an int (math.MaxInt + 1) are ignored
var (
	_ = time.Date(9223372036854775808, 2, 29, 3, 4, 5, 0, time.UTC)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 9223372036854775808, time.UTC)
)
