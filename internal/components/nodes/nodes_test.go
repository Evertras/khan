package nodes

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func TestEmptyViewDoesntCrash(t *testing.T) {
	m := NewEmptyModel()

	str := m.View()

	assert.Empty(t, str)
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

	m := NewModelWithNodes(nodes)

	view := m.View()

	assert.Contains(t, view, "hello-node", "Missing first node")
	assert.Contains(t, view, "another-node", "Missing second node")
}
