package mode

import (
	"fmt"
	"log"
)

type _M int

const (
	// Test 测试模式
	Test _M = 1

	// Debug 调试模式
	Debug _M = 2

	// Release 发布模式
	Release _M = 3
)

var m = Test

func Set(s string) {
	switch s {
	default:
		panic(fmt.Errorf("Invalid Mode %s", s))
	case "test", "Test", "TEST":
		m = Test
	case "debug", "Debug", "DEBUG":
		m = Debug
	case "release", "Release", "RELEASE":
		m = Release
	}
	log.Println("Mode:", s)
}

func String() string {
	switch m {
	default:
		return "test"
	case Debug:
		return "debug"
	case Release:
		return "release"
	}
}

func Val() _M {
	return m
}
