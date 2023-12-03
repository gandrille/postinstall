package processors

import (
	"os/exec"
	"strings"
)

func runCmdAndGetFirstColumn(cmd *exec.Cmd, captureFirstLine bool) ([]string, error) {
	// fetching lines
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")

	// building result
	results := []string{}
	for cpt, line := range lines {
		if cpt != 0 || captureFirstLine {
			line = strings.ReplaceAll(line, "\t", " ")
			name := strings.Split(line, " ")[0]
			if name != "" {
				results = append(results, name)
			}
		}
	}
	return results, nil
}
