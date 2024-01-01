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
		return fmt.Sprintf("'%v'", strings.ReplaceAll(string(val), "'", "''"))
	case string:
		return fmt.Sprintf("'%v'", strings.ReplaceAll(val, "'", "''"))
	case *string:
		if val != nil {
			return fmt.Sprintf("'%v'", strings.ReplaceAll(*val, "'", "''"))
		}
		return "NULL"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(val)
	case *int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64:
		if val != nil {
			switch val := val.(type) {
			case *int:
				return fmt.Sprint(*val)
			case *int8:
				return fmt.Sprint(*val)
			case *int16:
				return fmt.Sprint(*val)
			case *int32:
				return fmt.Sprint(*val)
			case *int64:
				return fmt.Sprint(*val)
			case *uint:
				return fmt.Sprint(*val)
			case *uint8:
				return fmt.Sprint(*val)
			case *uint16:
				return fmt.Sprint(*val)
			case *uint32:
				return fmt.Sprint(*val)
			case *uint64:
				return fmt.Sprint(*val)
			case *float32:
				return fmt.Sprint(*val)
			case *float64:
				return fmt.Sprint(*val)
			}
		}
		return "NULL"
	case time.Time:
		return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
	case *time.Time:
		if val != nil {
			return fmt.Sprintf("'%v'", val.Format("2006-01-02 15:04:05"))
		}
		return "NULL"
	default:
		panic("unknown val type")
	}
}
