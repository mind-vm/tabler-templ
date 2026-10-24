// Data of Tabler's activity and task pages (shared/data/people.json, tasks.json and the activity feed of
// preview/dist/activity.html), written once by a script.

package main

import (
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/blocks/activity"
	"github.com/mind-vm/tabler-templ/blocks/tasks"
)

// cPerson is an entry of people.json: the full name and the photo (empty for initials).
type cPerson struct{ Name, Photo string }

var cPeople = []cPerson{
	{"Paweł Kuna", "000m"},
	{"Jeffie Lewzey", "052f"},
	{"Mallory Hulme", "002m"},
	{"Dunn Slane", "003m"},
	{"Emmy Levet", "000f"},
	{"Maryjo Lebarree", "001f"},
	{"Egan Poetz", "004m"},
	{"Kellie Skingley", "002f"},
	{"Christabel Charlwood", "003f"},
	{"Haskel Shelper", "063m"},
	{"Lorry Mion", "006m"},
	{"Leesa Beaty", "004f"},
	{"Perren Keemar", "007m"},
	{"Sunny Airey", ""},
	{"Geoffry Flaunders", "009m"},
	{"Thatcher Keel", "010m"},
	{"Dyann Escala", "005f"},
	{"Avivah Mugleston", "006f"},
	{"Arlie Armstead", ""},
	{"Tessie Curzon", "008f"},
	{"Flossi Uttley", "009f"},
	{"Cesya Spritt", "010f"},
	{"Johnnie Gilby", "011m"},
	{"Ban Rzehor", "012m"},
	{"Carroll Erat", "011f"},
	{"Marsha Labat", "012f"},
	{"Elston Muffett", "013m"},
	{"Leigha Gorce", "013f"},
	{"Tallie Bettis", ""},
	{"Merrily Garforth", "015f"},
	{"Errol Blackley", ""},
	{"Ninon Don", "016f"},
	{"Delaney Cairney", "015m"},
	{"Gratia Gooley", "017f"},
	{"Odelinda McCosh", "018f"},
	{"Wilburt Siegertsz", "016m"},
	{"Julietta Coke", "019f"},
	{"Portie Christou", "017m"},
	{"Emmott Dowsett", "018m"},
	{"Rooney Cassy", "019m"},
	{"Haze Hubbert", "020m"},
	{"Mata Codlin", "021m"},
	{"Parker Oaten", "022m"},
	{"Johannes Paternoster", "023m"},
	{"Cary Baleine", "024m"},
	{"Riane Milward", "020f"},
	{"Reynold Indgs", "025m"},
	{"Parke Moneypenny", "026m"},
	{"Sandi Keys", "021f"},
	{"Peria Errichiello", "022f"},
}

func (p cPerson) avatar() avatar.Props {
	if p.Photo != "" {
		return avatar.Props{Src: avatarURL(p.Photo)}
	}
	return avatar.Props{Initials: cInitials(p.Name)}
}

// cUsers are n people from offset, like the preview's PeopleAvatarList.
func cUsers(offset, n int) []avatar.Props {
	var out []avatar.Props
	for _, p := range cPeople[offset : offset+n] {
		out = append(out, p.avatar())
	}
	return out
}

func cAssignee(id int) (*avatar.Props, string) {
	a := cPeople[id-1].avatar()
	return &a, cPeople[id-1].Name
}

func cTask(t tasks.Task, assignedTo int) tasks.Task {
	if assignedTo > 0 {
		t.Assignee, t.AssigneeName = cAssignee(assignedTo)
	}
	return t
}

