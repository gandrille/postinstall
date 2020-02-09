package functions

import (
	"fmt"
	"strings"

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
	onlySkip := true
	for _, exe := range executors {

		res := exe()

		msg := res.Message()
		if strings.Contains(msg, "\n") {
			msg = "\n" + msg
			msg = strings.ReplaceAll(msg, "\n", "\n* ")
		}

		if !res.IsSuccess() {
			result.PrintRed(msg)
			return result.NewError(title + " " + msg)
		} else {
			if !res.IsUnchanged() {
				onlySkip = false
			}
			fmt.Println("[" + strings.ToUpper(res.Status().String()) + "] " + msg)
		}
	}

	if onlySkip {
		return result.NewUnchanged(title)
	} else {
		return result.NewUpdated(title)
	}
}
