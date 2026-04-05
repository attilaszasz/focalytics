package progress

type Sink interface {
	Publish(event Event) error
}
