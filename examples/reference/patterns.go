package reference

// The explanation on top of each page: what the pattern is, when to use it, the rules and the block that
// implements it, and a link to the concept. It does not repeat the rules (concept §11). Rule numbers
// stand only here, in .ref-explain: the product pages themselves show none (GS-71).

type block struct {
	Package string // e.g. "blocks/detail"
	Parts   string // e.g. "Section, EditButton, Dialog"
}

type pattern struct {
	ID      string // B-01 …
	Name    string
	Path    string
	Summary string // one sentence for the overview
	What    string
	When    string
	Try     string // what to try on the page
	Rules   []string
	Blocks  []block
	Link    string // the section of the concept
	Note    string // an extra remark, e.g. on examples/crud
}

var patterns = []pattern{
	{ID: "B-01", Name: "Rahmen", Path: pathFrame,
		Summary: "Kopf, gruppierte Seitenleiste, grauer Inhalt und Panel rechts; schmal als Menü und Sheet.",
		What:    "Der Rahmen um jede Seite: Kopf mit Brotkrumen und den Knöpfen für das Panel, links die Seitenleiste in Gruppen, in der Mitte der Inhalt auf grauem Grund, rechts das Panel mit Reitern.",
		When:    "Auf jeder Büro-Seite. Nur Anmeldung und mobile Portale haben keinen Rahmen.",
		Try:     "Öffnen Sie oben rechts „Notizen“ oder „KI-Assistent“: Das Panel gleitet herein und bleibt beim Seitenwechsel offen. Machen Sie das Fenster schmal: Die Leiste wird zum Menü, das Panel zum Sheet.",
		Rules:   []string{"GS-10", "GS-11", "GS-12", "GS-13", "GS-14", "GS-15", "GS-20", "GS-21", "GS-22"},
		Blocks:  []block{{"blocks/shell", "Page"}, {"sidebar", "Sidebar, Group, Item"}, {"blocks/sidepanel", "Aside, Buttons, Script, Read"}, {"themes", "Graphite"}},
		Link:    pathConcept + "#3-application-layout-rahmen"},
	{ID: "B-02", Name: "Start", Path: pathStart,
		Summary: "Was heute zu tun ist: Kacheln mit Zahlen, die auf gefilterte Listen führen, darunter „Brennt“.",
		What:    "Die erste Seite nach der Anmeldung. Sie zeigt, was heute zu tun ist, nicht wozu die App da ist.",
		When:    "Als Startseite jeder App. Jede Zahl ist ein Link auf die Liste mit dem passenden Filter.",
		Try:     "Klicken Sie auf eine Kachel: Die Liste öffnet sich gefiltert, der Filter steht in der Adresse.",
		Rules:   []string{"GS-30"},
		Blocks:  []block{{"blocks/metric", "Card"}, {"table", "Table"}},
		Link:    pathConcept + "#5-seitentypen"},
	{ID: "B-03", Name: "Liste", Path: pathCustomers,
		Summary: "Suche, ein Filter, 25 Zeilen und Seiten; Status als Pille, Ampel als Punkt mit Wort, Zustand in der Adresse.",
		What:    "Eine Liste von Datensätzen mit Titel und Hauptaktion, Suche, einem Filter, Tabelle und Seiten.",
		When:    "Für jeden Bereich der Seitenleiste mit vielen Datensätzen. Sortiert wird nach dem, was zu tun ist (hier: überfällige Posten oben).",
		Try:     "Suchen Sie „muehl“: Es findet auch „Mühlhausen“ und „Mühlenweg“. Blättern, filtern, dann Zurück im Browser. Am Handy scrollen die Kunden in ihrer Insel, die Angebote werden zu Karten; ihre Frist zeigt die Restzeit.",
		Rules:   []string{"GS-31", "GS-50", "GS-51", "GS-52", "GS-53", "GS-54"},
		Blocks:  []block{{"blocks/listpage", "New, Paginate, Card"}, {"table", "Table, Props.Mobile"}, {"deadline", "Deadline"}, {"status", "Status"}, {"badge", "Badge"}},
		Link:    pathConcept + "#71-anzeigen",
		Note:    "Die Suche läuft hier im Speicher. In den Produkten rechnet SQL Suche, Filter, Seite und Summen (GS-51)."},
	{ID: "B-04", Name: "Detailseite", Path: customerURL("K-1001"),
		Summary: "Inseln zum Lesen, je Insel höchstens ein „Bearbeiten“ mit Dialog, Fehler im Dialog, Unterliste mit „Alle N“, Datei hinzufügen.",
		What:    "Ein Datensatz als weiße Inseln auf grauem Grund. Jede Insel zeigt nur an; „Bearbeiten“ öffnet einen Dialog nur mit ihren Feldern.",
		When:    "Für jeden Datensatz. Zuerst die Inseln, an denen man arbeitet, dann was aus anderen Apps kommt (hier: offene Posten).",
		Try:     "Leeren Sie im Dialog „Stammdaten“ den Namen und speichern Sie: Der Dialog öffnet sich wieder, mit Meldung und Ihren Werten, das Feld ist markiert.",
		Rules:   []string{"GS-32", "GS-40", "GS-41", "GS-42", "GS-43", "GS-44", "GS-45", "GS-46"},
		Blocks:  []block{{"blocks/detail", "Section, Fields, Field, EditButton, Dialog, Script"}, {"upload", "Button, Drop, Script"}},
		Link:    pathConcept + "#6-inseln-und-dialoge",
		Note:    "Die Produkte bearbeiten nie direkt in der Seite oder in einer Tabellenzeile. Das Beispiel examples/crud dieser Bibliothek zeigt das Bearbeiten in der Zeile als allgemeines Muster; die Produkte nutzen es nicht."},
	{ID: "B-05", Name: "Liste bearbeiten", Path: offerURL("A-2026-001"),
		Summary: "Lange Liste nur lesen; „Bearbeiten“ öffnet den ganzseitigen Dialog, jede Zeile ein Formular mit aufklappbarem Langtext, neue Zeile unten, Summe.",
		What:    "Eine lange Liste auf einer Detailseite (die Positionen eines Angebots) wird nur gelesen. „Bearbeiten“ öffnet sie in einem Dialog über fast die ganze Seite; jede Zeile ist dort ein eigenes Formular.",
		When:    "Für Listen, die zum Datensatz gehören und viele Zeilen haben. Nach der Abgabe zeigt die Insel den Grund statt des Knopfs (Angebot A-2026-002).",
		Try:     "Ändern Sie eine Menge und drücken Sie Enter: Die Zeile wird gespeichert, Betrag und Summe stimmen wieder. Löschen fragt nach. Unten steht die nächste Positionsnummer schon da. Unter jeder Zeile klappt der Langtext auf; ist er zu lang, öffnet sich die Zeile nach dem Speichern mit dem Fehler.",
		Rules:   []string{"GS-55", "GS-56", "GS-57", "GS-58", "GS-59"},
		Blocks:  []block{{"blocks/rowedit", "Dialog, List, Head, Row, Actions, Save, Delete, Add, Details, Total, Locked, Script"}, {"blocks/detail", "Section, EditButton"}},
		Link:    pathConcept + "#72-bearbeiten",
		Note:    "Die Produkte bearbeiten nicht in der Tabelle der Seite. Das Beispiel examples/crud dieser Bibliothek zeigt das Bearbeiten in der Zeile als allgemeines Muster; die Produkte nutzen es nicht."},
	{ID: "B-06", Name: "Anlegen", Path: pathNew,
		Summary: "Eine eigene Seite mit wenigen Feldern; danach die neue Detailseite mit Meldung.",
		What:    "Ein neuer Datensatz entsteht auf einer eigenen Seite, nicht in einem Dialog über der Liste.",
		When:    "Für jedes Anlegen. Nur die Felder, die dafür nötig sind; alles Weitere auf der Detailseite.",
		Try:     "Legen Sie einen Kunden ohne Ort an: Die Seite zeigt den Fehler am Feld. Mit Ort führt Speichern auf die neue Detailseite, mit Meldung und nächstem Schritt.",
		Rules:   []string{"GS-33", "GS-60"},
		Blocks:  []block{{"blocks/create", "Form, Footer"}, {"blocks/detail", "Section"}, {"form", "Field, Input"}},
		Link:    pathConcept + "#5-seitentypen"},
	{ID: "B-07", Name: "Zustände", Path: pathStates,
		Summary: "Meldung mit nächstem Schritt, leerer Zustand, leere Suche, Fehlerseiten, „nur lesen“ und „Demo-Instanz“, Platzhalter beim Laden, Aktion mit Grund und mit getippter Bestätigung.",
		What:    "Was eine Seite zeigt, wenn etwas geschehen ist, nichts da ist, etwas fehlt oder noch lädt.",
		When:    "Überall: nach jeder Aktion eine Meldung, bei leeren Listen eine Erklärung, Fehler im Rahmen.",
		Try:     "Schalten Sie „Nur lesen“ ein und gehen Sie auf eine Detailseite: Alle Knöpfe zum Bearbeiten sind weg. Nehmen Sie die Abgabe ohne Grund zurück: Der Dialog öffnet sich wieder mit dem Fehler am Feld.",
		Rules:   []string{"GS-44", "GS-60", "GS-61", "GS-62", "GS-63", "GS-64"},
		Blocks:  []block{{"blocks/shell", "Flash, Notice"}, {"blocks/detail", "Dialog"}, {"empty", "Empty"}, {"placeholder", "Placeholder"}},
		Link:    pathConcept + "#8-rückmeldung-und-zustände"},
	{ID: "B-08", Name: "Marke und Form", Path: pathBrand,
		Summary: "Theme-Tokens, Signalfarben mit Kontrast, Markierungen, Schrift, Knopfarten, Pillen, Ampel, Zahlen- und Datumsformat.",
		What:    "Die Bausteine der Form: neutraler Grund, Farbe nur mit Bedeutung, Geist als Schrift, Zahlen und Datum deutsch.",
		When:    "Immer. Was hier steht, gilt in jedem Theme gleich.",
		Try:     "Vergleichen Sie die Knopfarten: gefüllt ist nur die Hauptaktion.",
		Rules:   []string{"GS-07", "GS-15", "GS-70", "GS-71", "GS-72", "GS-73", "GS-74", "GS-75", "GS-76"},
		Blocks:  []block{{"themes", "Graphite, Signals, SignalStyles, Contrast"}, {"mark", "Text, Badge"}, {"button", "Button"}, {"badge", "Badge"}, {"status", "Status"}},
		Link:    pathConcept + "#9-sprache-und-form"},
	{ID: "B-09", Name: "KI-Assistent", Path: pathAssistant,
		Summary: "Der Assistent als Reiter im Panel rechts: liest, füllt Formulare vor, speichert nie.",
		What:    "Der KI-Assistent ist ein Reiter im Panel rechts. Er liest mit den Rechten der Person, zeigt, was er gelesen hat, und füllt Formulare vor; speichern tut der Mensch.",
		When:    "In Büro-Oberflächen, wenn die Firma ihn eingeschaltet hat. Die App ist ohne ihn vollständig nutzbar.",
		Try:     "Fragen Sie „Trage für die Wohnungsgesellschaft Nordstadt die Anschrift Hafenstraße 12 ein“: Der Dialog öffnet sich ausgefüllt. Fragen Sie nach der IBAN: abgelehnt.",
		Rules:   []string{"GS-13", "KA-01", "KA-02", "KA-03", "KA-15", "KA-20", "KA-28"},
		Blocks:  []block{{"blocks/sidepanel", "Aside, Buttons"}, {"blocks/assistant", "Conversation, Catalog, Stream, ActionScript, Transcript, German"}},
		Link:    pathConceptAssistant + "#6-wo-er-erscheint",
		Note:    "Die Regeln KA-xx stehen im Konzept Assistent. Ein Produkt zum Selbstbetrieb bietet den Assistenten nur mit einem Anbieter an, den der Betreiber gewählt hat; vertrauliche Daten sind nie Werkzeug."},
	{ID: "B-10", Name: "Vorschlag prüfen", Path: requestURL("R-2026-001"),
		Summary: "KI als Fachteil auf einer Arbeitsseite in drei Spalten: Vorschlag mit Quellen und Gründen, die Person entscheidet.",
		What:    "Eine KI bereitet die Arbeit vor, hier die Antwort auf eine Anfrage; der Mensch prüft, ändert und entscheidet. Die Arbeitsseite zeigt in drei Spalten, was geschah, den Vorschlag und den Zusammenhang.",
		When:    "Wo ein Agent einen Schritt vorschlägt, der nach außen wirkt. Er sendet nie selbst; Verwerfen fragt nach dem Grund.",
		Try:     "Ändern Sie die Antwort und senden Sie sie: Der Verlauf zeigt sie als geändert. Oder verwerfen Sie ohne Grund: Der Dialog öffnet sich mit dem Fehler am Feld wieder.",
		Rules:   []string{"GS-34", "GS-44", "GS-52", "GS-76", "KA-01", "KA-02", "KA-03", "KA-04"},
		Blocks:  []block{{"blocks/detail", "Section, Dialog"}, {"mark", "Text"}, {"deadline", "Deadline"}, {"blocks/assistant", "Level, German"}},
		Link:    pathConcept + "#5-seitentypen",
		Note:    "Bearbeiten in der Seite gilt hier nur für die Arbeit selbst, die Antwort. Die Inseln mit dem Zusammenhang lesen nur."},
	{ID: "B-11", Name: "Öffentlich und Externe", Path: pathExternal,
		Summary: "Derselbe Rahmen mit kurzer Leiste für Externe; ohne Anmeldung „Anmelden“ als Hauptaktion im Kopf.",
		What:    "Öffentliche Seiten und ein Bereich für Externe (Kunden, Bieter) nutzen denselben Rahmen wie das Büro, mit einer kurzen Leiste von ein bis drei Einträgen und ohne Panel.",
		When:    "Für alles, was Externe am Schreibtisch sehen. Portale am Handy haben eigene Regeln.",
		Try:     "Ohne Anmeldung stehen rechts im Kopf „Registrieren“ und „Anmelden“. Melden Sie sich an: Die Leiste zeigt „Meine Angebote“, der Kopf Ihren Namen und „Abmelden“.",
		Rules:   []string{"GS-16", "GS-10", "GS-11"},
		Blocks:  []block{{"blocks/shell", "Page"}, {"sidebar", "Sidebar, Group, Item"}, {"button", "Button"}},
		Link:    pathConcept + "#3-application-layout-rahmen"},
}

var overviewText = pattern{ID: "Übersicht", Name: "Referenz der Gestaltung",
	What: "Diese Referenz zeigt jedes Muster des Konzepts Gestaltung mit den echten Bausteinen von tabler-templ. Ändert sich ein Baustein, ändert sich das Beispiel mit.",
	When: "Zum Nachschlagen beim Bauen und zum Vergleichen der Produkte. Die Beispieldaten sind erfunden; Speichern wirkt nur bis zum Neustart.",
	Link: pathConcept + "#11-beispiele-die-referenz-app-in-tabler-templ"}

func patternOf(key string) (pattern, bool) {
	for _, p := range patterns {
		if p.ID == key {
			return p, true
		}
	}
	return pattern{}, false
}
