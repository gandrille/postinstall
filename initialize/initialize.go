package initialize

import (
	"fmt"

	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/initialize/functions"
)

// Run executes all the functions
func Run() {
	list := getInitFunctions()
	var results []result.Result
	for i, element := range list {
		result.PrintInfo(element.Infos().Title)
		result := element.Run()
		results = append(results, result)
		result.Print()
		if i != len(list)-1 {
			fmt.Printf("\n")
		}
	}

	result.PrintInfo("\nSummary\n")
	result.NewSet(results, "Initialization process").Print()
}

// Describe function
func Describe() {
	list := getInitFunctions()
	for i, element := range list {
		element.Infos().Describe()
		if i != len(list)-1 {
			fmt.Printf("\n")
		}
	}
}

func getInitFunctions() []functions.Function {
	var list []functions.Function

	list = append(list, functions.FreedesktopFunction{})
	list = append(list, functions.XfceDesktopFunction{})
	list = append(list, functions.GnomeDesktopFunction{})
	list = append(list, functions.BashrcFunction{})
	list = append(list, functions.ZimFunction{})
	list = append(list, functions.MvnSettingsFunction{})
	list = append(list, functions.BluemanFunction{})
	list = append(list, functions.SdkManFunction{})
	list = append(list, functions.SSHFunction{})
	list = append(list, functions.FirefoxFunction{})

	return list
}
