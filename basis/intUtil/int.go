package intUtil

func Int32ToPointer(input int32) *int32 {
	tmp := new(int32)
	*tmp = input
	return tmp
}
