package validation

import (
	"fmt"
	"time"
)

type dateTimeValidator struct {
	name string
}

var (
	validLayouts = []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05 -07:00",
	}
)

func NewDateTimeValidator(name string) Validator {
	return &dateTimeValidator{
		name: name,
	}
}

func (x *dateTimeValidator) Validate(args map[string][]string) (err error, res interface{}) {
	values, ok := args[x.name]

	if !ok {
		return
	}

	if len(values) > 1 {
		return fmt.Errorf("%s is not a valid datetime", x.name), nil
	}

	if res = ParseDateTime(values[0]); res == nil {
		return fmt.Errorf("%s is not a valid datetime", x.name), nil
	}

	return
}

func isValidDateTime(candidate string) bool {
	return ParseDateTime(candidate) != nil
}

func ParseDateTime(candidate interface{}) *time.Time {
	res, isTime := candidate.(time.Time)

	if isTime {
		return &res
	}

	str, isString := candidate.(string)

	if !isString {
		return nil
	}

	for _, layout := range validLayouts {
		dt, err := time.Parse(layout, str)
		if err == nil {
			return &dt
		}
	}

	return nil
}

type midnightValidator struct {
	name string
}

func NewMidnightValidator(name string) Validator {
	return &midnightValidator{name: name}
}

func (x *midnightValidator) Validate(args map[string][]string) (err error, res interface{}) {
	values, ok := args[x.name]

	if !ok {
		return
	}

	if len(values) > 1 {
		return fmt.Errorf("%s is not a valid datetime", x.name), nil
	}

	periodErr, period := NewPeriodValidator().Validate(args)

	if theDate := ParseDateTime(values[0]); theDate != nil &&
		periodErr == nil &&
		(period != nil && period != "hour") {

		if !isMidnight(theDate.UTC()) {
			return fmt.Errorf("%s must be midnight", x.name), nil
		}
	}

	return
}

func isMidnight(t time.Time) bool {
	hour, min, sec := t.Clock()
	return (hour == 0 && min == 0 && sec == 0)
}

type timespanValidator struct {
	length int
}

func NewTimespanValidator(length int) Validator {
	return &timespanValidator{length: length}
}

func (x *timespanValidator) Validate(args map[string][]string) (err error, res interface{}) {
	_, startAt := NewDateTimeValidator("start_at").Validate(args)
	_, endAt := NewDateTimeValidator("end_at").Validate(args)
	_, period := NewPeriodValidator().Validate(args)

	if startAt != nil && endAt != nil && (period != nil && period != "hour") {
		hours := endAt.(*time.Time).UTC().Sub(startAt.(*time.Time).UTC()).Hours()
		if hours < float64(24*7) {
			return fmt.Errorf("The minimum timespan for a query is %v days", x.length), nil
		}
		res = hours / 24
	}

	return
}

type mondayValidator struct {
	name string
}

func NewMondayValidator(name string) Validator {
	return &mondayValidator{name: name}
}

func (x *mondayValidator) Validate(args map[string][]string) (err error, res interface{}) {
	_, date := NewDateTimeValidator(x.name).Validate(args)
	_, period := NewPeriodValidator().Validate(args)

	if (period != nil && period == "week") &&
		date != nil &&
		date.(*time.Time).UTC().Weekday() != time.Monday {
		return fmt.Errorf("%v must be a Monday but was %v", x.name, date), nil
	}

	return
}

type monthValidator struct {
	name string
}

func NewMonthValidator(name string) Validator {
	return &monthValidator{name: name}
}

func (x *monthValidator) Validate(args map[string][]string) (err error, res interface{}) {
	_, date := NewDateTimeValidator(x.name).Validate(args)
	_, period := NewPeriodValidator().Validate(args)

	if (period != nil && period == "month") &&
		date != nil &&
		date.(*time.Time).UTC().Day() != 1 {
		return fmt.Errorf("%v must be a first of the month but was %v", x.name, date), nil
	}

	return
}
