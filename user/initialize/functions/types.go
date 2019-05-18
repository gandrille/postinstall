package functions

import (
	"fmt"

	"github.com/gandrille/go-commons/result"
)

// Function interface
type Function interface {
	Infos() FunctionInfos
	Run() result.Result
}

// FunctionInfos type
type FunctionInfos struct {
	Title            string
	ShortDescription string
	LongDescription  string
}

func (info FunctionInfos) Describe() {
	result.Describe(info.Title, info.ShortDescription)
	if info.LongDescription != "" {
		fmt.Println(info.LongDescription)
	}
}

type Executor func() result.Result

func execute(title string, executors ...Executor) result.Result {
	for _, exe := range executors {
		if res := exe(); !res.IsSuccess() {
			result.PrintRed(res.Message())
			return result.Failure(title + " " + res.Message())
		} else {
			fmt.Println(res.Message())
		}
	}
	return result.Success(title)
}
