package grecover

import (
	"log"
	"testing"
)

func TestCatchStack(t *testing.T) {
	t.Log("test grecover")
	// bridge logger
	LogEntry = LoggerFunc(log.Println)
	testSay()

	log.Println("ok")
}

func testSay() {
	defer CheckPanic()

	log.Println("11111")
	s := []string{
		"a", "b", "c",
	}

	// mock slice panic
	log.Println("d: ", s[3])
}

/**
=== RUN   TestCatchStack
    recover_test.go:9: test grecover
2021/05/07 21:49:05 11111
2021/05/07 21:49:05 panic error:  runtime error: index out of range [3] with length 3
2021/05/07 21:49:05 full stack:  goroutine 18 [running]:
runtime/debug.Stack(0x114fb28, 0xc00014c080, 0x2)
	/usr/local/go/src/runtime/debug/stack.go:24 +0x9f
github.com/daheige/tigago/grecover.CatchStack(...)
	/Users/heige/web/go/tigago/grecover/recover.go:38
github.com/daheige/tigago/grecover.CheckPanic()
	/Users/heige/web/go/tigago/grecover/recover.go:31 +0xd3
panic(0x1136f20, 0xc000148048)
	/usr/local/go/src/runtime/panic.go:965 +0x1b9
github.com/daheige/tigago/grecover.testSay()
	/Users/heige/web/go/tigago/grecover/recover_test.go:26 +0x92
github.com/daheige/tigago/grecover.TestCatchStack(0xc000126300)
	/Users/heige/web/go/tigago/grecover/recover_test.go:12 +0x97
testing.tRunner(0xc000126300, 0x114fac0)
	/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1238 +0x2b3

2021/05/07 21:49:05 ok
--- PASS: TestCatchStack (0.00s)
PASS
*/
