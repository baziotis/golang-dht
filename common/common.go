package common

import (
	"strings"
)

const (
	STARTING_PORT = 12000
	NSERVERS      = 4
)

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(cond bool) {
	if cond == false {
		panic("Assertion failed")
	}
}

func GetComponents(buffer string) []string {
	components := strings.Split(string(buffer), ":")
	Assert(len(components) >= 1)
	return components
}
