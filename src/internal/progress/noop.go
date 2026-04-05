package progress

type NoopSink struct{}

func (NoopSink) Publish(Event) error {
	return nil
}
