package fix

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"
)

// BenchmarkFloatAppend-24    	37752706	        30.91 ns/op
func BenchmarkFloatAppend(b *testing.B) {
	v := 123.456
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.AppendFloat(make([]byte, 0, 64), v, 'f', -1, 64)
	}
}

// BenchmarkFormatFloat-24    	25714983	        45.97 ns/op
func BenchmarkFormatFloat(b *testing.B) {
	v := 123123213.12345678
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(strconv.FormatFloat(v, 'f', -1, 64))
	}
}

// BenchmarkFormatToFloat-24   a         			 7569190	         105.0 ns/op
// BenchmarkFormatToFloat-24   0         			 540723948	         2.209 ns/op
// BenchmarkFormatToFloat-24   123         			 439286569	         2.725 ns/op
// BenchmarkFormatToFloat-24   0.12345678  			 162624674	         7.719 ns/op
// BenchmarkFormatToFloat-24   0.131212212 			 153940886	         7.817 ns/op
// BenchmarkFormatToFloat-24   12312312312.1312 	 123302217	         10.51 ns/op
// BenchmarkFormatToFloat-24   a					 58579447	         15.25 ns/op
// BenchmarkFormatToFloat-24   12312312312.131212212 123302217	         39.66 ns/op
func BenchmarkFormatToFloat(b *testing.B) {
	v := []byte("12312312312.131212212")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bytesToFloat(v)
	}
}

// BenchmarkFormatToFloatStrConv-24    	58696927	        33.25 ns/op
func BenchmarkFormatToFloatStrConv(b *testing.B) {
	v := []byte("a")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = strconv.ParseFloat(string(v), 64)
	}
}

type testFloatCase struct {
	v         []byte
	e         float64
	err       string
	errCustom string
}

func TestFormatFloatToStr(t *testing.T) {
	vv := []testFloatCase{
		{[]byte("0"), 0, "", ""},
		{[]byte("123"), 123, "", ""},
		{[]byte("123.456"), 123.456, "", ""},
		{[]byte("0.456"), 0.456, "", ""},
		{[]byte("0.1"), 0.1, "", ""},
		{[]byte("0.0000001"), 0.0000001, "", ""},
		{[]byte("-0.0000001"), -0.0000001, "", ""},
		{[]byte("0"), 0, "", ""},
		{[]byte("0.1"), 0.1, "", ""},
	}
	for _, c := range vv {
		t.Run(fmt.Sprintf("float case %+v", c.e), func(t *testing.T) {
			t.Logf("case %+v", c)
			v := []byte(strconv.FormatFloat(c.e, 'f', -1, 64))
			if !bytes.Equal(floatToBytes(c.e), c.v) {
				t.Errorf("got %v, want %v", v, c.v)
			}
		})
	}

}

type testUint64Case struct {
	v         []byte
	e         uint64
	err       string
	errCustom string
}

func TestFormatUintToStr(t *testing.T) {
	vv := []testUint64Case{
		{[]byte("0"), 0, "", ""},
		{[]byte("123"), 123, "", ""},
		{[]byte("-123"), 0,
			"strconv.ParseUint: parsing \"-123\": invalid syntax", "invalid input: non-numeric character"},
		{[]byte(""), 0, "strconv.ParseUint: parsing \"\": invalid syntax", "invalid input: empty byte slice"},
	}
	for _, c := range vv {
		t.Run(fmt.Sprintf("uint case %d", c.e), func(t *testing.T) {
			t.Logf("case %+v", c)
			v, err := strconv.ParseUint(string(c.v), 10, 64)
			if err != nil && c.err != err.Error() {
				t.Errorf("got error %+v, want %v", err, c.err)
			}
			i := NewUint(0)
			err2 := i.FromBytes(c.v)
			if err2 != nil && c.errCustom != err2.Error() {
				t.Errorf("bytesToUint: got error %+v, want %v", err2, c.errCustom)
			} else if v != i.value {
				t.Errorf("bytesToUint: got %v, want %v", v, c.v)
			}
			if c.err != "" {
				return
			}

			intBytes := NewUint(c.e)
			if !bytes.Equal(intBytes.ToBytes(), c.v) {
				t.Errorf("uintToBytes failed: got %v, want %v", intBytes, c.v)
			}
		})
	}
}

type testIntCase struct {
	v         []byte
	e         int
	err       string
	errCustom string
}

