package joblist

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"
	"golang.org/x/sync/errgroup"

	"github.com/evertras/khan/internal/repository"
)

func refreshJobsCmd() tea.Msg {
	c := repository.GetNomadClient()

	jobs, _, err := c.Jobs().List(&api.QueryOptions{})

	if err != nil {
		return errMsg(fmt.Errorf("failed to get job list: %w", err))
	}

	return jobs
}

func garbageCollectCmd() tea.Msg {
	client := repository.GetNomadClient()

	err := client.System().GarbageCollect()

	if err != nil {
		return errMsg(err)
	}

	// TODO: Figure out nicer load notifications
	time.Sleep(time.Millisecond * 200)

	return refreshJobsCmd()
}

func stopSelectedCmd(ids []string) func() tea.Msg {
	return func() tea.Msg {
		eg := errgroup.Group{}

		for _, id := range ids {
			eg.Go(func() error {
				client := repository.GetNomadClient()

				_, _, err := client.Jobs().Deregister(id, false, &api.WriteOptions{})

				return err
			})
		}

		err := eg.Wait()

		if err != nil {
			return errMsg(err)
		}

		// TODO: Figure out nicer load notifications
		time.Sleep(time.Millisecond * 200)

		return refreshJobsCmd()
	}
}

type logStreamingMsg struct {
	jobID     string
	allocID   string
	taskGroup string
	task      string
	logs      <-chan *api.StreamFrame
	errs      <-chan error
}

func showLogsForJobCmd(jobID string) func() tea.Msg {
	return func() tea.Msg {
		client := repository.GetNomadClient()

		job, _, err := client.Jobs().Info(jobID, &api.QueryOptions{})

		if err != nil {
			return errMsg(fmt.Errorf("failed to get job info for %q: %w", jobID, err))
		}

		if len(job.TaskGroups) == 0 || len(job.TaskGroups[0].Tasks) == 0 {
			return errMsg(fmt.Errorf("job %q has no tasks", jobID))
		}

		taskGroup := ""

		if job.TaskGroups[0].Name != nil {
			taskGroup = *job.TaskGroups[0].Name
		}

		// Default to first regardless of lifecycle hook
		taskName := job.TaskGroups[0].Tasks[0].Name

		// Try and find one without a lifecycle hook to use
		for _, task := range job.TaskGroups[0].Tasks {
			if task.Lifecycle == nil {
				taskName = task.Name
				break
			}
		}

		allocs, _, err := client.Jobs().Allocations(jobID, false, &api.QueryOptions{})

		if err != nil {
			return errMsg(fmt.Errorf("failed to get alloc list for job %q: %w", jobID, err))
		}

		numAllocs := len(allocs)

		if numAllocs == 0 {
			return errMsg(fmt.Errorf("no allocs found for job"))
		}

		randomAlloc := allocs[rand.Intn(numAllocs)]

		alloc, _, err := client.Allocations().Info(randomAlloc.ID, &api.QueryOptions{})

		if err != nil {
			return errMsg(fmt.Errorf("failed to get alloc %q: %w", randomAlloc.ID, err))
		}

		logMu.Lock()
		if logCancel != nil {
			logCancel()
		}

		logCtx, logCancel = context.WithCancel(context.Background())

		streamCh, errCh := client.AllocFS().Logs(alloc, true, taskName, "stdout", "start", 0, logCtx.Done(), &api.QueryOptions{})
		logMu.Unlock()

		return logStreamingMsg{
			jobID:     jobID,
			allocID:   alloc.ID,
			taskGroup: taskGroup,
			task:      taskName,
			logs:      streamCh,
			errs:      errCh,
		}
	}
}

func inspectJobCmd(jobID string) func() tea.Msg {
	return func() tea.Msg {
		client := repository.GetNomadClient()

		job, _, err := client.Jobs().Info(jobID, &api.QueryOptions{})

		if err != nil {
			return errMsg(err)
		}

		return job
	}
}
