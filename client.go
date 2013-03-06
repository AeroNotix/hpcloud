package hpcloud

import (
	"fmt"
)

func (a Access) IsFailed() bool {
	return a.Fail != nil
}

func (a Access) Describe() string {
	return fmt.Sprintf(
		"Code: %d\nDetails: %s\nMessage: %s\n",
		a.Fail.Code(), a.Fail.Details(), a.Fail.Message(),
	)
}
