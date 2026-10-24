package detail

import (
	"context"
	"encoding/json"
	"io"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/kit"
)

// ScriptConfig names what the script and the server exchange. The zero value uses the defaults; an app
// with existing names (form fields, a global variable) sets its own.
type ScriptConfig struct {
	DialogField string // hidden field with the dialog's ID in each of its forms, default "_dialog"
	ReturnField string // hidden field with the page to return to after an error, default "_return"
	FormField   string // hidden field with the form's kind (FormAttribute), default "_form"
	// FormAttribute marks a form's kind (e.g. "customer-edit"); after an error the script fills the form
	// of that kind whose hidden and read-only fields match the values. Default "data-form".
	FormAttribute string
	// StateVar is the global variable the page sets after a failed save:
	// {"id": dialog ID, "error": message, "values": {field: [values]}}. Default "kitDialogState".
	StateVar string
	// ValuesAttribute on a button that opens a shared dialog (data-bs-target) carries the values to put
	// into its form first, as JSON. Default "data-kit-values".
	ValuesAttribute string
}

func (c ScriptConfig) withDefaults(ctx context.Context) map[string]string {
	m := kit.T(ctx)
	return map[string]string{
		"dialogField":     kit.Or(c.DialogField, "_dialog"),
		"returnField":     kit.Or(c.ReturnField, "_return"),
		"formField":       kit.Or(c.FormField, "_form"),
		"formAttribute":   kit.Or(c.FormAttribute, "data-form"),
		"stateVar":        kit.Or(c.StateVar, "kitDialogState"),
		"valuesAttribute": kit.Or(c.ValuesAttribute, "data-kit-values"),
		"cancel":          m.Cancel,
		"sectionHint":     m.SectionHint,
	}
}

// Script renders the dialogs' script with its configuration; put it after Tabler's script. It exposes
// window.kitDialog.open(element) (opens the dialog an element is in) and
// window.kitDialog.matchForm(forms, values, kind).
func Script(c ScriptConfig) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		cfg, err := json.Marshal(c.withDefaults(ctx)) // escapes <, > and &: cannot end the script
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, csp.Open(ctx, "script")+"(function(cfg){"+script+"})("+string(cfg)+");</script>")
		return err
	})
}

// Styles renders the block's CSS (Tabler variables only); the shell block includes it.
func Styles() templ.Component {
	return csp.Style(css)
}

const css = `.kit-section{container-type:inline-size}
.kit-section-head{display:flex;align-items:center;justify-content:space-between;gap:.75rem;padding:1rem 1.25rem 0}
.kit-section-head+.card-body{padding-top:.75rem}
.kit-section-head+.table-responsive{margin-top:.75rem;border-top:var(--tblr-border-width) solid var(--tblr-border-color)}
.kit-section-title{min-width:0;margin:0;overflow-wrap:anywhere;font-size:1rem;font-weight:600}
.kit-section-title small{margin-left:.375rem;font-size:.8125rem;font-weight:400;color:var(--tblr-secondary)}
.kit-section-head>.kit-section-action,.kit-section-head>form,.kit-section-head>.kit-section-actions{flex-shrink:0}
.kit-section-actions{display:flex;flex-wrap:wrap;justify-content:flex-end;gap:.5rem}
.kit-fields{display:grid;grid-template-columns:minmax(7.5rem,auto) 1fr;gap:.5rem 1.25rem;margin:0}
.kit-fields dt{font-weight:400;color:var(--tblr-secondary)}
.kit-fields dd{min-width:0;margin:0;overflow-wrap:anywhere}
.kit-fields dd.kit-field-empty{color:var(--tblr-secondary)}
@container (max-width:320px){.kit-section-action-text{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0 0 0 0);white-space:nowrap;border:0}.kit-section-action{padding-left:.5rem;padding-right:.5rem}.kit-section-action .ti{margin-right:0}.kit-fields{grid-template-columns:max-content minmax(0,1fr);column-gap:.75rem}}
.kit-dialog-context{margin-bottom:.125rem;font-size:.8125rem;color:var(--tblr-secondary)}
.kit-dialog-footer{justify-content:flex-end;gap:.5rem}
.kit-dialog-footer>*{margin:0}
.kit-dialog-hint{margin-right:auto!important;font-size:.8125rem;color:var(--tblr-secondary)}
form.kit-dialog-invalid{outline:2px solid rgba(var(--tblr-danger-rgb),.25);outline-offset:-2px;border-radius:var(--tblr-border-radius)}`

