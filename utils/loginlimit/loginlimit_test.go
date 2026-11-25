package loginlimit

import (
	"testing"
	"time"
)

func TestLimiterBlocksAfterMaxFailures(t *testing.T) {
	l := New(3, time.Minute)
	key := "127.0.0.1|admin"

	for i := 0; i < 3; i++ {
		if !l.Allowed(key) {
			t.Fatalf("attempt %d should be allowed", i+1)
		}
		l.Fail(key)
	}

	if l.Allowed(key) {
		t.Fatal("should be blocked after max failures")
	}
	if l.Remaining(key) <= 0 {
		t.Fatal("expected positive remaining duration")
	}

	l.Reset(key)
	if !l.Allowed(key) {
		t.Fatal("reset should clear the counter")
	}
}

func TestLimiterExpires(t *testing.T) {
	l := New(1, 20*time.Millisecond)
	key := "k"

	l.Fail(key)
	if l.Allowed(key) {
		t.Fatal("should be blocked")
	}
	time.Sleep(30 * time.Millisecond)
	if !l.Allowed(key) {
		t.Fatal("should be allowed after window expires")
	}
}
