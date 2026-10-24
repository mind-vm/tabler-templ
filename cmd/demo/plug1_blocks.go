package main

import (
	"math"

	"github.com/mind-vm/tabler-templ/blocks/maps"
	"github.com/mind-vm/tabler-templ/mapview"
)

// Data of the map pages, from the Tabler preview's maps.astro and shared/components/cards/Map*.astro.

var berlinRoute = []mapview.LngLat{{Lng: 13.3695, Lat: 52.5251}, {Lng: 13.3777, Lat: 52.5208}, {Lng: 13.388, Lat: 52.5186},
	{Lng: 13.3989, Lat: 52.5176}, {Lng: 13.405, Lat: 52.5192}, {Lng: 13.4132, Lat: 52.5219}, {Lng: 13.4215, Lat: 52.5165},
	{Lng: 13.43, Lat: 52.5105}, {Lng: 13.4397, Lat: 52.505}}

var berlinArea = [][2]float64{{13.37, 52.53}, {13.42, 52.535}, {13.445, 52.52}, {13.43, 52.505}, {13.385, 52.503}, {13.365, 52.515}, {13.37, 52.53}}

var mapCities = []mapview.Bubble{{Name: "Berlin", Lng: 13.405, Lat: 52.52, Value: 820}, {Name: "London", Lng: -0.1276, Lat: 51.5072, Value: 1240},
	{Name: "New York", Lng: -74.006, Lat: 40.7128, Value: 1580}, {Name: "Tokyo", Lng: 139.6917, Lat: 35.6895, Value: 960},
	{Name: "São Paulo", Lng: -46.6333, Lat: -23.5505, Value: 540}, {Name: "Sydney", Lng: 151.2093, Lat: -33.8688, Value: 310},
	{Name: "Lagos", Lng: 3.3792, Lat: 6.5244, Value: 420}, {Name: "Mumbai", Lng: 72.8777, Lat: 19.076, Value: 1100}}

type mapPlace struct {
	Name, Text, Photo string
	Lng, Lat          float64
	Open              bool
}

var mapPlaces = []mapPlace{
	{Name: "Museum Island", Text: "Five museums on the Spree, open daily.", Photo: "a-visit-to-the-bookstore.jpg", Lat: 52.5169, Lng: 13.4019, Open: true},
	{Name: "Bellevue Palace", Text: "Residence by the Tiergarten park.", Photo: "book-on-the-grass.jpg", Lat: 52.5176, Lng: 13.3531},
	{Name: "Weissensee", Text: "A lake with a beach in the north-east.", Photo: "brainstorming-session-with-creative-designers.jpg", Lat: 52.5546, Lng: 13.4633},
}

var globeHub = mapview.LngLat{Lng: -0.1276, Lat: 51.5072}

var globeDestinations = []struct {
	Name string
	At   mapview.LngLat
}{{Name: "Tokyo", At: mapview.LngLat{Lng: 139.6917, Lat: 35.6895}}, {Name: "Mumbai", At: mapview.LngLat{Lng: 72.8777, Lat: 19.076}},
	{Name: "Singapore", At: mapview.LngLat{Lng: 103.8198, Lat: 1.3521}}, {Name: "Nairobi", At: mapview.LngLat{Lng: 36.8219, Lat: -1.2921}}}

func globeArcs() []mapview.Arc {
	var arcs []mapview.Arc
	for _, d := range globeDestinations {
		arcs = append(arcs, mapview.Arc{From: globeHub, To: d.At})
	}
	return arcs
}

// storePoints are stores scattered around a few European cities, with the preview's fixed offsets.
func storePoints() []mapview.LngLat {
	cities := []mapview.LngLat{{Lng: 13.405, Lat: 52.52}, {Lng: 2.3522, Lat: 48.8566}, {Lng: -0.1276, Lat: 51.5072}, {Lng: 12.4964, Lat: 41.9028},
		{Lng: -3.7038, Lat: 40.4168}, {Lng: 21.0122, Lat: 52.2297}, {Lng: 16.3738, Lat: 48.2082}, {Lng: 4.9041, Lat: 52.3676}}
	var out []mapview.LngLat
	for c, p := range cities {
		for i := 0; i < 5+c%4; i++ {
			out = append(out, mapview.LngLat{Lng: p.Lng + math.Sin(float64(c*7+i*3))*0.6, Lat: p.Lat + math.Cos(float64(c*5+i*2))*0.4})
		}
	}
	return out
}

