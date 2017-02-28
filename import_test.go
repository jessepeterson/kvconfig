package kvconfig

import "testing"

func TestSimpleStructImport1(t *testing.T) {
	type TestStruct struct {
		TestString    string  `kvconfig:"test_string"`
		TestInt       int     `kvconfig:"test_int"`
		TestPtrString *string `kvconfig:"test_ptr_string"`
		TestPtrInt    *int    `kvconfig:"test_ptr_int"`
	}

	ts := TestStruct{}

	kv := &MapStrStr{
		"test_string_0":     "test",
		"test_int_0":        "1",
		"test_ptr_string_0": "test2",
		"test_ptr_int_0":    "2",
	}

	Import(kv, &ts)

	if ts.TestString != "test" {
		t.Errorf("TestStruct.TestString = %q; wanted %q", ts.TestString, "test")
	}

	if ts.TestInt != 1 {
		t.Errorf("TestStruct.TestInt = %q; wanted %q", ts.TestInt, 1)
	}

	if ts.TestPtrString == nil {
		t.Error("ts.TestPtrString == nil")
	} else if *(ts.TestPtrString) != "test2" {
		t.Errorf("*(TestStruct.TestPtrString) = %q; wanted %q", *(ts.TestPtrString), "test2")
	}

	if ts.TestPtrInt == nil {
		t.Error("ts.TestPtrInt == nil")
	} else if *(ts.TestPtrInt) != 2 {
		t.Errorf("*(TestStruct.TestPtrInt) = %q; wanted %q", *(ts.TestPtrInt), 2)
	}
}
