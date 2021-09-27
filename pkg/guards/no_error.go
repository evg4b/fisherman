package guards

// NoError panic when passed error is not nil
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}
