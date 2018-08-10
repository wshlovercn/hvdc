package utils

import (
	"errors"
	"runtime/debug"
	"github.com/golang/glog"
)

func RecoverPrint() {
	var err error
	if r := recover(); r !=nil{
		switch x := r.(type) {
		case string:
			err = errors.New(x)
			break
		case error:
			err = x
			break
		default:
			err = errors.New("Unknown panic")
			break
		}

		glog.Error("Panic :", err.Error())
		glog.Error(string(debug.Stack()))
	}
}

