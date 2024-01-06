package document

import (
	"time"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

type Status int

const (
	StatusDraft   Status = 0
	StatusPosted  Status = 1
	StatusDeleted Status = 2
)

var ValidStatuses = []Status{
	StatusDraft,
	StatusPosted,
	StatusDeleted,
}

type Document struct {
	ID            uint64
	Number        string
	Date          time.Time
	Status        Status
	BasisDocument *Document
	Company       reference.Company
	AppendUser    *reference.User
	ChangeUser    *reference.User
}
