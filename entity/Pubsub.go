package entity

type (
	PublishParam struct {
		TopicName string
		Qos       byte
		Retained  bool
		Message   []byte
	}
)
