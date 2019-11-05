package middleware

func paginate(pageNumber, pageSize int, arr []interface{}) (ret []interface{}, exists bool) {
	start := pageNumber * pageSize
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
