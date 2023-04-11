package wrapper

import (
	"log"
	"testing"
)

func mockChanRecovery() {
	if err := recover(); err != nil {
		log.Printf("exec recover:%v\n", err)
	}
}

func TestWrapperChan(t *testing.T) {
	num := 10
	c := num + 2
	chWrap := New(ChWrapper, WithBufCap(c), WithRecover(mockChanRecovery))
	chWrap.Wrap(func() {
		log.Println("1111")
	})

	for i := 0; i < num; i++ {
		// The method of copying is used here to avoid the i
		// in the wrap func being the same variable
		index := i
		chWrap.Wrap(func() {
			log.Printf("current index: %d\n", index)
		})
	}

	chWrap.WrapWithRecover(func() {
		log.Println(2222)
		panic("mock panic test")
	})
	chWrap.Wait()
}

/**
=== RUN   TestWrapperChan
2023/04/11 17:23:55 2222
2023/04/11 17:23:55 current index: 3
2023/04/11 17:23:55 1111
2023/04/11 17:23:55 current index: 0
2023/04/11 17:23:55 current index: 1
2023/04/11 17:23:55 current index: 2
2023/04/11 17:23:55 current index: 8
2023/04/11 17:23:55 current index: 4
2023/04/11 17:23:55 current index: 7
2023/04/11 17:23:55 current index: 9
2023/04/11 17:23:55 current index: 5
2023/04/11 17:23:55 current index: 6
2023/04/11 17:23:55 exec recover:mock panic test
--- PASS: TestWrapperChan (0.00s)
PASS
*/
