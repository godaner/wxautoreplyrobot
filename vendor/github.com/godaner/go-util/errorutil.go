package go_util

//@OutOfDate
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}