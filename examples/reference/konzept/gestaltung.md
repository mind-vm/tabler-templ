# Gestaltung: Rahmen, Seiten, Inseln und Tabellen

| | |
|---|---|
| **Status** | v0.2 · 15.11.2026 · Entschieden (Benutzer 15.11.2026: alle Empfehlungen GSO-01 bis GSO-07; Verbund G-66): die geltenden Regeln der Oberfläche als ein Konzept GS-01 bis GS-75, für alle Produkte nutzbar; Beispiele als lebende Referenz-App in tabler-templ (Brief). [Historie](../historie/docs/konzepte/gestaltung.md) |
| **Anlass** | Wunsch des Benutzers vom 15.11.2026: aus Seitenleiste, Application Layout, „Inseln bearbeiten“ und dem Bearbeiten von Tabellen ein Design- und UX-Konzept mit Beispielen machen, das sich nach und nach überarbeiten, standardisieren und zwischen den Produkten austauschen lässt; die Beispiele als lebende Referenz-App in tabler-templ. |
| **Gilt zusammen mit** | [Navigation](navigation.md) (NV-01 bis NV-09, G-62), [Auftrag ADR-0013](../../apps/auftrag/docs/adr/0013-inseln-mit-dialogen.md) (Inseln mit Dialogen, E-45, E-70), [Verbund](verbund.md) G-41, G-44, G-62, G-63; [`libs/rahmen`](../../libs/rahmen/README.md) (Gerüst, Theme, Panel, Fehlerseiten), [`libs/liste`](../../libs/liste/README.md); Einsatz [ux-disponent](../../apps/einsatz/docs/ux-disponent.md) (UX-01 bis UX-14); Brief [gestaltung-referenz-app](../briefs/gestaltung-referenz-app.md) |
| **Präfix** | GS-xx (Regeln der Gestaltung), GSO-xx (offene Punkte) |
| **Entscheidung** | Benutzer 15.11.2026: „alle Empfehlungen annehmen“ entscheidet GSO-01 bis GSO-07 jeweils wie empfohlen (§13); Verbund G-66. |

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
| GS-01 | **Gegenstand** sind die Büro-Oberflächen (am Schreibtisch, Maus und Tastatur, auch schmal am Handy). Die mobilen Portale (`/m/`, offline, ein Daumen) haben eigene Regeln und stehen nur kurz in §10. |
| GS-02 | **Für alle Produkte:** Die Regeln sind produktneutral geschrieben. Was nur hier gilt (Paketnamen, Apps, Nummern anderer Dokumente), steht in der Spalte oder dem Absatz „hier“. Ein anderes Produkt übernimmt die GS-Regeln mit derselben Nummer und ersetzt nur „hier“. |
| GS-03 | **Quelle der Wahrheit:** Die Regel steht in diesem Dokument. Die Referenz-App (§11) zeigt sie, die Bausteine in tabler-templ und `libs/rahmen` setzen sie um. Weicht eins davon ab, gilt die Regel, und der Rest wird nachgezogen (oder die Regel wird per PR geändert). |

## 2. Grundsätze

| Nr. | Grundsatz | Warum |
|---|---|---|
| GS-05 | **Lesen vor Bearbeiten.** Eine Seite zeigt zuerst an. Bearbeiten ist ein eigener, sichtbarer Schritt (Dialog oder eigene Seite). | Versehentliche Änderungen kommen seltener vor, und man sieht eindeutig, wo sich was ändern lässt (ADR-0013, abgelehnte Option „inline“). |
| GS-06 | **Eine Aktion je Ort.** Eine Insel hat höchstens einen Knopf „Bearbeiten“, eine Seite oben höchstens eine Hauptaktion („Neuer Kunde“). Weitere Aktionen stehen im Dialog oder im Menü „…“. | Man muss nicht suchen, und die Seite bleibt ruhig. |
| GS-07 | **Farbe nur mit Bedeutung.** Der Grund ist neutral (Graphite: schwarz-weiß). Farbig sind nur die Signalfarben: Grün (gut), Gelb (Achtung), Rot (Problem), Blau (Info, heute). Sie gelten in jedem Theme gleich. | Rot fällt auf, weil sonst nichts bunt ist (`libs/rahmen`, Theme). |
| GS-08 | **Die Rolle blendet aus, der Use Case prüft.** Die Oberfläche zeigt nur, was die Rolle darf. Geprüft wird trotzdem auf dem Server. Gibt es für einen Abschnitt keinen Use Case zum Ändern, gibt es auch keinen Knopf. | Auf ein „Bearbeiten“ soll nie ein 403 folgen (G-41, NV-04, ADR-0013). |
| GS-09 | **Server rendert, Formulare schicken ab.** Seiten entstehen auf dem Server (templ), Formulare arbeiten mit POST und Rücksprung (PRG), wenig eigenes JavaScript. Nichts kommt von einem CDN. | Ein Binary, schnell, offline-fest, prüfbar (`libs/rahmen` `TestNichtsVonAussen`). |

