package flow

import (
	"fmt"
	"net/http"
	"slices"
)

// Register registers and configures an approval flow with the given id, store
// and steps. A function for retrieving a flow for a given object is returned.
func Register[T fmt.Stringer](id string, st Store, steps [][]Step[T]) func(T) Flow[T] {
	return func(o T) Flow[T] {
		return Flow[T]{
			ID:     id,
			Object: o,
			store:  st,
			steps:  steps,
		}
	}
}

type Flow[T fmt.Stringer] struct {
	ID     string
	Object T
	store  Store
	steps  [][]Step[T]
}

// Start starts a flow. From this point on the flow is in progress until
// completed or cleared. If the receiver is an already started flow,
// the flow state is reset to the start (even if previously completed).
func (f Flow[T]) Start() error {
	var steps []string
	for _, ss := range f.steps {
		for _, s := range ss {
			steps = append(steps, s.ID)
		}
	}
	if err := f.store.EnsureFlow(f.ID, f.Object.String(), steps); err != nil {
		return err
	}
	// Not implemented.
	return nil
}

// Clear cancels a flow. If the receiver is an already cleared flow or not
// started flow, this is a no-op.
func (f Flow[T]) Clear() error {
	if err := f.store.DeleteFlow(f.ID, f.Object.String()); err != nil {
		return err
	}
	// Not implemented.
	return nil
}

// InProgress reports whether the current flow is in progress.
// Not started, cleared and completed flows are not in progress.
func (f Flow[T]) InProgress() bool {
	steps, err := f.store.GetSteps(f.ID, f.Object.String())
	if err != nil {
		// FIXME.
		panic(err)
	}
	remaining := slices.DeleteFunc(steps, func(s stepStatus) bool {
		return s.done
	})
	return len(remaining) > 0
}

// NextSteps return the next actionable steps for the flow.
// No steps are returned for not started, cleared and completed flows.
func (Flow[T]) NextSteps() []Step[T] {
	// Not implemented.
	return nil
}

type Step[T fmt.Stringer] struct {
	ID      string
	Message string
	object  T
	Handle  func(object T, approved bool) error
}

// Approve approves the step and progresses the flow to the next steps.
// The step handler is called with a true approved value.
// If an error is returned, the flow state is not changed.
func (s Step[T]) Approve(u User, rationale string) error {
	if s.Handle != nil {
		if err := s.Handle(s.object, true); err != nil {
			return err
		}
	}
	// Not implemented.
	return nil
}

// Reject rejects the step and progresses the flow to the previous steps.
// The step handler is called with a false approved value.
// If an error is returned, the flow state is not changed.
func (s Step[T]) Reject(u User, rationale string) error {
	if s.Handle != nil {
		if err := s.Handle(s.object, false); err != nil {
			return err
		}
	}
	// Not implemented.
	return nil
}

// Started returns all started flows for the given object type.
func Started[T fmt.Stringer]() []Flow[T] {
	// Not implemented.
	return nil
}

// User must be implemented by approvers of the flows.
type User interface {
	// ID returns an identifier for the user.
	ID() string
	// CanHandleFlowStep reports whether a user can approve or reject
	// the given flow, step and object.
	CanHandleFlowStep(flowID, stepID, object string) bool
}

// Store is used to persist the flow state.
type Store interface {
	EnsureFlow(flowID, object string, steps []string) error
	DeleteFlow(flowID, object string) error
	GetSteps(flowID, object string) ([]stepStatus, error)
	UpdateSteps([]stepStatus) error
}

type stepStatus struct {
	id        string
	rationale string
	done      bool
}

func HTTPHandler[T fmt.Stringer](user func(req *http.Request) User, flow func(T) Flow[T]) http.Handler {
	// Not implemented.
	return nil
}
