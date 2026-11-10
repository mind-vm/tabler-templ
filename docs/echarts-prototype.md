# Prototyp: ECharts statt ApexCharts

Stand: 2026-11-10 · Branch `proto/echarts`

> **Umgesetzt:** Das Paket `chart` zeichnet inzwischen selbst mit ECharts; ApexCharts ist entfernt (siehe CHANGELOG, „Unreleased“). Die Vergleichsseite `/charts-echarts` gibt es deshalb nicht mehr. Dieses Dokument hält fest, was der Prototyp ergab.

## Warum

ApexCharts steht seit v5 unter der Community-Lizenz. Ab **2 Mio. USD** Jahresumsatz ist eine kommerzielle Lizenz nötig, und dabei zählen verbundene Firmen mit. Zum Weitergeben braucht man nach einem Forenbeitrag (Sekundärquelle) sogar eine OEM-Lizenz. tabler-templ ist eine Bibliothek, die ApexCharts in jede App mitliefert. Darum betrifft die Lizenz alle Apps, die tabler-templ nutzen. Apache ECharts steht unter der Apache-2.0-Lizenz und wird von der Apache Software Foundation gepflegt.

## Was der Prototyp zeigt

- Das neue Paket `echart` nimmt **unverändert `chart.Props`** an. Die Go-API kann also bleiben, nur die Ausgabe ändert sich.
- Die Vergleichsseite `/charts-echarts` zeigt alle 41 Demo-Diagramme nebeneinander, links ApexCharts 7, rechts ECharts 6.0.0. Alle 41 werden gezeichnet, ohne Konsolenfehler.
- Etwa 900 Zeilen: `echart/option.go` (Umsetzung von Props in ECharts-Optionen) und `echart/echart.templ` (Start-Skript).

| Bereich | Ergebnis |
| --- | --- |
| Line, Area, Bar, gestapelt, horizontal, Stepline, gestrichelt | fast identisch |
| Pie, Donut | identisch |
| Polar Area | als Balken in Polarkoordinaten, mit Ringen; sieht aus wie bei Apex |
| Radial Bar | als Gauge-Ringe; Abstände und Strichstärke anders, Legende vorhanden |
| Radar, Scatter, Bubble, Box Plot, Candlestick | gleichwertig |
| Heatmap, Treemap | gleichwertig (Abstufungen werden in Go berechnet) |
| Timeline (rangeBar) | über eine `custom`-Serie mit eigenem Renderer im Skript |
| Funnel | ECharts zeichnet einen Trichter (Trapeze), Apex zentrierte Balken |
| Zwei y-Achsen, Annotationen (x/y) | gleichwertig; Achsentitel werden teilweise abgeschnitten |
| Brush, Zoom-Toolbar, synchronisierte Gruppen | funktioniert (Brush über einen Slider auf einer versteckten Achse) |
| Dunkelmodus / Theme-Wechsel | funktioniert: CSS-Variablen werden aufgelöst und die Optionen neu gesetzt |
| Legende ein/aus (`Toggle`) | über `legendToggleSelect` |

## Was fehlt oder anders ist

- **`RawProps` und `Props.Options`** geben rohe ApexCharts-Optionen weiter. Das kann man nicht übertragen. Das ist der eigentliche **Breaking Change**. Betroffen sind `blocks/metric/widgets.templ` (Gauge) und `chart.Radial`.
- **Größe:** `echarts.min.js` hat 1092 KB (355 KB gzip), `apexcharts.min.js` 876 KB (248 KB gzip). Mit einem eigenen Build (nur die genutzten Diagrammtypen, Canvas-Renderer) wäre ECharts kleiner. Das bedeutet aber einen eigenen Build-Schritt.
- **Canvas statt SVG:** Bei ECharts sind Farben keine CSS-Werte. Das Skript wandelt jede `var()`- und `color-mix()`-Farbe in `rgba()` um. Ein SVG-Renderer ist mit `renderer: 'svg'` möglich, aber nicht getestet.
- Nicht übernommen: `XFormatter` auf Zeitachsen, `YTooltip`, `HidePoints`. Die Ausrichtung der Achsentitel ist nicht exakt nachgebaut.
- Nicht getestet: Barrierefreiheit (`aria`-Option von ECharts), htmx-Austausch, Druck, Leistung bei vielen Diagrammen.

## Aufwand bis zum Austausch

1. `echart` zu `chart` machen; `RawProps` und `Options` auf ECharts-Optionen umstellen (Breaking Change, Hauptversion).
2. `blocks/metric` (Gauge), `blocks/crypto` und `blocks/assistant` anpassen, `chart.Radial` neu bauen.
3. Feinarbeit am Aussehen: Funnel als Balken, Ringe beim Radial Bar, Achsentitel.
4. Versionen in `versions/versions.go`, CDN und Vendoring, README, CHANGELOG, Tests.

Schätzung: 2 bis 4 Tage, vor allem für Schritt 3 und die Tests.
