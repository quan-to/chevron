package SLog

import (
	"fmt"
	"testing"
)

// Crash tests

func TestSetTestMode(t *testing.T) {
	SetTestMode()
	// Test Mode should be all unset

	if debugEnabled {
		t.Fatalf("Debug is set to true! Should be false")
	}

	if warnEnabled {
		t.Fatalf("Warn is set to true! Should be false")
	}

	if errorEnabled {
		t.Fatalf("Error is set to true! Should be false")
	}

	if infoEnabled {
		t.Fatalf("Info is set to true! Should be false")
	}

	UnsetTestMode()
	// Test Mode should be all set

	if !debugEnabled {
		t.Fatalf("Debug is set to false! Should be true")
	}

	if !warnEnabled {
		t.Fatalf("Warn is set to false! Should be true")
	}

	if !errorEnabled {
		t.Fatalf("Error is set to false! Should be true")
	}

	if !infoEnabled {
		t.Fatalf("Info is set to false! Should be true")
	}
}

func TestDebug(t *testing.T) {
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestError(t *testing.T) {
	Error("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestWarn(t *testing.T) {
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestInfo(t *testing.T) {
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLog(t *testing.T) {
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLogNoFormat(t *testing.T) {
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestFatal(t *testing.T) {
	assertPanic(t, func() {
		Fatal("Test Fatal")
	}, "Fatal should panic")
}

func TestScope(t *testing.T) {
	scoped := Scope("test-scope")
	if scoped.scope != "test-scope" {
		t.Fatalf("Expected test-scope got %s", scoped.scope)
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	if !debugEnabled {
		t.Fatalf("Debug is set to false! Should be true")
	}
	SetDebug(false)
	if debugEnabled {
		t.Fatalf("Debug is set to true! Should be false")
	}
}

func TestSetError(t *testing.T) {
	SetError(true)
	if !errorEnabled {
		t.Fatalf("Error is set to false! Should be true")
	}
	SetError(false)
	if errorEnabled {
		t.Fatalf("Error is set to true! Should be false")
	}
}

func TestSetInfo(t *testing.T) {
	SetInfo(true)
	if !infoEnabled {
		t.Fatalf("Info is set to false! Should be true")
	}
	SetInfo(false)
	if infoEnabled {
		t.Fatalf("Info is set to true! Should be false")
	}
}

func TestSetWarn(t *testing.T) {
	SetWarning(true)
	if !warnEnabled {
		t.Fatalf("Warn is set to false! Should be true")
	}
	SetWarning(false)
	if warnEnabled {
		t.Fatalf("Warn is set to true! Should be false")
	}
}

type test struct{}

func (test) String() string {
	return "test"
}

func TestAsString(t *testing.T) {
	var tstringcast StringCast

	tstringcast = &test{}

	tests := []interface{}{
		"string",
		123456,
		123456.1,
		true,
		map[string]string{},
		[]int{1, 2, 3, 4, 5},
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
		fmt.Errorf("error format"),
		tstringcast,
	}

	outputs := make([]string, len(tests))

	for i, v := range tests { // Fill tests
		outputs[i] = fmt.Sprint(v) // Should be same output
	}

	for i, v := range tests {
		s := asString(v)
		if s != outputs[i] {
			t.Errorf("#%d expected %s got %s.", i, outputs[i], s)
		}
	}
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}
