// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

/*
Package posixsignal provides a listener for a posix signal. By default
it listens for SIGINT and SIGTERM, but others can be chosen in NewPosixSignalManager.
When ShutdownFinish is called it exits with os.Exit(0)
*/
package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/marmotedu/iam/pkg/shutdown"
)

// Name defines shutdown manager name.
const Name = "PosixSignalManager"

// PosixSignalManager implements ShutdownManager interface that is added
// to GracefulShutdown. Initialize with NewPosixSignalManager.
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager.
// As arguments you can provide os.Signal-s to listen to, if none are given,
// it will default to SIGINT and SIGTERM.
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

// GetName returns name of this ShutdownManager.
func (posixSignalManager *PosixSignalManager) GetName() string {
	return Name
}

// Start starts listening for posix signals.
func (posixSignalManager *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, posixSignalManager.signals...)

		// Block until a signal is received.
		<-c

		gs.StartShutdown(posixSignalManager)
	}()

	return nil
}

// ShutdownStart does nothing.
func (posixSignalManager *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (posixSignalManager *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
