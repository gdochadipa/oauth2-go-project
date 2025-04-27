package util

import "github.com/golang/glog"

func TerminateIfError(err error) {
	if err != nil {
		glog.Fatalf("%v", err)
	}
}
