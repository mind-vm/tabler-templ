package main

import (
	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/blocks/activity"
	"github.com/mind-vm/tabler-templ/blocks/crm"
	"github.com/mind-vm/tabler-templ/blocks/dashboard"
)

// The data of the dashboard (index) and the CRM dashboard, as the Tabler preview shows it.

var aCommits = []dashboard.Commit{
	{Avatar: avatar.Props{Src: avatarURL("000m")}, Text: "Fix dart Sass compatibility (#29755)", Date: "28 Nov 2025"},
	{Avatar: avatar.Props{Src: avatarURL("052f")}, Text: "Change deprecated html tags to text decoration classes (#29604)", Date: "27 Nov 2025"},
	{Avatar: avatar.Props{Src: avatarURL("002m")}, Text: "justify-content:between \u21d2 justify-content:space-between (#29734)", Date: "26 Nov 2025"},
	{Avatar: avatar.Props{Src: avatarURL("003m")}, Text: "Update change-version.js (#29736)", Date: "26 Nov 2025"},
	{Avatar: avatar.Props{Src: avatarURL("000f")}, Text: "Regenerate package-lock.json (#29730)", Date: "25 Nov 2025"},
}

var aVisits = []dashboard.PageVisit{
	{Path: "/", Visitors: "4,896", Unique: "3,654", Bounce: "82.54%", Trend: []float64{17, 24, 20, 10, 5, 1, 4, 18, 13}},
	{Path: "/form-elements.html", Visitors: "3,652", Unique: "3,215", Bounce: "76.29%", Trend: []float64{13, 11, 19, 22, 12, 7, 14, 3, 21}},
	{Path: "/index.html", Visitors: "3,256", Unique: "2,865", Bounce: "72.65%", Trend: []float64{10, 13, 10, 4, 17, 3, 23, 22, 19}},
	{Path: "/icons.html", Visitors: "986", Unique: "865", Bounce: "44.89%", Trend: []float64{6, 15, 13, 13, 5, 7, 17, 20, 19}},
	{Path: "/docs/", Visitors: "912", Unique: "822", Bounce: "41.12%", Trend: []float64{2, 11, 15, 14, 21, 20, 8, 23, 18, 14}},
	{Path: "/accordion.html", Visitors: "855", Unique: "798", Bounce: "32.65%", Trend: []float64{22, 12, 7, 14, 3, 21, 8, 23, 18, 14}},
}

var aSocial = []dashboard.BarRow{
	{Text: "Instagram", Value: "3,550", Percent: 71},
	{Text: "Twitter", Value: "1,798", Percent: 35.96},
	{Text: "Facebook", Value: "1,245", Percent: 24.9},
	{Text: "TikTok", Value: "986", Percent: 19.72},
	{Text: "Pinterest", Value: "854", Percent: 17.080000000000002},
	{Text: "VK", Value: "650", Percent: 13},
	{Text: "LinkedIn", Value: "420", Percent: 8.4},
}

var aTasks = []dashboard.Task{
	{Name: "Extend the data model.", Checked: true, Due: "December 08, 2024", Done: 2, Total: 7, Comments: 3, Avatar: avatar.Props{Src: avatarURL("000m")}},
	{Name: "Verify the event flow.", Checked: false, Due: "January 01, 2024", Done: 0, Total: 5, Comments: 0, Avatar: avatar.Props{Src: avatarURL("052f")}},
	{Name: "Database backup and maintenance", Checked: false, Due: "January 01, 2024", Done: 0, Total: 5, Comments: 0, Avatar: avatar.Props{Src: avatarURL("002m")}},
	{Name: "Identify the implementation team.", Checked: true, Due: "September 01, 2024", Done: 6, Total: 10, Comments: 12, Avatar: avatar.Props{Src: avatarURL("003m")}},
	{Name: "Define users and workflow", Checked: false, Due: "January 01, 2024", Done: 0, Total: 5, Comments: 0, Avatar: avatar.Props{Src: avatarURL("000f")}},
	{Name: "Check Pull Requests", Checked: true, Due: "July 17, 2024", Done: 2, Total: 9, Comments: 3, Avatar: avatar.Props{Src: avatarURL("001f")}},
}

