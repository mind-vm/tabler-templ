// Data of the plugin pages, written once from Tabler's shared/data (photos.json, maps-vector.json, maps.json).

package main

import (
	"github.com/mind-vm/tabler-templ/mapview"
	"github.com/mind-vm/tabler-templ/vectormap"
)

// lightboxPhotos are the horizontal photos of photos.json: file and title.
var lightboxPhotos = [][2]string{
	{"beautiful-blonde-woman-relaxing-with-a-can-of-coke-on-a-tree-stump-by-the-beach.jpg", "Beautiful blonde woman relaxing with a can of coke on a tree stump by the beach"},
	{"brainstorming-session-with-creative-designers.jpg", "Brainstorming session with creative designers"},
	{"finances-us-dollars-and-bitcoins-currency-money.jpg", "Finances - US Dollars and Bitcoins - Currency - Money"},
	{"group-of-people-brainstorming-and-taking-notes-2.jpg", "Group of people brainstorming and taking notes"},
	{"blue-sofa-with-pillows-in-a-designer-living-room-interior.jpg", "Blue sofa with pillows in a designer living room interior"},
	{"home-office-desk-with-macbook-iphone-calendar-watch-and-organizer.jpg", "Home office desk with Macbook, iPhone, calendar, watch & organizer"},
	{"young-woman-working-in-a-cafe.jpg", "Young woman working in a cafe"},
	{"everything-you-need-to-work-from-your-bed.jpg", "Everything you need to work from your bed"},
	{"young-entrepreneur-working-from-a-modern-cafe.jpg", "Young Entrepreneur Working from a Modern Cafe"},
	{"finances-us-dollars-and-bitcoins-currency-money-2.jpg", "Finances - US Dollars and Bitcoins - Currency - Money"},
	{"group-of-people-sightseeing-in-the-city.jpg", "Group of people sightseeing in the city"},
	{"color-palette-guide-sample-colors-catalog-.jpg", "Color Palette Guide. Sample Colors Catalog."},
	{"stylish-workplace-with-computer-at-home.jpg", "Stylish workplace with computer at home"},
	{"pink-desk-in-the-home-office.jpg", "Pink desk in the home office"},
	{"young-woman-sitting-on-the-sofa-and-working-on-her-laptop.jpg", "Young woman sitting on the sofa and working on her laptop"},
	{"coffee-on-a-table-with-other-items.jpg", "Coffee on a table with other items"},
	{"young-entrepreneur-working-from-a-modern-cafe-2.jpg", "Young Entrepreneur Working from a Modern Cafe"},
	{"soft-photo-of-woman-on-the-bed-with-the-book-and-cup-of-coffee-in-hands.jpg", "Soft photo of woman on the bed with the book and cup of coffee in hands"},
	{"fairy-lights-at-the-beach-in-bulgaria.jpg", "Fairy lights at the beach in Bulgaria"},
	{"woman-working-on-laptop-at-home-office.jpg", "Woman working on laptop at home office"},
	{"modern-home-office.jpg", "Modern Home Office"},
	{"blond-using-her-laptop-at-her-bedroom.jpg", "Blond using her laptop at her bedroom"},
	{"tropical-palm-leaves-floral-pattern-background.jpg", "Tropical palm leaves, floral pattern background"},
	{"woman-read-book-and-drink-coffee.jpg", "Woman read book and drink coffee"},
	{"book-on-the-grass.jpg", "Book on the grass"},
	{"colorful-exotic-flowers-and-greenery.jpg", "Colorful Exotic Flowers and Greenery"},
	{"workplace-with-laptop-on-table-at-home.jpg", "Workplace with laptop on table at home"},
	{"city-lights-reflected-in-the-water-at-night.jpg", "City lights reflected in the water at night"},
	{"cryptocurrency-bitcoin-coins.jpg", "Cryptocurrency Bitcoin coins"},
	{"woman-read-book-and-drink-coffee-2.jpg", "Woman read book and drink coffee"},
	{"friends-at-a-restaurant-drinking-wine.jpg", "Friends at a restaurant drinking wine"},
	{"beautiful-blonde-woman-on-a-wooden-pier-by-the-lake.jpg", "Beautiful blonde woman on a wooden pier by the lake"},
	{"white-apple-imac-computer-with-elephant-mousepad.jpg", "White Apple iMac computer with elephant mousepad"},
	{"workplace-with-laptop-on-table-at-home-2.jpg", "Workplace with laptop on table at home"},
	{"working-in-a-restaurant-macbook-cheese-cake-and-cup-of-coffee.jpg", "Working in a restaurant: Macbook, Cheese Cake and Cup of Coffee"},
	{"young-woman-sitting-on-the-sofa-and-working-on-her-laptop-2.jpg", "Young woman sitting on the sofa and working on her laptop"},
	{"brainstorming-session-with-creative-designers-2.jpg", "Brainstorming session with creative designers"},
	{"woman-drinking-hot-tea-in-her-home-office.jpg", "Woman drinking hot tea in her home office"},
	{"stylish-workspace-with-macbook-pro.jpg", "Stylish workspace with Macbook Pro"},
	{"people-by-a-banquet-table-full-with-food.jpg", "People by a banquet table full with food"},
	{"woman-working-on-a-laptop-while-enjoying-a-breakfast-coffee-and-chocolate-in-bed.jpg", "Woman working on a laptop while enjoying a breakfast coffee and chocolate in bed"},
	{"businesswoman-working-at-her-laptop.jpg", "Businesswoman Working at Her Laptop"},
	{"finances-us-dollars-and-bitcoins-currency-money-5.jpg", "Finances - US Dollars and Bitcoins - Currency - Money"},
	{"cup-of-coffee-on-table-in-cafe.jpg", "Cup of coffee on table in cafe"},
	{"young-woman-sitting-on-the-sofa-and-working-on-her-laptop-3.jpg", "Young woman sitting on the sofa and working on her laptop"},
	{"man-looking-out-to-sea.jpg", "Man looking out to sea"},
	{"cup-of-coffee-on-table-in-cafe-2.jpg", "Cup of coffee on table in cafe"},
	{"workplace-with-laptop-on-table-at-home-3.jpg", "Workplace with laptop on table at home"},
	{"businesswoman-working-at-her-laptop-2.jpg", "Businesswoman Working at Her Laptop"},
	{"cup-of-coffee-and-an-open-book.jpg", "Cup of Coffee and an Open Book"},
	{"group-of-people-brainstorming-and-taking-notes-4.jpg", "Group of people brainstorming and taking notes"},
	{"young-woman-sitting-on-the-sofa-and-working-on-her-laptop-4.jpg", "Young woman sitting on the sofa and working on her laptop"},
	{"making-magic-with-fairy-lights.jpg", "Making Magic with Fairy Lights"},
	{"woman-working-on-a-laptop-while-enjoying-a-breakfast-coffee-and-chocolate-in-bed-2.jpg", "Woman working on a laptop while enjoying a breakfast coffee and chocolate in bed"},
	{"woman-drinking-tea-and-reading-book.jpg", "Woman drinking tea and reading book"},
	{"a-woman-works-at-a-desk-with-a-laptop-and-a-cup-of-coffee.jpg", "A woman works at a desk with a laptop and a cup of coffee"},
	{"everything-you-need-to-work-from-your-bed-2.jpg", "Everything you need to work from your bed"},
	{"workplace-with-laptop-on-table-at-home-4.jpg", "Workplace with laptop on table at home"},
	{"stylish-workspace-with-macbook-pro-2.jpg", "Stylish workspace with Macbook Pro"},
	{"netflix-drug-lords-from-narcos.jpg", "Netflix Drug Lords From Narcos"},
	{"stone-texture-high-resolution-background-2.jpg", "Stone texture - high resolution background"},
	{"a-visit-to-the-bookstore.jpg", "A visit to the bookstore"},
	{"books-and-purple-flowers-on-a-wooden-stool-by-the-bed.jpg", "Books and purple flowers on a wooden stool by the bed"},
	{"people-watching-a-presentation-in-a-room.jpg", "People watching a presentation in a room"},
	{"home-office-laptop-organizer-and-cup-of-coffee.jpg", "Home Office - Laptop - organizer & cup of coffee"},
}

