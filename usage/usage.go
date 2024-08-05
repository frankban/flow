package usage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/frankban/flow"
)

type Order struct {
	ID         string
	CRRequired bool
}

func (o Order) String() string {
	return o.ID
}

// Defining and registering a flow.
var st = struct {
	flow.Store
}{}
var CR = flow.Register("CR", st, [][]flow.Step[*Order]{{
	{
		ID:      "step-1",
		Message: "This is a message to be displayed to users",
		Handle: func(object *Order, approved bool) error {
			// What happens to the order when this step is approved or rejected?
			return nil
		},
	},
}, {
	{
		ID:      "step-2",
		Message: "This is a message to be displayed to users",
	},
	{
		ID:      "step-3",
		Message: "This is a message to be displayed to users",
		Handle: func(object *Order, approved bool) error {
			// What happens to the order when this step is approved or rejected?
			return nil
		},
	},
}, {
	{
		ID:      "step-4",
		Message: "This is a message to be displayed to users",
	},
}})

// Starting/stopping a flow for a specific object.
func saveDraftOrder(o *Order) error {
	var err error
	if o.CRRequired {
		err = CR(o).Start()
	} else {
		err = CR(o).Clear()
	}
	return err
}

// Gating an action with an approval flow.
func publishOrder(o *Order) error {
	if CR(o).InProgress() {
		return errors.New("approvals required")
	}
	return nil
}

// Listing next steps required for approvals for a given object type.
func orderApprovals() {
	for _, f := range flow.Started[*Order]() {
		for _, s := range f.NextSteps() {
			fmt.Printf("flow %s step %s\n", f.ID, s.Message)
		}
	}
}

// Approving or rejecting a step.
func handleCR(o *Order, stepID string, user flow.User, approved bool, rationale string) error {
	f := CR(o)
	for _, s := range f.NextSteps() {
		if s.ID == stepID {
			if !user.CanHandleFlowStep(f.ID, s.ID, f.Object.String()) {
				return errors.New("forbidden")
			}
			var err error
			if approved {
				err = s.Approve(user, rationale)
			} else {
				err = s.Reject(user, rationale)
			}
			if err != nil {
				return err
			}
		}
	}
	return errors.New("invalid step " + stepID)
}

// HTTP handler actually provided by the library.
func registerHanldlers() {
	user := func(req *http.Request) flow.User {
		return nil
	}
	mux := http.NewServeMux()
	mux.Handle("/flow/cr/", flow.HTTPHandler(user, CR))
}

// Questions:
// What happens to already started flows if the flow definition changes?
// - If the flow was already started, it can have its own definition persisted.
// The messaging/notification should be handled by the library or by provided Handlers?
// How does the order state come into play?