var crmDeals = []crm.Deal{
	{Company: "Skyline Retail", Contact: "Hannah Cruz", Stage: "Negotiation", StageColor: "yellow", Value: "$48,000", Probability: 80, Owner: avatar.Props{Initials: "OB", Color: "primary"}},
	{Company: "CloudBase Inc", Contact: "David Park", Stage: "Proposal", StageColor: "purple", Value: "$32,500", Probability: 60, Owner: avatar.Props{Initials: "ML", Color: "purple"}},
	{Company: "NovaBuild Co", Contact: "Priya Sharma", Stage: "Qualified", StageColor: "azure", Value: "$27,000", Probability: 45, Owner: avatar.Props{Initials: "SR", Color: "green"}},
	{Company: "MetroPay Ltd", Contact: "Tom Ellis", Stage: "Proposal", StageColor: "purple", Value: "$21,800", Probability: 55, Owner: avatar.Props{Initials: "JO", Color: "yellow"}},
	{Company: "Greenridge Partners", Contact: "Amy Tan", Stage: "Negotiation", StageColor: "yellow", Value: "$18,400", Probability: 70, Owner: avatar.Props{Initials: "OB", Color: "primary"}},
	{Company: "Orion Logistics", Contact: "Ben Okafor", Stage: "Qualified", StageColor: "azure", Value: "$14,200", Probability: 35, Owner: avatar.Props{Initials: "LM", Color: "red"}},
}

var crmStages = []crm.Stage{
	{Label: "Leads", Percent: 100, Color: "primary", Count: "1240", Amount: "$2.4M"},
	{Label: "Qualified", Percent: 75, Color: "azure", Count: "680", Amount: "$1.6M"},
	{Label: "Proposal", Percent: 52, Color: "purple", Count: "310", Amount: "$940K"},
	{Label: "Negotiation", Percent: 34, Color: "yellow", Count: "148", Amount: "$580K"},
	{Label: "Closed won", Percent: 18, Color: "green", Count: "74", Amount: "$312K"},
}

func crmEvent(ic, color, text, time string) activity.Item {
	return activity.Item{Avatar: avatar.Props{Icon: ic, Color: color}, Text: templ.Raw(templ.EscapeString(text)), Time: time}
}

var crmActivity = []activity.Item{
	crmEvent("ti-git-merge", "primary", "Olivia Bennett moved Skyline Retail to Negotiation", "4 min ago"),
	crmEvent("ti-notes", "purple", "Marcus Lee added note to CloudBase opportunity", "22 min ago"),
	crmEvent("ti-check", "green", "Sara Mitchell completed task: Schedule NovaBuild demo", "1 hour ago"),
	crmEvent("ti-alert-triangle", "red", "Churn alert triggered for Orion Logistics account", "2 hours ago"),
	crmEvent("ti-file-text", "yellow", "James Okonkwo created proposal for MetroPay renewal", "Yesterday"),
}

var crmTodos = []dashboard.Todo{
	{ID: "task-due-0", Text: "Send Q2 proposal to Skyline Retail", Color: "red", When: "Today"},
	{ID: "task-due-1", Text: "Follow up with CloudBase - 2nd touch", Color: "red", When: "Today"},
	{ID: "task-due-2", Text: "Schedule demo for NovaBuild Co", Color: "yellow", When: "Tomorrow"},
	{ID: "task-due-3", Text: "Update pipeline notes - Orion Logistics", Color: "secondary", When: "Tomorrow"},
	{ID: "task-due-4", Text: "Prepare renewal deck for MetroPay", Color: "yellow", When: "May 4"},
	{ID: "task-due-5", Text: "Review churn signals - 3 at-risk accounts", Color: "red", When: "May 5"},
}

var crmLeaders = []dashboard.Leader{
	{Avatar: avatar.Props{Src: avatarURL("000m"), Color: "primary"}, Name: "Olivia Bennett", Percent: 92, Color: "primary", Value: "$148K", Detail: "18 deals - 92%"},
	{Avatar: avatar.Props{Src: avatarURL("052f"), Color: "purple"}, Name: "Marcus Lee", Percent: 78, Color: "purple", Value: "$124K", Detail: "15 deals - 78%"},
	{Avatar: avatar.Props{Src: avatarURL("002m"), Color: "green"}, Name: "Sara Mitchell", Percent: 65, Color: "green", Value: "$98K", Detail: "12 deals - 65%"},
	{Avatar: avatar.Props{Src: avatarURL("003m"), Color: "yellow"}, Name: "James Okonkwo", Percent: 54, Color: "yellow", Value: "$82K", Detail: "10 deals - 54%"},
	{Avatar: avatar.Props{Src: avatarURL("000f"), Color: "red"}, Name: "Lena Muller", Percent: 40, Color: "red", Value: "$61K", Detail: "8 deals - 40%"},
}

var crmRisks = []crm.Risk{
	{Name: "Orion Logistics", Note: "Last contact: 18 days ago", Level: "High", Color: "red", Value: "$24K"},
	{Name: "Summit Media", Note: "Last contact: 22 days ago", Level: "High", Color: "red", Value: "$18K"},
	{Name: "BluePeak Tech", Note: "Last contact: 12 days ago", Level: "Medium", Color: "yellow", Value: "$31K"},
	{Name: "Greenridge Partners", Note: "Last contact: 8 days ago", Level: "Low", Color: "green", Value: "$12K"},
	{Name: "MetroPay Ltd", Note: "Last contact: 5 days ago", Level: "Medium", Color: "yellow", Value: "$24K"},
}
