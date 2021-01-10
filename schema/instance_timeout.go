package schema

import (
	"log"
	"time"
)

const (
	TimeoutCreate  = "create"
	TimeoutRead    = "read"
	TimeoutUpdate  = "update"
	TimeoutDelete  = "delete"
	TimeoutDefault = "default"
)

func timeoutKeys() []string {
	return []string{
		TimeoutCreate,
		TimeoutRead,
		TimeoutUpdate,
		TimeoutDelete,
		TimeoutDefault,
	}
}

// could be time.Duration, int64 or float64
func DefaultTimeout(tx interface{}) *time.Duration {
	var td time.Duration
	switch raw := tx.(type) {
	case time.Duration:
		return &raw
	case int64:
		td = time.Duration(raw)
	case float64:
		td = time.Duration(int64(raw))
	default:
		log.Printf("[WARN] Unknown type in DefaultTimeout: %#v", tx)
	}
	return &td
}

type InstanceTimeout struct {
	Create, Read, Update, Delete, Default *time.Duration
}
