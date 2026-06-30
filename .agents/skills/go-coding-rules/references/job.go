package references

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"

	"github.com/reshap0318/go-project/internal/helpers"
	"github.com/reshap0318/go-project/internal/services"
)

// ============================================================
// 00_jobs.go — Jobs struct, type constants, Register
// ============================================================

// Task type constants — used by both scheduler (cmd/worker/main.go)
// and server mux (Register). ALWAYS define here, never inline strings.
const (
	TypeClearTmp = "cleartmp"
	// TypeMyJob = "myjob"  ← add new types here
)

// Jobs holds all background job handlers.
// Follows the same single-struct pattern as Handlers.
type Jobs struct {
	svcs *services.Services
}

// NewJobs creates a new Jobs instance.
func NewJobs(svcs *services.Services) *Jobs {
	return &Jobs{svcs: svcs}
}

// Register wires all job handlers to the asynq ServeMux.
// EVERY type constant MUST have a corresponding HandleFunc here.
func (j *Jobs) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeClearTmp, j.HandleClearTmp)
	// mux.HandleFunc(TypeMyJob, j.HandleMyJob)  ← add alongside the constant above
}

// ============================================================
// {feature}_job.go — one handler per file
// ============================================================

// HandleClearTmp processes the cleartmp job.
// Deletes files in storage/tmp older than TMP_AGE_HOURS (default 24h).
func (j *Jobs) HandleClearTmp(ctx context.Context, t *asynq.Task) error {
	dir := helpers.GetEnv("TMP_DIR", "storage/tmp")
	ageHours := helpers.GetEnvInt("TMP_AGE_HOURS", 24)
	threshold := time.Now().Add(-time.Duration(ageHours) * time.Hour)

	j.svcs.Logger.LogStart("HandleClearTmp", "Scanning %s (files older than %dh)", dir, ageHours)

	deleted, errCount := 0, 0

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			errCount++
			return nil // skip bad entries, don't abort the walk
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			errCount++
			return nil
		}

		if info.ModTime().After(threshold) {
			return nil
		}

		if err := helpers.DeleteFile(path); err != nil {
			j.svcs.Logger.LogWarn("HandleClearTmp", "Failed to delete %s: %v", path, err)
			errCount++
			return nil
		}

		deleted++
		j.svcs.Logger.LogStep("HandleClearTmp", "Deleted: %s", path)
		return nil
	})

	if err != nil {
		j.svcs.Logger.LogEndWithError("HandleClearTmp", "Walk failed: %v", err)
		return fmt.Errorf("%s: walk %s: %w", TypeClearTmp, dir, err)
	}

	j.svcs.Logger.LogEnd("HandleClearTmp", "Done — deleted: %d, errors: %d", deleted, errCount)
	return nil
}
