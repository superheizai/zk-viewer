package base

import "errors"

var (
	// ErrClosed is the error resulting if the pool is Closed via pool.Close().

	ResponseErrorCode   = 5000
	NoClusterFoundError = errors.New("没有对应的集群")
	TimeOutError        = errors.New("访问超时")
	NoReulstFoundError = errors.New("no result found")

)