## 3. Application Layout (Rahmen)

| Nr. | Teil | Regel | hier |
|---|---|---|---|
| GS-10 | **Kopf** | 56 px hoch. Links die Marke der App, rechts Notizen, Assistent, Rechner und ganz rechts der App-Umschalter (mit der Zahl „dringend“). Keine eigene Navigation im Kopf. | `rahmen.Seite`, `rahmen.Umschalter` (G-63), Rechner RA-12 |
| GS-11 | **Seitenleiste** | 240 px links, gruppiert nach §4, je Eintrag ein Symbol und ein Wort. Beim Wechsel über die Breite lg kein Übergang (der Inhalt springt, statt hinter die Leiste zu gleiten). | `sidebar.Group` aus tabler-templ, `rahmen.Symbol` |
| GS-12 | **Inhalt** | Grauer Grund, Breite höchstens 1240 px, linksbündig neben der Leiste. Oben der Seitentitel (bei Detailseiten mit Brotkrumen „Bereich › Liste › Datensatz“) und rechts davon die Hauptaktion (GS-06). Darunter die Meldungszeile (GS-60), dann die Inseln. | `theme.css` (`--grund`, `--linie`) |
| GS-13 | **Panel rechts** | Arbeitsbereich neben der Seite mit Reitern (Notizen, Assistent). Er bleibt beim Seitenwechsel offen, der Zustand steht im Cookie. Kein Formular zum Bearbeiten im Panel. | `rahmen.Panel` (RA-11, Auftrag E-70) |
| GS-14 | **Handy und schmal** | Unter lg wird die Leiste zum Menü (gleiche Reihenfolge, NV-07), das Panel zum Sheet über der Seite (öffnet sich nach einem Seitenwechsel nicht von selbst), ein Dialog zum Vollbild. Die Seite scrollt nie seitlich; breite Tabellen scrollen in ihrer Insel. Felder haben 16 px Schrift, damit das Handy nicht zoomt. | `tools/e2e-staging/layout.mjs` prüft das |
| GS-15 | **Theme** | Ein Theme für den ganzen Verbund, an einer Stelle wählbar. Farben, Ecken und Linien sind Tokens; Signalfarben gelten in jedem Theme (GS-07). | `VERBUND_THEME`, Standard Graphite |

## 4. Navigation

Die Regeln stehen in [Navigation](navigation.md) §2 und gelten hier unverändert. Kurz:

| Nr. | Regel |
|---|---|
| GS-20 | „Start“ steht oben ohne Gruppentitel (NV-01). Danach kommen Gruppen in fester Reihenfolge: „Für mich“/„Unterwegs“, höchstens zwei Gruppen Tagesgeschäft, „Erledigen“, „Verwaltung“ (NV-02). In „Verwaltung“ erst die Stammdaten, dann Einstellungen, Importieren, Daten exportieren (NV-03). |
| GS-21 | Höchstens etwa 5 bis 7 Einträge je Gruppe, keine Unterbereiche mit „Zurück“ (NV-06). Gleiches heißt in allen Apps gleich (NV-05). Jeder Eintrag hat einen stabilen Anker für Touren und Tests (NV-08). |
| GS-22 | Die Leiste ist nur für die Hauptbereiche da. Von einem Datensatz zu einem anderen führt ein Link in der Insel oder in der Tabellenzeile, nie ein Eintrag in der Leiste. |

