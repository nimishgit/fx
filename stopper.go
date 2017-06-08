// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package fx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Stopper interface {
	TakeOver() (chan<- Reason, error)
}

type Reason string

type stopper struct {
	trigger       chan Reason
	signalTrigger chan os.Signal
}

func newStopper() *stopper {
	return &stopper{}
}

func (s *stopper) start() error {
	if s.trigger != nil {
		return nil
	}
	s.trigger = make(chan Reason, 1)
	s.signalTrigger = make(chan os.Signal, 1)
	signal.Notify(s.signalTrigger, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-s.signalTrigger
			if sig == nil { // exit on channel closed.
				break
			}
			s.trigger <- Reason(sig.String())
		}
	}()
	return nil
}

func (s *stopper) stop() error {
	if s.signalTrigger != nil {
		signal.Stop(s.signalTrigger)
	}
	return nil
}

func (s *stopper) stopped() <-chan Reason {
	return s.trigger
}

func (s *stopper) TakeOver() (chan<- Reason, error) {
	if s.trigger != nil {
		return nil, fmt.Errorf("stopper already started")
	}
	s.trigger = make(chan Reason, 1)
	return s.trigger, nil
}
