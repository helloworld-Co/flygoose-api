package tools

// Find获取一个切片并在其中查找元素。如果找到它，它将返回它的索引，否则它将返回-1和一个错误的bool。
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
