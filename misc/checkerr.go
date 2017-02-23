package misc

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}


