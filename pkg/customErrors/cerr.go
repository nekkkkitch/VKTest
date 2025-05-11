package cerr

import "fmt"

type CustomError error

var (
	ErrSubClosed CustomError = fmt.Errorf("SubPub hub is closed")
	ErrNoTopic   CustomError = fmt.Errorf("Topic does not exists")
)