## 5. Seitentypen

| Nr. | Typ | Aufbau | Vorbild |
|---|---|---|---|
| GS-30 | **Start** | Was heute zu tun ist, nicht wozu die App da ist. Kacheln mit Zahlen, und jede Zahl führt auf die gefilterte Liste. Darunter Arbeitslisten („Brennt“, „Diese Woche“). | Einsatz UX-01, Konto „Mein Tag“ |
| GS-31 | **Liste** | Titel und Hauptaktion („Neuer …“), Suche, ein Filter, Tabelle, Seiten (§7.1). | Auftrag „Kunden“, E-49 |
| GS-32 | **Detailseite** | Kopf mit Name, Nummer und Status, darunter die Inseln in fester Reihenfolge: zuerst, woran man arbeitet, dann was aus anderen Apps kommt (§6). Eine Karte „Als Nächstes“ zeigt, was zu tun ist, wenn es einen Ablauf gibt. | Auftrag „Auftrag“, Einsatz UX-03, UX-07 |
| GS-33 | **Anlegen** | Eine eigene Seite (`/kunden/neu`), nicht als Dialog über der Liste. Sie hat nur die Felder, die zum Anlegen nötig sind, und führt nach dem Speichern auf die neue Detailseite. | ADR-0013 Punkt 4 |
| GS-34 | **Arbeitsseite** | Ausnahme für Arbeit am Stück (Plantafel, Freigabe, Erfassen, Importieren): eigenes Layout erlaubt, Rahmen und Regeln §8 und §9 gelten trotzdem. | ADR-0013 |

## 6. Inseln und Dialoge

| Nr. | Regel | Quelle |
|---|---|---|
| GS-40 | **Insel:** eine weiße Karte auf grauem Grund mit Überschrift. Der Inhalt ist nur zum Lesen da, meist als Paare aus Bezeichnung und Wert in zwei Spalten (am Handy eine). Fehlende Werte zeigen „–“. | ADR-0013 Punkt 1, 5 |
| GS-41 | **Bearbeiten:** rechts oben in der Insel steht genau ein Knopf „Bearbeiten“. Gibt es keinen passenden Use Case, entfällt er. | ADR-0013 Punkt 1, GS-08 |
| GS-42 | **Dialog:** ein zentrierter Dialog mit genau den Feldern dieser Insel. Der Titel heißt wie die Insel, darunter steht klein, worum es geht („A-2026-001 · Wohnungsgesellschaft Nordstadt“). Unten links „Abbrechen“, rechts „Speichern“. Gespeichert wird nur diese Insel. | ADR-0013 Punkt 2 |
| GS-43 | **Fehler im Dialog:** Scheitert das Speichern, öffnet die Seite den Dialog wieder, mit der Meldung oben und den eingegebenen Werten; das fehlerhafte Feld ist markiert. | ADR-0013 Punkt 2 |
| GS-44 | **Kleine Ergänzungen** an einem bestehenden Datensatz (Nachtrag, Zuweisung, PIN) sind Dialoge seiner Detailseite, mit einem Knopf in der Insel, zu der sie gehören. | ADR-0013 Punkt 4 |
| GS-45 | **Keine Drawer, kein Bearbeiten in der Insel.** Felder stehen nie direkt in der Leseansicht. Das Panel rechts ist kein Formular (GS-13). | ADR-0013, E-70 |
| GS-46 | **Unterlisten in einer Insel** zeigen höchstens 10 Zeilen und darunter „Alle N“, das auf die Liste mit Filter führt. | Verbund G-44 |

## 7. Tabellen

### 7.1 Anzeigen

