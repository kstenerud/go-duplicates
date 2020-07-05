// duplicates examines an abitrary object and reports any duplicate pointers it
// finds (where more than one pointer is pointing to the same object).
package duplicates

import (
	"reflect"
)

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
	finder := NewDuplicateFinder()
	finder.ScanObject(value)
	return finder.DuplicatePointers
}

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

// TypedPointerOf gets the typed pointer of the object that value references.
// Note: value must be addressable or else this function will panic.
func TypedPointerOfRV(rv reflect.Value) TypedPointer {
	return TypedPointer{
		Type:    rv.Type(),
		Pointer: rv.Pointer(),
	}
}

// DuplicateFinder scans objects for pointers and keeps track of them so that
// any duplicates can be found.
type DuplicateFinder struct {
	DuplicatePointers map[TypedPointer]bool
}

func NewDuplicateFinder() *DuplicateFinder {
	_this := &DuplicateFinder{}
	_this.Init()
	return _this
}

func (_this *DuplicateFinder) Init() {
	_this.DuplicatePointers = make(map[TypedPointer]bool)
}

// Record a pointer, returning true if it has been recorded before.
// This method panics if the pointer's Kind is not Chan, Func, Map, Ptr, Slice,
// or UnsafePointer.
func (_this *DuplicateFinder) CheckPtrAlreadyFound(pointer reflect.Value) (alreadyExists bool) {
	typedPtr := TypedPointerOfRV(pointer)
	if _, ok := _this.DuplicatePointers[typedPtr]; ok {
		_this.DuplicatePointers[typedPtr] = true
		return true
	}

	_this.DuplicatePointers[typedPtr] = false
	return false
}

func (_this *DuplicateFinder) ScanObject(object interface{}) {
	_this.scanRV(reflect.ValueOf(object))
}

func (_this *DuplicateFinder) scanRV(value reflect.Value) {
	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return
		}
		elem := value.Elem()
		if !isSearchableKind(elem.Kind()) {
			return
		}
		_this.scanRV(elem)
	case reflect.Ptr:
		if value.IsNil() {
			return
		}
		if _this.CheckPtrAlreadyFound(value) {
			return
		}
		elem := value.Elem()
		if !isSearchableKind(elem.Type().Kind()) {
			return
		}
		_this.scanRV(elem)
	case reflect.Map:
		if value.IsNil() {
			return
		}
		if _this.CheckPtrAlreadyFound(value) {
			return
		}
		if !isSearchableKind(value.Type().Elem().Kind()) {
			return
		}
		iter := mapRange(value)
		for iter.Next() {
			_this.scanRV(iter.Value())
		}
	case reflect.Slice:
		if value.IsNil() {
			return
		}
		if _this.CheckPtrAlreadyFound(value) {
			return
		}
		if !isSearchableKind(value.Type().Elem().Kind()) {
			return
		}
		count := value.Len()
		for i := 0; i < count; i++ {
			_this.scanRV(value.Index(i))
		}
	case reflect.Array:
		if !isSearchableKind(value.Type().Elem().Kind()) {
			return
		}
		count := value.Len()
		for i := 0; i < count; i++ {
			_this.scanRV(value.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanAddr() {
				field = field.Addr()
			}
			if isSearchableKind(field.Kind()) {
				_this.scanRV(field)
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
