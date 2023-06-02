package iteration

func Repeat(chars string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated += chars
	}
	return repeated
}
