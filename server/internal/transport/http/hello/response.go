package hello

type GetInput struct{}

type GetOutput struct {
	Body GreetingResponse
}

type GreetingResponse struct {
	Message string `json:"message" doc:"Greeting text." example:"Hello from Netstamp"`
	Service string `json:"service" doc:"Service name that generated the greeting." example:"netstamp-api"`
}
