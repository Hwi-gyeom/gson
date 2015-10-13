//  Copyright (c) 2015 Couchbase, Inc.

package gson

import "strconv"
import "fmt"
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
		n += normalizeFloat(obj, code[n:], config.nk)
		code[n] = Terminator
		n++
		return n

	case float32:
		n := 0
		code[n] = TypeNumber
		n++
		n += normalizeFloat(float64(value), code[n:], config.nk)
		code[n] = Terminator
		n++
		return n

	case int64:
		n := 0
		code[n] = TypeNumber
		n++
		n += normalizeFloat(obj, code[n:], config.nk)
		code[n] = Terminator
		n++
		return n

	case uint64:
		n := 0
		code[n] = TypeNumber
		n++
		n += normalizeFloat(int64(value), code[n:], config.nk)
		code[n] = Terminator
		n++
		return n

	case int:
		n := 0
		code[n] = TypeNumber
		n++
		n += normalizeFloat(int64(value), code[n:], config.nk)
		code[n] = Terminator
		n++
		return n

	case json.Number:
		n := 0
		code[n] = TypeNumber
		n++
		f, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			panic(err)
		}
		n += normalizeFloat(f, code[n:], config.nk)
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
		if config.doMissing && MissingLiteral.Equal(value) {
			code[0], code[1] = TypeMissing, Terminator
			return 2
		}
		n := 0
		code[n] = TypeString
		n++
		n += suffixEncodeString(str2bytes(value), code[n:])
		code[n] = Terminator
		n++
		return n

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

		keys := config.pools.keysPool.Get().([]string)
		defer config.pools.keysPool.Put(keys[:0])
		for _, key := range sortProps(value, keys) {
			n += gson2collate(key, code[n:], config)        // encode key
			n += gson2collate(value[key], code[n:], config) // encode value
		}
		code[n] = Terminator
		n++
		return n
	}
	panic(fmt.Errorf("collate invalid golang type %T", obj))
}

func collate2gson(code []byte, config *Config) (interface{}, int) {
	if len(code) == 0 {
		return nil, 0
	}

	var scratch [64]byte
	n := 1
	switch code[0] {
	case TypeMissing:
		return MissingLiteral, 2

	case TypeNull:
		return nil, 2

	case TypeTrue:
		return true, 2

	case TypeFalse:
		return false, 2

	case TypeNumber:
		m := getDatum(code[n:])
		f := denormalizeFloat(code[n:n+m-1], config.nk) // -1 to skip terminator
		return f, n + m

	case TypeString:
		s := make([]byte, encodedStringSize(code[n:]))
		x, y := suffixDecodeString(code[n:], s)
		return bytes2str(s[:y]), n + x

	case TypeBinary:
		m := getDatum(code[n:])
		bs := make([]byte, m-1)
		copy(bs, code[n:n+m-1])
		return bs, n + m

	case TypeArray:
		var arr []interface{}
		if config.arrayLenPrefix {
			if code[n] != TypeLength {
				panic("collate decode expected array length prefix")
			}
			n++
			m := getDatum(code[n:])
			_, y := collated2Int(code[n:n+m], scratch[:])
			ln, err := strconv.Atoi(bytes2str(scratch[:y]))
			if err != nil {
				panic(err)
			}
			arr = make([]interface{}, 0, ln)
			n += m
		} else {
			arr = make([]interface{}, 0, 8)
		}
		for code[n] != Terminator {
			item, y := collate2gson(code[n:], config)
			arr = append(arr, item)
			n += y
		}
		return arr, n + 1 // +1 to skip terminator

	case TypeObj:
		obj := make(map[string]interface{})
		if config.propertyLenPrefix {
			if code[n] != TypeLength {
				panic("collate decode expected object length prefix")
			}
			n++
			m := getDatum(code[n:])
			_, y := collated2Int(code[n:n+m], scratch[:])
			_, err := strconv.Atoi(bytes2str(scratch[:y])) // just skip
			if err != nil {
				panic(err)
			}
			n += m
		}
		for code[n] != Terminator {
			key, m := collate2gson(code[n:], config)
			n += m
			value, m := collate2gson(code[n:], config)
			obj[key.(string)] = value
			n += m
		}
		return obj, n + 1 // +1 to skip terminator
	}
	panic("collate decode invalid binary")
}

// get the collated datum based on Terminator and return the length
// of the datum.
func getDatum(code []byte) int {
	var i int
	var b byte
	for i, b = range code {
		if b == Terminator {
			return i + 1
		}
	}
	panic("collate decode terminator not found")
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
