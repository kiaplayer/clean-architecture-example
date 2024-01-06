package reference

type Status int

const (
	StatusActive  Status = 0
	StatusDeleted Status = 1
)

var ValidStatuses = []Status{
	StatusActive,
	StatusDeleted,
}

type Reference struct {
	ID     uint64
	Name   string
	Status Status
}
