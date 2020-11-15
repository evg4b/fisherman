package utils

func HandleCriticalError(err error) {
	if err != nil {
		panic(err)
	}
}
