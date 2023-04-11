package wrapper

import (
	"log"
	"testing"
)

func mockRecovery() {
	if err := recover(); err != nil {
		log.Printf("exec recover:%v\n", err)
	}
}

func TestWrapper(t *testing.T) {
	wg := New(ChWrapper, WithRecover(mockRecovery))
	wg.Wrap(func() {
		log.Println("1111")
	})

	num := 10 * 100
	for i := 0; i < num; i++ {
		// The method of copying is used here to avoid the i
		// in the wrap func being the same variable
		index := i
		wg.Wrap(func() {
			log.Printf("current index: %d\n", index)
		})
	}

	wg.WrapWithRecover(func() {
		panic("mock panic test")
	})

	wg.Wait()
}

/*
$ go test -v
2023/04/11 17:36:00 current index: 708
2023/04/11 17:36:00 current index: 703
2023/04/11 17:36:00 current index: 684
2023/04/11 17:36:00 current index: 711
2023/04/11 17:36:00 current index: 459
2023/04/11 17:36:00 current index: 718
--- PASS: TestWrapper (4.06s)
PASS
*/
