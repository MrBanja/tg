package tg

var token string

func SetToken(t string) {
	token = t
}

func ToOptional[T any](t T) *T {
	return &t
}