| Nr. | Regel | Quelle |
|---|---|---|
| GS-50 | **Liste:** oben eine Textsuche (jedes Wort muss vorkommen, ohne Groß- und Kleinschreibung, ä = ae), daneben höchstens ein einfacher Filter, unten Seiten mit 25 Zeilen. Suche, Filter und Seite stehen in der URL, damit Zurück, Lesezeichen und Links funktionieren. | Auftrag E-49, `libs/liste` |
| GS-51 | **Rechnen in der Datenbank:** Suche, Filter, Seite, Summen und Gesamtzahl rechnet SQL in wenigen festen Abfragen je Seite, nie eine je Zeile. | G-44, ADR-0011 |
| GS-52 | **Spalten:** Was man erkennt, steht vorn (Name, Ort), Nummern dahinter. Name oder Nummer ist der Link zur Detailseite. Zahlen und Beträge stehen rechtsbündig mit gleich breiten Ziffern, Einheiten in eigener Spalte oder hinter der Zahl. Status als Pille, Ampel als Punkt mit Wort. | Einsatz UX-04, UX-06; `theme.css` |
| GS-53 | **Sortierung** nach dem, was Handlung braucht (Fälligkeit, „Frist fehlt“ oben), nicht nach Anlagedatum. Eine Liste zum Abarbeiten hat Schnellfilter statt vieler Spalten. | Einsatz UX-05, Ausschreibungen AS002 |
| GS-54 | **Breite Tabellen** scrollen seitlich in ihrer Insel. Die Seite selbst scrollt nie seitlich (GS-14). | |

### 7.2 Bearbeiten

| Nr. | Regel | Quelle |
|---|---|---|
| GS-55 | **Lange Listen im ganzseitigen Dialog.** Auf der Detailseite wird die Liste (z. B. Positionen eines Angebots) nur gelesen. „Bearbeiten“ öffnet einen Dialog über fast die ganze Seite, unter lg als echtes Vollbild. | ADR-0013 Punkt 3; hier `modal-ganzseitig` in Auftrag |
| GS-56 | **Jede Zeile ein Formular.** Im Dialog ist jede Zeile ein eigenes kleines Formular mit festen Spalten. Rechts stehen „Speichern“ (Haken) und „Löschen“ (Papierkorb). Enter speichert die Zeile. Ein Hinweis oben sagt: „Jede Zeile wird für sich gespeichert.“ | hier `liste-bearbeiten` in Auftrag `seite_angebot.templ` |
| GS-57 | **Neue Zeile unten** als leeres Formular mit Vorschlag (z. B. nächste Positionsnummer) und „Hinzufügen“. | Auftrag „Positionen bearbeiten“ |
| GS-58 | **Löschen fragt nach** („Position 01.01.0010 löschen?“), über `data-bestaetigen` am Knopf (heute die Nachfrage des Browsers). Nach dem Löschen sagt die Meldung, was gelöscht wurde. | hier Auftrag `formulare.js`; GS-60 |
| GS-59 | **Summen und berechnete Spalten** (Betrag = Menge × Preis) zeigt der Dialog nur an und rechnet sie nach jedem Speichern neu. Schreibgeschützte Listen (z. B. nach Abgabe eines Angebots) zeigen den Grund statt der Knöpfe. | Auftrag „Positionen bearbeiten“ |

## 8. Rückmeldung und Zustände

| Nr. | Regel | Quelle |
|---|---|---|
| GS-60 | **Meldung nach jeder Aktion:** eine wegklickbare Zeile über dem Inhalt, die sagt, was geschehen ist, und den nächsten Schritt als Link anbietet („Einsatz disponiert. Zur Plantafel“). Grün für Erfolg, Rot für Fehler. | Einsatz UX-06, `shell.Flash` |
| GS-61 | **Leerer Zustand:** sagt, warum hier nichts steht und was hilft, mit einem Knopf für den nächsten Schritt. Eine leere Suche bietet „Suche zurücksetzen“ an. | Einsatz UX-09 |
| GS-62 | **Fehlerseiten** (nicht gefunden, keine Rechte, Fehler) stehen im Rahmen mit Seitenleiste und erklären in einem Satz, was los ist. Technische Einzelheiten gehören ins Log, nie in den Browser. | `libs/rahmen` RA-01 bis RA-05 |
| GS-63 | **Nur lesen:** Ist die Firma gesperrt oder gekündigt, zeigt der Rahmen einen Hinweis mit Link zum Export. Die Knöpfe zum Bearbeiten entfallen. | Datenhoheit DH-12 |
| GS-64 | **Laden:** Teile, die länger brauchen (Zahlen der Seitenleiste, „dringend“ im Umschalter), kommen nach dem Seitenaufbau nach; die Seite selbst steht sofort. | [Laden im Hintergrund](laden-im-hintergrund.md) |

