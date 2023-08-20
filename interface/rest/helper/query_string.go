package helper

import (
	"strconv"
)

func ParseNumberParameter(parameter string) (uint64, error) {
	parsedParameter, err := strconv.ParseUint(parameter, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsedParameter, nil
}
