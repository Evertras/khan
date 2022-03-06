package nodes

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"

	"github.com/evertras/khan/internal/screens"
)

func TestEmptyViewDoesntPanic(t *testing.T) {
	m := New(screens.Size{})

	m.View()
}

func TestModelWithNodesShowsAllNodeNames(t *testing.T) {
	nodes := []*api.NodeListStub{
		{
			Name: "hello-node",
		},
		{
			Name: "another-node",
		},
	}

	m := New(screens.Size{})

	updated, _ := m.Update(nodes)

	m = updated.(Model)

	view := m.View()

	assert.Contains(t, view, "hello-node", "Missing first node")
	assert.Contains(t, view, "another-node", "Missing second node")
}
