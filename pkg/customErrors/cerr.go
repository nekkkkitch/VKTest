package cerr

import "fmt"

type CustomError error

var (
	ErrSubClosed    CustomError = fmt.Errorf("SubPub hub is closed")
	ErrNoTopic      CustomError = fmt.Errorf("Topic does not exists")
	ErrEmptyRequest CustomError = fmt.Errorf("Request is empty")
	ErrEmptyTopic   CustomError = fmt.Errorf("Topic is empty")
	ErrEmptyMessage CustomError = fmt.Errorf("Message is empty")
)
