package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// An assistant beside an app may act on the page the reader works in: open a page, open a form and
// fill it in, or point at an element. It never saves; the reader checks the form and saves it.
//
// The agent only asks for an action by name. A Catalog, the app's allowlist, turns the request into an
// Action or refuses it: unknown names, other roles, fields the form does not have, IDs that are not
// IDs, paths outside the app, sensitive fields (bank details, e-mail, tax numbers) and values that look
// like them. Stream.Action sends the checked action as an event; ActionScript runs it on the page.

// Action kinds.
const (
	ActionOpen = "open" // go to Path
	ActionFill = "fill" // open the form Form on Path and fill in Fields; never submit
	ActionShow = "show" // scroll to the element Anchor on Path and highlight it
)

// ActionSpec is one entry of a Catalog.
type ActionSpec struct {
	Name        string // what the agent asks for: "order-new"
	Title       string // short, for the card in the conversation: "New order"
	Description string // for the agent's prompt
	Kind        string // ActionOpen, ActionFill or ActionShow
	// Path is the page the action runs on: "/orders/{id}". Each {param} is filled from Params; the
	// parameter "path" (ParamPath) carries a whole path, which must lie in Catalog.Pages.
	Path   string
	Form   string   // ActionFill: the form's data-form attribute
	Anchor string   // ActionShow: the element's data-tour attribute
	Params []string // the placeholders in Path: IDs, or ParamPath
	Fields []string // ActionFill: the form fields (name attributes) the agent may fill
	Roles  []string // who may get the action; empty: everyone
}

// ParamPath is the parameter that carries a whole path ("open any page of the app").
const ParamPath = "path"

// Action is a checked action as the page gets it. Build it with Catalog.Check, not by hand.
type Action struct {
	ID     string            `json:"id"` // the card's element, set by Stream.Action
	Name   string            `json:"name"`
	Kind   string            `json:"kind"`
	Title  string            `json:"title"`
	Path   string            `json:"path"`
	Form   string            `json:"form,omitempty"`
	Anchor string            `json:"anchor,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Catalog is an app's allowlist of actions.
type Catalog struct {
	Actions []ActionSpec
	// Pages are the pages a ParamPath may open, without the leading slash: "orders" allows /orders and
	// /orders/{id}; "crm/customers" allows /crm/customers and /crm/customers/{id}.
	Pages []string
	// SensitiveFields matches field names the assistant never fills; nil: DefaultSensitiveFields.
	SensitiveFields *regexp.Regexp
	// SensitiveValues matches values it never puts in any field; nil: DefaultSensitiveValues.
	SensitiveValues *regexp.Regexp
	// MaxValue is the longest value of a field; 0: 500 characters.
	MaxValue int
	// RequireRoles: every action names its roles. Problems reports an action without roles and Check
	// refuses it, instead of offering it to everyone. For apps whose rights are per role.
	RequireRoles bool
}

// ActionOptions configure ActionScript for a request (WithActionOptions).
type ActionOptions struct {
	// FormAttribute marks the forms a fill action finds by its Form, default "data-form"; set it like
	// detail.ScriptConfig.FormAttribute when the app uses another attribute.
	FormAttribute string
}

func (o ActionOptions) formAttribute() string {
	if validFormAttribute.MatchString(o.FormAttribute) {
		return o.FormAttribute
	}
	return "data-form"
}

type actionOptionsKey struct{}

// WithActionOptions sets the options of ActionScript for the request.
func WithActionOptions(ctx context.Context, o ActionOptions) context.Context {
	return context.WithValue(ctx, actionOptionsKey{}, o)
}

// ActionOptionsFrom returns the request's options, the zero value if none were set.
func ActionOptionsFrom(ctx context.Context) ActionOptions {
	o, _ := ctx.Value(actionOptionsKey{}).(ActionOptions)
	return o
}

// DefaultSensitiveFields: bank details, e-mail addresses, tax and routing numbers. A wrongly filled IBAN
// is the costliest damage a prompt injection can do; people enter these themselves.
var DefaultSensitiveFields = regexp.MustCompile(`(?i)iban|bic|swift|e-?mail|password|passwort|tax|steuer|vat|ust-?id|leitweg`)

// DefaultSensitiveValues: values that look like an IBAN or an e-mail address, whatever the field.
var DefaultSensitiveValues = regexp.MustCompile(`(?i)\b[A-Z]{2}[0-9]{2}(?:[ ]?[A-Z0-9]){11,30}\b|[^\s@]+@[^\s@]+\.[^\s@]+`)

var (
	validID     = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)
	validName   = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
	validAnchor = regexp.MustCompile(`^[a-z][a-z0-9-]*(\.[a-z0-9-]+)*$`)
	validField  = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,63}$`)
	validPage   = regexp.MustCompile(`^[a-z][a-z0-9-]*(/[a-z][a-z0-9-]*)*$`)

	validFormAttribute = regexp.MustCompile(`^data-[a-z][a-z0-9-]*$`)
)

