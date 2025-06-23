package envvar

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type syncwriter struct {
	buf bytes.Buffer
	m   sync.Mutex
}

func (w *syncwriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.buf.Write(b)
}

// NOTE: don't make this test parallel!
// this test deliberately doesn't run in CI since it's darned slow.
func Test_LogBatching(t *testing.T) {
	const key = "TEST_LOG_BATCHING"
	// we can't use pkg/tests.SkipLongf() or pkg/tests.Skipf, since they both import pkg/envvar
	if testing.Short() {
		t.Skipf("SKIP: %s: this is slow and sleeps for 70 seconds!", t.Name())
	}
	if b, _ := lookupBool(key); !b {
		t.Skipf("SKIP: %s: this test is too slow to run by default. set the environment variable %s=true to enable", t.Name(), key)
	}
	t.Log("starting log batching test - this is slow and sleeps for 70 seconds!")
	w := new(syncwriter)

	// we don't want a race where we move the logger underneath the batching goroutine, so we need the lock.
	logMux.Lock()
	oldLogger := logger
	logger = zerolog.New(w)
	logMux.Unlock()
	defer func() {
		logMux.Lock()
		logger = oldLogger
		logMux.Unlock()
	}()

	for i := 'A'; i <= 'Z'; i++ {
		key := fmt.Sprintf("_ENV_TEST_TEST_LOGBATCHING_%c", i)
		GetString(key, fmt.Sprintf("%c", i))
	}
	for i := 7; i > 0; i-- {
		t.Logf("%2d seconds remaining\n", 10*i)
		time.Sleep(10 * time.Second)
	}
	w.m.Lock()
	got := w.buf.String()
	w.m.Unlock()
	if got == "" {
		t.Fatalf("expected at least one batch of logs, but got none")
	}
	if lines := strings.Count(got, "\n"); lines > 5 {
		t.Fatalf("expected logs to be batched into 5 or fewer entries, but got %d", lines)
	}

}
