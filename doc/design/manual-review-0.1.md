# Manual review for 0.1

One line per page of the manual, in reading order. Tick a page when you have read it against the current xts, and write what has to change in the notes column. Generated pages (Commands reference, changelog) are not listed, they come from commands.xml and changelog.xml.

Checks per page: every listing runs with the current xts or is marked as a fragment; command and attribute names are the current ones; claims about defaults and behaviour match the implementation; the page links to an example in xts-examples where one exists. `rake doccheck` covers the links.

Local preview: `cd doc/manual && hugo server`, then open http://localhost:1313/ next to this file.

Tracked in issue #33.

## Layout Guide

| Done | Page | URL | Notes |
|---|---|---|---|
| [x] | `manual/_index.md` (Layout Guide) | https://doc.speedata.de/xts/manual/ | |
| [x ] | `manual/getting-started/_index.md` (Getting Started) | https://doc.speedata.de/xts/manual/getting-started/ | |
| [x] | `manual/core-concepts/_index.md` (Core Concepts) | https://doc.speedata.de/xts/manual/core-concepts/ | |
| [x] | `manual/core-concepts/how-it-works.md` (How It Works) | https://doc.speedata.de/xts/manual/core-concepts/how-it-works/ | |
| [x] | `manual/core-concepts/grid.md` (The Grid) | https://doc.speedata.de/xts/manual/core-concepts/grid/ | |
| [x] | `manual/core-concepts/placing-objects.md` (Placing Objects) | https://doc.speedata.de/xts/manual/core-concepts/placing-objects/ | |
| [x] | `manual/core-concepts/positioning-areas.md` (Positioning Areas) | https://doc.speedata.de/xts/manual/core-concepts/positioning-areas/ | |
| [x] | `manual/text-and-styling/_index.md` (Text & Styling) | https://doc.speedata.de/xts/manual/text-and-styling/ | |
| [x] | `manual/text-and-styling/fonts.md` (Fonts) | https://doc.speedata.de/xts/manual/text-and-styling/fonts/ | Leading rausgeworfen (erledigt) |
| [x] | `manual/text-and-styling/text-formatting.md` (Text Formatting) | https://doc.speedata.de/xts/manual/text-and-styling/text-formatting/ | Das hier scheint nicht zu funktionieren: ```<Span class="highlight"> <Value>highlighted</Value>    </Span> ``` mit background-color: yellow; -> wenn Zeit ist mal prüfen |
| [x] | `manual/text-and-styling/css-html.md` (CSS and HTML) | https://doc.speedata.de/xts/manual/text-and-styling/css-html/ | Supported CSS properties ist sicher nicht mehr aktuell. Ist aber eine Doppelung mit dem Reference Kapitel, oder? |
| [x] | `manual/text-and-styling/lists.md` (Lists) | https://doc.speedata.de/xts/manual/text-and-styling/lists/ | |
| [x] | `manual/tables/_index.md` (Tables) | https://doc.speedata.de/xts/manual/tables/ | |
| [x] | `manual/tables/tables.md` (Working with Tables) | https://doc.speedata.de/xts/manual/tables/tables/ | |
| [x] | `manual/images-and-graphics/_index.md` (Images & Graphics) | https://doc.speedata.de/xts/manual/images-and-graphics/ | |
| [x] | `manual/images-and-graphics/images.md` (Images) | https://doc.speedata.de/xts/manual/images-and-graphics/images/ | Query image dimensions for dynamic layouts -> Kurze Erklärung |
| [x ] | `manual/images-and-graphics/boxes-and-shapes.md` (Boxes and Shapes) | https://doc.speedata.de/xts/manual/images-and-graphics/boxes-and-shapes/ | |
| [x] | `manual/images-and-graphics/barcodes.md` (Barcodes) | https://doc.speedata.de/xts/manual/images-and-graphics/barcodes/ | |
| [x] | `manual/page-layout/_index.md` (Page Layout) | https://doc.speedata.de/xts/manual/page-layout/ | |
| [x] | `manual/page-layout/master-pages.md` (Master Pages) | https://doc.speedata.de/xts/manual/page-layout/master-pages/ | Hier noch das Zusammenspiel zwischen CSS Page setup und XTS masterpage beschreiben (siehe page-hooks) |
| [x] | `manual/page-layout/page-hooks.md` (Page Hooks) | https://doc.speedata.de/xts/manual/page-layout/page-hooks/ | Ich glaube, das CSS-Beispiel wäre besser auf der Seite davor aufgehoben. Außerdem könnte ein Bild erscheinen, wo die CSS-Margins überhaupt liegen. |
| [x] | `manual/page-layout/multi-page.md` (Multi-Page Content) | https://doc.speedata.de/xts/manual/page-layout/multi-page/ | |
| [x] | `manual/advanced/_index.md` (Advanced Topics) | https://doc.speedata.de/xts/manual/advanced/ | |
| [x] | `manual/advanced/pdf-options.md` (PDF Options) | https://doc.speedata.de/xts/manual/advanced/pdf-options/ | |
| [x] | `manual/advanced/slates.md` (Slates) | https://doc.speedata.de/xts/manual/advanced/slates/ | |
| [x] | `manual/advanced/colors.md` (Colors) | https://doc.speedata.de/xts/manual/advanced/colors/ | Hier müsste noch ein Verweis auf die neue Möglichkeit, mit ```@-bag-color gold {    value: #FFC72C;}```eine Frage zu bestimmen |
| [x] | `manual/accessibility/_index.md` (Accessible PDF) | https://doc.speedata.de/xts/manual/accessibility/ | A complete minimal example: das erste pdfua="2" steht ohne Kommentar dort. Beschreiben! (kurz) |
| [x] | `manual/running-xts/_index.md` (Running XTS) | https://doc.speedata.de/xts/manual/running-xts/ | |
| [x] | `manual/running-xts/command-line.md` (Command Line) | https://doc.speedata.de/xts/manual/running-xts/command-line/ | |
| [x] | `manual/running-xts/configuration.md` (Configuration) | https://doc.speedata.de/xts/manual/running-xts/configuration/ | |
| [x] | `manual/running-xts/file-organization.md` (File Organization) | https://doc.speedata.de/xts/manual/running-xts/file-organization/ | CSS path resolution -> letzte Wort in dem Abschnitt (directory) müsste CSS file sein, oder?  Und: funktioniert das mit dem --extra-dir auch für CSS-Dateien (wie fonts oder background-imags) |
| [x] | `manual/running-xts/schema-validation.md` (Schema Validation) | https://doc.speedata.de/xts/manual/running-xts/schema-validation/ | |
| [x] | `manual/running-xts/quality-assurance.md` (Quality Assurance) | https://doc.speedata.de/xts/manual/running-xts/quality-assurance/ | |
| [x] | `manual/running-xts/versions.md` (Versions) | https://doc.speedata.de/xts/manual/running-xts/versions/ | |