// script is the body of function(cfg){…}. Tabler's script opens and closes the modals: window.tabler since
// Tabler 1.x, window.bootstrap with Bootstrap's own script.
const script = `"use strict";var bs=window.tabler&&window.tabler.Modal?window.tabler:window.bootstrap;
function hidden(f,n,v){if(f.elements.namedItem(n))return;var h=document.createElement("input");h.type="hidden";h.name=n;h.value=v;f.appendChild(h)}
document.querySelectorAll(".modal[data-kit-dialog]").forEach(function(d){
var forms=d.querySelectorAll("form");
forms.forEach(function(f){hidden(f,cfg.dialogField,d.id);hidden(f,cfg.returnField,location.pathname+location.search);var k=f.getAttribute(cfg.formAttribute);if(k)hidden(f,cfg.formField,k);
if(d.classList.contains("kit-dialog-large"))f.addEventListener("submit",function(){try{sessionStorage.setItem("kitDialogOpen",d.id)}catch(_){}})});
var f=forms.length===1?forms[0]:null;var rows=f?Array.from(f.children).filter(function(e){return!(e.tagName==="INPUT"&&e.type==="hidden")}):[];var row=rows[rows.length-1];
var onlyButtons=row&&row.children.length>0&&Array.from(row.children).every(function(e){return e.tagName==="BUTTON"});
if(onlyButtons&&row.querySelector(":scope > .btn-primary")){if(!f.id)f.id=d.id+"-form";var foot=document.createElement("div");foot.className="modal-footer kit-dialog-footer";
var hint=document.createElement("span");hint.className="kit-dialog-hint";if(d.hasAttribute("data-kit-section"))hint.textContent=cfg.sectionHint;foot.appendChild(hint);
var c=document.createElement("button");c.type="button";c.className="btn";c.setAttribute("data-bs-dismiss","modal");c.textContent=cfg.cancel;foot.appendChild(c);
Array.from(row.children).forEach(function(b){b.setAttribute("form",f.id);b.classList.remove("w-100");foot.appendChild(b)});row.remove();d.querySelector(".modal-content").appendChild(foot)}
d.addEventListener("shown.bs.modal",function(){if(d.classList.contains("kit-dialog-large"))return;var first=d.querySelector("[data-kit-focus], input:not([type=hidden]):not([disabled]):not([readonly]), select, textarea");if(first&&first.focus)first.focus()})});
function open(el){var d=el&&el.closest?el.closest(".modal"):null;if(!d||!bs)return false;bs.Modal.getOrCreateInstance(d).show();return true}
function matchForm(forms,values,kind){var cands=Array.from(forms).filter(function(f){return!kind||f.getAttribute(cfg.formAttribute)===kind});var one=function(v){return Array.isArray(v)?v[0]:v};
for(var i=0;i<cands.length;i++){var f=cands[i];var keys=Array.from(f.querySelectorAll("input[type=hidden], input[readonly]")).filter(function(e){return e.name&&e.name.charAt(0)!=="_"});
if(keys.length&&keys.every(function(e){return values[e.name]!==undefined&&String(one(values[e.name]))===e.value}))return f}return cands[0]||null}
window.kitDialog={open:open,matchForm:matchForm};
function restore(z){var d=document.getElementById(z.id);if(!d)return;var vals=z.values||{};var f=matchForm(d.querySelectorAll("form"),vals,vals[cfg.formField]?vals[cfg.formField][0]:"");
if(f){Object.keys(vals).forEach(function(n){if(n.charAt(0)==="_")return;var fld=f.elements.namedItem(n);if(!fld)return;var list=vals[n];var els=fld.tagName?[fld]:Array.from(fld);
els.forEach(function(e){if(e.type==="checkbox"||e.type==="radio")e.checked=list.indexOf(e.value)>=0;else if(e.type!=="file"&&e.type!=="password")e.value=list[0]||""})});if(z.error)f.classList.add("kit-dialog-invalid")}
open(d.firstElementChild||d)}
var state=window[cfg.stateVar];if(state)restore(state);else{try{var o=sessionStorage.getItem("kitDialogOpen");sessionStorage.removeItem("kitDialogOpen");var od=o&&document.getElementById(o);if(od&&bs)bs.Modal.getOrCreateInstance(od).show()}catch(_){}}
document.addEventListener("click",function(e){var b=e.target.closest?e.target.closest("["+cfg.valuesAttribute+"][data-bs-target]"):null;if(!b)return;var d=document.querySelector(b.getAttribute("data-bs-target"));var f=d&&d.querySelector("form");if(!f)return;
var vals={};try{vals=JSON.parse(b.getAttribute(cfg.valuesAttribute))||{}}catch(_){return}var err=d.querySelector(".kit-dialog-error");if(err)err.remove();
Object.keys(vals).forEach(function(n){var fld=f.elements.namedItem(n);if(fld&&fld.tagName)fld.value=vals[n]===null||vals[n]===undefined?"":String(vals[n])})},true);`
