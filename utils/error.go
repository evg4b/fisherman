package utils

// HandleCriticalError check error and when error is not nil call panic.
func HandleCriticalError(err error) {
	if err != nil {
		panic(err)
	}
}
