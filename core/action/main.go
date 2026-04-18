package action

import "symlinker/core/msg"

type Runner struct {
    IsDryRun  bool
    IsVerbose bool
}

func (r *Runner) Do(description string, action func() error) error {
    if r.IsVerbose || r.IsDryRun {
        msg.Info("Action: " + description)
    }

    if r.IsDryRun {
        msg.Warn("-> Skipping (Dry Run)")
        return nil
    }
    return action()
}
