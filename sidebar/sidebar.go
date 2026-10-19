package sidebar

import "net/http"

// CookieName holds the desktop state the script writes when the sidebar is toggled: "collapsed" or
// "expanded". Read it with Collapsed so the server renders the sidebar the way the user left it.
const CookieName = "kit_sidebar"

// Collapsed reports whether the user collapsed the sidebar to icons.
func Collapsed(r *http.Request) bool {
	c, err := r.Cookie(CookieName)
	return err == nil && c.Value == "collapsed"
}

// css is the sidebar's layout. It only uses Tabler's variables and its own (--kit-sidebar-…), so the
// app's theme decides colors:
//
//	--kit-sidebar-width            expanded width (15rem)
//	--kit-sidebar-width-collapsed  icon-only width (3.75rem)
//	--kit-sidebar-bg               background (Tabler's surface)
//	--kit-sidebar-hover            background of a hovered link
//	--kit-sidebar-active-bg        background of the current page's link
//	--kit-sidebar-active-color     text of the current page's link
const css = `:root{--kit-sidebar-width:15rem;--kit-sidebar-width-collapsed:3.75rem}
.kit-sidebar-inner{display:flex;flex-direction:column;height:100%;min-height:0;width:100%;padding:.75rem;background:var(--kit-sidebar-bg,var(--tblr-bg-surface));border-right:var(--tblr-border-width) solid var(--tblr-border-color)}
.kit-sidebar-header{display:flex;align-items:center;gap:.5rem;min-width:0;margin-bottom:1rem}
.kit-sidebar-brand{display:flex;align-items:center;gap:.625rem;min-width:0;flex:1 1 auto;padding:.25rem .375rem;color:inherit;text-decoration:none}
.kit-sidebar-brand:hover{color:inherit;text-decoration:none}
.kit-sidebar-brand-text{display:flex;flex-direction:column;min-width:0;line-height:1.25}
.kit-sidebar-brand-title{font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kit-sidebar-brand-subtitle{font-size:.75rem;color:var(--tblr-secondary);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kit-sidebar-content{display:flex;flex-direction:column;flex:1 1 auto;min-height:0;overflow-x:hidden;overflow-y:auto}
.kit-sidebar-group{display:flex;flex-direction:column}
.kit-sidebar-group-end{margin-top:auto}
.kit-sidebar-label{padding:1rem .625rem .375rem;font-size:.75rem;font-weight:500;color:var(--tblr-secondary);white-space:nowrap;overflow:hidden}
.kit-sidebar-menu{display:flex;flex-direction:column;gap:2px;margin:0;padding:0;list-style:none}
.kit-sidebar-link{display:flex;align-items:center;gap:.625rem;height:2.25rem;padding:0 .625rem;border-radius:var(--tblr-border-radius);color:var(--tblr-body-color);white-space:nowrap;text-decoration:none}
.kit-sidebar-link:hover{background:var(--kit-sidebar-hover,rgba(var(--tblr-body-color-rgb),.05));color:var(--tblr-body-color);text-decoration:none}
.kit-sidebar-link.active{background:var(--kit-sidebar-active-bg,rgba(var(--tblr-body-color-rgb),.07));color:var(--kit-sidebar-active-color,var(--tblr-body-color));font-weight:500}
.kit-sidebar-link>.ti{flex:0 0 auto}
.kit-sidebar-text{flex:1 1 auto;min-width:0;overflow:hidden;text-overflow:ellipsis}
.kit-sidebar-badge{margin-left:auto}
.kit-sidebar-separator{height:var(--tblr-border-width);margin:.75rem .625rem;background:var(--tblr-border-color)}
.kit-sidebar-footer{margin-top:.5rem;padding-top:.625rem;border-top:var(--tblr-border-width) solid var(--tblr-border-color)}
.kit-sidebar-user{display:flex;align-items:center;gap:.625rem;width:100%;padding:.5rem .625rem;border:0;border-radius:var(--tblr-border-radius);background:none;color:inherit;text-align:left}
.kit-sidebar-user:hover,.kit-sidebar-user[aria-expanded=true]{background:var(--kit-sidebar-hover,rgba(var(--tblr-body-color-rgb),.05))}
.kit-sidebar-avatar{display:inline-flex;flex:0 0 auto;align-items:center;justify-content:center;width:2rem;height:2rem;border-radius:999px;background:rgba(var(--tblr-body-color-rgb),.08);font-size:.75rem;font-weight:600}
.kit-sidebar-user-text{display:flex;flex:1 1 auto;flex-direction:column;min-width:0;line-height:1.3}
.kit-sidebar-user-name{font-weight:500;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kit-sidebar-user-detail{font-size:.75rem;color:var(--tblr-secondary);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.kit-sidebar-footer .dropdown-menu{min-width:15rem}
@media (min-width:992px){
.kit-sidebar{position:fixed;top:0;bottom:0;left:0;z-index:1030;display:flex;width:var(--kit-sidebar-width);transition:width .15s ease}
.kit-sidebar~.page-wrapper{margin-left:var(--kit-sidebar-width);transition:margin-left .15s ease}
.kit-sidebar[data-collapsed]{width:var(--kit-sidebar-width-collapsed)}
.kit-sidebar[data-collapsed]~.page-wrapper{margin-left:var(--kit-sidebar-width-collapsed)}
.kit-sidebar[data-collapsed] .kit-sidebar-inner{padding:.75rem .5rem}
.kit-sidebar[data-collapsed] :is(.kit-sidebar-text,.kit-sidebar-brand-text,.kit-sidebar-user-text,.kit-sidebar-hide-collapsed){display:none}
.kit-sidebar[data-collapsed] .kit-sidebar-label{height:0;padding:0;margin:.5rem 0;visibility:hidden}
.kit-sidebar[data-collapsed] :is(.kit-sidebar-link,.kit-sidebar-brand,.kit-sidebar-user){justify-content:center;padding-left:0;padding-right:0}
.kit-sidebar[data-collapsed] .kit-sidebar-link{position:relative}
.kit-sidebar[data-collapsed] .kit-sidebar-badge{position:absolute;top:.25rem;right:.25rem;min-width:0;width:.5rem;height:.5rem;padding:0;font-size:0;background:var(--tblr-red)!important}
}
@media (max-width:991.98px){.kit-sidebar.offcanvas-lg{--tblr-offcanvas-width:17rem}.kit-sidebar-inner{border-right:0}}
@media print{.kit-sidebar{display:none!important}.kit-sidebar~.page-wrapper{margin-left:0!important}}`

// script toggles the desktop state, keeps the tooltips (title) of collapsed links and writes the cookie.
// The mobile sheet is Tabler's own offcanvas.
const script = `(function(){document.addEventListener("click",function(e){var b=e.target.closest&&e.target.closest("[data-kit-sidebar-toggle]");if(!b)return;var s=document.getElementById(b.getAttribute("data-kit-sidebar-toggle"));if(!s)return;var c=!s.hasAttribute("data-collapsed");s.toggleAttribute("data-collapsed",c);b.setAttribute("aria-expanded",String(!c));s.querySelectorAll("[data-kit-tooltip]").forEach(function(l){if(c)l.setAttribute("title",l.getAttribute("data-kit-tooltip"));else l.removeAttribute("title")});document.cookie="` + CookieName + `="+(c?"collapsed":"expanded")+"; path=/; max-age=31536000; samesite=lax"})})();`
