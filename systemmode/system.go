package systemmode

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/systemmode/processors"
)

// Run executes the system installer
func Run() {
	// gets user choices
	lines, err := askLines()
	if err != nil {
		result.PrintError(err.Error())
		os.Exit(1)
	}

	// run
	t0 := time.Now()
	results := runAll(lines)
	duration := time.Since(t0).Truncate(time.Second).String()

	// summary
	result.PrintInfo("\nSummary")
	results.SetMessage(results.Message() + " (" + duration + ")")
	results.Print()
}

func runAll(lines []string) result.Set {
	var results result.Set
	tot := strconv.Itoa(len(lines))

	for i, line := range lines {
		if len(line) != 0 {
			// preparation and message
			tokens := strings.Split(line, " ")
			cmd := processors.GetProcessor(tokens[0])
			args := tokens[1:]
			fmt.Println("")
			result.PrintInfo(strconv.Itoa(i+1) + "/" + tot + " " + cmd.Describe(args))

			// run
			runner := func() result.Result { return cmd.Run(args) }
			result := result.Run(runner)

			// register
			results.Add(result)
			result.Print()
		}
	}

	return results
}

// Describe prints what the system installer does
func Describe() {
	describeAll()
}
