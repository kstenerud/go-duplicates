// duplicates examines an abitrary object and reports any duplicate pointers it
// finds (where more than one pointer is pointing to the same object).
package duplicates

import (
	"reflect"
)

// TypedPointer is a pointer value with an associated type. Typing is necessary
// because the first field of a struct will have the same address as the struct
// itself
type TypedPointer struct {
	Type    reflect.Type
	Pointer uintptr
}

// TypedPointerOf gets the typed pointer of an arbitrary object.
// Note: value must be addressable or else this function will panic.
func TypedPointerOf(value interface{}) TypedPointer {
	return TypedPointerOfRV(reflect.ValueOf(value))
}

func TypedPointerOfRV(rv reflect.Value) TypedPointer {
	return TypedPointer{
		Type:    rv.Type(),
		Pointer: rv.Pointer(),
	}
}

// FindDuplicatePointers walks an object and its contents looking for pointer
// values that are used multiple times, marking slices, maps, pointers, and
// struct fields that point to the same instance of an object. Both exported and
// unexported fields are examined and followed.
//
// The returned duplicatePtrs will map to true for every duplicate pointer found.
// Non-duplicates will either not be present in the map, or will map to false.
// Either way, duplicatePtrs[myTypedPtr] will return true if and only if
// myTypedPtr represents a duplicate pointer.
func FindDuplicatePointers(value interface{}) (duplicatePtrs map[TypedPointer]bool) {
	duplicatePtrs = make(map[TypedPointer]bool)
	findDuplicatePtrsInValue(reflect.ValueOf(value), duplicatePtrs)
	return
}

func findDuplicatePtrsInValue(value reflect.Value, foundPtrs map[TypedPointer]bool) {
	switch value.Kind() {
	case reflect.Interface:
		if !value.IsNil() {
			value = value.Elem()
			if isSearchableKind(value.Kind()) {
				findDuplicatePtrsInValue(value, foundPtrs)
			}
		}
	case reflect.Ptr:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				findDuplicatePtrsInValue(value.Elem(), foundPtrs)
			}
		}
	case reflect.Map:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				iter := value.MapRange()
				for iter.Next() {
					findDuplicatePtrsInValue(iter.Value(), foundPtrs)
				}
			}
		}
	case reflect.Slice:
		if !value.IsNil() && !checkPtrAlreadyFound(value, foundPtrs) {
			if isSearchableKind(value.Type().Elem().Kind()) {
				count := value.Len()
				for i := 0; i < count; i++ {
					findDuplicatePtrsInValue(value.Index(i), foundPtrs)
				}
			}
		}
	case reflect.Array:
		if isSearchableKind(value.Type().Elem().Kind()) {
			count := value.Len()
			for i := 0; i < count; i++ {
				findDuplicatePtrsInValue(value.Index(i), foundPtrs)
			}
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanAddr() {
				findDuplicatePtrsInValue(field.Addr(), foundPtrs)
			} else if isSearchableKind(field.Kind()) {
				findDuplicatePtrsInValue(field, foundPtrs)
			}
		}
	}
}

const searchableKinds uint = (uint(1) << reflect.Interface) |
	(uint(1) << reflect.Ptr) |
	(uint(1) << reflect.Slice) |
	(uint(1) << reflect.Map) |
	(uint(1) << reflect.Array) |
	(uint(1) << reflect.Struct)

func isSearchableKind(kind reflect.Kind) bool {
	return searchableKinds&(uint(1)<<kind) != 0
}

func checkPtrAlreadyFound(value reflect.Value, foundPtrs map[TypedPointer]bool) (alreadyExists bool) {
	ptr := TypedPointerOfRV(value)
	if _, ok := foundPtrs[ptr]; ok {
		foundPtrs[ptr] = true
		return true
	}

	foundPtrs[ptr] = false
	return false
}