// ErrAction is wrapped by every refusal of Check, so callers can tell a refused action from other
// errors. The message says why, for the agent: it can correct the request or tell the reader.
var ErrAction = errors.New("action refused")

func refuse(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrAction, fmt.Sprintf(format, args...))
}

func (c *Catalog) sensitiveFields() *regexp.Regexp {
	if c.SensitiveFields != nil {
		return c.SensitiveFields
	}
	return DefaultSensitiveFields
}

func (c *Catalog) sensitiveValues() *regexp.Regexp {
	if c.SensitiveValues != nil {
		return c.SensitiveValues
	}
	return DefaultSensitiveValues
}

func (c *Catalog) maxValue() int {
	if c.MaxValue > 0 {
		return c.MaxValue
	}
	return 500
}

func (c *Catalog) spec(role, name string) (ActionSpec, error) {
	i := slices.IndexFunc(c.Actions, func(a ActionSpec) bool { return a.Name == name })
	if i < 0 {
		return ActionSpec{}, refuse("there is no action %q (available: %s)", name, strings.Join(c.Names(role), ", "))
	}
	a := c.Actions[i]
	if !c.offers(a, role) {
		return ActionSpec{}, refuse("action %q is not available to the role %q", name, role)
	}
	return a, nil
}

// Check turns an agent's request into an Action, or says why not. values holds the parameters of the
// path and the form fields. An empty field value means "leave as is" and is dropped, so the assistant
// cannot clear what the reader entered.
func (c *Catalog) Check(role, name string, values map[string]string) (Action, error) {
	a, err := c.spec(role, name)
	if err != nil {
		return Action{}, err
	}
	out := Action{Name: a.Name, Kind: a.Kind, Title: a.Title, Path: a.Path, Form: a.Form, Anchor: a.Anchor}
	for _, p := range a.Params {
		v := strings.TrimSpace(values[p])
		switch {
		case p == ParamPath && !c.pagePath(v):
			return Action{}, refuse("%q is not a page of this app", v)
		case p != ParamPath && !validID.MatchString(v):
			return Action{}, refuse("%s %q is not a valid ID", p, v)
		}
		if p == ParamPath && a.Kind == ActionOpen {
			out.Title = v // "Open page: /orders/A-1" says more than the spec's title
		}
		out.Path = strings.ReplaceAll(out.Path, "{"+p+"}", v)
	}
	names := make([]string, 0, len(values))
	for k := range values {
		names = append(names, k)
	}
	sort.Strings(names) // the first refusal does not depend on map order
	for _, f := range names {
		v := strings.TrimSpace(values[f])
		if slices.Contains(a.Params, f) {
			continue
		}
		if c.sensitiveFields().MatchString(f) {
			return Action{}, refuse("the assistant does not fill in %q; the reader enters it in the form", f)
		}
		if !slices.Contains(a.Fields, f) {
			return Action{}, refuse("%s has no field %q (fields: %s)", name, f, strings.Join(a.Fields, ", "))
		}
		if v == "" {
			continue
		}
		if c.sensitiveValues().MatchString(v) {
			return Action{}, refuse("the value for %q looks like an IBAN or an e-mail address; the reader enters such values", f)
		}
		if len(v) > c.maxValue() {
			return Action{}, refuse("the value for %q is longer than %d characters", f, c.maxValue())
		}
		if out.Fields == nil {
			out.Fields = map[string]string{}
		}
		out.Fields[f] = v
	}
	return out, nil
}

func (c *Catalog) pagePath(p string) bool {
	for _, page := range c.Pages {
		rest, ok := strings.CutPrefix(p, "/"+page)
		if ok && (rest == "" || strings.HasPrefix(rest, "/") && validID.MatchString(rest[1:])) {
			return true
		}
	}
	return false
}

