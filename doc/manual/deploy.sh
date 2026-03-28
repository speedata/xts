#!/bin/bash

set -o nounset
set -o errexit

export HUGO_ENV=production
export HUGO=/usr/local/bin/hugo147_9

FONT_DIR=/usr/local/fonts/xtsdoc

a=$(mktemp -d)
cd $a

{
    echo $a
    git clone ~patrick/git/xts.git
    cd xts/doc/manual

    # -- Proprietary content injection -----------------------

    # 1. Copy font files
    cp "$FONT_DIR"/*.woff2 static/fonts/

    # 2. Append @font-face rules and variable overrides to custom.css
    cat >> assets/css/custom.css << 'CSSEOF'

/* ibm-plex-mono-regular - latin */
@font-face {
  font-display: swap;
  font-family: 'IBM Plex Mono';
  font-style: normal;
  font-weight: 400;
  size-adjust: 105%;
  src: url('../fonts/ibm-plex-mono-v15-latin-regular.woff2') format('woff2');
}
/* ibm-plex-mono-italic - latin */
@font-face {
  font-display: swap;
  font-family: 'IBM Plex Mono';
  font-style: italic;
  font-weight: 400;
  src: url('../fonts/ibm-plex-mono-v15-latin-italic.woff2') format('woff2');
}
/* ibm-plex-mono-600 - latin */
@font-face {
  font-display: swap;
  font-family: 'IBM Plex Mono';
  font-style: normal;
  font-weight: 600;
  src: url('../fonts/ibm-plex-mono-v15-latin-600.woff2') format('woff2');
}
/* ibm-plex-mono-600italic - latin */
@font-face {
  font-display: swap;
  font-family: 'IBM Plex Mono';
  font-style: italic;
  font-weight: 600;
  src: url('../fonts/ibm-plex-mono-v15-latin-600italic.woff2') format('woff2');
}
@font-face {
  font-family: 'SanaSans';
  font-weight: 600;
  src: url('../fonts/sansbold.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSans';
  font-weight: 600;
  font-style: italic;
  src: url('../fonts/sansbolditalic.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSansAlt';
  font-weight: 600;
  src: url('../fonts/altbold.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSansAlt';
  font-weight: 600;
  font-style: italic;
  src: url('../fonts/altbolditalic.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSansAlt';
  font-style: italic;
  src: url('../fonts/altitalic.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSans';
  font-style: italic;
  src: url('../fonts/sansitalic.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSans';
  src: url('../fonts/sansregular.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSansAlt';
  src: url('../fonts/altregular.woff2') format('woff2');
  font-display: swap;
}
@font-face {
  font-family: 'SanaSansAlt';
  font-weight: 500;
  src: url('../fonts/altmedium.woff2') format('woff2');
  font-display: swap;
}

/* Override font variables with actual fonts */
:root {
  --md-text-font: "SanaSansAlt", system-ui, sans-serif;
  --md-code-font: "IBM Plex Mono", ui-monospace, monospace;
}
CSSEOF

    # 3. Inject Matomo tracking
    cat > layouts/partials/custom/head-end.html << 'MATOMOEOF'
{{- if hugo.IsProduction }}
<!-- Matomo -->
<script>
  var _paq = window._paq = window._paq || [];
  _paq.push(['trackPageView']);
  _paq.push(['enableLinkTracking']);
  (function() {
    var u="https://piwik.speedata.de/";
    _paq.push(['setTrackerUrl', u+'matomo.php']);
    _paq.push(['setSiteId', '8']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
    g.async=true; g.src=u+'matomo.js'; s.parentNode.insertBefore(g,s);
  })();
</script>
<!-- End Matomo Code -->
{{end -}}
MATOMOEOF

    # -- Build and deploy ------------------------------------
    $HUGO --baseURL "https://doc.speedata.de/xts/"
    rsync -avx --delete public/ /var/www/speedata.de/doc/xts/

} > /tmp/mkxts-docs.txt

cd ~/
