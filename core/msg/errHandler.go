package msg

import (
	"fmt"
	"os"
)

func Fail(err error){
	if err != nil{
	Error(err.Error())
		os.Exit(1)
	}
}
func Wrap(msg string, err error) error{
	if err == nil {return nil}
	return fmt.Errorf("%s: %w", msg,err)
}
