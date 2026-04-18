package action

import "symlinker/core/msg"

type Runner struct {
    IsDryRun  bool
    IsVerbose bool
}

func (r *Runner) Do(description string, action func() error) error {
    if r.IsVerbose || r.IsDryRun {
        msg.Info(description)
    }

    if r.IsDryRun {
        return nil
    }
    return action()
}
