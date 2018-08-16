package main

// Constants for format detector
const (
	firstFormat  = "first_format"
	secondFormat = "second_format"
)

const (
	firstFormatLayout  = `Jan 2, 2006 at 3:04:05pm (UTC)`
	secondFormatLayout = `2006-01-02T15:04:05Z`
)

// Regex strings for format detector
/*
const (
	firstVersionRegexp  = `(?m)^(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[ ](([1-9])|(1[0-9])|(2[0-9])|(3[0-1]))[,][ ](20[0-9][0-9])( at )([1-9]|1[0-2])[:]([0-5][1-9])[:]([0-5][1-9])(am|pm) (\(UTC\))$`
	secondVersionRegexp = `(?m)^(([0-2][0-9][0-9][0-9])-(0[1-9]|1[0-2])-(0[0-9]|[1-2][0-9]|3[0-1])T(0[1-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))(Z)$`
)
*/
