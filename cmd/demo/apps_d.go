package main

import "github.com/a-h/templ"

// appsDPages are the chat, email, music and wizard pages (apps_d.templ).
func appsDPages() map[string]templ.Component {
	return map[string]templ.Component{
		"chat":        chatPage(),
		"email-inbox": emailInboxPage(),
		"emails":      emailsPage(),
		"music":       musicPage(),
		"wizard":      wizardPage(),
	}
}
