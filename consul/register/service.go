package main

import "errors"

// Service Define a service interface
type Service interface {
	// Add calculate a+b
	Add(a, b int) int

	// Subtract calculate a-b
	Subtract(a, b int) int

	// Multiply calculate a*b
	Multiply(a, b int) int

	// Divide calculate a/b
	Divide(a, b int) (int, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

// service 实现
type ArithmeticService struct {
}

func (ArithmeticService) Add(a, b int) int {
	return a + b
}

func (ArithmeticService) Subtract(a, b int) int {
	return a - b
}

func (ArithmeticService) Multiply(a, b int) int {
	return a * b
}

func (ArithmeticService) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("the dividend can not be zero!")
	}
	return a / b, nil
}

func (ArithmeticService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
