package errors

import (
	"errors"
	"fmt"
)

//ErrCartActionNotRecognized the structure and the information that it needs to produce a valid err
type ErrCartActionNotRecognized struct {
	Action string
}

func (err *ErrCartActionNotRecognized) Error() string {
	return fmt.Sprintf("the cart action %s was not recognized, please try with a valid action", err.Action)
}

//ErrNoCart is returned when there is no cart for
//the provided order. Typically this happens when
//there aren't any products associated with the
//order itself.
type ErrNoCart struct {
	OrderID int64
}

func (err *ErrNoCart) Error() string {
	return fmt.Sprintf("no cart was returned for order %d, please add a product to the order in order to initialize one", err.OrderID)
}

var (
	//ErrCartOrderNotProvided is returned when there is no order provided during the any process that my need it
	ErrCartOrderNotProvided = errors.New("no order was provided while trying to query one that requires it, please provide a valid order")
)
