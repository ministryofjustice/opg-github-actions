package semver

import (
	"math/big"
	"strconv"
)

var (
	defaultSemver string = "0.0.0"
)

func bigIntOrNil(str string) *big.Int {
	b, err := new(big.Int).SetString(str, 0)
	if err != true {
		return nil
	}
	return b
}

func intOrNil(str string) *uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return nil
	}
	return &i
}
