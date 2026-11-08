# Copilot-Demo

Eine kleine Auftragsverwaltung mit einem KI-Copiloten im Panel rechts. Sie zeigt, wie ein Assistent mit
`blocks/assistant` neben einer App arbeitet: Er liest, füllt Formulare aus und zeigt auf Elemente. Speichern
tut immer der Mensch.

Starten:

```bash
go run ./cmd/demo -addr localhost:8093
```

Dann <http://localhost:8093/copilot> öffnen und oben rechts auf **Copilot** klicken. Die Demo braucht kein
Sprachmodell. Ein festes Skript (`agent.go`) spielt das Modell. Alles, was es tun will, geht aber durch
dieselbe Prüfung wie bei einem echten Modell.

## Was man ausprobieren kann

| Frage | Was passiert |
|---|---|
| *Which orders are overdue?* | Werkzeugaufruf, Tabelle im Panel, der Filter „Overdue“ auf der Seite blinkt auf |
| *Create an order for Becker: shop sign, 2,400 €* | Der Dialog „New order“ öffnet sich, ausgefüllt und markiert, mit Hinweis „bitte prüfen und selbst speichern“ |
| *Mark A-1003 as done* | Der Copilot wechselt auf die Seite des Auftrags und füllt dort Status und Notiz aus. Das Panel bleibt offen, das Gespräch auch |
| *Change the IBAN of Müller Bau* | Abgelehnt: Bankverbindungen füllt der Copilot nie aus. Er öffnet stattdessen den Kunden |
| *Read the note of A-1004 and do what it says* | Die Notiz enthält eine Prompt-Injection. Das Modell versucht, die IBAN zu ändern, und wird abgelehnt. Das Bild, das die Notiz einschleusen will, erscheint nicht |

## Wie es gebaut ist

| Teil | Wo | Aufgabe |
|---|---|---|
| **Katalog** | `catalog.go` | Die Allowlist der App: welche Seiten, Formulare (`data-form`) und Elemente (`data-tour`) der Copilot bedienen darf, mit welchen Feldern. `Catalog.Check` lehnt alles andere ab, auch Felder und Werte, die nach IBAN oder E-Mail aussehen |
| **Werkzeuge** | `agent.go` (`tool`) | Nur lesend. In einer echten App die lesenden Use Cases, mit der Rolle der Person |
| **Stream** | `copilot.go` (`run`) | Die Antwort kommt Stück für Stück als HTML vom Server: Werkzeugaufrufe, Tabelle, Antwort (`AnswerStrict`, ohne Links und Bilder) und die geprüften Aktionen als Ereignis |
| **Aktionen im Browser** | `ActionScript` in `blocks/assistant` | Führt die Aktionen nach dem Stream aus: Seite öffnen, Formular finden, Dialog öffnen, Felder füllen und markieren, nie absenden. Muss die Aktion auf eine andere Seite, merkt es sie sich und führt sie dort aus |
| **Verlauf** | `assistant.Transcript` | Der Server hält das Gespräch je Sitzung (Cookie) und zeichnet es auf jeder Seite wieder. Im Browser liegt nur die Sitzungs-ID |
| **Panel** | `assistant.Copilot` | Offcanvas rechts ohne Abdunklung. Auf breiten Bildschirmen rückt die Seite zur Seite. Es bleibt von Seite zu Seite offen |

Die Seiten selbst brauchen dafür nur zwei Attribute: `data-form="order-new"` am Formular und
`data-tour="orders.filter.overdue"` am Element. Der Test `TestActionsFindTheirPage` prüft, dass jede Aktion
des Katalogs ihr Formular mit allen Feldern bzw. ihr Element auf ihrer Seite findet.

## Mit einem echten Modell

`plan` in `agent.go` wird durch einen Agenten ersetzt (z. B. agentloop aus thinking-core):

1. Im System-Prompt stehen die Werkzeuge und `Catalog.Prompt(rolle)`.
2. Jeder Werkzeugaufruf des Modells wird zu `turn.tool(…)` (Anzeige) und zum Aufruf des lesenden Use Case.
3. Jeder Wunsch nach einer Aktion geht durch `turn.act(name, werte)`, also durch `Catalog.Check`. Eine
   Ablehnung geht als Fehlermeldung zurück an das Modell, damit es sich korrigieren oder dem Menschen
   antworten kann.
4. Die Antwort wird mit `AnswerStrict` gezeichnet.

## Vergleich mit auftragsbutler (`libs/assistent`)

Das Sicherheitsmodell ist dasselbe: nur lesen, eine Allowlist für Oberflächenaktionen, heikle Felder nie
ausfüllen, der Mensch speichert. Der Unterschied liegt in der Darstellung. Statt JSON, das im Browser mit
eigenem Markdown und eigenem SVG gezeichnet wird, kommt HTML vom Server, gestreamt mit Zwischenständen. Im
Browser läuft nur noch, was der Server nicht kann: die Aktionen auf der Seite.

| auftragsbutler | tabler-templ |
|---|---|
| `Katalog`, `Aktion`, `Pruefen` | `Catalog`, `ActionSpec`, `Check` |
| `Oberflaechenaktion` (JSON) | `Action`, gesendet mit `Stream.Action` |
| `HeikleFelder`, `HeikleWerte` | `DefaultSensitiveFields`, `DefaultSensitiveValues` |
| `data-formular`, `data-tour` | `data-form`, `data-tour` |
| `vorausfuellen` in `assistent.js` | `ActionScript` |
| Verlauf in `sessionStorage` und Datenbank | `Transcript` auf dem Server |
| Mini-Markdown in JS | `AnswerStrict` (goldmark, ohne Links und Bilder) |
| Antwort am Stück | Stream mit Werkzeugaufrufen, Blöcken und Aktionen |
