package cluster

type Action uint8

const (
	PUT Action = iota
	DEL
)

func NewEvent(data *Data, action Action) *Event {
	return &Event{
		Data:   data,
		Action: action,
	}
}

type Event struct {
	*Data
	Action
}
