package def

import (
	"encoding/json"
	"log"
	"testing"
)

func TestDef(t *testing.T) {
	obj := EmptyObject{}
	log.Println("empty object: ", obj)

	b, err := json.Marshal(obj)
	log.Println("empty object to json: ", string(b), "err: ", err)

	arr := EmptyArray{}
	log.Println("empty array: ", arr)
	b, err = json.Marshal(arr)
	log.Println("empty array to json: ", string(b), "err: ", err)

	m := H{
		"a": 1,
		"b": "hello",
		"c": 1.234,
	}

	b, err = json.Marshal(m)
	log.Println("map to json: ", string(b), "err: ", err)
}

/**
=== RUN   TestDef
2021/05/05 18:17:54 empty object:  {}
2021/05/05 18:17:54 empty object to json:  {} err:  <nil>
2021/05/05 18:17:54 empty array:  []
2021/05/05 18:17:54 empty array to json:  [] err:  <nil>
2021/05/05 18:17:54 map to json:  {"a":1,"b":"hello","c":1.234} err:  <nil>
--- PASS: TestDef (0.00s)
PASS
*/
