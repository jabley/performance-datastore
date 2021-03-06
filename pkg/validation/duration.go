package validation

import (
	"fmt"
)

type durationValidator struct{}

func NewDurationValidator() Validator {
	return &durationValidator{}
}

func (x *durationValidator) Validate(args map[string][]string) (err error, res interface{}) {
	values, durationOk := args["duration"]

	_, periodOk := args["period"]
	_, startAtOk := args["start_at"]
	_, endAtOk := args["end_at"]

	if durationOk && startAtOk && endAtOk {
		return fmt.Errorf(`Absolute and relative time cannot be requested at the same time - either ask for 'start_at' and 'end_at', or ask for 'start_at'/'end_at' with 'duration'`), nil
	}

	if startAtOk && !(durationOk || endAtOk) {
		return fmt.Errorf(`Use of 'start_at' requires 'end_at' or 'duration'`), nil
	}

	if endAtOk && !(durationOk || startAtOk) {
		return fmt.Errorf(`Use of 'end_at' requires 'start_at' or 'duration'`), nil
	}

	if durationOk {
		if !periodOk {
			return fmt.Errorf(`If 'duration' is requested (for relative time), 'period' is required - please add a period (like 'day', 'month' etc)`), nil
		}
		if len(values) > 1 {
			return fmt.Errorf("duration should be a single argument but received %v", len(values)), nil
		}
		if values[0] == "0" {
			return fmt.Errorf("duration must be positive"), nil
		}

	}

	if startAtOk && endAtOk {

	}

	return
}