## 9. Sprache und Form

| Nr. | Regel | Quelle |
|---|---|---|
| GS-70 | **Begriffe der Nutzer:** die Wörter aus dem Glossar der App (`CONTEXT.md`), keine Fachsprache der Entwicklung. Gleiches heißt überall gleich. | NV-05, Einsatz UX-13 |
| GS-71 | **Keine internen Kennungen:** keine Regelnummern (E-xx, G-xx, §), keine IDs in Texten und Meldungen. Gesetze („§ 14 UStG“) dürfen bleiben. | Auftrag E-52, `rahmen.OhneRegelnummern`, Architekturtest |
| GS-72 | **Knöpfe sagen, was passiert:** „Speichern“, „Angebot anlegen“, „Freigeben“, nicht „OK“ oder „Weiter“. Die Hauptaktion ist gefüllt, alle anderen sind umrandet oder schlicht. Gefährliche Aktionen sind rot und fragen nach. | |
| GS-73 | **Zahlen und Datum** im deutschen Format: 1.234,50 €, TT.MM.JJJJ, Uhrzeit 07:30. Zeiträume als „01.04.2026 – 31.03.2028“. Listen werden mit Komma verbunden. | Auftrag E-52 |
| GS-74 | **Links** im Fließtext und in Tabellen sind dezent unterstrichen (der Akzent ist schwarz wie der Text). | `libs/rahmen` Theme |
| GS-75 | **Beispieldaten** ohne erkennbare Branche (Musterdienst GmbH, Wohnungsgesellschaft Nordstadt), in Referenz-App, Mockups und Screenshots gleich. | Verbund G-60 |

## 10. Mobile Portale (nur kurz)

Die mobilen Portale (`/m/` in Auftrag, Einsatz, Zeit) nutzen dieselben Werte (Akzent, Schrift, Ecken, Signalfarben)
aus eigenem CSS ohne Tabler. Sonst sind sie anders gebaut: ein großer Knopf je Schritt, offline, kein Rahmen, keine
Seitenleiste. Eigene Regeln dafür sind offen (GSO-07).

## 11. Beispiele: die Referenz-App in tabler-templ

Wunsch des Benutzers (15.11.2026): Das Konzept soll lebendig sein. Dafür kommt in tabler-templ eine Referenz-App,
die alle Muster dieses Dokuments mit den echten Bausteinen zeigt und jeweils kurz erklärt. Sie ist kein Mockup:
Ändert sich ein Baustein, ändert sich das Beispiel mit, und Screenshot-Tests merken eine Abweichung. Den Auftrag dazu
beschreibt der Brief [gestaltung-referenz-app](../briefs/gestaltung-referenz-app.md).

| Beispiel | zeigt | Regeln |
|---|---|---|
| B-01 Rahmen | Kopf, gruppierte Seitenleiste, Inhalt, Panel rechts mit Reitern; schmal als Menü und Sheet | GS-10 bis GS-15, GS-20 bis GS-22 |
| B-02 Start | Kacheln mit Zahlen, die auf gefilterte Listen führen, Arbeitsliste „Brennt“ | GS-30 |
| B-03 Liste | Suche, Filter, 25 Zeilen, Seiten, Status-Pillen, Ampel, Zahlen rechtsbündig, Zustand in der URL | GS-31, GS-50 bis GS-54 |
| B-04 Detailseite | Kopf, Inseln mit Bearbeiten, Dialog, Fehler im Dialog, Unterliste mit „Alle N“ | GS-32, GS-40 bis GS-46 |
| B-05 Liste bearbeiten | ganzseitiger Dialog, Zeile speichern, löschen mit Nachfrage, neue Zeile mit Vorschlag, Summe | GS-55 bis GS-59 |
| B-06 Anlegen | eigene Seite, wenige Felder, danach die Detailseite mit Meldung | GS-33, GS-60 |
| B-07 Zustände | Meldung mit nächstem Schritt, leerer Zustand, leere Suche, Fehlerseite, „nur lesen“, Platzhalter beim Laden | GS-60 bis GS-64 |
| B-08 Marke und Form | Theme-Tokens, Signalfarben, Schrift, Knopfarten, Pillen, Ampel, Zahlen- und Datumsformat | GS-07, GS-15, GS-70 bis GS-75 |

