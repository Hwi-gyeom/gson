package gson

import "fmt"
import "testing"

var _ = fmt.Sprintf("dummy")

// All test cases are folded into cbor_value_test.go, contains only few
// missing testcases (if any) and benchmarks.

func BenchmarkVal2CborNull(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(nil)

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborTrue(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(true)

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborFalse(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(false)
	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborUint8(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(uint8(255))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborInt8(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(int8(-128))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborUint16(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(uint16(65535))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborInt16(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(int16(-32768))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborUint32(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(uint32(4294967295))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborInt32(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(int32(-2147483648))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborUint64(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(uint64(18446744073709551615))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborInt64(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(int64(-2147483648))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborFlt32(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(float32(10.2))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborFlt64(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue(float64(10.2))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborTBytes(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue([]byte("hello world"))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborText(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 128), 0)
	val := config.NewValue("hello world")

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborArr0(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 1024), 0)
	val := config.NewValue(make([]interface{}, 0))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborArr5(b *testing.B) {
	value := interface{}([]interface{}{5, 5.0, "hello world", true, nil})
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 1024), 0)
	val := config.NewValue(value)

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborMap0(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 1024), 0)
	val := config.NewValue(make([][2]interface{}, 0))

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborMap5(b *testing.B) {
	value := interface{}([][2]interface{}{
		[2]interface{}{"key0", 5}, [2]interface{}{"key1", 5.0},
		[2]interface{}{"key2", "hello world"},
		[2]interface{}{"key3", true}, [2]interface{}{"key4", nil},
	})
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 1024), 0)
	val := config.NewValue(value)

	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}

func BenchmarkVal2CborTyp(b *testing.B) {
	config := NewDefaultConfig()
	cbr := config.NewCbor(make([]byte, 10*1024), 0)
	jsn := config.NewJson(testdataFile("testdata/typical.json"), -1)
	_, value := jsn.Tovalue()
	val := config.NewValue(value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val.Tocbor(cbr.Reset(nil))
	}
}
