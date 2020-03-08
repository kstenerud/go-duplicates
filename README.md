Go-Duplicates
=============

The duplicates module examines an abitrary object and reports any duplicate
pointers it finds (where more than one pointer is pointing to the same object).


Usage
-----

```golang
type SomeStruct struct {
	Name         string
	NameAlias    *string
	recursive    *SomeStruct
	RandomValues []interface{}
}

func describeDuplicates(dups map[duplicates.TypedPointer]bool) string {
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
		describeDuplicates(duplicates.FindDuplicatePointers(v)))

	// -----------------------------------

	v.NameAlias = &v.Name
	fmt.Printf("NameAlias points to Name: %v\n",
		describeDuplicates(duplicates.FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.NameAlias = &v.Name
	v.recursive = v
	fmt.Printf("NameAlias points to Name and recursive points to self: %v\n",
		describeDuplicates(duplicates.FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.RandomValues = append(v.RandomValues, v)
	fmt.Printf("RandomValues contains pointer to self: %v\n",
		describeDuplicates(duplicates.FindDuplicatePointers(v)))

	// -----------------------------------

	v = &SomeStruct{Name: "My name"}
	v.RandomValues = append(v.RandomValues, &v.NameAlias)
	fmt.Printf("RandomValues contains pointer to NameAlias: %v\n",
		describeDuplicates(duplicates.FindDuplicatePointers(v)))
}
```

#### Output:

```
No duplicates: 
NameAlias points to Name: {*string 824634402880}
NameAlias points to Name and recursive points to self: {*duplicates.SomeStruct 824634403008}, {*string 824634403008}
RandomValues contains pointer to self: {*duplicates.SomeStruct 824634403072}
RandomValues contains pointer to NameAlias: {**string 824634403152}
```


License
-------

MIT License:

Copyright 2020 Karl Stenerud

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
