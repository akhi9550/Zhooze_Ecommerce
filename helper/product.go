package helper

import (
	"errors"
	"strconv"
)

func StringToUInt(str string) (uint, error) {
	if str == "" {
		return 0, errors.New("empty string ")
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}