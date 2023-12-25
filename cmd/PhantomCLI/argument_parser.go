package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrMissingRequired = errors.New("missing required argument")

type ArgumentParser struct {
	parsed map[string]string

	lastParseName    string
	lastParseSuccess bool
}

func ParseArgs(args []string) *ArgumentParser {
	var a = &ArgumentParser{
		parsed: make(map[string]string),
	}

	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		a.parsed[arg[1:]] = args[i+1]
	}

	return a
}

func (a *ArgumentParser) GetParam(name, short string) (result string) {
	a.lastParseName = name
	result, a.lastParseSuccess = a.parsed[short]
	if !a.lastParseSuccess {
		result, a.lastParseSuccess = a.parsed[name]
	}
	return
}

func (a *ArgumentParser) Required() *ArgumentParser {
	if !a.lastParseSuccess {
		panic(fmt.Errorf("%w : %s", ErrMissingRequired, a.lastParseName))
	}
	return a
}

func (a *ArgumentParser) String(name, short string, result *string) *ArgumentParser {
	*result = a.GetParam(name, short)
	return a
}

func (a *ArgumentParser) Int(name, short string, result *int) *ArgumentParser {

	var res = a.GetParam(name, short)
	if res != "" {
		var err error
		*result, err = strconv.Atoi(res)
		if err != nil {
			panic(err)
		}
	}

	return a
}
