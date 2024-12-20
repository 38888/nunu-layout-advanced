package pagination

func GetPageOffset(pageNum, pageSize int) int {
	return (pageNum - 1) * pageSize
}
func GetCurrent(current int64) int {
	if current <= 0 {
		return 1
	}
	return int(current)
}
func GetPageSize(pageSize int64) int {
	if pageSize <= 0 {
		return 20
	}
	return int(pageSize)
}
