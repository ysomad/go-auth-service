package errors

import "errors"

var (
	ErrCacheDuplicate = errors.New("given key already exist in cache")
	ErrCacheNotFound  = errors.New("given key not found in cache")
)