Jede Seite der Referenz-App hat oben einen kurzen Text: was das Muster ist, wann man es nimmt, welche GS-Nummern
gelten, welcher Baustein es umsetzt. Die Texte verweisen auf dieses Dokument und wiederholen es nicht.

## 12. Weiterentwickeln und austauschen

- **Eine Regel ändern** ist ein PR an diesem Dokument (neue Version, Eintrag in der Historie). Im selben PR oder
  direkt danach wird die Referenz-App nachgezogen; welche Apps folgen, steht im PR.
- **Nummern bleiben stabil.** Ersetzte Regeln kommen in den Anhang „Abgelöst“, ihre Nummer bleibt im Text stehen.
- **Andere Produkte** übernehmen die Regeln mit Nummer und Version („gilt Gestaltung v0.1 GS-40 bis GS-46“) und
  ersetzen nur die Spalte „hier“. Bis v1.0 liegt das Dokument hier, danach der produktneutrale Teil in tabler-templ (GSO-01).

## 13. Entschiedene Punkte

Alle nach der Empfehlung entschieden (Benutzer 15.11.2026).

| Nr. | Frage | Optionen | Empfehlung | Entscheidung |
|---|---|---|---|---|
| GSO-01 | **Wo liegt das Konzept für alle Produkte?** | (a) hier, andere Produkte verlinken; (b) in tabler-templ neben der Referenz-App (`docs/`), hier bleibt nur „hier“; (c) eigenes Repo | **(a) jetzt, (b) ab v1.0.** Erst hier abstimmen, solange die Regeln aus diesem Verbund kommen. Mit der Referenz-App wandert der produktneutrale Teil nach tabler-templ, das alle Produkte schon mit Version einbinden. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-02 | **Sprache der Referenz-App** | (a) Texte und Beispieldaten deutsch, Code und README englisch wie tabler-templ; (b) alles englisch | **(a):** Die Produkte sind deutsch, Beispiele sollen aussehen wie sie; tabler-templ bleibt im Code englisch. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-03 | **Liste bearbeiten als Baustein** (`liste-bearbeiten`, `modal-ganzseitig` gibt es nur in Auftrag) | (a) mit der Referenz-App als Block in tabler-templ (z. B. `blocks/rowedit`), Auftrag stellt um; (b) in `libs/rahmen`; (c) bleibt in Auftrag | **(a):** Die Referenz-App braucht ihn ohnehin, und die nächste App mit langer Liste kann ihn dann gleich nutzen. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-04 | **Inline-Bearbeiten im Beispiel `examples/crud` von tabler-templ** widerspricht GS-45 | (a) bleibt als Beispiel für die Bibliothek, die Referenz-App sagt deutlich, dass die Produkte es nicht nutzen; (b) entfernen | **(a):** tabler-templ ist eine allgemeine Bibliothek; die Referenz-App zeigt unsere Muster. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-05 | **Prüfbare Regeln als Test** (z. B. keine Drawer in Büro-Oberflächen, Insel höchstens ein „Bearbeiten“, keine Seite scrollt seitlich) | (a) Architektur- und Layout-Tests je Regel, wo es einfach geht; (b) nur Review | **(a):** wie bei Navigation (`navigation.go`) und Regelnummern; je Test die GS-Nummer im Namen. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-06 | **Verbund-Regel** | (a) nach der Entscheidung G-66 „Gestaltung nach dem Konzept Gestaltung“ in `verbund.md`; (b) ohne G-Nummer | **(a)**, damit der Verweis in den `kurz.md` der Apps kurz bleibt. | **Benutzer 15.11.2026: nach der Empfehlung.** |
| GSO-07 | **Regeln für die mobilen Portale** | (a) eigener Abschnitt in diesem Dokument (GS-8x); (b) eigenes Konzept | **(a) später**, wenn ein zweites Portal neu gebaut wird. | **Benutzer 15.11.2026: nach der Empfehlung.** |
