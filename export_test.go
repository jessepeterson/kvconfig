package kvconfig

import "testing"

func TestSimpleStruct1(t *testing.T) {
	type TestStruct struct {
		TestString string `config:"test_string"`
		TestInt    int    `config:"test_int"`
	}

	kv := NewMap()

	ts := TestStruct{
		TestString: "test",
		TestInt:    1,
	}

	Export(&ts, kv)

	testTable := map[string]string{
		"test_string_0": "test",
		"test_int_0":    "1",
	}

	for k, v := range testTable {
		av := kv.Get(k)
		if av != v {
			t.Errorf("kv.Get(%q) = %q; wanted %q", k, av, v)
		}
	}
}