// vectorMaps are the maps of maps-vector.json by id.
var vectorMaps = map[string]vectormap.Props{
	"empty":         {Map: "world", Label: "Empty map"},
	"world":         {Map: "world", Label: "World map", Values: worldValues},
	"world-markers": {Map: "world", Label: "Map with markers", Color: "red", Markers: []vectormap.Marker{{Name: "Russia", Lat: 61.524, Lng: 105.3188}, {Name: "Canada", Lat: 56.1304, Lng: -106.3468}, {Name: "Greenland", Lat: 71.7069, Lng: -42.6043}, {Name: "Egypt", Lat: 26.8206, Lng: 30.8025}, {Name: "Brazil", Lat: -14.235, Lng: -51.9253}, {Name: "China", Lat: 35.8617, Lng: 104.1954}, {Name: "United States", Lat: 37.0902, Lng: -95.7129}, {Name: "Norway", Lat: 60.472024, Lng: 8.468946}, {Name: "Ukraine", Lat: 48.379433, Lng: 31.16558}}},
	"world-lines":   {Map: "world", Label: "Map with markers and lines", Markers: []vectormap.Marker{{Name: "Russia", Lat: 61.524, Lng: 105.3188}, {Name: "Canada", Lat: 56.1304, Lng: -106.3468}, {Name: "Greenland", Lat: 71.7069, Lng: -42.6043}, {Name: "Egypt", Lat: 26.8206, Lng: 30.8025}, {Name: "Brazil", Lat: -14.235, Lng: -51.9253}, {Name: "China", Lat: 35.8617, Lng: 104.1954}, {Name: "United States", Lat: 37.0902, Lng: -95.7129}, {Name: "Norway", Lat: 60.472024, Lng: 8.468946}, {Name: "Ukraine", Lat: 48.379433, Lng: 31.16558}}, Lines: []vectormap.Line{{From: "Russia", To: "Greenland"}, {From: "Russia", To: "United States"}, {From: "Russia", To: "Canada"}, {From: "Brazil", To: "Norway"}, {From: "Brazil", To: "Ukraine"}, {From: "Brazil", To: "Egypt"}, {From: "Brazil", To: "China"}}},
	"world-zoom":    {Map: "world", Label: "Map with zoom", Color: "green", Zoom: true, Values: map[string]float64{"PL": 8, "DE": 6, "FR": 5, "ES": 3, "IT": 4, "GB": 7, "SE": 2, "UA": 3, "RO": 2, "CZ": 5}},
	"world-colors":  {Map: "world", Label: "Map with a color scale", Legend: true, Colors: []string{"var(--tblr-red)", "var(--tblr-yellow)", "var(--tblr-green)"}, Values: worldValues},
	"world-select":  {Map: "world", Label: "Map with selection", Select: "multiple", Selected: []string{"PL", "FR", "BR", "AU"}},
}

