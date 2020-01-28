package middleware

import (
	"errors"
	"reflect"
	"sort"
)

var (
	ErrTypeUnsupported = errors.New("internal type not supported")
	ErrEmptyOrdering   = errors.New("ordering by an empty string is invalid")
	ErrMultiTypeArray  = errors.New("only homogenously typed arrays supported")
)

type genericSortableArray struct {
	arr        []interface{}
	fieldIndex int
	fieldType  reflect.Type
}

func (s genericSortableArray) Len() int {
	return len(s.arr)
}

func (s genericSortableArray) Swap(i, j int) {
	temp := s.arr[j]
	s.arr[j] = s.arr[i]
	s.arr[i] = temp
}

func (s genericSortableArray) Less(i, j int) bool {
	iValue := reflect.ValueOf(s.arr[i]).Field(s.fieldIndex)
	jValue := reflect.ValueOf(s.arr[j]).Field(s.fieldIndex)

	less, err := genericLess(s.fieldType, iValue, jValue)
	if err != nil {
		panic("cannot compare types, " + err.Error())
	}

	return less
}

func genericLess(fieldType reflect.Type, a, b reflect.Value) (bool, error) {
	if fieldType.Kind() == reflect.Ptr {
		if a.IsNil() {
			return true, nil
		}
		if b.IsNil() {
			return false, nil
		}

		fieldType = fieldType.Elem()
		a = a.Elem()
		b = b.Elem()
	}

	switch fieldType.Kind() {
	case reflect.String:
		aVal := a.String()
		bVal := b.String()
		return aVal < bVal, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		aVal := a.Int()
		bVal := b.Int()
		return aVal < bVal, nil
	case reflect.Float32, reflect.Float64:
		aVal := a.Float()
		bVal := b.Float()
		return aVal < bVal, nil
	}
	return false, ErrTypeUnsupported
}

func isGenericLessSupported(fieldType reflect.Type) bool {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	switch fieldType.Kind() {
	case reflect.String:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func order(orderBy string, direction orderDirection, arr []interface{}) (err error) {
	if len(arr) == 0 {
		return nil
	}

	if orderBy == "" {
		// We don't want to expose fields that are not exposed through
		// JSON tags
		return ErrEmptyOrdering
	}

	// Catch errors thrown by reflect package or sort.Sort()
	defer func() {
		r := recover()
		if r != nil {
			switch typedPanicMessage := r.(type) {
			case error:
				err = errors.New("Error thrown while ordering: " + typedPanicMessage.Error())
			case string:
				err = errors.New("Error thrown while ordering: " + typedPanicMessage)
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	// Get the type of the array, get the index and type of the json-tagged field
	var arrayType reflect.Type = reflect.TypeOf(arr[0])
	var fieldType reflect.Type
	var fieldIndex int

	for i := 0; i < arrayType.NumField(); i++ {
		tags := arrayType.Field(i).Tag
		jsonTag := tags.Get("json")
		if jsonTag == orderBy {
			fieldIndex = i
			fieldType = arrayType.Field(i).Type
			break
		}
	}

	// Make sure the type is supported by our generic Less function
	if !isGenericLessSupported(fieldType) {
		return ErrTypeUnsupported
	}

	// Make sure the array is homogoneously typed
	for i := range arr {
		if reflect.TypeOf(arr[i]) != arrayType {
			return ErrMultiTypeArray
		}
	}

	// Get our generic sortable array
	sortableArr := sort.Interface(genericSortableArray{
		arr,
		fieldIndex,
		fieldType,
	})

	// Reverse the direction if needed, then sort.
	if direction == OrderDescending {
		sortableArr = sort.Reverse(sortableArr)
	}
	sort.Sort(sortableArr)
	return nil
}
