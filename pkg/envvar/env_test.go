package envvar

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func set(key string, val interface{}) { os.Setenv(key, format(val)) }
func unset(key string)                { os.Unsetenv(key) }
func TestMain(m *testing.M) {
	logger = zerolog.New(io.Discard)
	os.Exit(m.Run())
}
func format(v interface{}) string {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("unknown type %T", v))
	case time.Time:
		b, err := v.MarshalText()
		if err != nil {
			panic(err)
		}
		return string(b)
	case net.IP:
		b, err := v.MarshalText()
		if err != nil {
			panic(err)
		}
		return string(b)
	case time.Duration:
		return v.String()
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case []string:
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
}

func Test_Bool_True(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_BOOL_TRUE"
	for _, s := range []string{"1", "t", "T", "TRUE", "true", "True"} {
		set(k, s)
		if b, err := lookupBool(k); err != nil || !b {
			t.Fail()
		}
		unset(s)
	}

}

func TestBool_False(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_BOOL_FALSE"
	for _, s := range []string{"0", "f", "F", "FALSE", "false", "False"} {
		set(k, s)
		if b, err := lookupBool(k); err != nil || b {
			t.Fail()
		}

		unset(k)
	}
}

func Test_Duration(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_DURATION"

	for _, tt := range []time.Duration{
		10 * time.Nanosecond,
		5 * time.Millisecond,
		-3 * time.Second,
		24 * time.Hour,
	} {
		set(k, tt)
		if got := GetDuration(k, 0); got != tt {
			t.Fatal("getduration")
		}
	}
	for _, bad := range []string{"twenty-four-hours", "auasda", ""} {
		set(k, bad)
		if got, err := lookupDuration(k); got != 0 || err == nil {
			t.Fatal("getduration err - expected no error")
		}
	}
}

func Test_Time(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_TIME"
	someDay := time.Date(2002, time.January, 1, 2, 3, 4, 5, time.UTC)
	set(k, someDay)
	if got := GetTime(k, time.Now()); !got.Equal(someDay) {
		t.Fatal("get time")
	}
}

func TestInt(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_INT"
	if GetInt(k, 2) != 2 {
		t.Fatal("default failed")
	}
	if !testing.Short() {
		for i := 0; i < 0xff_ff; i++ {
			want := rand.Int()
			if i%2 == 0 {
				want = -want
			}
			set(k, want)
			if GetInt(k, 0) != want {
				t.Fatal("getint")
			}
			unset(k)
		}
	}
	for _, s := range []string{"foo", "2.1", "a9sdasd", "--"} {
		set(k, s)
		if got := GetInt(k, -1); got != -1 {
			t.Fatalf("expected -1 from %v, but got %v", s, got)
		}
	}

}

func TestFloat(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_FLOAT"

	unset(k)
	if GetFloat(k, 2.2) != 2.2 {
		t.Fatal("default failed")
	}
	for _, s := range []string{"foo", "", "a9sdasd", "--"} {
		set(k, s)
		if got := GetFloat(k, 0.0); got != 0.0 {
			t.Fatalf("expected 0 from %v, but got %v", s, got)
		}
	}
	if !testing.Short() {
		for i := 0; i < 0xff_ff; i++ {
			want := rand.ExpFloat64()
			if i%2 == 0 {
				want = -want
			}
			set(k, want)
			if GetFloat(k, 0.0) != want {
				t.Fatal("getint")
			}
			unset(k)
		}
	}

}

func TestString(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_String"

	unset(k)
	hamlet := "Alas, poor Yorick! I knew him, Horatio."
	if GetString(k, hamlet) != hamlet {
		t.Fatal("default failed")
	}
	for _, s := range []string{"foo", "", "a9sdasd", "--"} {
		set(k, s)
		if got := GetString(k, hamlet); got != s {
			t.Fatalf("expected %s from %v, but got hamlet", s, got)
		}
	}
}

func TestPanics(t *testing.T) {
	t.Parallel()
	const k = "_ENV_TEST_TEST_PANIC"
	for name, f := range map[string]func(){
		"bool-f":   func() { MustGetBool(k) },
		"int-f":    func() { MustGetInt(k) },
		"float-f":  func() { MustGetFloat(k) },
		"string-f": func() { MustGetString(k) },
		"dur-f":    func() { MustGetDuration(k) },
		"time-f":   func() { MustGetTime(k) },
		"ip-f":     func() { MustGetIP(k) },
	} {
		name, t := name, t
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if p := recover(); p == nil {
					t.Fatal("should have panicked")
				}
			}()
			f()
		})
	}
}
