package hello

type Greeting struct {
	message string
	service string
}

func NewGreeting(service string) Greeting {
	return Greeting{
		message: "Hello from NetStamp",
		service: service,
	}
}

func (g Greeting) Message() string {
	return g.message
}

func (g Greeting) Service() string {
	return g.service
}
