package hello

import (
	"context"

	domainhello "github.com/yorukot/netstamp/internal/domain/hello"
)

type Service struct {
	serviceName string
}

func NewService(serviceName string) *Service {
	return &Service{serviceName: serviceName}
}

func (s *Service) GetGreeting(ctx context.Context) (GreetingResult, error) {
	if err := ctx.Err(); err != nil {
		return GreetingResult{}, err
	}

	greeting := domainhello.NewGreeting(s.serviceName)
	return GreetingResult{
		Message: greeting.Message(),
		Service: greeting.Service(),
	}, nil
}
