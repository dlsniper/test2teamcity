package stdlib

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dlsniper/test2teamcity/teamcity"
)

type Message struct {
	Time                                  time.Time
	Elapsed                               float64
	Action, Package, Output, Status, Test string
}

var testPrefixes = []string{"--- PASS: ", "--- FAIL: ", "--- SKIP: ", "--- BENCH: "}

func removePrefixes(msg string) string {
	result := msg
	for _, v := range testPrefixes {
		result = strings.Replace(result, v, "", 1)
	}
	return result
}

func ProcessStdLib(line string, w io.Writer) {
	test := &Message{}
	err := json.Unmarshal([]byte(line), test)
	if err != nil {
		return
	}

	testTime := test.Time.Format(teamcity.TeamCityTs)
	testName := test.Package
	if test.Test != "" {
		testName += "." + test.Test
	}
	testName = teamcity.Escape(testName)
	testOutput := removePrefixes(teamcity.Escape(test.Output))

	matchers := []struct {
		prefix string
		action string
		output func()
	}{
		{
			prefix: "",
			action: "run",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testStarted flowId='%s' timestamp='%s' name='%s' captureStandardOutput='false']\n", testName, testTime, testName)
			},
		},
		{
			prefix: "",
			action: "skip",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testIgnored flowId='%s' timestamp='%s' name='%s']\n", testName, testTime, testName)
			},
		},
		{
			prefix: "",
			action: "pass",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testFinished flowId='%s' timestamp='%s' name='%s' duration='%d']\n", testName, testTime, testName, (time.Duration(test.Elapsed)*time.Second)/time.Millisecond)
			},
		},
		{
			prefix: "",
			action: "fail",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testFailed flowId='%s' timestamp='%s' name='%s' details='%s']\n", testName, testTime, testName, testOutput)
			},
		},
		{
			prefix: "    ",
			action: "output",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testStdOut flowId='%s' timestamp='%s' name='%s' out='%s']\n", testName, testTime, testName, strings.Replace(testOutput, "    ", "", 1))
			},
		},
		{
			prefix: "=== RUN   ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "=== PAUSE ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "=== CONT  ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "--- PASS: ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "--- FAIL: ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "--- SKIP: ",
			action: "output",
			output: func() {},
		},
		{
			prefix: "FAIL",
			action: "output",
			output: func() {},
		},
		{
			prefix: "?   \t",
			action: "output",
			output: func() {},
		},
		{
			prefix: "testing: warning: no tests to run",
			action: "output",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testStdErr flowId='%s' timestamp='%s' name='%s' out='%s']\n", testName, testTime, testName, testOutput)
			},
		},
		{
			prefix: "exit status ",
			action: "output",
			output: func() {
				fmt.Fprintf(w, "##teamcity[testStdErr flowId='%s' timestamp='%s' name='%s' out='%s']\n", testName, testTime, testName, testOutput)
			},
		},
	}
	for _, m := range matchers {
		if m.action != "" && m.action == test.Action && m.prefix != "" && strings.HasPrefix(testOutput, m.prefix) {
			m.output()
			break
		} else if m.action != "" && m.action == test.Action && m.prefix == "" {
			m.output()
			break
		}
		if m.prefix != "" && strings.HasPrefix(testOutput, m.prefix) {
			m.output()
			break
		}
	}
}

