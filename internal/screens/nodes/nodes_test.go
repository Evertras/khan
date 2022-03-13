package nodes

import (
	"testing"

	"github.com/evertras/khan/internal/screens"
)

func TestEmptyViewDoesntPanic(t *testing.T) {
	m := New(screens.Size{})

	m.View()
}
