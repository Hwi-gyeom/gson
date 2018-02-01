// transform golang native value into binary collated encoding.
// cnf: NumberKind, doMissing, arrayLenPrefix, propertyLenPrefix

package gson

import "fmt"
import "strconv"
import "encoding/json"

func gson2collate(obj interface{}, code []byte, config *Config) int {
	if obj == nil {
		code[0], code[1] = TypeNull, Terminator
		return 2
	}

	switch value := obj.(type) {
	case bool:
		if value {
			code[0] = TypeTrue
		} else {
			code[0] = TypeFalse
		}
		code[1] = Terminator
		return 2

	case float64:
		n := 0
		code[n] = TypeNumber
		n++
		n += collateFloat64(value, code[n:])
		code[n] = Terminator
		n++
		return n

	case float32:
		n := 0
		code[n] = TypeNumber
		n++
		n += collateFloat64(float64(value), code[n:])
		code[n] = Terminator
		n++
		return n

	case int64:
		n := 0
		code[n] = TypeNumber
		n++
		n += collateInt64(value, code[n:])
		code[n] = Terminator
		n++
		return n

	case uint64:
		n := 0
		code[n] = TypeNumber
		n++
		n += collateUint64(value, code[n:])
		code[n] = Terminator
		n++
		return n

	case int:
		n := 0
		code[n] = TypeNumber
		n++
		n += collateInt64(int64(value), code[n:])
		code[n] = Terminator
		n++
		return n

	case json.Number:
		if isnegative(value) {
			n := 0
			code[n] = TypeNumber
			n++
			n += collateInt64Str(string(value), code[n:])
			code[n] = Terminator
			n++
			return n
		}
		n := 0
		code[n] = TypeNumber
		n++
		n += collateUint64Str(string(value), code[n:])
		code[n] = Terminator
		n++
		return n

	case Missing:
		if config.doMissing && MissingLiteral.Equal(string(value)) {
			code[0], code[1] = TypeMissing, Terminator
			return 2
		}
		panic("collate missing not configured")

	case string:
		return collateString(value, code, config)

	case []byte:
		n := 0
		code[n] = TypeBinary
		n++
		m := copy(code[n:], value)
		n += m
		code[n] = Terminator
		n++
		return n

	case []interface{}:
		n := 0
		code[n] = TypeArray
		n++
		if config.arrayLenPrefix {
			n += collateLength(len(value), code[n:])
		}
		for _, val := range value {
			n += gson2collate(val, code[n:], config)
		}
		code[n] = Terminator
		n++
		return n

	case map[string]interface{}:
		n := 0
		code[n] = TypeObj
		n++
		if config.propertyLenPrefix {
			n += collateLength(len(value), code[n:])
		}

		poolobj := config.pools.keysPool.Get()
		keys := poolobj.([]string)
		defer config.pools.keysPool.Put(poolobj)

		for _, key := range sortProps1(value, keys) {
			n += collateString(key, code[n:], config)       // encode key
			n += gson2collate(value[key], code[n:], config) // encode value
		}
		code[n] = Terminator
		n++
		return n
	}
	panic(fmt.Errorf("collate invalid golang type %T", obj))
}

func collateLength(length int, code []byte) (n int) {
	var num [64]byte
	code[n] = TypeLength
	n++
	bs := strconv.AppendInt(num[:0], int64(length), 10)
	n += collateInt(bs, code[n:])
	code[n] = Terminator
	n++
	return n
}

func collateString(str string, code []byte, config *Config) (n int) {
	if config.doMissing && MissingLiteral.Equal(str) {
		code[0], code[1] = TypeMissing, Terminator
		return 2
	}
	code[n] = TypeString
	n++
	n += suffixEncodeString(str2bytes(str), code[n:])
	code[n] = Terminator
	n++
	return n
}
