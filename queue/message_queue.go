package controller

type MessageQueue interface {
	Subscribe(string, func([]byte)) error
	Publish(string, []byte) error
}
