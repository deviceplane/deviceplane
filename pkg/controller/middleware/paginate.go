package middleware

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrPageNotFound = errors.New("requested page not found")

func simplePaginate(start, pageSize int, arr []interface{}) (ret []interface{}, exists bool) {
	end := start + pageSize

	maxEnd := len(arr)
	if end > maxEnd {
		end = maxEnd
	}

	maxStart := len(arr) - 1
	if start > maxStart || start < 0 {
		return nil, false
	}

	return arr[start:end], true
}

func paginateAfter(after, paginateOn string, pageSize int, arr []interface{}) (ret []interface{}, err error) {
	if len(arr) == 0 {
		return arr, nil
	}

	if paginateOn == "" {
		// We don't want to expose fields that are not exposed through
		// JSON tags
		return nil, ErrEmptyOrdering
	}

	// Catch errors thrown by reflect package or sort.Sort()
	defer func() {
		r := recover()
		if r != nil {
			switch typedPanicMessage := r.(type) {
			case error:
				err = errors.New("Error thrown while pagianting: " + typedPanicMessage.Error())
			case string:
				err = errors.New("Error thrown while pagianting: " + typedPanicMessage)
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	// Make sure the array is homogoneously typed
	var arrayType reflect.Type = reflect.TypeOf(arr[0])
	for i := range arr {
		if reflect.TypeOf(arr[i]) != arrayType {
			return nil, ErrMultiTypeArray
		}
	}

	// Get the index and type of the json-tagged field
	var fieldIndex int
	for i := 0; i < arrayType.NumField(); i++ {
		tags := arrayType.Field(i).Tag
		jsonTag := tags.Get("json")
		if jsonTag == paginateOn {
			fieldIndex = i
			break
		}
	}

	// Paginate
	var startIndex int
	var found bool
	if after == "" {
		startIndex = 0
		found = true
	} else {
		for i := range arr {
			// Find index where `after`'s value equals the `orderby` field's value
			if after == fmt.Sprint(reflect.ValueOf(arr[i]).Field(fieldIndex).Interface()) {
				startIndex = i
				found = true
				break
			}
		}
	}
	if !found {
		return nil, ErrPageNotFound
	}

	page, exists := simplePaginate(startIndex, pageSize, arr)
	if !exists {
		return nil, ErrPageNotFound
	}

	return page, nil
}