var cColumns = []tasks.Column{
	{Name: "Upcoming", Tasks: []tasks.Task{
		cTask(tasks.Task{Name: "Improve email marketing strategy", Image: photo("projects/dashboard-1.png"), Users: cUsers(4, 1), Due: "August 1, 2023", Likes: 7, Comments: 2, Priority: "Medium"}, 5),
		cTask(tasks.Task{Name: "Develop new product video", Description: "Dedicated form for a category of users that will perform actions.", Users: cUsers(1, 1), Due: "July 15, 2023", Likes: 6, Priority: "High"}, 2),
		cTask(tasks.Task{Name: "Conduct user interviews for new feature", Color: "yellow", Due: "June 20, 2023", Favorite: true, Likes: 34, Comments: 4, Priority: "Low"}, 3),
		cTask(tasks.Task{Name: "Create social media content calendar", Description: "Plan and schedule social media posts for Q3", Users: cUsers(5, 2), Due: "September 5, 2023", Likes: 12, Comments: 8, Priority: "Medium"}, 6),
		cTask(tasks.Task{Name: "Review and update documentation", Users: cUsers(3, 1), Due: "July 25, 2023", Comments: 5, Priority: "Medium"}, 4),
		cTask(tasks.Task{Name: "Plan quarterly team meeting", Users: cUsers(8, 3), Due: "August 20, 2023", Favorite: true, Likes: 9, Priority: "Low"}, 7),
		cTask(tasks.Task{Name: "Redesign product landing page", Description: "Complete redesign based on user feedback", Users: cUsers(6, 2), Due: "September 1, 2023", Likes: 23, Comments: 15, Priority: "High"}, 8),
	}},
	{Name: "In Progress", Tasks: []tasks.Task{
		cTask(tasks.Task{Name: "Implement new analytics tracking", Color: "green", Due: "July 1, 2023", Comments: 11, Subtasks: []tasks.Subtask{{Name: "Implement new designs", Done: true}, {Name: "Usability testing", Done: false}, {Name: "Design navigation changes", Done: false}}, Priority: "Low"}, 2),
		cTask(tasks.Task{Name: "Design new marketing campaign", Description: "Dedicated form for a category of users that will perform actions.", Starred: true, Users: cUsers(7, 2), Due: "August 15, 2023", Likes: 11, Comments: 6, Priority: "High"}, 3),
		cTask(tasks.Task{Name: "Conduct A/B testing on landing page", Image: photo("projects/dashboard-3.png"), Due: "June 30, 2023", Comments: 11, Priority: "Low"}, 5),
		cTask(tasks.Task{Name: "Implement user authentication system", Description: "Add OAuth and JWT authentication", Color: "red", Users: cUsers(4, 3), Due: "July 20, 2023", Likes: 14, Comments: 18, Priority: "High"}, 9),
		cTask(tasks.Task{Name: "Write blog post about product launch", Users: cUsers(9, 1), Due: "August 5, 2023", Likes: 10, Comments: 7, Priority: "Medium"}, 10),
		cTask(tasks.Task{Name: "Update customer onboarding flow", Starred: true, Users: cUsers(2, 2), Due: "July 12, 2023", Likes: 8, Comments: 12, Priority: "Medium"}, 11),
		cTask(tasks.Task{Name: "Setup automated testing pipeline", Color: "blue", Users: cUsers(10, 4), Due: "August 25, 2023", Likes: 16, Comments: 20, Subtasks: []tasks.Subtask{{Name: "Configure CI/CD", Done: true}, {Name: "Write unit tests", Done: false}, {Name: "Setup code coverage", Done: false}}, Priority: "High"}, 12),
		cTask(tasks.Task{Name: "Review customer support tickets", Users: cUsers(11, 1), Due: "July 8, 2023", Comments: 4, Priority: "Low"}, 13),
	}},
	{Name: "Completed", Tasks: []tasks.Task{
		cTask(tasks.Task{Name: "Optimize website performance", Users: cUsers(7, 1), Due: "June 15, 2023", Likes: 1, Priority: "Low"}, 3),
		cTask(tasks.Task{Name: "Develop mobile app prototype", Description: "This content is a little bit longer.", Users: cUsers(0, 5), Due: "August 10, 2023", Favorite: true, Likes: 7, Subtasks: []tasks.Subtask{{Name: "Find out the old contract documents", Done: true}, {Name: "Organize meeting sales associates to understand need in detail", Done: false}, {Name: "Make sure to cover every small details", Done: false}}, Priority: "Medium"}, 5),
		cTask(tasks.Task{Name: "Conduct user research interviews", DescriptionContent: templ.Raw("Publishing industries for previewing layouts and visual <a href=\"#\">#family</a> 🔥"), Users: cUsers(10, 2), Due: "July 20, 2023", Likes: 4, Priority: "Low"}, 5),
		cTask(tasks.Task{Name: "Migrate database to new server", Description: "Complete database migration with zero downtime", Users: cUsers(12, 5), Due: "June 10, 2023", Likes: 18, Comments: 25, Priority: "Medium"}, 14),
		cTask(tasks.Task{Name: "Launch new marketing website", Image: photo("projects/dashboard-2.png"), Users: cUsers(14, 3), Due: "May 30, 2023", Favorite: true, Likes: 45, Comments: 30, Priority: "High"}, 15),
		cTask(tasks.Task{Name: "Complete annual security audit", Color: "red", Users: cUsers(15, 4), Due: "May 25, 2023", Likes: 12, Comments: 22, Priority: "High"}, 16),
		cTask(tasks.Task{Name: "Design new company logo", Users: cUsers(16, 2), Due: "April 15, 2023", Likes: 21, Comments: 9, Priority: "Low"}, 17),
		cTask(tasks.Task{Name: "Write technical documentation", Description: "Document all API endpoints and integrations", Users: cUsers(17, 1), Due: "June 5, 2023", Likes: 7, Comments: 14, Subtasks: []tasks.Subtask{{Name: "Document REST API", Done: true}, {Name: "Document GraphQL API", Done: true}, {Name: "Create code examples", Done: false}}, Priority: "Medium"}, 18),
		cTask(tasks.Task{Name: "Organize team building event", Users: cUsers(18, 6), Due: "April 20, 2023", Favorite: true, Likes: 52, Comments: 35, Priority: "Low"}, 19),
	}},
}