var stores = []maps.Store{
	{Name: "Tabler Store Mitte", Address: "Friedrichstraße 112, 10117 Berlin", Distance: "0.4 km", Open: true, Lng: 13.3889, Lat: 52.5246},
	{Name: "Tabler Store Alexanderplatz", Address: "Karl-Liebknecht-Straße 9, 10178 Berlin", Distance: "1.2 km", Open: true, Lng: 13.4094, Lat: 52.5219},
	{Name: "Tabler Store Kreuzberg", Address: "Oranienstraße 34, 10999 Berlin", Distance: "2.6 km", Open: true, Lng: 13.4197, Lat: 52.5009},
	{Name: "Tabler Store Charlottenburg", Address: "Kantstraße 54, 10627 Berlin", Distance: "4.8 km", Lng: 13.3087, Lat: 52.5063},
	{Name: "Tabler Store Prenzlauer Berg", Address: "Kastanienallee 21, 10435 Berlin", Distance: "2.1 km", Open: true, Lng: 13.4098, Lat: 52.5385},
	{Name: "Tabler Store Friedrichshain", Address: "Warschauer Straße 58, 10243 Berlin", Distance: "3.9 km", Lng: 13.4492, Lat: 52.5083},
}

func homes() []maps.Home {
	return []maps.Home{
		{Title: "Loft in Mitte", Price: "$420k", Beds: 2, Baths: 1, Area: "78 m²", Photo: photo("photos/contemporary-black-and-white-home-decor.jpg"), Lng: 13.3989, Lat: 52.5246},
		{Title: "Flat by the park", Price: "$365k", Beds: 2, Baths: 1, Area: "64 m²", Photo: photo("photos/blue-sofa-with-pillows-in-a-designer-living-room-interior.jpg"), Lng: 13.4336, Lat: 52.5287},
		{Title: "Studio in Kreuzberg", Price: "$240k", Beds: 1, Baths: 1, Area: "38 m²", Photo: photo("photos/elegant-home-office-with-golden-accessories.jpg"), Lng: 13.4107, Lat: 52.4993},
		{Title: "Family home", Price: "$690k", Beds: 4, Baths: 2, Area: "142 m²", Photo: photo("photos/home-office-desk-with-macbook-iphone-calendar-watch-and-organizer.jpg"), Lng: 13.3458, Lat: 52.5083},
	}
}

var serverRegions = []maps.Region{
	{Name: "Frankfurt", Code: "eu-central", State: maps.Operational, Latency: "24 ms", Uptime: "99.99%", Lng: 8.6821, Lat: 50.1109},
	{Name: "London", Code: "eu-west", State: maps.Operational, Latency: "31 ms", Uptime: "99.98%", Lng: -0.1276, Lat: 51.5072},
	{Name: "Virginia", Code: "us-east", State: maps.Degraded, Latency: "212 ms", Uptime: "99.71%", Lng: -77.4874, Lat: 39.0438},
	{Name: "São Paulo", Code: "sa-east", State: maps.Operational, Latency: "48 ms", Uptime: "99.95%", Lng: -46.6333, Lat: -23.5505},
	{Name: "Singapore", Code: "ap-southeast", State: maps.Outage, Latency: "-", Uptime: "98.20%", Lng: 103.8198, Lat: 1.3521},
	{Name: "Sydney", Code: "ap-south", State: maps.Operational, Latency: "36 ms", Uptime: "99.97%", Lng: 151.2093, Lat: -33.8688},
}

var userCountries = []maps.Country{
	{Name: "United States", Flag: "us", Users: "12,480", Share: 32}, {Name: "Brazil", Flag: "br", Users: "7,120", Share: 18},
	{Name: "Germany", Flag: "de", Users: "5,930", Share: 15}, {Name: "India", Flag: "in", Users: "4,360", Share: 11},
	{Name: "Australia", Flag: "au", Users: "3,540", Share: 9}, {Name: "Poland", Flag: "pl", Users: "2,210", Share: 6},
}
