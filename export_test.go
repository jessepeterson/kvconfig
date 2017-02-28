package kvconfig

import "testing"

func TestSimpleStructExport1(t *testing.T) {
	type TestStruct struct {
		TestString    string  `kvconfig:"test_string"`
		TestInt       int     `kvconfig:"test_int"`
		TestPtrString *string `kvconfig:"test_ptr_string"`
		TestPtrInt    *int    `kvconfig:"test_ptr_int"`
	}

	testStr := "test2"
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
		"test_ptr_string_0": "test2",
		"test_ptr_int_0":    "2",
	}

	for k, v := range testTable {
		if _, ok := kv.Lookup(k); !ok {
			t.Errorf("kv.Exists(%q) != true", k)
		} else if av := kv.Get(k); av != v {
			t.Errorf("kv.Get(%q) = %q; wanted %q", k, av, v)
		}
	}
}
