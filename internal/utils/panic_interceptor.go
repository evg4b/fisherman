package utils

func PanicInterceptor(action func(interface{})) {
	if recovered := recover(); recovered != nil {
		action(recovered)
	}
}
