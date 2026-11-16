package detail

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
)

// State is what the page tells the dialog script after a failed save: which dialog opens again, the
// message and the values entered. Render it with the page that answers the failed POST (usually with
// status 422); Script then opens the dialog, fills in the values of the form that was sent and marks it.
type State struct {
	ID     string     // the dialog's ID
	Error  string     // the message on top of the dialog
	Values url.Values // the values sent, usually r.PostForm; see StateValues
	// Kind is the form's kind (ScriptConfig.FormAttribute) when the dialog holds several forms, e.g. the
	// rows of a list dialog; with the form's hidden fields in Values, the script finds the row that failed.
	Kind string
}

// StateValues are the values of a form to put back: everything but the script's own fields (they start
// with "_"), passwords and the fields named in skip (e.g. a CSRF token).
func StateValues(form url.Values, skip ...string) url.Values {
	out := url.Values{}
	for k, v := range form {
		if strings.HasPrefix(k, "_") || strings.Contains(strings.ToLower(k), "password") || contains(skip, k) {
			continue
		}
		out[k] = v
	}
	return out
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

// StateScript sets the state for Script (window.kitDialogState, or cfg.StateVar) with the nonce of
// templ.WithNonce. Nothing when ID is empty.
func StateScript(s State, cfg ...ScriptConfig) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if s.ID == "" {
			return nil
		}
		var c ScriptConfig
		if len(cfg) > 0 {
			c = cfg[0]
		}
		names := c.withDefaults(ctx)
		values := map[string][]string{}
		for k, v := range s.Values {
			values[k] = v
		}
		if s.Kind != "" {
			values[names["formField"]] = []string{s.Kind}
		}
		b, err := json.Marshal(map[string]any{"id": s.ID, "error": s.Error, "values": values}) // escapes <, > and &
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, csp.Open(ctx, "script")+"window["+jsString(names["stateVar"])+"]="+string(b)+";</script>")
		return err
	})
}

func jsString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