## Programming

| Done | Page | URL | Notes |
|---|---|---|---|
| [x] | `programming/_index.md` (Programming) | https://doc.speedata.de/xts/programming/ | Blockquote sieht doof aus (border-left und dicke Anführungszeichen) |
| [x] | `programming/execution-model.md` (Execution model) | https://doc.speedata.de/xts/programming/execution-model/ | |
| [x] | `programming/data-and-actions.md` (Data vs. action) | https://doc.speedata.de/xts/programming/data-and-actions/ | Ist das erste Beispiel unter The "boundary is the absence of a type" korrekt? Ebenso das Beispiel in The bridge from data to a command, ist da nicht ein Columns zu viel in der Schachteln? |
| [ ] | `programming/values-and-types.md` (Values and types) | https://doc.speedata.de/xts/programming/values-and-types/ | |
| [ ] | `programming/variables.md` (Variables) | https://doc.speedata.de/xts/programming/variables/ | |
| [ ] | `programming/functions.md` (Functions) | https://doc.speedata.de/xts/programming/functions/ | |
| [ ] | `programming/templates.md` (Templates) | https://doc.speedata.de/xts/programming/templates/ | |
| [ ] | `programming/control-flow.md` (Control flow) | https://doc.speedata.de/xts/programming/control-flow/ | |
| [ ] | `programming/records-and-dispatch.md` (Records and dispatch) | https://doc.speedata.de/xts/programming/records-and-dispatch/ | |
| [ ] | `programming/xpath-basics.md` (XPath basics) | https://doc.speedata.de/xts/programming/xpath-basics/ | |
| [ ] | `programming/xpath.md` (XPath in XTS) | https://doc.speedata.de/xts/programming/xpath/ | |
| [ ] | `programming/maps-and-arrays.md` (Maps and arrays) | https://doc.speedata.de/xts/programming/maps-and-arrays/ | |
| [ ] | `programming/data-files.md` (Data files) | https://doc.speedata.de/xts/programming/data-files/ | |

## Reference (handwritten pages)

| Done | Page | URL | Notes |
|---|---|---|---|
| [ ] | `reference/_index.md` (Reference Manual) | https://doc.speedata.de/xts/reference/ | |
| [ ] | `reference/defaults.md` (Defaults) | https://doc.speedata.de/xts/reference/defaults/ | Ist CSS defaults noch aktuell? -> Prüfen |
| [ ] | `reference/cli.md` (CLI Reference) | https://doc.speedata.de/xts/reference/cli/ | |
| [ ] | `reference/units.md` (Units) | https://doc.speedata.de/xts/reference/units/ | |
| [ ] | `reference/xpath-functions.md` (XPath Functions) | https://doc.speedata.de/xts/reference/xpath-functions/ | |
| [ ] | `reference/css-properties.md` (CSS Properties) | https://doc.speedata.de/xts/reference/css-properties/ | Diese Seite müsste man mit der Fähigkeit von htmlbag synchron halten |
