package kvconfig

import "testing"

func TestSimpleStruct1(t *testing.T) {
	type TestStruct struct {
		TestString    string  `config:"test_string"`
		TestInt       int     `config:"test_int"`
		TestPtrString *string `config:"test_ptr_string"`
		TestPtrInt    *int    `config:"test_ptr_int"`
	}

	testStr := "testptrstr"
	testInt := 2

	ts := TestStruct{
		TestString:    "test",
		TestInt:       1,
		TestPtrString: &testStr,
		TestPtrInt:    &testInt,
	}

	kv := NewMap()

	Export(&ts, kv)

	testTable := map[string]string{
		"test_string_0":     "test",
		"test_int_0":        "1",
		"test_ptr_string_0": "testptrstr",
		"test_ptr_int_0":    "2",
	}

	for k, v := range testTable {
		if !kv.Exists(k) {
			t.Errorf("kv.Exists(%q) != true", k)
		} else if av := kv.Get(k); av != v {
			t.Errorf("kv.Get(%q) = %q; wanted %q", k, av, v)
		}
	}
}
