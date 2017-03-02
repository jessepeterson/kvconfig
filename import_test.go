package kvconfig

import "testing"

func TestSimpleStructImport1(t *testing.T) {
	type TestSubStruct struct {
		TestSubString string `kvconfig:"test_sub_string"`
		TestSubInt    int    `kvconfig:"test_sub_int"`
	}

	type TestStruct struct {
		TestString    string  `kvconfig:"test_string"`
		TestInt       int     `kvconfig:"test_int"`
		TestPtrString *string `kvconfig:"test_ptr_string"`
		TestPtrInt    *int    `kvconfig:"test_ptr_int"`
		SubStructs    []*TestSubStruct
	}

	ts := TestStruct{}

	kv := &MapStrStr{
		"test_string_0":     "test",
		"test_int_0":        "1",
		"test_ptr_string_0": "test2",
		"test_ptr_int_0":    "2",
		"test_sub_string_0": "test3",
		"test_sub_int_0":    "3",
	}

	Import(kv, &ts)

	if ts.TestString != "test" {
		t.Errorf("TestStruct.TestString = %q; wanted %q", ts.TestString, "test")
	}

	if ts.TestInt != 1 {
		t.Errorf("TestStruct.TestInt = %d; wanted %d", ts.TestInt, 1)
	}

	if ts.TestPtrString == nil {
		t.Error("TestStruct.TestPtrString == nil")
	} else if *(ts.TestPtrString) != "test2" {
		t.Errorf("*(TestStruct.TestPtrString) = %q; wanted %q", *(ts.TestPtrString), "test2")
	}

	if ts.TestPtrInt == nil {
		t.Error("TestStruct.TestPtrInt == nil")
	} else if *(ts.TestPtrInt) != 2 {
		t.Errorf("*(TestStruct.TestPtrInt) = %q; wanted %q", *(ts.TestPtrInt), 2)
	}

	if ts.SubStructs == nil {
		t.Error("TestStruct.SubStructs == nil")
	} else if len(ts.SubStructs) < 1 {
		t.Error("len(TestStruct.SubStructs) < 1")
	} else if ts.SubStructs[0] == nil {
		t.Error("TestStruct.SubStructs[0] == nil")
	} else {
		if ts.SubStructs[0].TestSubInt != 3 {
			t.Errorf("TestStruct.SubStructs[0].TestSubInt) = %d; wanted %d", ts.SubStructs[0].TestSubInt, 3)
		}
		if ts.SubStructs[0].TestSubString != "test3" {
			t.Errorf("TestStruct.SubStructs[0].TestSubString) = %q; wanted %q", ts.SubStructs[0].TestSubString, "test3")
		}
	}

}
