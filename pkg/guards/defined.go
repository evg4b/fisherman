package guards

func ShouldBeDefined(object interface{}, message string) {
	if object == nil {
		panic(message)
	}
}
