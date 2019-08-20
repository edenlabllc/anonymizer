package gormlogger

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"github.com/rs/zerolog/log"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		duration := float64(v[2].(time.Duration).Nanoseconds()/1e4) / 100.0
		var (
			sql             string
			formattedValues []string
		)
		for _, value := range v[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			if indirectValue.IsValid() {
				value = indirectValue.Interface()
				if t, ok := value.(time.Time); ok {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
				} else if b, ok := value.([]byte); ok {
					if str := string(b); isPrintable(str) {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
					} else {
						formattedValues = append(formattedValues, "'<binary>'")
					}
				} else if r, ok := value.(driver.Valuer); ok {
					if value, err := r.Value(); err == nil && value != nil {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					} else {
						formattedValues = append(formattedValues, "NULL")
					}
				} else {
					switch value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
						formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
					default:
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				}
			} else {
				formattedValues = append(formattedValues, "NULL")
			}
		}

		if numericPlaceHolderRegexp.MatchString(v[3].(string)) {
			sql = v[3].(string)
			for index, value := range formattedValues {
				placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
				sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
			}
		} else {
			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(v[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		log.Debug().Str("module", "gorm").Str("type", "sql").Str("took", fmt.Sprintf("%.1f ms", duration)).Str("sql", sql).Msg("")
	}
	if v[0] == "log" {
		log.Debug().Str("module", "gorm").Str("type", "log").Str("sql", fmt.Sprintf("%+v", v[2])).Msg("")
	}
}
