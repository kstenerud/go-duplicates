// +build go1.12

package duplicates

import (
	"reflect"
)

func mapRange(v reflect.Value) *reflect.MapIter {
	return v.MapRange()
}
