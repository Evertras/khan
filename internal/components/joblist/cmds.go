package joblist

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

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

	time.Sleep(time.Millisecond * 20)

	return refreshJobsCmd()
}