// libreMaps are the maps of maps.json by id (the data keeps latitude first; LngLat takes longitude first).
var libreMaps = map[string]struct {
	Props   mapview.Props
	Card    bool
	Markers []mapview.MarkerProps
}{
	"simple":  {Props: mapview.Props{ID: "simple", Label: "Simple map", Center: &mapview.LngLat{Lng: 13.4049, Lat: 52.518827}, Zoom: 13, Controls: "true"}, Card: false, Markers: []mapview.MarkerProps{}},
	"light":   {Props: mapview.Props{ID: "light", Label: "Street map style", Center: &mapview.LngLat{Lng: -0.2416782, Lat: 51.5285582}, Zoom: 11, Style: "liberty"}, Card: false, Markers: []mapview.MarkerProps{}},
	"markers": {Props: mapview.Props{ID: "markers", Label: "Map with markers", Center: &mapview.LngLat{Lng: 0, Lat: 0}, Zoom: 1, Style: "positron"}, Card: true, Markers: []mapview.MarkerProps{{Kind: mapview.Dot, Lng: -58.666667, Lat: -34.58333333, Popup: "Buenos Aires"}, {Kind: mapview.Dot, Lng: 16.366667, Lat: 48.2, Popup: "Vienna"}, {Kind: mapview.Dot, Lng: 116.383333, Lat: 39.91666667, Popup: "Beijing"}, {Kind: mapview.Dot, Lng: 149.133333, Lat: -35.26666667, Popup: "Canberra"}, {Kind: mapview.Dot, Lng: 58.583333, Lat: 23.61666667, Popup: "Muscat"}, {Kind: mapview.Dot, Lng: -77, Lat: 38.883333, Popup: "Washington"}}},
	"card":    {Props: mapview.Props{ID: "card", Label: "Full-card map", Center: &mapview.LngLat{Lng: 0, Lat: 0}, Zoom: 1, Style: "bright", DarkStyle: "fiord"}, Card: true, Markers: []mapview.MarkerProps{}},
}