var cActivity = []activity.Item{
	{Avatar: avatar.Props{Src: avatarURL("052f")}, Text: templ.Raw("<strong>Jeffie Lewzey</strong> commented on your post: <strong>Not everything is what it seems</strong>"), Time: "1 day ago", Unread: true},
	{Avatar: avatar.Props{Src: avatarURL("002m")}, Text: templ.Raw("It's <strong>Mallory Hulme</strong>'s birthday today"), Time: "now", Unread: true},
	{Avatar: avatar.Props{Src: avatarURL("003m")}, Text: templ.Raw("<strong>Dunn Slane</strong> posted a status: <strong>Trying to figure things out one step at a time.</strong>"), Time: "now", Unread: true},
	{Avatar: avatar.Props{Src: avatarURL("000f")}, Text: templ.Raw("<strong>Emmy Levet</strong> created a new project: <strong>Smart Morning Alarm</strong>"), Time: "4 days ago", Unread: true},
	{Avatar: avatar.Props{Src: avatarURL("001f")}, Text: templ.Raw("<strong>Maryjo Lebarree</strong> liked your photo"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("004m")}, Text: templ.Raw("<strong>Egan Poetz</strong> registered a new client: <strong>Trilia Solutions</strong>"), Time: "1 day ago"},
	{Avatar: avatar.Props{Src: avatarURL("002f")}, Text: templ.Raw("<strong>Kellie Skingley</strong> closed a new deal on project: <strong>PPAP Campaign Strategy</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("003f")}, Text: templ.Raw("<strong>Christabel Charlwood</strong> created a new project: <strong>Wikibox Knowledge Platform</strong>"), Time: "4 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("063m")}, Text: templ.Raw("<strong>Haskel Shelper</strong> changed project status of <strong>Tabler Icons</strong> from <strong>open</strong> to <strong>closed</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("006m")}, Text: templ.Raw("<strong>Lorry Mion</strong> liked the <strong>Tabler UI Kit</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("004f")}, Text: templ.Raw("<strong>Leesa Beaty</strong> shared a new video: <strong>Design Tips for Busy Creators</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("007m")}, Text: templ.Raw("<strong>Perren Keemar</strong> and <strong>3 others</strong> followed you"), Time: "now"},
	{Avatar: avatar.Props{Initials: "SA"}, Text: templ.Raw("<strong>Sunny Airey</strong> uploaded <strong>3 new photos</strong> to category: <strong>Creative Inspirations</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("009m")}, Text: templ.Raw("<strong>Geoffry Flaunders</strong> made a <strong>$10 donation</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("010m")}, Text: templ.Raw("<strong>Thatcher Keel</strong> created a profile"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("005f")}, Text: templ.Raw("<strong>Dyann Escala</strong> hosted the event: <strong>Tabler UI Anniversary</strong>"), Time: "4 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("006f")}, Text: templ.Raw("<strong>Avivah Mugleston</strong> mentioned you in <strong>Best of 2020 Highlights</strong>"), Time: "now"},
	{Avatar: avatar.Props{Initials: "AA"}, Text: templ.Raw("<strong>Arlie Armstead</strong> sent a review request to <strong>Amanda Blake</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("008f")}, Text: templ.Raw("<strong>Tessie Curzon</strong> posted a photo from her vacation in <strong>Tuscany</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("009f")}, Text: templ.Raw("<strong>Flossi Uttley</strong> shared an article: <strong>Top 10 productivity hacks</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("010f")}, Text: templ.Raw("<strong>Cesya Spritt</strong> commented on your photo"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("011m")}, Text: templ.Raw("<strong>Johnnie Gilby</strong> joined the group: <strong>UX Designers Hub</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("012m")}, Text: templ.Raw("<strong>Ban Rzehor</strong> uploaded new design mockups"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("011f")}, Text: templ.Raw("<strong>Carroll Erat</strong> reacted with ❤️ to your story"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("012f")}, Text: templ.Raw("<strong>Marsha Labat</strong> posted a video: <strong>Building Better Dashboards</strong>"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("013m")}, Text: templ.Raw("<strong>Elston Muffett</strong> joined your webinar: <strong>Intro to TypeScript</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("013f")}, Text: templ.Raw("<strong>Leigha Gorce</strong> commented: <strong>This is so helpful!</strong>"), Time: "4 days ago"},
	{Avatar: avatar.Props{Initials: "TB"}, Text: templ.Raw("<strong>Tallie Bettis</strong> updated her profile photo"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("015f")}, Text: templ.Raw("<strong>Merrily Garforth</strong> posted a job listing: <strong>Frontend Developer</strong>"), Time: "now"},
	{Avatar: avatar.Props{Initials: "EB"}, Text: templ.Raw("<strong>Errol Blackley</strong> started following your project: <strong>Clean Code Patterns</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("016f")}, Text: templ.Raw("<strong>Ninon Don</strong> gave feedback on your latest design"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("015m")}, Text: templ.Raw("<strong>Delaney Cairney</strong> shared a playlist: <strong>Coding Flow</strong>"), Time: "4 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("017f")}, Text: templ.Raw("<strong>Gratia Gooley</strong> joined the event: <strong>Night UI Meetup</strong>"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("018f")}, Text: templ.Raw("<strong>Odelinda McCosh</strong> shared a memory: <strong>2 years at Tabler!</strong>"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("016m")}, Text: templ.Raw("<strong>Wilburt Siegertsz</strong> retweeted your latest post"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("019f")}, Text: templ.Raw("<strong>Julietta Coke</strong> shared your story to their network"), Time: "3 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("017m")}, Text: templ.Raw("<strong>Portie Christou</strong> invited you to collaborate on her workspace"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("018m")}, Text: templ.Raw("<strong>Emmott Dowsett</strong> published a new blog post: <strong>Dark Mode for All</strong>"), Time: "2 days ago"},
	{Avatar: avatar.Props{Src: avatarURL("019m")}, Text: templ.Raw("<strong>Rooney Cassy</strong> mentioned you in a comment"), Time: "now"},
	{Avatar: avatar.Props{Src: avatarURL("020m")}, Text: templ.Raw("<strong>Haze Hubbert</strong> shared a new design: <strong>Modern Dashboard UI</strong>"), Time: "now"},
}

// cInitials are the first letters of the words of a name, like the preview's firstLetters.
func cInitials(name string) string {
	s := ""
	for _, w := range strings.Fields(name) {
		s += string([]rune(w)[:1])
	}
	return s
}
