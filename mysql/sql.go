package mysql

import (
	libJSON "encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tonny-zhang/goutils/num"
)

const defaultBufSize = 4096
const statusNoBackslashEscapes = true

func interpolateParams(query string, args ...any) (string, error) {
	// Number of ? should be same to len(args)
	if strings.Count(query, "?") != len(args) {
		err := fmt.Errorf("argv [%v] number error, excepted %d, actural %d", args, strings.Count(query, "?"), len(args))

		return query, err
	}

	var buf []byte
	var err error
	argPos := 0

	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf = append(buf, query[i:]...)
			break
		}
		buf = append(buf, query[i:i+q]...)
		i += q

		arg := args[argPos]
		argPos++

		if arg == nil {
			buf = append(buf, "NULL"...)
			continue
		}

		switch v := arg.(type) {
		case int:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int16:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case uint:
			buf = strconv.AppendUint(buf, uint64(v), 10)
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case uint64:
			// Handle uint64 explicitly because our custom ConvertValue emits unsigned values
			buf = strconv.AppendUint(buf, v, 10)
		case float64:
			buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
		case bool:
			if v {
				buf = append(buf, '1')
			} else {
				buf = append(buf, '0')
			}
		case time.Time:
			if v.IsZero() {
				buf = append(buf, "'0000-00-00'"...)
			} else {
				buf = append(buf, '\'')
				buf, err = appendDateTime(buf, v)
				if err != nil {
					return "", err
				}
				buf = append(buf, '\'')
			}
		case libJSON.RawMessage:
			buf = append(buf, '\'')
			if statusNoBackslashEscapes {
				buf = escapeBytesBackslash(buf, v)
			} else {
				buf = escapeBytesQuotes(buf, v)
			}
			buf = append(buf, '\'')
		case []byte:
			if v == nil {
				buf = append(buf, "NULL"...)
			} else {
				buf = append(buf, "_binary'"...)
				if statusNoBackslashEscapes {
					buf = escapeBytesBackslash(buf, v)
				} else {
					buf = escapeBytesQuotes(buf, v)
				}
				buf = append(buf, '\'')
			}
		case string:
			buf = append(buf, '\'')
			if statusNoBackslashEscapes {
				buf = escapeStringBackslash(buf, v)
			} else {
				buf = escapeStringQuotes(buf, v)
			}
			buf = append(buf, '\'')
		case []int:
			buf = append(buf, '(')
			buf = append(buf, num.Join(v, ",")...)
			buf = append(buf, ')')
		case []int16:
			buf = append(buf, '(')
			buf = append(buf, num.Join(v, ",")...)
			buf = append(buf, ')')
		case []string:
			buf = append(buf, '(')
			len := len(v)
			for i, str := range v {
				buf = append(buf, '\'')
				if statusNoBackslashEscapes {
					buf = escapeStringBackslash(buf, str)
				} else {
					buf = escapeStringQuotes(buf, str)
				}
				buf = append(buf, '\'')
				if i != len-1 {
					buf = append(buf, ',')
				}
			}

			buf = append(buf, ')')
		default:
			return query, fmt.Errorf("[%d]argv [%v:%s] no support", argPos, v, reflect.TypeOf(arg))
		}

	}
	if argPos != len(args) {
		return query, fmt.Errorf("args number [%d] error not [%d]", argPos, len(args))
	}
	return string(buf), nil
}
