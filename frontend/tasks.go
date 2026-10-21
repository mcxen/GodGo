package frontend

import (
	"context"
	"time"

	"godgo/backend/workflow"
)

func (a *App) startTask(task workflow.Task) {
	if a.running[task.ID] {
		return
	}
	a.running[task.ID] = true
	root := a.root
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		res := workflow.Execute(ctx, task, root)
		a.update(func() {
			delete(a.running, task.ID)
			if a.runs == nil {
				a.runs = map[string]workflow.Result{}
			}
			a.runs[task.ID] = res
			if a.setWin != nil {
				a.setWin.Invalidate()
			}
		})
	}()
}
