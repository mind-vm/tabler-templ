// Data of the marketing pages, written once from Tabler's testimonials.json, people.json and real-estate.json.

package main

import (
	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/blocks/marketing"
)

var mktTestimonials = []marketing.Testimonial{
	{Avatar: avatar.Props{Src: photo("avatars/052f.jpg")}, Name: "Jeffie Lewzey", Role: "Chemical Engineer", Text: "Ever since I started using Tabler, my HTML designing process has been significantly smoother. Its intuitive interface is a breath of fresh air!"},
	{Avatar: avatar.Props{Src: photo("avatars/002m.jpg")}, Name: "Mallory Hulme", Role: "Geologist IV", Text: "I'm beyond impressed with Tabler. This UI kit has transformed my workflow for the better. It's easy to use, streamlined, and remarkably efficient."},
	{Avatar: avatar.Props{Src: photo("avatars/003m.jpg")}, Name: "Dunn Slane", Role: "Research Nurse", Text: "Tabler's feature-rich package has drastically improved my website design process. It's an absolute game-changer in the UI kit landscape."},
	{Avatar: avatar.Props{Src: photo("avatars/000f.jpg")}, Name: "Emmy Levet", Role: "VP Product Management", Text: "Having used many UI kits in the past, I can confidently say that Tabler stands out. The level of customization it offers is unparalleled."},
	{Avatar: avatar.Props{Src: photo("avatars/001f.jpg")}, Name: "Maryjo Lebarree", Role: "Civil Engineer", Text: "What I love about Tabler is its sleek design and ease of use. It's really simplified my HTML designing work. Five stars from me!"},
	{Avatar: avatar.Props{Src: photo("avatars/004m.jpg")}, Name: "Egan Poetz", Role: "Research Nurse", Text: "I'd recommend Tabler to anyone, whether they're a seasoned pro or a beginner in HTML design. Its versatility and user-friendly nature make it a cut above the rest."},
	{Avatar: avatar.Props{Src: photo("avatars/002f.jpg")}, Name: "Kellie Skingley", Role: "Teacher", Text: "Navigating through Tabler has been a breeze! Its intuitive layout helps me save time and focus more on the creative aspect of my work."},
	{Avatar: avatar.Props{Src: photo("avatars/003f.jpg")}, Name: "Christabel Charlwood", Role: "Tax Accountant", Text: "With Tabler, designing HTML has never been more straightforward. The kit's sleek features and functionalities have taken my design work to new heights."},
	{Avatar: avatar.Props{Src: photo("avatars/063m.jpg")}, Name: "Haskel Shelper", Role: "Staff Scientist", Text: "Tabler has revolutionized the way I approach HTML design. Its array of tools and efficient interface make it a must-have for any designer."},
	{Avatar: avatar.Props{Src: photo("avatars/006m.jpg")}, Name: "Lorry Mion", Role: "Automation Specialist IV", Text: "I'm truly amazed by Tabler's capabilities. It's not just an HTML UI kit, it's a design powerhouse that has made my work a lot more efficient and enjoyable."},
	{Avatar: avatar.Props{Src: photo("avatars/004f.jpg")}, Name: "Leesa Beaty", Role: "Editor", Text: "As a web developer, Tabler has revolutionized my workflow. The HTML UI kit is exceptionally well-structured, making my projects much more efficient and visually appealing."},
	{Avatar: avatar.Props{Src: photo("avatars/007m.jpg")}, Name: "Perren Keemar", Role: "Analog Circuit Design manager", Text: "Tabler has transformed the way I approach web design. It's intuitive, easy to use, and has saved me countless hours of coding from scratch. A real game-changer!"},
	{Avatar: avatar.Props{Initials: "SA"}, Name: "Sunny Airey", Role: "Nuclear Power Engineer", Text: "The breadth of design options that Tabler provides is astounding. It has added a whole new level of professionalism to my websites and made my job a lot easier."},
	{Avatar: avatar.Props{Src: photo("avatars/009m.jpg")}, Name: "Geoffry Flaunders", Role: "Software Test Engineer II", Text: "I've been in the web development industry for years and I must say, Tabler is one of the best HTML UI kits I've come across. It's user-friendly and allows for a high degree of customization."},
	{Avatar: avatar.Props{Src: photo("avatars/010m.jpg")}, Name: "Thatcher Keel", Role: "VP Sales", Text: "Since I've started using Tabler, I've noticed a significant improvement in my work. It's streamlined, efficient, and helps me produce better results for my clients."},
	{Avatar: avatar.Props{Src: photo("avatars/005f.jpg")}, Name: "Dyann Escala", Role: "Mechanical Systems Engineer", Text: "Using Tabler has been a fantastic experience. The number of components it offers has allowed me to design more efficiently, and its consistent updates ensure that it always stays ahead of the curve."},
	{Avatar: avatar.Props{Src: photo("avatars/006f.jpg")}, Name: "Avivah Mugleston", Role: "Actuary", Text: "Tabler is a web designer's dream come true. The HTML UI kit offers so many features and functionalities that make my design process more streamlined and efficient."},
	{Avatar: avatar.Props{Initials: "AA"}, Name: "Arlie Armstead", Role: "VP Quality Control", Text: "I cannot recommend Tabler enough! This HTML UI kit is a must-have for any web developer looking to level up their design game."},
	{Avatar: avatar.Props{Src: photo("avatars/008f.jpg")}, Name: "Tessie Curzon", Role: "Research Nurse", Text: "From a beginner's perspective, Tabler has been a lifesaver. It has enabled me to create professional-quality websites without needing advanced coding skills."},
	{Avatar: avatar.Props{Src: photo("avatars/009f.jpg")}, Name: "Flossi Uttley", Role: "Engineer I", Text: "Tabler offers the perfect balance of simplicity and sophistication. It has greatly improved my workflow and I can't imagine working on web development projects without it."},
}

var mktProperties = []marketing.Property{
	{Name: "Cozy Cottage", Image: photo("marketing/photo-1558036117-15d82a90b9b1.jpg"), Price: "$250,000", Address: "123 Main Street, California USA", Beds: 3, Baths: 2},
	{Name: "Modern Townhouse", Image: photo("marketing/photo-1494526585095-c41746248156.jpg"), Price: "$450,000", Address: "456 Elm Avenue, California USA", Beds: 4, Baths: 3},
	{Name: "Spacious Villa", Image: photo("marketing/photo-1512917774080-9991f1c4c750.jpg"), Price: "$1,200,000", Address: "789 Oak Drive, California USA", Beds: 6, Baths: 5},
	{Name: "Rustic Farmhouse", Image: photo("marketing/photo-1564013799919-ab600027ffc6.jpg"), Price: "$350,000", Address: "321 Pine Lane, California USA", Beds: 2, Baths: 2},
	{Name: "Beachfront Condo", Image: photo("marketing/photo-1580587771525-78b9dba3b914.jpg"), Price: "$500,000", Address: "567 Sandy Beach Road, California USA", Beds: 1, Baths: 2},
	{Name: "Luxury Penthouse", Image: photo("marketing/photo-1592595896551-12b371d546d5.jpg"), Price: "$2,000,000", Address: "890 Highrise Avenue, California USA", Beds: 3, Baths: 4},
}
