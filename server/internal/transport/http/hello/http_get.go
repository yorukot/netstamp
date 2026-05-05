package hello

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) Get(ctx context.Context, _ *GetInput) (*GetOutput, error) {
	result, err := h.service.GetGreeting(ctx)
	if err != nil {
		return nil, huma.Error503ServiceUnavailable("request was cancelled")
	}

	return &GetOutput{
		Body: GreetingResponse{
			Message: result.Message,
			Service: result.Service,
		},
	}, nil
}
