package stdlib_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/dlsniper/test2teamcity/stdlib"
)

func TestCleanOutput(t *testing.T) {
	lines := []struct {
		in, out string
	}{
		// 1
		{
			in:  `{"Time":"2018-07-07T22:58:04.6231601+03:00","Action":"run","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoPassFunc"}`,
			out: "##teamcity[testStarted flowId='github.com/dlsniper/test2teamcity.TestDemoPassFunc' timestamp='2018-07-07T22:58:04.623' name='github.com/dlsniper/test2teamcity.TestDemoPassFunc' captureStandardOutput='false']\n",
		},
		// 2
		{
			in:  `{"Time":"2018-07-07T22:58:04.6241615+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoPassFunc","Output":"=== RUN   TestDemoPassFunc\n"}`,
			out: "",
		},
		// 3
		{
			in:  `{"Time":"2018-07-07T22:58:04.6241615+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoPassFunc","Output":"--- PASS: TestDemoPassFunc (0.00s)\n"}`,
			out: "",
		},
		// 4
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoPassFunc","Output":"    main_test.go:8: pass output here\n"}`,
			out: "##teamcity[testStdOut flowId='github.com/dlsniper/test2teamcity.TestDemoPassFunc' timestamp='2018-07-07T22:58:04.625' name='github.com/dlsniper/test2teamcity.TestDemoPassFunc' out='main_test.go:8: pass output here|n']\n",
		},
		// 5
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"pass","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoPassFunc","Elapsed":0}`,
			out: "##teamcity[testFinished flowId='github.com/dlsniper/test2teamcity.TestDemoPassFunc' timestamp='2018-07-07T22:58:04.625' name='github.com/dlsniper/test2teamcity.TestDemoPassFunc' duration='0']\n",
		},
		// 6
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"run","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoSkipFunc"}`,
			out: "##teamcity[testStarted flowId='github.com/dlsniper/test2teamcity.TestDemoSkipFunc' timestamp='2018-07-07T22:58:04.625' name='github.com/dlsniper/test2teamcity.TestDemoSkipFunc' captureStandardOutput='false']\n",
		},
		// 7
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoSkipFunc","Output":"=== RUN   TestDemoSkipFunc\n"}`,
			out: "",
		},
		// 8
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoSkipFunc","Output":"--- SKIP: TestDemoSkipFunc (0.00s)\n"}`,
			out: "",
		},
		// 9
		{
			in:  `{"Time":"2018-07-07T22:58:04.6251648+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoSkipFunc","Output":"    main_test.go:12: skip output here\n"}`,
			out: "##teamcity[testStdOut flowId='github.com/dlsniper/test2teamcity.TestDemoSkipFunc' timestamp='2018-07-07T22:58:04.625' name='github.com/dlsniper/test2teamcity.TestDemoSkipFunc' out='main_test.go:12: skip output here|n']\n",
		},
		// 10
		{
			in:  `{"Time":"2018-07-07T22:58:04.6261602+03:00","Action":"skip","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoSkipFunc","Elapsed":0}`,
			out: "##teamcity[testIgnored flowId='github.com/dlsniper/test2teamcity.TestDemoSkipFunc' timestamp='2018-07-07T22:58:04.626' name='github.com/dlsniper/test2teamcity.TestDemoSkipFunc']\n",
		},
		// 11
		{
			in:  `{"Time":"2018-07-07T22:58:04.6261602+03:00","Action":"run","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoFailFunc"}`,
			out: "##teamcity[testStarted flowId='github.com/dlsniper/test2teamcity.TestDemoFailFunc' timestamp='2018-07-07T22:58:04.626' name='github.com/dlsniper/test2teamcity.TestDemoFailFunc' captureStandardOutput='false']\n",
		},
		// 12
		{
			in:  `{"Time":"2018-07-07T22:58:04.6281609+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoFailFunc","Output":"=== RUN   TestDemoFailFunc\n"}`,
			out: "",
		},
		// 13
		{
			in:  `{"Time":"2018-07-07T22:58:04.6281609+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoFailFunc","Output":"--- FAIL: TestDemoFailFunc (0.00s)\n"}`,
			out: "",
		},
		// 14
		{
			in:  `{"Time":"2018-07-07T22:58:04.6291734+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoFailFunc","Output":"    main_test.go:16: fail output here\n"}`,
			out: "##teamcity[testStdOut flowId='github.com/dlsniper/test2teamcity.TestDemoFailFunc' timestamp='2018-07-07T22:58:04.629' name='github.com/dlsniper/test2teamcity.TestDemoFailFunc' out='main_test.go:16: fail output here|n']\n",
		},
		// 15
		{
			in:  `{"Time":"2018-07-07T22:58:04.6301596+03:00","Action":"fail","Package":"github.com/dlsniper/test2teamcity","Test":"TestDemoFailFunc","Elapsed":0}`,
			out: "##teamcity[testFailed flowId='github.com/dlsniper/test2teamcity.TestDemoFailFunc' timestamp='2018-07-07T22:58:04.630' name='github.com/dlsniper/test2teamcity.TestDemoFailFunc' details='']\n",
		},
		// 16
		{
			in:  `{"Time":"2018-07-07T22:58:04.6301596+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Output":"FAIL\n"}`,
			out: "",
		},
		// 17
		{
			in:  `{"Time":"2018-07-07T22:58:04.631161+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Output":"exit status 1\n"}`,
			out: "##teamcity[testStdErr flowId='github.com/dlsniper/test2teamcity' timestamp='2018-07-07T22:58:04.631' name='github.com/dlsniper/test2teamcity' out='exit status 1|n']\n",
		},
		// 18
		{
			in:  `{"Time":"2018-07-07T22:58:04.6321606+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Output":"FAIL\tgithub.com/dlsniper/test2teamcity\t0.079s\n"}`,
			out: "",
		},
		// 19
		{
			in:  `{"Time":"2018-07-07T22:58:04.6341617+03:00","Action":"fail","Package":"github.com/dlsniper/test2teamcity","Elapsed":0.082}`,
			out: "##teamcity[testFailed flowId='github.com/dlsniper/test2teamcity' timestamp='2018-07-07T22:58:04.634' name='github.com/dlsniper/test2teamcity' details='']\n",
		},
		// 20
		{
			in: `{"Time":"2018-07-08T15:19:01.0255244+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity","Output":"testing: warning: no tests to run\n"}`,
			out: "##teamcity[testStdErr flowId='github.com/dlsniper/test2teamcity' timestamp='2018-07-08T15:19:01.025' name='github.com/dlsniper/test2teamcity' out='testing: warning: no tests to run|n']\n",
		},
		// 21
		{
			in: `{"Time":"2018-07-08T16:05:20.9049863+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity/stdlib","Test":"TestCleanOutput/sub20","Output":"    --- PASS: TestCleanOutput/sub20 (0.00s)\n"}`,
			out: "",
		},
		// 22
		{
			in: `{"Time":"2018-07-08T16:05:20.9069869+03:00","Action":"output","Package":"github.com/dlsniper/test2teamcity/stdlib","Test":"TestCleanOutput/sub20","Output":"        stdlib_test.go:120: running test sub20\n"}`,
			out: "##teamcity[testStdOut flowId='github.com/dlsniper/test2teamcity.TestCleanOutput/sub20' timestamp='2018-07-08T16:05:20.906' name='github.com/dlsniper/test2teamcity.TestCleanOutput/sub20' out='stdlib_test.go:120: running test sub20|n']\n",
		},
	}

	for idx := range lines {
		if idx+1 != 22 {
			continue
		}
		testName := "sub" + strconv.Itoa(idx+1)
		t.Run(testName, func(t *testing.T) {
			t.Logf("running test %s\n", testName)
			line := lines[idx]
			out := bytes.NewBuffer(make([]byte, 0, 4096))
			stdlib.ProcessStdLib(line.in, out)

			result := out.String()
			if line.out != result {
				t.Logf("got : %s\n", result)
				t.Logf("want: %s\n", line.out)
				t.Fail()
			}
		})
	}
}
