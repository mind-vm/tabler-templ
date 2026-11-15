package sidepanel

import (
	"context"
	"io"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
)

func scriptElement() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, csp.Open(ctx, "script")+script+"</script>")
		return err
	})
}

const css = `.kit-sidepanel{--kit-sidepanel-width:25rem;position:fixed;top:0;right:0;bottom:0;z-index:1040;display:flex;flex-direction:column;width:var(--kit-sidepanel-width);max-width:100vw;background:var(--tblr-bg-surface);border-left:var(--tblr-border-width) solid var(--tblr-border-color);transform:translateX(0);transition:transform .2s ease}
.kit-sidepanel.kit-sidepanel-out{transform:translateX(100%)}
.kit-sidepanel[hidden],.kit-sidepanel-content[hidden],.kit-sidepanel-tab-head[hidden]{display:none}
@media (min-width:1200px){body.kit-sidepanel-open .page{margin-right:var(--kit-sidepanel-width,25rem)}.page{transition:margin-right .2s ease}}
@media (max-width:1199.98px){.kit-sidepanel{box-shadow:-8px 0 24px rgba(0,0,0,.18)}.kit-sidepanel[data-kit-from-server]{display:none}html:has(body.kit-sidepanel-open),body.kit-sidepanel-open{overflow:hidden}.kit-sidepanel-content{overscroll-behavior:contain}.kit-sidepanel-head{padding-top:max(.5rem,env(safe-area-inset-top));padding-right:max(.75rem,env(safe-area-inset-right))}.kit-sidepanel-tab,.kit-sidepanel-head .btn{min-height:44px;min-width:44px}@supports (height:100dvh){.kit-sidepanel{bottom:auto;height:100dvh}}}
@media (max-width:575.98px){.kit-sidepanel{left:0;width:auto;max-width:none;border-left:0}.kit-sidepanel-head{padding-left:max(1rem,env(safe-area-inset-left))}.kit-sidepanel-tab{padding:0 .5rem}}
@media (prefers-reduced-motion:reduce){.kit-sidepanel,.page{transition:none}}
.kit-sidepanel-head{display:flex;align-items:center;gap:.5rem;padding:.5rem .75rem .5rem 1rem;border-bottom:var(--tblr-border-width) solid var(--tblr-border-color)}
.kit-sidepanel-tab-head{display:contents}
.kit-sidepanel-tabs{display:flex;gap:.25rem;align-self:stretch}
.kit-sidepanel-tab{display:inline-flex;align-items:center;gap:.375rem;padding:0 .625rem;border:0;border-bottom:2px solid transparent;background:none;color:var(--tblr-secondary);font:inherit;font-weight:500;cursor:pointer}
.kit-sidepanel-tab:hover,.kit-sidepanel-tab[aria-selected=true]{color:var(--tblr-body-color)}
.kit-sidepanel-tab[aria-selected=true]{border-bottom-color:currentColor}
.kit-sidepanel-content{flex:1;min-height:0;display:flex;flex-direction:column;overflow-y:auto}`

// script: the header buttons open the panel on their tab, a click on the open tab's button closes it.
// The state goes into the cookie so the server renders the next page the same way; on small screens the
// panel stays closed after a page change (unless a form in a tab with data-kit-sidepanel-forms was sent).
const script = `(function(){"use strict";
var panel=document.getElementById("kit-sidepanel");if(!panel)return;
var small=window.matchMedia("(max-width: 1199.98px)"),KEEP="kitSidepanelKeep";
function active(){return panel.dataset.active}
function isOpen(){return !panel.hidden&&!panel.classList.contains("kit-sidepanel-out")}
function has(t){return !!panel.querySelector('[data-kit-sidepanel-content="'+t+'"]')}
function remember(){document.cookie="kit-sidepanel="+(isOpen()?active():"closed")+"; path=/; max-age=31536000; samesite=lax"}
function draw(){var t=active();
panel.querySelectorAll("[data-kit-sidepanel-content]").forEach(function(e){e.hidden=e.dataset.kitSidepanelContent!==t});
panel.querySelectorAll("[data-kit-sidepanel-show]").forEach(function(b){b.setAttribute("aria-selected",String(b.dataset.kitSidepanelShow===t))});
panel.querySelectorAll("[data-kit-sidepanel-only]").forEach(function(e){e.hidden=e.dataset.kitSidepanelOnly!==t});
document.body.classList.toggle("kit-sidepanel-open",isOpen());
document.querySelectorAll("[data-kit-sidepanel-tab]").forEach(function(k){k.setAttribute("aria-expanded",String(isOpen()&&k.dataset.kitSidepanelTab===t))})}
var vv=window.visualViewport;
function fit(){var kb=vv&&small.matches&&isOpen()&&Math.abs(vv.scale-1)<0.01&&(vv.height<window.innerHeight-1||vv.offsetTop>0);
if(kb){panel.style.top=vv.offsetTop+"px";panel.style.bottom="auto";panel.style.height=vv.height+"px"}else{panel.style.removeProperty("top");panel.style.removeProperty("bottom");panel.style.removeProperty("height")}}
if(vv){vv.addEventListener("resize",fit);vv.addEventListener("scroll",fit)}
function show(t,on,focus){if(t&&has(t))panel.dataset.active=t;panel.removeAttribute("data-kit-from-server");
if(on){if(panel.hidden){panel.classList.add("kit-sidepanel-out");panel.hidden=false;void panel.offsetWidth}panel.classList.remove("kit-sidepanel-out")}
else if(!panel.hidden){panel.classList.add("kit-sidepanel-out");var done=function(){if(panel.classList.contains("kit-sidepanel-out"))panel.hidden=true};
if(window.matchMedia("(prefers-reduced-motion: reduce)").matches)done();else{panel.addEventListener("transitionend",done,{once:true});setTimeout(done,400)}}
draw();fit();remember();
if(on){panel.dispatchEvent(new CustomEvent("kit-sidepanel:shown",{detail:{tab:active(),focus:!!focus}}));
if(focus&&!small.matches){var f=panel.querySelector('[data-kit-sidepanel-content="'+active()+'"][data-kit-sidepanel-forms] textarea');if(f)f.focus({preventScroll:true})}}}
document.querySelectorAll("[data-kit-sidepanel-tab]").forEach(function(k){k.addEventListener("click",function(e){var t=k.dataset.kitSidepanelTab;if(!has(t))return;e.preventDefault();
if(isOpen()&&active()===t)show(null,false);else show(t,true,true)})});
panel.querySelectorAll("[data-kit-sidepanel-show]").forEach(function(b){b.addEventListener("click",function(){show(b.dataset.kitSidepanelShow,true,true)})});
var x=document.getElementById("kit-sidepanel-close");if(x)x.addEventListener("click",function(){show(null,false)});
panel.addEventListener("keydown",function(e){if(e.key==="Escape"&&small.matches)show(null,false)});
document.addEventListener("submit",function(e){if(e.defaultPrevented||!e.target.closest||!e.target.closest("[data-kit-sidepanel-forms]"))return;try{sessionStorage.setItem(KEEP,"1")}catch(_){}});
var keep=false;try{keep=sessionStorage.getItem(KEEP)==="1";sessionStorage.removeItem(KEEP)}catch(_){}
if(!panel.hidden&&small.matches&&!keep){panel.hidden=true;draw()}
panel.removeAttribute("data-kit-from-server");
window.kitSidepanel={open:function(t){show(t,true,false)},close:function(){show(null,false)},active:active,isOpen:isOpen}})();`
