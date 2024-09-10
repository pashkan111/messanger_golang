package utils

import "time"

func ParseTimeFromString(t string) (*time.Time, error) {
	parsedTime, err := time.ParseInLocation(
		"2006-01-02 15:04:05-07", t, time.UTC,
	)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
