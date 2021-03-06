// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package timeout

import (
	"time"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/hashing"
	"github.com/ava-labs/gecko/utils/timer"
	"github.com/ava-labs/gecko/utils/wrappers"
	"github.com/prometheus/client_golang/prometheus"
)

// Manager registers and fires timeouts for the snow API.
type Manager struct{ tm timer.AdaptiveTimeoutManager }

// Initialize this timeout manager.
func (m *Manager) Initialize(
	namespace string,
	registerer prometheus.Registerer,
) error {
	return m.tm.Initialize(
		time.Second,
		500*time.Millisecond,
		2,
		time.Millisecond,
		namespace,
		registerer,
	)
}

// Dispatch ...
func (m *Manager) Dispatch() { m.tm.Dispatch() }

// Register request to time out unless Manager.Cancel is called
// before the timeout duration passes, with the same request parameters.
func (m *Manager) Register(validatorID ids.ShortID, chainID ids.ID, requestID uint32, timeout func()) time.Time {
	return m.tm.Put(createRequestID(validatorID, chainID, requestID), timeout)
}

// Cancel request timeout with the specified parameters.
func (m *Manager) Cancel(validatorID ids.ShortID, chainID ids.ID, requestID uint32) {
	m.tm.Remove(createRequestID(validatorID, chainID, requestID))
}

func createRequestID(validatorID ids.ShortID, chainID ids.ID, requestID uint32) ids.ID {
	p := wrappers.Packer{Bytes: make([]byte, wrappers.IntLen)}
	p.PackInt(requestID)

	return ids.NewID(hashing.ByteArraysToHash256Array(validatorID.Bytes(), chainID.Bytes(), p.Bytes))
}