// offers: role may get a; without roles everyone, unless RequireRoles.
func (c *Catalog) offers(a ActionSpec, role string) bool {
	if len(a.Roles) == 0 {
		return !c.RequireRoles
	}
	return slices.Contains(a.Roles, role)
}

// Names are the actions available to role, sorted.
func (c *Catalog) Names(role string) []string {
	var n []string
	for _, a := range c.Actions {
		if c.offers(a, role) {
			n = append(n, a.Name)
		}
	}
	sort.Strings(n)
	return n
}

// Prompt describes the actions available to role for the agent's prompt, one line each.
func (c *Catalog) Prompt(role string) string {
	var b strings.Builder
	for _, a := range c.Actions {
		if !c.offers(a, role) {
			continue
		}
		fmt.Fprintf(&b, "- %s: %s", a.Name, a.Description)
		if f := append(slices.Clone(a.Params), a.Fields...); len(f) > 0 {
			fmt.Fprintf(&b, " (fields: %s)", strings.Join(f, ", "))
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Problems checks the catalog itself; call it from a test of the app. It reports duplicate or invalid
// names, missing titles, unknown kinds, a form or anchor that does not fit the kind, parameters missing
// from the path, invalid or sensitive fields and invalid pages.
func (c *Catalog) Problems() []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }
	for _, p := range c.Pages {
		if !validPage.MatchString(p) {
			add("page %q: lower-case segments of letters, digits and '-', no leading or trailing slash", p)
		}
	}
	seen := map[string]bool{}
	for _, a := range c.Actions {
		if !validName.MatchString(a.Name) {
			add("action %q: name of lower-case letters, digits and '-'", a.Name)
		}
		if seen[a.Name] {
			add("action %q twice", a.Name)
		}
		seen[a.Name] = true
		if strings.TrimSpace(a.Title) == "" {
			add("action %q: no title", a.Name)
		}
		if c.RequireRoles && len(a.Roles) == 0 {
			add("action %q: no role", a.Name)
		}
		switch a.Kind {
		case ActionOpen:
			if a.Form != "" || a.Anchor != "" || len(a.Fields) > 0 {
				add("action %q: an open action has no form, anchor or fields", a.Name)
			}
		case ActionFill:
			if !validName.MatchString(a.Form) {
				add("action %q: data-form %q invalid", a.Name, a.Form)
			}
			if a.Anchor != "" {
				add("action %q: a fill action has no anchor", a.Name)
			}
		case ActionShow:
			if !validAnchor.MatchString(a.Anchor) {
				add("action %q: data-tour %q invalid", a.Name, a.Anchor)
			}
			if a.Form != "" || len(a.Fields) > 0 {
				add("action %q: a show action has no form or fields", a.Name)
			}
		default:
			add("action %q: unknown kind %q", a.Name, a.Kind)
		}
		if !strings.HasPrefix(a.Path, "/") && a.Path != "{"+ParamPath+"}" {
			add("action %q: path %q does not start with /", a.Name, a.Path)
		}
		for _, p := range a.Params {
			if !strings.Contains(a.Path, "{"+p+"}") {
				add("action %q: parameter %q missing from the path", a.Name, p)
			}
		}
		for _, f := range append(slices.Clone(a.Params), a.Fields...) {
			if !validField.MatchString(f) {
				add("action %q: field %q invalid", a.Name, f)
			}
		}
		for _, f := range a.Fields {
			if c.sensitiveFields().MatchString(f) {
				add("action %q: field %q is sensitive and is never filled", a.Name, f)
			}
		}
	}
	return errs
}

// Action sends a checked action: its card into SlotBlocks (with a button to run it again) and the
// action itself as an event, which ActionScript runs on the page once the stream has ended. id names
// the card; it must be unique on the page.
func (s *Stream) Action(ctx context.Context, id string, a Action) error {
	a.ID = id
	if err := s.Send(ctx, SlotBlocks, id, "", ActionCard(a)); err != nil {
		return err
	}
	return s.event("action", a)
}

func (s *Stream) event(typ string, args any) error {
	b, err := json.Marshal(map[string]any{"type": typ, "args": args})
	if err != nil {
		return err
	}
	return s.write(b)
}

func (a Action) json() string {
	b, _ := json.Marshal(a)
	return string(b)
}
