# Gestaltung: Rahmen, Seiten, Inseln und Tabellen

| | |
|---|---|
| **Status** | v0.14 · 16.11.2026 · Umsetzungsstand tabler-templ v0.22.0 (GS-09, GS-56 „hier“), keine neue Regel. [Historie](../historie/docs/konzepte/gestaltung.md) |
| **Anlass** | Wunsch des Benutzers vom 15.11.2026: aus Seitenleiste, Application Layout, „Inseln bearbeiten“ und dem Bearbeiten von Tabellen ein Design- und UX-Konzept mit Beispielen machen, das sich nach und nach überarbeiten, standardisieren und zwischen den Produkten austauschen lässt; die Beispiele als lebende Referenz-App in tabler-templ. |
| **Gilt zusammen mit** | [Navigation](navigation.md) (NV-01 bis NV-09, G-62), [Auftrag ADR-0013](../../apps/auftrag/docs/adr/0013-inseln-mit-dialogen.md) (Inseln mit Dialogen, E-45, E-70), [Verbund](verbund.md) G-41, G-44, G-62, G-63; [`libs/rahmen`](../../libs/rahmen/README.md) (Gerüst, Theme, Panel, Fehlerseiten), [`libs/liste`](../../libs/liste/README.md); Einsatz [ux-disponent](../../apps/einsatz/docs/ux-disponent.md) (UX-01 bis UX-14); Brief [gestaltung-referenz-app](../briefs/gestaltung-referenz-app.md) |
| **Präfix** | GS-xx (Regeln der Gestaltung), GSO-xx (offene Punkte) |
| **Entscheidung** | Benutzer 15.11.2026: „alle Empfehlungen annehmen“ entscheidet GSO-01 bis GSO-07, „gerne eintragen“ GSO-08 bis GSO-14 (Fragen aus tabler-templ #60), „Ja, ich folge deinen Empfehlungen zu Auftragsbutler 474“ (in der Session tabler-templ) GSO-15 bis GSO-44, „dann (a)“ (16.11.2026) GSO-45, „Ja, ich folge den Empfehlungen“ (16.11.2026) GSO-46 und GSO-47, „Ja, ich folge der Empfehlung“ (16.11.2026) GSO-48, „Ja, weiter so“ (16.11.2026) GSO-49, „Ja, am besten so machen mit dem Handy als Karten.“ (16.11.2026, in der Session von Vergabe) GSO-50, jeweils wie empfohlen (§13 bis §16); Verbund G-66. |

## 0. Kurzfassung

Dieses Dokument fasst zusammen, was in den Apps des Verbunds schon gilt, und schreibt es so, dass auch andere
Produkte (studio-apps, Vergabe, myoperator) es übernehmen können. Neue Regeln führt es nicht ein. Wo eine Regel schon
eine Nummer hat (NV-, E-, UX-, RA-), steht sie dabei. Die Quelle bleibt das dort genannte Dokument.

1. **Rahmen:** Kopf oben (56 px), Seitenleiste links (240 px), Inhalt auf grauem Grund, Panel rechts für Notizen und
   Assistent. Auf dem Handy wird die Leiste zum Menü und das Panel zum Sheet (§3).
2. **Navigation:** „Start“ oben, dann Gruppen in fester Reihenfolge: „Für mich“, Tagesgeschäft, „Erledigen“,
   „Verwaltung“ (§4, NV-01 bis NV-09).
3. **Fünf Seitentypen:** Start, Liste, Detailseite, Anlegen, Arbeitsseite (§5).
4. **Inseln:** Eine Detailseite besteht aus weißen Inseln, die nur anzeigen. Jede hat höchstens einen Knopf
   „Bearbeiten“, der einen Dialog nur mit den Feldern dieser Insel öffnet. Es gibt keine Drawer und kein Bearbeiten
   direkt in der Seite (§6, ADR-0013).
5. **Tabellen:** Eine Liste hat oben eine Suche und einen Filter, zeigt 25 Zeilen und unten die Seiten. Der Zustand
   steht in der URL. Eine lange Liste auf einer Detailseite wird nur gelesen. „Bearbeiten“ öffnet sie in einem
   ganzseitigen Dialog, in dem jede Zeile ein eigenes Formular mit Speichern und Löschen ist und unten eine neue
   Zeile wartet (§7).
6. **Rückmeldung:** Nach jeder Aktion kommt eine Meldung mit dem nächsten Schritt. Ein Dialog mit Fehler öffnet
   sich wieder, die eingegebenen Werte bleiben stehen. Ein leerer Zustand erklärt sich und nennt den nächsten
   Schritt (§8).
7. **Sprache und Form:** Die Oberfläche nutzt die Begriffe der Nutzer und zeigt keine Regelnummern. Datum als
   TT.MM.JJJJ. Farbe nur, wenn sie etwas bedeutet (§9).
8. **Entschieden** (Benutzer 15.11.2026, alle Empfehlungen GSO-01 bis GSO-07, G-66): Konzept hier bis v1.0, dann nach tabler-templ; Referenz-App deutsch; „Liste bearbeiten“ wird ein Block in tabler-templ; prüfbare Regeln bekommen Tests.
9. **Beispiele:** eine lebende Referenz-App in tabler-templ, gebaut aus denselben Bausteinen wie die Produkte. Zu
   jedem Muster erklärt ein kurzer Text, welche GS-Regel es zeigt (§11, Brief).

## 1. Zweck und Geltung

| Nr. | Regel |
|---|---|
| GS-01 | **Gegenstand** sind die Büro-Oberflächen (am Schreibtisch, Maus und Tastatur, auch schmal am Handy), auch öffentliche Seiten und Bereiche für Externe am Schreibtisch (GS-16, GSO-39). Die Portale am Handy (`/m/`, offline, ein Daumen) haben eigene Regeln und stehen nur kurz in §10 (GSO-22). |
| GS-02 | **Für alle Produkte:** Die Regeln sind produktneutral geschrieben. Was nur hier gilt (Paketnamen, Apps, Nummern anderer Dokumente), steht in der Spalte oder dem Absatz „hier“. Ein anderes Produkt übernimmt die GS-Regeln mit derselben Nummer und ersetzt nur „hier“. |
| GS-03 | **Quelle der Wahrheit:** Die Regel steht in diesem Dokument. Die Referenz-App (§11) zeigt sie, die Bausteine in tabler-templ und `libs/rahmen` setzen sie um. Weicht eins davon ab, gilt die Regel, und der Rest wird nachgezogen (oder die Regel wird per PR geändert). |

## 2. Grundsätze

| Nr. | Grundsatz | Warum |
|---|---|---|
| GS-05 | **Lesen vor Bearbeiten.** Eine Seite zeigt zuerst an. Bearbeiten ist ein eigener, sichtbarer Schritt (Dialog oder eigene Seite). | Versehentliche Änderungen kommen seltener vor, und man sieht eindeutig, wo sich was ändern lässt (ADR-0013, abgelehnte Option „inline“). |
| GS-06 | **Eine Aktion je Ort.** Eine Insel hat höchstens einen Knopf „Bearbeiten“, eine Seite oben höchstens eine Hauptaktion („Neuer Kunde“). Weitere Aktionen stehen im Dialog oder im Menü „…“. | Man muss nicht suchen, und die Seite bleibt ruhig. |
| GS-07 | **Farbe nur mit Bedeutung.** Der Grund ist neutral (Graphite: schwarz-weiß). Farbig sind nur die Signalfarben: Grün (gut), Gelb (Achtung), Rot (Problem), Blau (Info, heute); sie gelten in jedem Theme gleich. Dazu kommen benannte Markierungen mit eigenem Token je Theme: „intern, nie gesendet“ und „von außen“ (GS-76). Frei wählbare Farben (z. B. Labels) nur aus Signalfarben und Grau. **Kontrast** in jedem Theme nach WCAG 2.1 AA: Text ≥ 4,5:1, Punkte, Linien und Flächen ≥ 3:1. Signalfarben sind nie Textfarbe, Text nimmt ihre dunklen Töne (Grün `#1E7537`, Gelb `#8A5A00`, Rot `#B42318`, Blau `#1D4ED8`); Gelb für Punkte `#B57F00`. Markierungen: „intern“ Punkt `#7C3AED`, Text `#6D28D9`; „von außen“ Punkt `#B45309`, Text `#92400E`. Punkte im Dunkelmodus: Blau `#3B82F6`, „intern“ `#8B5CF6`, „von außen“ `#D97706` (GSO-25, GSO-31; Werte entschieden von Claude im Auftrag des Benutzers, 15.11.2026, nach Vorschlag und Test von tabler-templ). | Rot fällt auf, weil sonst nichts bunt ist (`libs/rahmen`, Theme); BITV 2.0 für öffentliche Auftraggeber. Ein Test in tabler-templ `themes` und `libs/rahmen` prüft alle Paare, auch Text in einer Pille auf dem grauen Grund (getönte Fläche). Hier nutzt `libs/rahmen/static/theme.css` noch Grün `#1F7A3A` auf eigener Fläche (4,9:1); es folgt beim Umstellen auf die Werte von tabler-templ. |
| GS-08 | **Die Rolle blendet aus, der Use Case prüft.** Die Oberfläche zeigt nur, was die Rolle darf. Geprüft wird trotzdem auf dem Server. Gibt es für einen Abschnitt keinen Use Case zum Ändern, gibt es auch keinen Knopf. | Auf ein „Bearbeiten“ soll nie ein 403 folgen (G-41, NV-04, ADR-0013). |
| GS-09 | **Server rendert, Formulare schicken ab.** Seiten entstehen auf dem Server (templ), Formulare arbeiten mit POST und Rücksprung (PRG), wenig eigenes JavaScript. Nichts kommt von einem CDN. Eigenes JavaScript nur mit Grund, eingebettet und getestet; ein solcher Grund ist Kryptografie beim Nutzer (WebCrypto, der Schlüssel verlässt den Browser nie; GSO-38). Statische Dateien haben eine Version im Pfad und kommen dauerhaft aus dem Cache (immutable); Icons nur als Untersatz der verwendeten; Schriften nur des eigenen Themes (GSO-48). Ein Klick lädt die ganze Seite neu, kein Austausch nur des Inhalts (hx-boost) (GSO-49). | Ein Binary, schnell, offline-fest, prüfbar (`libs/rahmen` `TestNichtsVonAussen`). tabler-templ `static.New`, `Files.URL`, `static.Head` und `cmd/static` (Untersatz, `-check` im Test); hier `libs/rahmen` (RA-15, RA-16); ab v0.22 `static.WithFilledIcons` (Webfont der gefüllten Icons aus dem Bündel) und `themes.Styles` (CSS aller Themes in `kit.css`), `kit.css` steht nach Tabler (`rahmen.KitCSS`, `TestGS07_KitCSSNachTabler`); davor überschrieb Tabler Signalfarben und Fokusring. |

## 3. Application Layout (Rahmen)

| Nr. | Teil | Regel | hier |
|---|---|---|---|
| GS-10 | **Kopf** | 56 px hoch. Links der Knopf für die Leiste und bei Detailseiten die Brotkrumen (die Marke steht oben in der Leiste, GS-11; ohne Leiste links im Kopf), rechts Notizen, Assistent und Rechner, soweit das Produkt sie hat, und ganz rechts der App-Umschalter (mit der Zahl „dringend“) nur im Verbund (GSO-23). Keine eigene Navigation im Kopf. App-Umschalter und Rechner gehören dem Produkt, nicht tabler-templ (GSO-11). | `rahmen.Seite`, `rahmen.Umschalter` (G-63), Rechner RA-12; die Knöpfe zeigen unter sm nur das Symbol, Name als `aria-label` und `title` (`libs/rahmen`) |
| GS-11 | **Seitenleiste** | 240 px links. Oben die Marke: im Verbund Logo und Name der Firma, darunter der Name der App, ohne Logo das Zeichen der App, ohne Anmeldung „Auftragsbutler“ (G-67); andere Produkte ihre eigene. Darunter die Einträge, gruppiert nach §4, je Eintrag ein Symbol und ein Wort. Beim Wechsel über die Breite lg kein Übergang (der Inhalt springt, statt hinter die Leiste zu gleiten). | `sidebar.Group` aus tabler-templ, `rahmen.Symbol` |
| GS-12 | **Inhalt** | Grauer Grund, Breite höchstens 1240 px, linksbündig neben der Leiste. Oben der Seitentitel (bei Detailseiten mit Brotkrumen „Bereich › Liste › Datensatz“) und rechts davon die Hauptaktion (GS-06). Darunter die Meldungszeile (GS-60), dann die Inseln. | `theme.css` (`--grund`, `--linie`) |
| GS-13 | **Panel rechts** | Arbeitsbereich neben der Seite mit Reitern (Notizen, Assistent). Er bleibt beim Seitenwechsel offen, der Zustand steht im Cookie. Kein Formular zum Bearbeiten im Panel. Ein Knopf im Kopf blendet es ein und aus; es gleitet herein, steht ab der Breite xl angedockt neben der Seite und darunter als Sheet über ihr (GS-14). Verlangt der Browser weniger Bewegung (`prefers-reduced-motion`), erscheint es ohne Animation (GSO-10). Der Reiter „Notizen“ zeigt die Notizen des aktuellen Datensatzes, markiert als „intern“ (GSO-33). | tabler-templ `blocks/sidepanel`; hier `rahmen.Panel` darauf (RA-11, Auftrag E-70) |
| GS-14 | **Handy und schmal** | Unter lg wird die Leiste zum Menü (gleiche Reihenfolge, NV-07), das Panel zum Sheet über der Seite (öffnet sich nach einem Seitenwechsel nicht von selbst), ein Dialog zum Vollbild. Die Seite scrollt nie seitlich; breite Tabellen scrollen in ihrer Insel. Felder haben 16 px Schrift, damit das Handy nicht zoomt. Layout-Tests prüfen 1400, 1366, 1280 und 390 px Breite (GSO-21). | `tools/e2e-staging/layout.mjs` prüft das; `detail.Dialog` unter lg Vollbild ab tabler-templ v0.21 |
| GS-15 | **Theme** | Ein Theme je Produkt (eigene Marke), im Verbund eines für alle Apps an einer Stelle wählbar. Ein Produkt zum Selbstbetrieb lässt den Betreiber das Theme je Instanz wählen, Standard Graphite. Farben, Ecken und Linien sind Tokens; gleich bleiben in jedem Theme Aufbau, Muster, Signalfarben und Kontrast (GS-07; GSO-15, GSO-24). | hier `VERBUND_THEME`, Standard Graphite |
| GS-16 | **Ohne Anmeldung** | Öffentliche Seiten und Bereiche für Externe nutzen denselben Rahmen mit kurzer Leiste. Ohne Anmeldung steht rechts im Kopf „Anmelden“ als Hauptaktion, „Registrieren“ schlicht daneben (GSO-39, Beispiel B-11). | |
| GS-17 | **Benutzermenü** | Unten in der Leiste, in jeder App gleich: Kürzel, Name und Rolle; im Menü „Meine Daten“ (im Verbund die Seite in Konto) und „Abmelden“. Eigene Einträge einer App stehen davor (z. B. Benutzerwechsel der Demo). | `rahmen.Benutzer` (RA-18) |

## 4. Navigation

Die Regeln stehen in [Navigation](navigation.md) §2 und gelten hier unverändert. Kurz:

| Nr. | Regel |
|---|---|
| GS-20 | „Start“ steht oben ohne Gruppentitel (NV-01). Danach kommen Gruppen in fester Reihenfolge: „Für mich“/„Unterwegs“, höchstens zwei Gruppen Tagesgeschäft, „Erledigen“, „Verwaltung“ (NV-02). In „Verwaltung“ erst die Stammdaten, dann Einstellungen, Importieren, Daten exportieren (NV-03). |
| GS-21 | Höchstens etwa 5 bis 7 Einträge je Gruppe, keine Unterbereiche mit „Zurück“ (NV-06). Gleiches heißt in allen Apps gleich (NV-05). Jeder Eintrag hat einen stabilen Anker für Touren und Tests (NV-08). |
| GS-22 | Die Leiste ist nur für die Hauptbereiche da. Von einem Datensatz zu einem anderen führt ein Link in der Insel oder in der Tabellenzeile, nie ein Eintrag in der Leiste. |
| GS-23 | Gruppen des Tagesgeschäfts dürfen in der Reihenfolge eines echten Ablaufs stehen (Buchhalter: Erfassen, Abschließen); nur dann tragen die Einträge Nummern (GSO-17). |

## 5. Seitentypen

| Nr. | Typ | Aufbau | Vorbild |
|---|---|---|---|
| GS-30 | **Start** | Was heute zu tun ist, nicht wozu die App da ist. Kacheln mit Zahlen, und jede Zahl führt auf die gefilterte Liste. Darunter Arbeitslisten („Brennt“, „Diese Woche“). | Einsatz UX-01, Konto „Mein Tag“ |
| GS-31 | **Liste** | Titel und Hauptaktion („Neuer …“), Suche, ein Filter, Tabelle, Seiten (§7.1). | Auftrag „Kunden“, E-49 |
| GS-32 | **Detailseite** | Kopf mit Name, Nummer (wenn das Produkt eine fachliche Nummer hat, nie eine interne ID; GSO-43) und Status, darunter die Inseln in fester Reihenfolge: zuerst, woran man arbeitet, dann was aus anderen Apps kommt (§6). Eine Karte „Als Nächstes“ zeigt, was zu tun ist, wenn es einen Ablauf gibt. **Arbeitsstand und Fassung:** Wird ein Stand veröffentlicht, zeigt ein Hinweis oben „Es gibt Änderungen gegenüber Fassung N“ mit der Hauptaktion „Neue Fassung veröffentlichen“ (mit Grund, GS-44), die Fassungen stehen als Insel (GSO-40). **Einstellungen** sind eine Detailseite aus Inseln, kein großes Formular (GSO-36). | Auftrag „Auftrag“, Einsatz UX-03, UX-07 |
| GS-33 | **Anlegen** | Eine eigene Seite (`/kunden/neu`), nicht als Dialog über der Liste. Sie hat nur die Felder, die zum Anlegen nötig sind, und führt nach dem Speichern auf die neue Detailseite. | ADR-0013 Punkt 4 |
| GS-34 | **Arbeitsseite** | Ausnahme für Arbeit am Stück (Plantafel, Freigabe, Erfassen, Importieren, ein Vorgang im Support): eigenes Layout erlaubt, auch drei Spalten. Bearbeiten in der Seite nur für die Arbeit selbst (Antwort, Vorschlag, Notiz); Inseln mit Kontext folgen GS-41 bis GS-45. Rahmen und Regeln §8 und §9 gelten trotzdem (GSO-21, GSO-28). | ADR-0013; Support E-72 |

## 6. Inseln und Dialoge

| Nr. | Regel | Quelle |
|---|---|---|
| GS-40 | **Insel:** eine weiße Karte auf grauem Grund mit Überschrift. Der Inhalt ist nur zum Lesen da, meist als Paare aus Bezeichnung und Wert in zwei Spalten (am Handy eine). Fehlende Werte zeigen „–“. | ADR-0013 Punkt 1, 5 |
| GS-41 | **Bearbeiten:** rechts oben in der Insel steht genau ein Knopf „Bearbeiten“. Gibt es keinen passenden Use Case, entfällt er. | ADR-0013 Punkt 1, GS-08 |
| GS-42 | **Dialog:** ein zentrierter Dialog mit genau den Feldern dieser Insel. Der Titel heißt wie die Insel, darunter steht klein, worum es geht („A-2026-001 · Wohnungsgesellschaft Nordstadt“). Unten links „Abbrechen“, rechts „Speichern“. Gespeichert wird nur diese Insel. | ADR-0013 Punkt 2 |
| GS-43 | **Fehler im Dialog:** Scheitert das Speichern, öffnet die Seite den Dialog wieder, mit der Meldung oben und den eingegebenen Werten; das fehlerhafte Feld ist markiert. | ADR-0013 Punkt 2 |
| GS-44 | **Kleine Ergänzungen** an einem bestehenden Datensatz (Nachtrag, Zuweisung, PIN, Datei hochladen mit Datei, Art und Bezug) sind Dialoge seiner Detailseite, mit einem Knopf in der Insel, zu der sie gehören (GSO-41). **Aktionen mit Grund** (rechtlich wirksam, Entscheidung, Zurücksetzen) sind ein kleiner Dialog mit Pflichtfeld und rotem Knopf; eine unumkehrbare, folgenreiche Aktion verlangt im Dialog eine getippte Bestätigung (z. B. die Adresse; GSO-26). **Datei wählen:** ein Knopf „Datei wählen“ mit dem Dateinamen daneben; am Handy (unter lg) daneben „Foto aufnehmen“ (Kamera), nur wo ein Foto als Beleg taugt (GSO-46). **Ablage:** Dateien lassen sich auf die Insel oder Liste ziehen, zu der sie gehören; braucht der Upload mehr als die Datei (Art, Bezug), öffnet die Ablage den Dialog mit der Datei, sonst lädt sie gleich hoch. Beim Ziehen zeigt die Fläche „Zum Hochladen loslassen“; der Knopf bleibt der sichtbare Weg, die Ablage ist eine Abkürzung (GSO-47). | ADR-0013 Punkt 4; tabler-templ `upload` ab v0.21 |
| GS-45 | **Keine Drawer, kein Bearbeiten in der Insel.** Felder stehen nie direkt in der Leseansicht. Das Panel rechts ist kein Formular (GS-13). | ADR-0013, E-70 |
| GS-46 | **Unterlisten in einer Insel** zeigen höchstens 10 Zeilen und darunter „Alle N“, das auf die Liste mit Filter führt. Das gilt für Listen anderer Datensätze (Aufträge eines Kunden). Eigene Zeilen des Datensatzes (Positionen eines Angebots) zeigt die Insel ganz; sie scrollen in der Insel (GS-54, GSO-09). Das gilt auch für ein unveränderliches Protokoll; „Alle N“ führt dann zur vollständigen Liste mit Suche (GSO-42). | Verbund G-44 |

## 7. Tabellen

### 7.1 Anzeigen

| Nr. | Regel | Quelle |
|---|---|---|
| GS-50 | **Liste:** oben eine Textsuche (jedes Wort muss vorkommen, ohne Groß- und Kleinschreibung, ä = ae), daneben höchstens ein einfacher Filter, unten Seiten mit 25 Zeilen. Eine Arbeitsliste darf daneben Schnellfilter haben („Nur meine“, GS-53); weitere Merkmale gehen in die Suche (`label:rechnung`). Jede Liste, die über 25 Zeilen wachsen kann, hat Suche und Seiten. Suche, Filter und Seite stehen in der URL, damit Zurück, Lesezeichen und Links funktionieren (GSO-19, GSO-30). | Auftrag E-49, `libs/liste`; tabler-templ `listpage.Contains`, `listpage.Fold` (ab v0.19) |
| GS-51 | **Rechnen in der Datenbank:** Suche, Filter, Seite, Summen und Gesamtzahl rechnet SQL in wenigen festen Abfragen je Seite, nie eine je Zeile. | G-44, ADR-0011 |
| GS-52 | **Spalten:** Was man erkennt, steht vorn (Name, Ort), Nummern dahinter. Name oder Nummer ist der Link zur Detailseite. Zahlen und Beträge stehen rechtsbündig mit gleich breiten Ziffern, Einheiten in eigener Spalte oder hinter der Zahl. Status als Pille, Ampel als Punkt mit Wort. Eine Frist ist eine Ampel, deren Wort die Restzeit ist („noch 58 Min.“, „überfällig seit 20 Std.“; GSO-35). | Einsatz UX-04, UX-06; `theme.css`; tabler-templ `deadline` ab v0.21 |
| GS-53 | **Sortierung** nach dem, was Handlung braucht (Fälligkeit, „Frist fehlt“ oben), nicht nach Anlagedatum. Eine Liste zum Abarbeiten hat Schnellfilter statt vieler Spalten. | Einsatz UX-05, Ausschreibungen AS002 |
| GS-54 | **Breite Tabellen** scrollen seitlich in ihrer Insel oder werden unter 768 px zu Karten (je Tabelle gewählt; GSO-21). Die Seite selbst scrollt nie seitlich (GS-14). Gilt auch für „Liste bearbeiten“ (GS-56, GSO-50). | Buchhalter OF-04 |

### 7.2 Bearbeiten

| Nr. | Regel | Quelle |
|---|---|---|
| GS-55 | **Lange Listen im ganzseitigen Dialog.** Auf der Detailseite wird die Liste (z. B. Positionen eines Angebots) nur gelesen. „Bearbeiten“ öffnet einen Dialog über fast die ganze Seite, unter lg als echtes Vollbild. | ADR-0013 Punkt 3; hier `modal-ganzseitig` in Auftrag |
| GS-56 | **Jede Zeile ein Formular.** Jede Liste mit mehr als einem Feld je Zeile wird so bearbeitet, auch eine kurze (Labels). Im Dialog ist jede Zeile ein eigenes kleines Formular mit festen Spalten. Rechts stehen „Speichern“ (Haken) und „Löschen“ (Papierkorb). Enter speichert die Zeile. Ein Hinweis oben sagt: „Jede Zeile wird für sich gespeichert.“ **Ausnahme:** Eine kurze Liste (bis etwa 10 Zeilen), deren Zeilen nur zusammen gültig sind (Gewichte = 100 %), ist ein Formular im normalen Dialog (GS-42; GSO-27). **Lange Texte:** Eine Zeile darf eine aufklappbare zweite Zeile im selben Formular haben, mit Textfeldern über die ganze Breite (Langtext, Anforderungen). Sie ist zu und nach einem Fehler in einem ihrer Felder offen; Speichern und Löschen bleiben in der ersten Zeile; Enter in einem Textfeld speichert nicht, in der ersten Zeile weiter schon (GSO-45). **Am Handy als Karten (je Liste gewählt):** Unter 768 px wird jede Zeile eine Karte. Jedes Feld steht mit seiner Bezeichnung untereinander. Speichern und Löschen stehen unten in der Karte, Enter speichert weiter. Details lassen sich in der Karte zuklappen. Die neue Zeile ist eine Karte am Ende. Darüber bleibt die Liste mit festen Spalten (GS-54, GSO-50). | tabler-templ `blocks/rowedit` (ab v0.18; aufklappbare Zeile `rowedit.Details` und feste Aktionsspalte ab v0.20); hier Auftrag (alle Listen-Dialoge); `rowedit.ListProps.Mobile` und `rowedit.Cell` (Karten am Handy) ab v0.22, Auftrag nutzt sie noch nicht |
| GS-57 | **Neue Zeile unten** als leeres Formular mit Vorschlag (z. B. nächste Positionsnummer) und „Hinzufügen“. | Auftrag „Positionen bearbeiten“ |
| GS-58 | **Löschen fragt nach** („Position 01.01.0010 löschen?“), über `data-bestaetigen` am Knopf (heute die Nachfrage des Browsers). Nach dem Löschen sagt die Meldung, was gelöscht wurde. Löschen mit Grund oder unumkehrbaren Folgen ist ein Dialog nach GS-44. | hier Auftrag `formulare.js`; GS-60 |
| GS-59 | **Summen und berechnete Spalten** (Betrag = Menge × Preis) zeigt der Dialog nur an und rechnet sie nach jedem Speichern neu. Schreibgeschützte Listen (z. B. nach Abgabe eines Angebots) zeigen den Grund statt der Knöpfe. | Auftrag „Positionen bearbeiten“ |

## 8. Rückmeldung und Zustände

| Nr. | Regel | Quelle |
|---|---|---|
| GS-60 | **Meldung nach jeder Aktion:** eine wegklickbare Zeile über dem Inhalt, die sagt, was geschehen ist, und den nächsten Schritt als Link anbietet („Einsatz disponiert. Zur Plantafel“). Grün für Erfolg, Rot für Fehler. | Einsatz UX-06, `shell.Flash` |
| GS-61 | **Leerer Zustand:** sagt, warum hier nichts steht und was hilft, mit einem Knopf für den nächsten Schritt. Eine leere Suche bietet „Suche zurücksetzen“ an. | Einsatz UX-09 |
| GS-62 | **Fehlerseiten** (nicht gefunden, keine Rechte, Fehler) stehen im Rahmen mit Seitenleiste und erklären in einem Satz, was los ist. Technische Einzelheiten gehören ins Log, nie in den Browser. | `libs/rahmen` RA-01 bis RA-05 |
| GS-63 | **Nur lesen:** Ist die Firma gesperrt oder gekündigt, zeigt der Rahmen einen Hinweis mit Link zum Export. Die Knöpfe zum Bearbeiten entfallen. Ebenso, an derselben Stelle über der Meldungszeile und in Info-Farbe: der Hinweis „Demo-Instanz“, wenn alle Daten erfunden sind (GSO-44). | Datenhoheit DH-12 |
| GS-64 | **Laden:** Teile, die länger brauchen (Zahlen der Seitenleiste, „dringend“ im Umschalter), kommen nach dem Seitenaufbau nach; die Seite selbst steht sofort. | [Laden im Hintergrund](laden-im-hintergrund.md) |

## 9. Sprache und Form

| Nr. | Regel | Quelle |
|---|---|---|
| GS-70 | **Begriffe der Nutzer:** die Wörter aus dem Glossar der App (`CONTEXT.md`), keine Fachsprache der Entwicklung. Gleiches heißt überall gleich. | NV-05, Einsatz UX-13 |
| GS-71 | **Keine internen Kennungen:** keine Regelnummern (E-xx, G-xx, §), keine IDs in Texten und Meldungen. Gesetze („§ 14 UStG“) dürfen bleiben. Ausnahme: die Erklärtexte der Referenz-App nennen die GS-Nummern (§11, GSO-13). | Auftrag E-52, `rahmen.OhneRegelnummern`, Architekturtest |
| GS-72 | **Knöpfe sagen, was passiert:** „Speichern“, „Angebot anlegen“, „Freigeben“, nicht „OK“ oder „Weiter“. Die Hauptaktion ist gefüllt, alle anderen sind umrandet oder schlicht. Gefährliche Aktionen sind rot und fragen nach; mit Grund oder unumkehrbar im Dialog nach GS-44 (GSO-26). | |
| GS-73 | **Zahlen und Datum** im deutschen Format: 1.234,50 €, TT.MM.JJJJ, Uhrzeit 07:30. Fristen und Zeitpunkte mit rechtlicher Wirkung als „05.12.2026, 22:00 Uhr“ (GSO-35). Zeiträume als „01.04.2026 – 31.03.2028“. Listen werden mit Komma verbunden. Ausnahme: Eingabefelder für Datum und Zeit sind die nativen Felder (`type=date`, `datetime-local`) im Format des Browsers (GSO-32). | Auftrag E-52 |
| GS-74 | **Links** im Fließtext und in Tabellen sind dezent unterstrichen (der Akzent ist schwarz wie der Text). | `libs/rahmen` Theme |
| GS-75 | **Beispieldaten** ohne erkennbare Branche (Musterdienst GmbH, Wohnungsgesellschaft Nordstadt), in Referenz-App, Mockups und Screenshots gleich. | Verbund G-60 |
| GS-76 | **Text von außen** (von Kunden, Lieferanten, Bietern) ist gekennzeichnet (Markierung „von außen“, GS-07), escaped und nie ein Link. Interne Notizen tragen die Markierung „intern, nie gesendet“ (GSO-31). | Support E-15, E-29; Baustein in tabler-templ |

## 10. Mobile Portale (nur kurz)

Die mobilen Portale (`/m/` in Auftrag, Einsatz, Zeit) nutzen dieselben Werte (Akzent, Schrift, Ecken, Signalfarben)
aus eigenem CSS ohne Tabler. Sonst sind sie anders gebaut: ein großer Knopf je Schritt, offline, kein Rahmen, keine
Seitenleiste. Eigene Regeln dafür kommen später, zusammen mit den Portalen für andere Zielgruppen am Handy (GSO-07, GSO-22).

## 11. Beispiele: die Referenz-App in tabler-templ

Wunsch des Benutzers (15.11.2026): Das Konzept soll lebendig sein. Dafür kommt in tabler-templ eine Referenz-App,
die alle Muster dieses Dokuments mit den echten Bausteinen zeigt und jeweils kurz erklärt. Sie ist kein Mockup:
Ändert sich ein Baustein, ändert sich das Beispiel mit, und Screenshot-Tests merken eine Abweichung. Den Auftrag dazu
beschreibt der Brief [gestaltung-referenz-app](../briefs/gestaltung-referenz-app.md).

**Gebaut** (tabler-templ PR #60, 15.11.2026): `examples/reference`, in der Demo unter `/reference`, mit den Blöcken
`blocks/rowedit` (Liste bearbeiten, B-05) und `blocks/sidepanel` (Panel rechts, B-01, B-09). Die App zeigt dieses
Dokument als Kopie mit Commit-Stand (`/reference/konzept`) und verlinkt in beide Richtungen. Ihre Seiten stehen in der
Seitenleiste in einer eigenen Gruppe „Referenz“ am Ende; das ist keine Gruppe der Produkte und keine Ausnahme von
GS-20 (GSO-08).

| Beispiel | zeigt | Regeln |
|---|---|---|
| B-01 Rahmen | Kopf, gruppierte Seitenleiste, Inhalt, Panel rechts mit Reitern; schmal als Menü und Sheet | GS-10 bis GS-15, GS-20 bis GS-22 |
| B-02 Start | Kacheln mit Zahlen, die auf gefilterte Listen führen, Arbeitsliste „Brennt“ | GS-30 |
| B-03 Liste | Suche, Filter, 25 Zeilen, Seiten, Status-Pillen, Ampel, Zahlen rechtsbündig, Zustand in der URL | GS-31, GS-50 bis GS-54 |
| B-04 Detailseite | Kopf, Inseln mit Bearbeiten, Dialog, Fehler im Dialog, Unterliste mit „Alle N“ | GS-32, GS-40 bis GS-46 |
| B-05 Liste bearbeiten | ganzseitiger Dialog, Zeile speichern, löschen mit Nachfrage, neue Zeile mit Vorschlag, Summe | GS-55 bis GS-59 |
| B-06 Anlegen | eigene Seite, wenige Felder, danach die Detailseite mit Meldung | GS-33, GS-60 |
| B-07 Zustände | Meldung mit nächstem Schritt, leerer Zustand, leere Suche, Fehlerseite, „nur lesen“ und „Demo-Instanz“, Platzhalter beim Laden, Aktion mit Grund und mit getippter Bestätigung | GS-44, GS-60 bis GS-64 |
| B-08 Marke und Form | Theme-Tokens, Signalfarben, Schrift, Knopfarten, Pillen, Ampel, Zahlen- und Datumsformat | GS-07, GS-15, GS-70 bis GS-75 |
| B-09 KI-Assistent | Panel mit Reiter Assistent, Beispielfragen, Assistent füllt einen Dialog, die Person speichert; verweigert, was er nicht füllen darf | GS-13, [Assistent](assistent.md) KA-01 bis KA-03, KA-15, KA-20 (GSO-10) |
| B-10 Vorschlag prüfen | KI als Fachteil: Vorschlag als Insel mit Quellen und Gründen, Entscheidungen als Knöpfe, Grund als Pflichtauswahl | GS-34, GS-44, Assistent KA-01 bis KA-04 sinngemäß (GSO-34) |
| B-11 Öffentlich und Externe | Rahmen ohne Anmeldung mit „Anmelden“, kurzer Bereich für Externe | GS-16 (GSO-39) |

Jede Seite der Referenz-App hat oben einen kurzen Text: was das Muster ist, wann man es nimmt, welche GS-Nummern
gelten, welcher Baustein es umsetzt. Die Texte verweisen auf dieses Dokument und wiederholen es nicht. Für sie gilt
GS-71 nicht; der Test `TestGS71_*` der Referenz-App prüft alles außerhalb davon (GSO-13).

## 12. Weiterentwickeln und austauschen

- **Eine Regel ändern** ist ein PR an diesem Dokument (neue Version, Eintrag in der Historie). Im selben PR oder
  direkt danach wird die Referenz-App nachgezogen, dazu ihre Kopie dieses Dokuments mit `scripts/reference-konzept.sh`
  in tabler-templ (GSO-14; die Kopie entfällt ab v1.0). Welche Apps folgen, steht im PR.
- **Nummern bleiben stabil.** Ersetzte Regeln kommen in den Anhang „Abgelöst“, ihre Nummer bleibt im Text stehen.
- **Andere Produkte** übernehmen die Regeln mit Nummer und Version („gilt Gestaltung v0.1 GS-40 bis GS-46“) und
  ersetzen nur die Spalte „hier“. Bis v1.0 liegt das Dokument hier, danach der produktneutrale Teil in tabler-templ (GSO-01).
- **Oberfläche über dem Theme** (grauer Grund, Inseln, Signalfarben, 1240 px, Links, 16 px in Feldern) steht heute in
  `libs/rahmen/static/theme.css`, die Referenz-App hat eine Kopie. Mit v1.0 kommt dieser produktneutrale Teil als
  eigene Ebene über den Themes nach tabler-templ; jedes Produkt bindet ihn mit einer Zeile ein (GSO-12).

## 13. Entschiedene Punkte

Alle nach der Empfehlung entschieden (Benutzer 15.11.2026). GSO-08 bis GSO-14 sind die Fragen aus tabler-templ PR #60.

| Nr. | Frage | Optionen | Empfehlung | Entscheidung |
|---|---|---|---|---|
| GSO-01 | **Wo liegt das Konzept für alle Produkte?** | (a) hier, andere Produkte verlinken; (b) in tabler-templ neben der Referenz-App (`docs/`), hier bleibt nur „hier“; (c) eigenes Repo | **(a) jetzt, (b) ab v1.0.** Erst hier abstimmen, solange die Regeln aus diesem Verbund kommen. Mit der Referenz-App wandert der produktneutrale Teil nach tabler-templ, das alle Produkte schon mit Version einbinden. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-02 | **Sprache der Referenz-App** | (a) Texte und Beispieldaten deutsch, Code und README englisch wie tabler-templ; (b) alles englisch | **(a):** Die Produkte sind deutsch, Beispiele sollen aussehen wie sie; tabler-templ bleibt im Code englisch. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-03 | **Liste bearbeiten als Baustein** (`liste-bearbeiten`, `modal-ganzseitig` gibt es nur in Auftrag) | (a) mit der Referenz-App als Block in tabler-templ (z. B. `blocks/rowedit`), Auftrag stellt um; (b) in `libs/rahmen`; (c) bleibt in Auftrag | **(a):** Die Referenz-App braucht ihn ohnehin, und die nächste App mit langer Liste kann ihn dann gleich nutzen. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-04 | **Inline-Bearbeiten im Beispiel `examples/crud` von tabler-templ** widerspricht GS-45 | (a) bleibt als Beispiel für die Bibliothek, die Referenz-App sagt deutlich, dass die Produkte es nicht nutzen; (b) entfernen | **(a):** tabler-templ ist eine allgemeine Bibliothek; die Referenz-App zeigt unsere Muster. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-05 | **Prüfbare Regeln als Test** (z. B. keine Drawer in Büro-Oberflächen, Insel höchstens ein „Bearbeiten“, keine Seite scrollt seitlich) | (a) Architektur- und Layout-Tests je Regel, wo es einfach geht; (b) nur Review | **(a):** wie bei Navigation (`navigation.go`) und Regelnummern; je Test die GS-Nummer im Namen. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-06 | **Verbund-Regel** | (a) nach der Entscheidung G-66 „Gestaltung nach dem Konzept Gestaltung“ in `verbund.md`; (b) ohne G-Nummer | **(a)**, damit der Verweis in den `kurz.md` der Apps kurz bleibt. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-07 | **Regeln für die mobilen Portale** | (a) eigener Abschnitt in diesem Dokument (GS-8x); (b) eigenes Konzept | **(a) später**, wenn ein zweites Portal neu gebaut wird. | **Benutzer 15.11.2026: nach der Empfehlung.** |

| GSO-08 | **Gruppe „Referenz“ in der Seitenleiste der Referenz-App** steht am Ende, nach GS-20 gehört dort „Verwaltung“ hin | (a) so lassen, §11 sagt, dass sie zur Referenz gehört; (b) umordnen | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-09 | **Gilt GS-46 für eigene Zeilen** eines Datensatzes (Positionen)? | (a) nein, nur für Listen anderer Datensätze; eigene Zeilen ganz, Scrollen in der Insel; (b) ja | **(a)**, so gebaut | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-10 | **KI-Assistent in §11 und Verhalten des Panels** | (a) B-09 aufnehmen, Satz zum Panel in GS-13; (b) ohne | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-11 | **App-Umschalter und Rechner** liegen in `libs/rahmen` | (a) bleiben beim Produkt, die Referenz zeigt nur die Knöpfe des Panels; (b) nach tabler-templ | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-12 | **Oberfläche über dem Theme** (`theme.css`) | (a) ab v1.0 als Ebene über den Themes nach tabler-templ; (b) bleibt je Produkt | **(a)**: Buchhalter, Support und Vergabe bekämen sie mit einer Zeile | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-13 | **Regelnummern in den Erklärtexten** der Referenz-App, obwohl GS-71 sie verbietet | (a) Ausnahme für die Erklärtexte; (b) keine Nummern | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-14 | **Kopie dieses Dokuments in der Referenz-App** | (a) zu jeder Regeländerung `scripts/reference-konzept.sh` in tabler-templ; ab v1.0 entfällt die Kopie; (b) nur Link auf GitHub | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |

## 14. Abgleich mit Buchhalter, Support und Vergabe (Bestandsaufnahme)

Entschieden nach der Empfehlung (Benutzer 15.11.2026, in der Session tabler-templ); die Regeln oben sind nachgezogen. Stand 15.11.2026, gelesen von der Session tabler-templ beim Bau der Referenz-App (PR #60), nichts geändert. Die
Produkte liegen in eigenen Repos (`mind-vm/buchhalter`, `mind-vm/support`, `mind-vm/vergabe`); jede Umstellung ist
ein PR dort, nach der Entscheidung hier.

| Regel | Buchhalter | Support | Vergabe |
|---|---|---|---|
| Theme (GS-15) | Tabler statt Graphite | Tabler | Tabler |
| Rahmen, Seitenleiste (GS-10, GS-11) | wie GS | wie GS | keine Seitenleiste |
| Gruppen der Leiste (GS-20) | nach dem Monatsablauf, Einträge nummeriert (Buchhalter OF-01) | – | – |
| Bearbeiten in der Seite statt Insel mit Dialog (GS-41 bis GS-45) | ja | ja | ja |
| Listen (GS-50) | ohne Suche und Seiten | drei Filter, 50 Zeilen (Support E-72) | ohne Suche und Seiten |
| Lange Listen (GS-55 bis GS-57) | ein großes Formular | – | ein großes Formular |

Muster, die dieses Konzept noch nicht kennt: eine Arbeitsseite mit drei Spalten (Support E-72), Tabellen, die unter
768 px zu Karten werden, und Prüfbreiten 1366 und 1280 px (Buchhalter OF-04, OF-06), ein Tastaturmodell (Buchhalter
OF-06.4), Portale für andere Zielgruppen (Vergabe: Bieter, öffentlich).

| Nr. | Frage | Optionen | Empfehlung | Entscheidung |
|---|---|---|---|---|
| GSO-15 | **Theme je Produkt** | (a) jedes Produkt wählt sein Theme (eigene Marke), alle GS-Regeln und die Ebene über dem Theme (GSO-12) gelten trotzdem; (b) alle auf Graphite | **(a):** Marke ist Sache des Produkts; GS-15 heißt dann „ein Theme je Produkt“. Gleich sind Aufbau, Muster und Signalfarben. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-16 | **Vergabe ohne Seitenleiste** | (a) die Büro-Oberfläche der Vergabe bekommt die Seitenleiste nach GS-11 und GS-20, ihre Portale für Bieter und Öffentlichkeit bleiben außen vor (GSO-22); (b) Ausnahme | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-17 | **Leiste nach einem Ablauf mit Nummern** (Buchhalter) | (a) als Regel GS-23 aufnehmen: Gruppen Tagesgeschäft dürfen in der Reihenfolge eines echten Ablaufs stehen, Nummern nur dann; (b) Buchhalter ordnet nach GS-20 um | **(a):** Der Monatsablauf ist eine echte Reihenfolge, und GS-20 erlaubt schon zwei Gruppen Tagesgeschäft. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-18 | **Bearbeiten in der Seite** (alle drei) | (a) nach und nach auf Inseln mit Dialog umstellen, je Seite beim nächsten Umbau, kein Großumbau; (b) Ausnahme je Produkt | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-19 | **Listen** (Filter, Seitengröße, fehlende Suche) | (a) GS-50 bleibt: eine Suche, ein Filter, 25 Zeilen; eine Arbeitsliste darf daneben Schnellfilter haben (GS-53, wie die Queue in Support), und jede Liste, die über 25 Zeilen wachsen kann, bekommt Suche und Seiten; (b) Seitengröße je Produkt | **(a)** | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-20 | **Lange Listen als ein großes Formular** (Buchhalter, Vergabe) | (a) auf `blocks/rowedit` umstellen (GS-55 bis GS-59); (b) bleibt | **(a)**, sobald tabler-templ mit dem Block veröffentlicht ist | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-21 | **Neue Muster aufnehmen** | (a) drei Spalten als Variante der Arbeitsseite (GS-34); Karten unter 768 px als Alternative zum Scrollen in der Insel, je Tabelle gewählt (GS-54); Prüfbreiten 1366 und 1280 px in die Layout-Tests (GS-14); das Tastaturmodell nach Lesen von Buchhalter OF-06.4 als eigener Abschnitt vorschlagen; (b) erst bei Bedarf | **(a)**, je Muster als Seite der Referenz-App | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-22 | **Portale für andere Zielgruppen** (Bieter, öffentlich, Kundenbereich) | (a) wie die mobilen Portale außerhalb der Büro-Regeln, später ein Abschnitt zusammen mit GSO-07; (b) gleich mitregeln | **(a) nur für Portale am Handy.** Portale für Externe am Schreibtisch regelt GSO-39. | **Benutzer 15.11.2026: nach der Empfehlung.** |

## 15. Fragen aus den Prüfungen von Support und Vergabe

Support und Vergabe haben ihre Oberfläche gegen v0.2 geprüft (Berichte `docs/reviews/gestaltung-2026-11-15.md` in
mind-vm/support#20 und mind-vm/vergabe#33, gesammelt weitergegeben von der Session tabler-templ). Die Fragen stehen
hier zusammengefasst, Doppelte zusammengelegt; in Klammern die Nummer im Bericht (Support A-x/B-x, Vergabe a-x/b-x).
Die Empfehlung folgt dem Bericht, wo nichts anderes steht. Entschieden nach der Empfehlung (Benutzer 15.11.2026, in der Session tabler-templ); die Regeln oben sind nachgezogen.

**Kontrast (zu GSO-25), nachgerechnet 15.11.2026:** Text in Pillen und Ampeln nutzt schon dunklere Töne und erreicht
WCAG AA (Grün `#1F7A3A` 4,9:1 auf seiner Fläche, Gelb `#8A5A00` 5,4:1, Rot `#B42318` 5,8:1). Als Textfarbe auf Weiß
fallen die Signalfarben selbst durch (Grün `#22A04B` 3,4:1, Gelb `#E0A100` 2,3:1, Rot `#E5484D` 3,9:1; Blau `#2563EB`
5,2:1 besteht). Für Punkte und Flächen (3:1) fällt nur Gelb durch, auf dem grauen Grund knapp auch Grün (3,1:1).

| Nr. | Frage | Empfehlung | Entscheidung |
|---|---|---|---|
| GSO-23 | **Kopf ohne Verbund** (Support A-1, Vergabe a-1): Produkte ohne Verbund haben keinen App-Umschalter, manche keinen Rechner und kein Panel. | GS-10: Notizen, Assistent, Rechner „soweit das Produkt sie hat“, der App-Umschalter nur im Verbund. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-24 | **Theme je Instanz** (Vergabe a-3): Vergabe läuft bei fremden Betreibern. | Zusammen mit GSO-15: Theme je Produkt; ein Produkt zum Selbstbetrieb lässt den Betreiber das Theme je Instanz wählen, Standard Graphite. Signalfarben bleiben fest. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-25 | **Kontrast** (Vergabe a-4): öffentliche Auftraggeber müssen BITV 2.0 / WCAG 2.1 AA erfüllen. | GS-07 verlangt in jedem Theme Text ≥ 4,5:1 und Punkte, Linien, Flächen ≥ 3:1. Signalfarben sind nie Textfarbe, Text nimmt die dunklen Töne. Gelb für Punkte dunkler (`#B57F00`, 3,5:1 auf Weiß), Grün prüfen. Ein Test in tabler-templ `themes` und in `libs/rahmen` prüft alle Paare; Vergabe löscht dann seine Korrektur in `vergabe.css`. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-26 | **Aktionen mit Pflichtgrund und gefährliche Aktionen** (Support B-4, Vergabe b-3) | Ergänzung von GS-44, GS-58 und GS-72: Eine rechtlich wirksame Aktion oder eine Aktion mit Grund ist ein kleiner Dialog mit Pflichtfeld und rotem Knopf; eine unumkehrbare, folgenreiche Aktion verlangt im Dialog eine getippte Bestätigung (z. B. die Adresse). `data-bestaetigen` bleibt für einfaches Löschen. Beispiel in B-07. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-27 | **Kurze Listen** (Support A-6: `rowedit` schon ab mehr als einem Feld; Vergabe a-5: Liste, die nur zusammen gültig ist, z. B. Gewichte = 100 %) | Beides: `rowedit` ist der Weg für jede Liste mit mehr als einem Feld je Zeile (GS-55, GS-56). Ausnahme: eine kurze Liste (bis etwa 10 Zeilen), deren Zeilen nur zusammen gültig sind, ist ein Formular im normalen Dialog (GS-42), mit genau diesem Grund. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-28 | **Arbeitsseite mit Bearbeiten in der Seite** (Support A-3, Vorgang als Dreispalter) | GS-34 sagt ausdrücklich: Bearbeiten in der Seite nur für die Arbeit selbst (Antwort, Vorschlag, Notiz); Inseln mit Kontext folgen GS-41 bis GS-45. Dreispalter als Variante zusammen mit GSO-21. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-29 | **Arbeitsliste als Start** (Support A-2) | Keine Änderung nötig: NV-01 erlaubt schon eine App ohne eigene Startseite, `/` führt dann auf den obersten Eintrag (wie Austausch „Läufe“). Zahlen kommen als Kopf der Liste (GS-30). Abweichend vom Bericht, der den Eintrag „Start“ umbenennen wollte. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-30 | **Mehr als ein Filter** (Support A-5) | Wie GSO-19: GS-50 bleibt. Status als Filter, „Nur meine“ als Schnellfilter über der Tabelle (GS-53), Label in der Suche (`label:rechnung`). | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-31 | **Benannte Markierungen und Text von außen** (Support A-4, B-3) | GS-07 kennt neben den Signalfarben benannte Markierungen mit eigenem Token je Theme: „intern, nie gesendet“ und „von außen“. Frei wählbare Farben (Labels) nur aus Signalfarben und Grau. Neue Regel in §9: Text von außen ist gekennzeichnet, escaped und nie ein Link; Baustein in tabler-templ. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-32 | **Datumsfelder** (Support A-7): das Feld zeigt das Format des Browsers | GS-73 nennt Eingabefelder als Ausnahme: das native Feld (`type=date`, `datetime-local`) ist am Handy die beste Eingabe. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-33 | **Notizen an einem Datensatz** (Support B-1) | Panel-Reiter „Notizen“ zeigt den aktuellen Datensatz (GS-13, wie Auftrag E-70), mit der Markierung „intern“ (GSO-31). | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-34 | **KI als Fachteil** (Support B-2): ein Agent entwirft, der Mensch entscheidet mit Pflichtgrund | Beispiel **B-10 „Vorschlag prüfen“** in der Referenz-App: Vorschlag als Insel mit Quellen und Gründen, Entscheidungen als Knöpfe, Grund als Pflichtauswahl; Assistent KA-01 bis KA-04 gelten sinngemäß. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-35 | **Fristen** (Support B-5: Restzeit; Vergabe a-6: „Uhr“) | GS-52: eine Frist ist eine Ampel, deren Wort die Restzeit ist („noch 58 Min.“, „überfällig seit 20 Std.“). GS-73: „TT.MM.JJJJ, HH:MM Uhr“ für Fristen und Zeitpunkte mit rechtlicher Wirkung, sonst „07:30“. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-36 | **Einstellungen** (Support B-6) | Inseln mit Dialogen wie eine Detailseite (GS-32), kein großes Formular. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-37 | **Assistent im Selbstbetrieb** (Vergabe a-2) | Ergänzung im Konzept [Assistent](assistent.md) zu KA-04 und KA-20: Ein Produkt zum Selbstbetrieb bietet den Assistenten nur mit einem vom Betreiber gewählten Anbieter an; vertrauliche Daten (Angebote vor der Öffnung) sind nie Werkzeug. Bis dahin hat Vergabe kein Panel. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-38 | **Eigenes JavaScript für Kryptografie** (Vergabe a-7, WebCrypto bei der Abgabe) | GS-09 nennt Kryptografie beim Nutzer als zulässigen Grund für eigenes JavaScript, eingebettet und getestet. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-39 | **Öffentliche Seiten und Bereich für Externe am Schreibtisch** (Vergabe b-1) | Derselbe Rahmen mit kurzer Leiste; ohne Anmeldung rechts im Kopf „Anmelden“ als Hauptaktion, „Registrieren“ schlicht. Beispiel **B-11** in der Referenz-App. Portale am Handy bleiben bei GSO-22. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-40 | **Arbeitsstand und veröffentlichte Fassung** (Vergabe b-2) | Muster als Unterfall von GS-32 und GS-60: Hinweis oben „Es gibt Änderungen gegenüber Fassung N“ mit der Hauptaktion „Neue Fassung veröffentlichen“ (mit Grund, GSO-26), Fassungen als Insel. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-41 | **Dateien hochladen** (Vergabe b-4) | Dialog „… hinzufügen“ an der Insel (GS-44) mit Datei, Art und Bezug; die Liste bleibt lesend, Entfernen mit Nachfrage. Geht in tabler-templ schon mit `detail.Dialog` und `multipart`. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-42 | **Unveränderliches Protokoll** (Vergabe b-5) | GS-46 gilt auch hier: die letzten 10, „Alle N“ führt zur vollständigen Liste mit Suche. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-43 | **Detailseite ohne fachliche Nummer** (Vergabe b-6) | GS-32: Nummer „wenn vorhanden“; eine interne ID nie (GS-71). | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-44 | **Hinweis „Demo-Instanz“** (Vergabe b-7) | Wie „nur lesen“ (GS-63): eine Zeile im Rahmen über der Meldungszeile, Farbe Info. | **Benutzer 15.11.2026: nach der Empfehlung.** |

## 16. Weitere Fragen

| Nr. | Frage | Optionen | Empfehlung | Entscheidung |
|---|---|---|---|---|
| GSO-45 | **Lange Texte in „Liste bearbeiten“** (Vergabe, vergabe#36): Eine Leistung im Leistungsverzeichnis hat 12 Felder, darunter Langtext und Anforderungen; als feste Spalten scrollt die Zeile weit seitlich. | (a) aufklappbare zweite Zeile im selben Formular mit Textfeldern über die ganze Breite, als Teil von GS-56; (b) lange Texte als einzeilige Felder in den Spalten; (c) lange Texte in einem eigenen Dialog je Zeile | **(a):** zu, nach einem Fehler in einem ihrer Felder offen; Speichern und Löschen in der ersten Zeile; Enter im Textfeld speichert nicht; Baustein `rowedit.Details` in tabler-templ, Beispiel in B-05 | **Benutzer 16.11.2026: (a).** |
| GSO-46 | **„Foto aufnehmen“ beim Hochladen** (Buchhalter, buchhalter#48, über tabler-templ) | (a) am Handy (unter lg) ein zweiter Knopf „Foto aufnehmen“ neben „Datei wählen“, nur wo ein Foto als Beleg taugt; (b) nur „Datei wählen“ | **(a):** spart am Handy einen Schritt, am Desktop ändert sich nichts | **Benutzer 16.11.2026: nach der Empfehlung.** |
| GSO-47 | **Dateien auf eine Liste oder Insel ziehen** (Buchhalter, buchhalter#48) | (a) ja, die Insel bzw. Liste, zu der die Datei gehört; mit Art oder Bezug öffnet sich der Dialog mit der Datei, sonst lädt sie gleich hoch; (b) nur über den Knopf im Dialog; (c) nur in einer eigenen Fläche im Dialog | **(a):** der Knopf bleibt der sichtbare Weg, das Ziehen ist eine Abkürzung für die Maus; beim Ziehen „Zum Hochladen loslassen“ | **Benutzer 16.11.2026: nach der Empfehlung.** |
| GSO-48 | **Ausliefern statischer Dateien** (Support: Icons flackerten, die ganze Icon-Webfont mit 506 kB und die Schriften aller Themes wurden geladen, `/static` ohne Cache-Control) | (a) Satz bei GS-09; (b) keine Regel, nur Doku in tabler-templ | **(a):** gemessen nach dem Vorgehen von `libs/rahmen` (RA-15, RA-16): 11,7 kB und 3 Anfragen statt 207 kB und 17, FCP 0,25 s statt 0,6 s | **Benutzer 16.11.2026: nach der Empfehlung.** |
| GSO-49 | **Navigation ohne vollen Seitenaufbau (hx-boost)** (Vergabe) | (a) nein, voller Seitenaufbau, schnell genug durch GS-09 (statische Dateien im Cache); (b) hx-boost als Option in tabler-templ `blocks/shell`, je Produkt einschaltbar, je Seite und Link abschaltbar; (c) hx-boost überall | **(a) jetzt, (b) erst, wenn nach GS-09 gemessen ein Seitenwechsel spürbar zu langsam ist:** im Auftragsbutler kostet ein Wechsel mit Cache nur das HTML, 11,7 kB, FCP 0,25 s. Bedingungen für ein späteres (b): Titel, Brotkrumen und aktiver Eintrag folgen; die Skripte von `detail`, `rowedit`, `upload` und `sidepanel` starten neu, ohne doppelte Listener; die Meldung per Cookie und der Fehlerfall mit 422 und offenem Dialog funktionieren; Verlauf und Zurück funktionieren; der Fokus springt auf den neuen Inhalt; Seiten mit WebCrypto sind ausgenommen | **Benutzer 16.11.2026: nach der Empfehlung.** |
| GSO-50 | **„Liste bearbeiten“ am Handy** (Vergabe): Der Dialog „Leistungsverzeichnis bearbeiten“ hat 10 feste Spalten plus aufklappbare Details (GSO-45); am Handy scrollt die Liste seitlich, die rechten Spalten sieht man erst nach dem Wischen. Gilt GS-54 auch hier? | (a) ja, unter 768 px Karten, je Liste gewählt; (b) nein, „Liste bearbeiten“ scrollt immer seitlich in ihrem Dialog | **(a):** Bei vielen Spalten ist seitliches Wischen am Handy kaum bedienbar; GS-54 macht es für Tabellen schon so. | **Benutzer 16.11.2026 in der Session von Vergabe: „Ja, am besten so machen mit dem Handy als Karten.“, nach der Empfehlung (a).** |
