package errors

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

//ErrNoProductFound ...
type ErrNoProductFound struct {
	ID uuid.UUID
}

func (err *ErrNoProductFound) Error() string {
	return fmt.Sprintf("no product return with the id %s", err.ID)
}

var (
	//ErrProductNotProvide is a generic error for one
	ErrProductNotProvided = errors.New("there was no product provided when one was required, please provide a proper product")
)
