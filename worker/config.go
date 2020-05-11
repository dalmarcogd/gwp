package worker

import "time"

//Apply the setting at worker
func (c *Config) Apply(w *Worker) {
	if c.k != nil && w != nil {
		c.k(w)
	}
}

//WithConcurrency generate a setup to worker with specify concurrency
func WithConcurrency(c int) Config {
	return Config{k: func(w *Worker) { w.Concurrency = c }}
}

//WithRestartAlways generate a setup to worker with restart always
func WithRestartAlways() Config {
	return Config{k: func(w *Worker) { w.RestartAlways = true }}
}

//WithTimeout generate a setup to worker with specify timeout
func WithTimeout(t time.Duration) Config {
	return Config{k: func(w *Worker) { w.Timeout = t }}
}

//WithDeadline generate a setup to worker with specify deadline
func WithDeadline(t time.Time) Config {
	return Config{k: func(w *Worker) { w.Deadline = t }}
}
