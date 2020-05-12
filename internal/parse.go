package internal

import (
	"runtime"
	"strconv"
	"time"
)

func ParseServerInfos(s Server) map[string]interface{} {
	infos := map[string]interface{}{
		"cpus":       strconv.Itoa(runtime.NumCPU()),
		"goroutines": strconv.Itoa(runtime.NumGoroutine()),
		"workers":    []map[string]interface{}{},
	}

	for _, w := range s.Workers() {
		finishedAt := ""
		if !w.FinishedAt.IsZero() {
			finishedAt = w.FinishedAt.Format(time.RFC3339)
		}
		infos["workers"] = append(infos["workers"].([]map[string]interface{}), map[string]interface{}{
			"id":             w.ID,
			"name":           w.Name,
			"concurrency":    w.Concurrency,
			"restart_always": w.RestartAlways,
			"restarts":       w.Restarts,
			"started_at":     w.StartAt.Format(time.RFC3339),
			"finished_at":    finishedAt,
			"status":         w.Status(),
		})
	}
	return infos
}
