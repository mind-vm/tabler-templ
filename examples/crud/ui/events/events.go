// Package events names the browser events the CRUD components fire (HX-Trigger) and listen to. They are
// the only coupling between the components: a component fires an event after a change and refreshes
// itself on the events it cares about, without knowing who fires them or who listens.
package events

const (
	Saved         = "contact:saved"    // created or edited; detail: the contact id
	Deleted       = "contact:deleted"  // detail: the contact id
	Restored      = "contact:restored" // undone deletion; detail: the contact id
	StatusChanged = "contact:status"   // the status cell saved; detail: the contact id
)

// ListChanged are the events after which a list of contacts may show other rows.
var ListChanged = []string{Saved, Deleted, Restored}

// Any are all events that change a contact.
var Any = []string{Saved, Deleted, Restored, StatusChanged}
