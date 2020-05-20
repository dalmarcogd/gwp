package worker

import (
	"context"
	"testing"
	"time"
)

func TestWithConcurrency(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})
	concurrency := WithConcurrency(2)
	concurrency.Apply(w)
	if w.Concurrency != 2 {
		t.Errorf("Expected 2 concurrency but received %v", w.Concurrency)
	}
}

func TestWithDeadline(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})
	now := time.Now().Add(2 * time.Second)
	deadline := WithDeadline(now)
	deadline.Apply(w)
	if w.Deadline != now {
		t.Errorf("Expected %v deadline but received %v", now, w.Deadline)
	}
}

func TestWithRestartAlways(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})
	restartAlways := WithRestartAlways()
	restartAlways.Apply(w)
	if !w.RestartAlways {
		t.Errorf("Expected restart always but received %t", w.RestartAlways)
	}
}

func TestWithTimeout(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})
	duration := 2 * time.Second
	timeout := WithTimeout(duration)
	timeout.Apply(w)
	if w.Timeout != duration {
		t.Errorf("Expected %v deadline but received %v", duration, w.Timeout)
	}
}

func TestWithCron(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})
	duration := 2 * time.Second
	cron := WithCron(duration)
	cron.Apply(w)
	if w.Cron != duration {
		t.Errorf("Expected %v deadline but received %v", duration, w.Cron)
	}
}

func Test_workerConfig_Apply(t *testing.T) {
	config := Config{}
	config.Apply(nil)
	config.Apply(NewWorker("test", func(ctx context.Context) error {
		return nil
	}))
}

func TestFullWith(t *testing.T) {
	w := NewWorker("test", func(ctx context.Context) error {
		return nil
	})

	duration := 5 * time.Second
	timeout := WithTimeout(duration)
	timeout.Apply(w)

	concurrency := WithConcurrency(8)
	concurrency.Apply(w)

	now := time.Now().Add(80 * time.Second)
	deadline := WithDeadline(now)
	deadline.Apply(w)

	restartAlways := WithRestartAlways()
	restartAlways.Apply(w)

	if w.Deadline != now {
		t.Errorf("Expected %v deadline but received %v", now, w.Deadline)
	}

	if !w.RestartAlways {
		t.Errorf("Expected restart always but received %t", w.RestartAlways)
	}

	if w.Concurrency != 8 {
		t.Errorf("Expected 8 concurrency but received %v", w.Concurrency)
	}

	if w.Timeout != duration {
		t.Errorf("Expected %v deadline but received %v", duration, w.Timeout)
	}
}
