package qb

import (
	"fmt"
	"strings"
	"time"
)

func normalize(val any) string {
	switch val := val.(type) {
	case nil:
		return "NULL"
	case []byte:
		return fmt.Sprintf("'%s'", val)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(val)
	case string:
		if val != "" {
			return fmt.Sprintf("'%v'", strings.ReplaceAll(val, "'", "''"))
		}
	case time.Time:
		return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64:
		if val != nil {
			return fmt.Sprint(val)
		}
	case *string:
		if val != nil && *val != "" {
			return fmt.Sprintf("'%v'", strings.ReplaceAll(*val, "'", "''"))
		}
	case *time.Time:
		if val != nil {
			return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
		}
	default:
		panic("unknown val type")
	}
	return ""
}
