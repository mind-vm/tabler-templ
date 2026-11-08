package copilot

import "github.com/mind-vm/tabler-templ/blocks/assistant"

// Paths of the app.
const (
	Home     = "/copilot" // the order list
	orders   = Home + "/orders"
	customer = Home + "/customers"
)

// Catalog is the assistant's allowlist: everything it may do on the pages of this app. It reads data
// through the agent's tools, which only read; it changes nothing itself. It opens pages, fills in
// forms for the reader to check and save, and points at elements.
//
// The forms carry data-form, the elements data-tour; the test checks that each action finds its form
// or element on its page.
var Catalog = &assistant.Catalog{
	Pages: []string{"copilot", "copilot/orders", "copilot/customers"},
	Actions: []assistant.ActionSpec{
		{Name: "order-new", Title: "New order", Kind: assistant.ActionFill, Path: Home, Form: "order-new",
			Description: "open the form for a new order and fill it in",
			Fields:      []string{"customer", "title", "amount", "due", "note"}},
		{Name: "order-edit", Title: "Edit order", Kind: assistant.ActionFill, Path: orders + "/{id}", Form: "order-edit",
			Description: "open an order and fill in its status, due date or note",
			Params:      []string{"id"}, Fields: []string{"status", "due", "note"}},
		{Name: "customer-edit", Title: "Edit customer", Kind: assistant.ActionFill, Path: customer + "/{id}", Form: "customer-edit",
			Description: "open a customer and fill in name or city (never e-mail or bank details)",
			Params:      []string{"id"}, Fields: []string{"name", "city"}},
		{Name: "overdue-filter", Title: "Overdue filter", Kind: assistant.ActionShow, Path: Home, Anchor: "orders.filter.overdue",
			Description: "point at the filter for overdue orders"},
		{Name: "page", Title: "Open page", Kind: assistant.ActionOpen, Path: "{path}", Params: []string{assistant.ParamPath},
			Description: "open a page of the app: " + Home + ", " + orders + "/{id}, " + customer + "/{id}"},
	},
}