func TestFormatIntToStr(t *testing.T) {
	vv := []testIntCase{
		{[]byte("0"), 0, "", ""},
		{[]byte("123"), 123, "", ""},
		{[]byte("-123"), -123, "", ""},
		{[]byte(""), 0, "strconv.Atoi: parsing \"\": invalid syntax", "invalid input: empty byte slice"},
	}
	for _, c := range vv {
		t.Run(fmt.Sprintf("uint case %d", c.e), func(t *testing.T) {
			v, err := strconv.Atoi(string(c.v))
			if err != nil && c.err != err.Error() {
				t.Errorf("got error %+v, want %v", err, c.err)
			}
			i, err2 := bytesToInt(c.v)
			if err2 != nil && c.errCustom != err2.Error() {
				t.Errorf("bytesToUint: got error %+v, want %v", err2, c.errCustom)
			} else if v != i {
				t.Errorf("bytesToUint: got %v, want %v", v, c.v)
			}
			if c.err != "" {
				return
			}

			intBytes := NewInt(c.e)
			if !bytes.Equal(intBytes.ToBytes(), c.v) {
				t.Errorf("uintToBytes failed: got %v, want %v", intBytes, c.v)
			}
		})
	}

}

func TestFormatToFloatStrConv(t *testing.T) {
	vv := []testFloatCase{
		{[]byte("0"), 0, "", ""},
		{[]byte("0.0000"), 0, "", ""},
		{[]byte("123"), 123, "", ""},
		{[]byte("123.456"), 123.456, "", ""},
		{[]byte("0.456"), 0.456, "", ""},
		{[]byte("0.1"), 0.1, "", ""},
		{[]byte("0.0000001"), 0.0000001, "", ""},
		{[]byte("-0.0000001"), -0.0000001, "", ""},
		{[]byte("0.123123"), 0.123123, "", ""},
		{[]byte("0.12312321312312"), 0.12312321312312, "", ""},
		{[]byte("0.12312999999999922"), 0.12312999999999922, "", ""},
		{[]byte("0.123129999999999221"), 0.123129999999999221, "", ""},
		{[]byte("1"), 1, "", ""},
		{[]byte("1797693134862315"), 1.797693134862315e+15, "", ""},
		{[]byte("17976931348623157"), 17976931348623157, "", ""},
		{[]byte("179769313486231574"), 179769313486231574, "", ""},
		{[]byte("3.402823466385288598e+10"), 3.402823466385288598e+10, "", ""},
		{[]byte("3.40282346638528859811704183484516925440e+10"), 3.40282346638528859811704183484516925440e+10, "", ""},
		{[]byte("3.40282346638528859811704183484516925440e+38"), math.MaxFloat32, "", ""},
		{[]byte("1.401298464324817e-45"), math.SmallestNonzeroFloat32, "", ""},
		{[]byte("1.79769313486231570814527423731704356798070e+308"), math.MaxFloat64, "", ""},
		{[]byte("1.79769313486231570814527423731704356798070e-308"), 1.79769313486231570814527423731704356798070e-308, "", ""},
		{[]byte("5e-324"), math.SmallestNonzeroFloat64, "", ""},
		{[]byte("179769313486231574112351123123"), 179769313486231574112351123123, "", ""},
		{[]byte(""), 0,
			"strconv.ParseFloat: parsing \"\": invalid syntax",
			"invalid syntax: empty string"},
		{[]byte("test"), 0,
			"strconv.ParseFloat: parsing \"test\": invalid syntax",
			"invalid syntax: invalid character"},
		{[]byte("."), 0,
			"strconv.ParseFloat: parsing \".\": invalid syntax",
			"invalid syntax: unparsable tail left"},
		{[]byte(".."), 0,
			"strconv.ParseFloat: parsing \"..\": invalid syntax",
			"invalid syntax"},
		{[]byte("1e1"), 10, "", ""},
		{[]byte("1e2"), 100, "", ""},
		{[]byte("1e-2"), 0.01, "", ""},
		{[]byte("1E-5"), 1e-05, "", ""},
		{[]byte("1E+16"), 1e+16, "", ""},
		{[]byte("1E-16"), 1e-16, "", ""},
		{[]byte("1E+32"), 1e+32, "", ""},
		{[]byte("1E-32"), 1e-32, "", ""},
		{[]byte("1E+38"), 1e+38, "", ""},
		{[]byte("1E300"), 1e+300, "", ""},
		{[]byte("1E-300"), 1e-300, "", ""},
		{[]byte("0."), 0, "", ""},
		{[]byte(".1"), 0.1, "", ""},
	}

	tstr := ""
	for _, c := range vv {
		t.Run(fmt.Sprintf("float case %s", c.v), func(t *testing.T) {
			str := fmt.Sprintf("%s,", string(c.v))
			v, err := strconv.ParseFloat(string(c.v), 64)
			if v != c.e {
				str += "-,"
				str += fmt.Sprintf("%f,", v)
				t.Errorf("got %v, want %v", v, c.e)
			} else {
				str += "+,"
				str += fmt.Sprintf("%v,", v)
			}
			if c.err != "" {
				if err == nil || c.err != err.Error() {
					t.Errorf("got %+v, want %v", err, c.err)
				}
			} else if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			v2, err := bytesToFloat(c.v)
			if v2 != c.e {
				t.Errorf("custom convert got %v, want %v", v2, c.e)
			}

			if c.err != "" {
				if err == nil || c.errCustom != err.Error() {
					t.Errorf("got %v, want %v", err, c.err)
				}
			} else if err != nil {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
	fmt.Println(tstr)

}

// Benchmark_UntToByteAppend-24    	147965528	         8.104 ns/op
func Benchmark_UntToByteAppend(b *testing.B) {
	v := uint64(1123123)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.AppendUint(make([]byte, 0, 20), v, 10)
	}
}

// Benchmark_UntToByteFormat-24    	74329178	        14.91 ns/op
func Benchmark_UntToByteFormat(b *testing.B) {
	v := uint64(1123123)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(strconv.FormatUint(v, 10))
	}
}

func Test_UntToByteAppend(t *testing.T) {
	vv := []uint64{0, 123, 123456, 1234567890, math.MaxUint64}
	for _, v := range vv {
		bytesVal := strconv.AppendUint(make([]byte, 0, 20), v, 10)
		if string(bytesVal) != strconv.FormatUint(v, 10) {
			t.Errorf("got %v, want %v", string(bytesVal), strconv.FormatUint(v, 10))
		}
	}
}

func Test_IntToByteAppend(t *testing.T) {
	vv := []int64{-1, 0, 123, 123456, 1234567890}
	for _, v := range vv {
		bytesVal := strconv.AppendInt(make([]byte, 0, 20), v, 10)
		if string(bytesVal) != strconv.FormatInt(v, 10) {
			t.Errorf("got %v, want %v", string(bytesVal), strconv.FormatInt(v, 10))
		}
	}
}

// Benchmark_IntToByte-24    	41440184	        26.98 ns/op
func Benchmark_IntToByte(b *testing.B) {
	vv := -123123123
	for i := 0; i < b.N; i++ {
		_ = strconv.AppendInt(make([]byte, 20), int64(vv), 10)
	}
}

// Benchmark_Itoa-24    	60535741	        18.13 ns/op
func Benchmark_Itoa(b *testing.B) {
	vv := -123123123
	for i := 0; i < b.N; i++ {
		_ = []byte(strconv.Itoa(vv))
	}
}

// Benchmark_ByteToInt-24    	283137664	         4.230 ns/op
func Benchmark_ByteToInt(b *testing.B) {
	vv := []byte("-123123123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bytesToInt(vv)
	}
}

// Benchmark_AtoiIntegerBytes-24    	161953503	         7.714 ns/op
func Benchmark_AtoiIntegerBytes(b *testing.B) {
	vv := []byte("-123123123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = strconv.Atoi(string(vv))
	}
}

func Test_BytesToInt(t *testing.T) {
	vv := [][]byte{
		[]byte(strconv.FormatInt(-500, 10)),
		[]byte(strconv.FormatInt(0, 10)),
		[]byte(strconv.FormatInt(123, 10)),
		[]byte(strconv.FormatInt(math.MaxInt, 10))}
	for _, v := range vv {
		bytesVal, err := bytesToInt(v)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
		in, _ := strconv.Atoi(string(v))
		if bytesVal != in {
			t.Errorf("got %v, want %v", bytesVal, string(v))
		}
	}
}

// Benchmark_TimeToBytes-24			21122396	54.07 ns/op
func Benchmark_TimeToBytes(b *testing.B) {
	v := time.Now()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = timeToBytes(v)
	}
}

// Benchmark_TimeFormatToBytes-24	10808502	109.3 ns/op
func Benchmark_TimeFormatToBytes(b *testing.B) {
	v := time.Now()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = []byte(v.Format(TimeLayout))
	}
}

func TestTimeToBytes(t *testing.T) {

	cases := []time.Time{
		time.Now(),
		time.Now().Add(time.Hour),
		time.Now().Add(time.Hour * 6516),
	}
	for _, v := range cases {
		bytesVal := timeToBytes(v)
		if string(bytesVal) != v.Format(TimeLayout) {
			t.Errorf("got %v, want %v", string(bytesVal), v.Format(TimeLayout))
		}
	}

	bytesVal := timeToBytes(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	if string(bytesVal) != "00001130-00:00:00.000" {
		t.Errorf("got %v, want %v", string(bytesVal), "00001130-00:00:00.000")
	}
}

func TestStringWrite(t *testing.T) {
	cases := []string{
		"",
		"as ",
	}
	casesB := [][]byte{
		[]byte(""),
		[]byte("as "),
	}
	for i, v := range cases {
		buff := bytes.NewBuffer(make([]byte, 0, 200))
		NewString(v).WriteBytes(buff)
		if !bytes.Equal(buff.Bytes(), casesB[i]) {
			t.Errorf("%s got %v, want %v", v, buff.Bytes(), casesB[i])
		}
	}
}
