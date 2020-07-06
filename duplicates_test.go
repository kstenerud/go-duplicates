package duplicates

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func mapDifference(a, b map[TypedPointer]bool) (difference []TypedPointer) {
	for k, _ := range a {
		if _, ok := b[k]; !ok {
			difference = append(difference, k)
		}
	}
	return
}

func assertDuplicates(t *testing.T, value interface{}, expectedDuplicates ...interface{}) {
	expected := make(map[TypedPointer]bool)
	for _, dup := range expectedDuplicates {
		expected[TypedPointerOf(dup)] = true
	}

	actual := make(map[TypedPointer]bool)
	for dup, found := range FindDuplicatePointers(value) {
		if found {
			actual[dup] = true
		}
	}

	actualMissing := mapDifference(expected, actual)
	if len(actualMissing) > 0 {
		t.Errorf("Expected but not found: %v", actualMissing)
	}
	expectedMissing := mapDifference(actual, expected)
	if len(expectedMissing) > 0 {
		t.Errorf("Found but unexpected: %v", expectedMissing)
	}

	finder := NewDuplicateFinder()
	finder.ScanForPointers(value)
	for _, duplicate := range expectedDuplicates {
		if !finder.IsDuplicatePointer(duplicate) {
			t.Errorf("Expected duplicate %v", duplicate)
		}
	}
}

func TestZeroValue(t *testing.T) {
	FindDuplicatePointers(nil)
}

func TestDuplicatesBasicTypes(t *testing.T) {
	var nullIntf interface{}
	var nullPtr *int
	var nullFunc func()
	var aChan chan int
	assertDuplicates(t, nil)
	assertDuplicates(t, nullIntf)
	assertDuplicates(t, nullPtr)
	assertDuplicates(t, nullFunc)
	assertDuplicates(t, aChan)
	assertDuplicates(t, func() {})
	assertDuplicates(t, true)
	assertDuplicates(t, int(1))
	assertDuplicates(t, int8(1))
	assertDuplicates(t, int16(1))
	assertDuplicates(t, int32(1))
	assertDuplicates(t, int64(1))
	assertDuplicates(t, uint(1))
	assertDuplicates(t, uint8(1))
	assertDuplicates(t, uint16(1))
	assertDuplicates(t, uint32(1))
	assertDuplicates(t, uint64(1))
	assertDuplicates(t, uintptr(1))
	assertDuplicates(t, float32(1))
	assertDuplicates(t, float64(1))
	assertDuplicates(t, complex64(complex(1, 1)))
	assertDuplicates(t, complex128(complex(1, 1)))
	assertDuplicates(t, "")
}

func TestDuplicatesSlice(t *testing.T) {
	v1 := 1
	v2 := 2
	slice := []*int{&v1, &v2, &v1}
	assertDuplicates(t, slice, &v1)
}

func TestDuplicatesArray(t *testing.T) {
	v1 := 1
	v2 := 2
	array := [3]*int{&v1, &v2, &v1}
	assertDuplicates(t, array, &v1)
}

func TestDuplicatesMap(t *testing.T) {
	v1 := "1"
	v2 := "2"
	v3 := "3"
	m := map[int]*string{1: &v1, 2: &v2, 3: &v1, 4: &v3, 5: &v3, 6: &v3}
	assertDuplicates(t, m, &v1, &v3)
}

func TestDuplicatesDeep(t *testing.T) {
	v1 := []byte{1, 1}
	v2 := []byte{2, 2}
	v3 := []byte{1, 1}

	m := map[interface{}]interface{}{
		1: &v1,
		2: &v3,
		"x": []interface{}{
			1,
			&v2,
			map[int]*[]byte{
				10: &v1,
				20: &v2,
			},
		},
	}
	assertDuplicates(t, m, &v2, &v1)
}

func TestDuplicatesRecursive(t *testing.T) {
	m := map[interface{}]interface{}{}
	m[1] = m
	assertDuplicates(t, m, m)
}

type DuplicatesTestStruct struct {
	A int
	B *DuplicatesTestStruct
	C []DuplicatesTestStruct
	D []interface{}
	e *DuplicatesTestStruct

	Bo    bool
	By    byte
	I     int
	I8    int8
	I16   int16
	I32   int32
	I64   int64
	U     uint
	U8    uint8
	U16   uint16
	U32   uint32
	U64   uint64
	F32   float32
	F64   float64
	Ar    [4]byte
	St    string
	Ba    []byte
	Sl    []interface{}
	M     map[interface{}]interface{}
	Pi    *int
	Time  time.Time
	URL   url.URL
	PTime *time.Time
	PURL  *url.URL
}

func TestDuplicatesStruct1(t *testing.T) {
	v := DuplicatesTestStruct{}
	v.A = 1
	v.B = &v

	assertDuplicates(t, &v, &v)
}

func TestDuplicatesStruct2(t *testing.T) {
	v := DuplicatesTestStruct{}
	v.C = []DuplicatesTestStruct{v}

	assertDuplicates(t, &v)
}

func TestDuplicatesStruct3(t *testing.T) {
	v := DuplicatesTestStruct{}
	v.D = []interface{}{v}

	assertDuplicates(t, &v)
}

func TestDuplicatesStruct4(t *testing.T) {
	v := DuplicatesTestStruct{}
	v.D = []interface{}{&v}

	assertDuplicates(t, &v, &v)
}

func TestDuplicatesStruct5(t *testing.T) {
	v := DuplicatesTestStruct{}
	v.e = &v

	assertDuplicates(t, &v, &v)
}

type SomeStruct struct {
	Name         string
	NameAlias    *string
	recursive    *SomeStruct
	RandomValues []interface{}
}

func describeDuplicates(dups map[TypedPointer]bool) string {
	str := strings.Builder{}
	isFirst := true
	for ptr, isDup := range dups {
		if isDup {
			if !isFirst {
				str.WriteString(", ")
			}
			isFirst = false
			str.WriteString(fmt.Sprintf("%v", ptr))
		}
	}
	return str.String()
}

func Demonstrate() {
	v := &SomeStruct{Name: "My name"}
	fmt.Printf("No duplicates: %v\n",
		describeDuplicates(FindDuplicatePointers(v)))

	// -----------------------------------

	v.NameAlias = &v.Name
	fmt.Printf("NameAlias points to Name: %v\n",
		describeDuplicates(FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.NameAlias = &v.Name
	v.recursive = v
	fmt.Printf("NameAlias points to Name and recursive points to self: %v\n",
		describeDuplicates(FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.RandomValues = append(v.RandomValues, v)
	fmt.Printf("RandomValues contains pointer to self: %v\n",
		describeDuplicates(FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.RandomValues = append(v.RandomValues, &v.NameAlias)
	fmt.Printf("RandomValues contains pointer to NameAlias: %v\n",
		describeDuplicates(FindDuplicatePointers(v)))
}

func TestDemonstrate(t *testing.T) {
	Demonstrate()
}
