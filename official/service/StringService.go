package service

// service
type StringService interface {
	UpperCase(string) (string, error)
	Count(string) int
}
