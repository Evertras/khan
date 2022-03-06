package logs

import (
	"context"
	"sync"
)

var (
	logCtx    context.Context
	logCancel context.CancelFunc
	logMu     sync.Mutex
)

func cancelExistingLogStream() {
	logMu.Lock()
	if logCancel != nil {
		logCancel()
		logCancel = nil
	}
	logMu.Unlock()
}

type logReceivedMsg struct {
	logStreamingMsg

	receivedLogChunk string
}
