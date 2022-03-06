package list

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/khan/internal/repository"
	"github.com/hashicorp/nomad/api"
	"golang.org/x/sync/errgroup"
)

func refreshJobsCmd() tea.Msg {
	c := repository.GetNomadClient()

	jobs, _, err := c.Jobs().List(&api.QueryOptions{})

	if err != nil {
		return fmt.Errorf("failed to get job list: %w", err)
	}

	return jobs
}

func garbageCollectCmd() tea.Msg {
	client := repository.GetNomadClient()

	err := client.System().GarbageCollect()

	if err != nil {
		return err
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
			return err
		}

		// TODO: Figure out nicer load notifications
		time.Sleep(time.Millisecond * 200)

		return refreshJobsCmd()
	}
}

func inspectJobCmd(jobID string) func() tea.Msg {
	return func() tea.Msg {
		client := repository.GetNomadClient()

		job, _, err := client.Jobs().Info(jobID, &api.QueryOptions{})

		if err != nil {
			return err
		}

		return job
	}
}
