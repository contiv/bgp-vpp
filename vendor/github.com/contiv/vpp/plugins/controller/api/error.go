// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

/********************************* Fatal Error ********************************/

// FatalError tells Controller to abort the event loop and stop the agent
// as soon as possible.
type FatalError struct {
	origErr error
}

// NewFatalError is the constructor for FatalError.
func NewFatalError(origErr error) error {
	return &FatalError{origErr: origErr}
}

// Error delegates the call to the underlying error.
func (e *FatalError) Error() string {
	return e.origErr.Error()
}

// GetOriginalError returns the underlying error.
func (e *FatalError) GetOriginalError() error {
	return e.origErr
}

/****************************** Abort Event Error *****************************/

// AbortEventError tells controller to abort the processing of the event
// (and for Update/RevertOnFailure to revert the changes). The agent does not have
// to restart but resync should be performed as soon as possible.
type AbortEventError struct {
	origErr error
}

// NewAbortEventError is the constructor for the AbortEventError.
func NewAbortEventError(origErr error) error {
	return &AbortEventError{origErr: origErr}
}

// Error delegates the call to the underlying error.
func (e *AbortEventError) Error() string {
	return e.origErr.Error()
}

// GetOriginalError returns the underlying error.
func (e *AbortEventError) GetOriginalError() error {
	return e.origErr
}
