package common

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Gansack(db *gorm.DB, cond map[string]interface{}) error {
	for k, value := range cond {
		splitKey := strings.Split(k, "__")

		if len(splitKey) == 2 {
			field := splitKey[0]
			value = mappingValue(value)

			switch splitKey[1] {
			case "gt":
				db = db.Where(field+" > ?", value)
			case "lt":
				db = db.Where(field+" < ?", value)
			case "gte":
				db = db.Where(field+" >= ?", value)
			case "lte":
				db = db.Where(field+" <= ?", value)
			case "ne":
				db = db.Where(field+" != ?", value)

			//	sorting
			case "asc":
				db = db.Order(fmt.Sprintf("%s ASC", field))
			case "desc":
				db = db.Order(fmt.Sprintf("%s DESC", field))
			default:
				return errors.New(fmt.Sprintf("field %s with operator %s not valid", field, splitKey[1]))
			}

			delete(cond, k)
		}
	}

	return nil
}

func mappingValue(a interface{}) interface{} {
	switch value := a.(type) {
	case string:
		return a
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return a
	case bool:
		return true
	case time.Time:
		return value.Format("2006-01-02 15:04:05")
	default:
		return a.(fmt.Stringer).String()
	}
}
