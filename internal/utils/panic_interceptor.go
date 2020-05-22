package utils

func PanicInterceptor(action func(any)) {
	if recovered := recover(); recovered != nil {
		action(recovered)
	}
}
