package hello

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) get(ctx context.Context, _ *getInput) (*getOutput, error) {
	result, err := h.service.GetGreeting(ctx)
	if err != nil {
		return nil, huma.Error503ServiceUnavailable("request was cancelled")
	}

	return &getOutput{
		Body: greetingResponse{
			Message: result.Message,
			Service: result.Service,
		},
	}, nil
}

type getInput struct{}

type getOutput struct {
	Body greetingResponse
}

type greetingResponse struct {
	Message string `json:"message" doc:"Greeting text." example:"Hello from Netstamp"`
	Service string `json:"service" doc:"Service name that generated the greeting." example:"netstamp-api"`
}
