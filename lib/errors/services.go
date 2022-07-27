package errors

import (
	"errors"
	"fmt"
)

var (
	//ErrNoAuthService provides a clean way to prevent auth service for throwing
	//exceptions during any initialization that might require it
	ErrNoAuthService = errors.New("no auth service was provided during service initialization, please provide one")
	//ErrNoPaypalService provides a clean way to prevent paypal service for throwing
	//exceptions during any initialization that might require it
	ErrNoPaypalService = errors.New("no paypal service was provided during service initialization, please provide one")
	//ErrNoOrderService provides a clean way to prevent order service for throwing
	//exceptions during any initialization that might require it
	ErrNoOrderService = errors.New("no order service was provided during service initialization, please provide one")
	//ErrNoProductService provides a clean way to prevent product service for throwing
	//exceptions during any initialization that might require it
	ErrNoProductService = errors.New("no product service was provided during service initialization, please provide one")
	//ErrNoCartService provides a clean way to prevent cart service for throwing
	//exceptions during any initialization that might require it
	ErrNoCartService = errors.New("no cart service was provided during service initialization, please provide one")
)

type ErrInvalidRequest struct {
	Fields map[string]string
}

func (err *ErrInvalidRequest) Error() (str string) {

	for key, prop := range err.Fields {
		str += fmt.Sprintf("\n%s:%s", key, prop)
	}

	return
}
