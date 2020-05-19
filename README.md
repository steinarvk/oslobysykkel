# oslobysykkel

En app skrevet av Steinar Kaldager som en demonstrasjon
av APIet for Oslo Bysykkel
( https://oslobysykkel.no/apne-data/sanntid ).

## Ferdigversjon 

En versjon av web-appen er hostet på følgende adresse:
https://oslobysykkel.app.kaldager.com

## Bruksanvisning

Appen skal være brukbar både på desktop og på mobil/tablet.

I øvre halvdel av appen finnes det et kart; i nedre
halvdel er det en tabell.

I tabellen finner du antall ledige sykler og antall
ledige plasser for hver rad.

Hver rad har også en lenke til Google Maps og til
Google Street View dersom en ønsker å orientere seg mer.

Kartet har ikoner med forskjellige farger avhengig av om
det finnes sykler og/eller plasser tilgjengelig.

Grønne markører betyr at begge deler finnes. Beige markører
med minus-symbol betyr at det er tomt for sykler. Oransje
markører med pluss-symbol betyr at det er fullt (eller
"tomt" for sykkelplasser). Røde markører betyr at det er
tomt for begge deler, altså sannsynligvis at noe er
i ustand.

En kan klikke på enten en rad i tabellen eller på en markør
på kartet for å fokusere på et nærområde. Kartet zoomer da
inn mot markøren, og radene i tabellen sorteres etter
distanse fra den.

Tabellen kan ellers sorteres etter ønsket kolonne ved å
klikke på kolonneoverskriften. Klikk to ganger for å
reversere sorteringen.

## Kildekode

git@git.kaldager.com:oslobysykkel.git
Tilgang til dette repoet er foreløpig etter avtale.
Send en RSA-pubkey for å få tilgang.

Alternativt er et arkiv av kildekoden tilgjengelig for
nedlasting på følgende adresse:
https://kaldager.com/download/oslobysykkel-app.tar.gz

## Instruksjoner for å bygge prosjektet

```
  $ go build
  $ docker build
```

Deretter kjøres prosjektet som en vanlig docker-container
som lytter på porten $PORT. Om $PORT ikke er satt til en
verdi lytter serveren på port 8080.

## Om implementasjonen

Implementasjonen er i Go med litt Javascript på frontend.

Kartet bruker biblioteket Leaflet og er basert på data
fra OpenStreetMap.

Appen fungerer som et lag med cache rundt det underliggende
APIet. Hver instanse av appen sender kun én forespørsel til
det underliggende APIet omtrent hvert tiende sekund, og
dette også kun hvis det faktisk finnes trafikk.

En app under tung trafikk belaster altså det underliggende
APIet minimalt, og sikrer så langt det er mulig god ytelse

Dersom det ikke finnes noe trafikk sender appen ingen
forespørsler til APIet etter oppstart. Dette minimerer
kostnadene ved å kjøre en instans som stort sett ikke får
trafikk.

## Biblioteker (skrevet av andre)

Se også vendor/-mappa.

github.com/umahmood/haversine: et bibliotek for å regne ut
avstander mellom geografiske koordinater.

leaflet: et Javascript-bibliotek for visualisering av kart.
https://leafletjs.com/

leaflet.awesome-markers: en utvidelse til leaflet som gir
flere muligheter for markørene.
https://github.com/lvoogdt/Leaflet.awesome-markers

jquery: generelle verktøy for browser-Javascript.
https://jquery.com/

ionicons: et ikon-bibliotek, brukt for markørene på kartet.
https://ionicons.com/

reset: en klassisk CSS-reset-stylesheet.
https://meyerweb.com/eric/tools/css/reset/

sorttable: et Javascript-bibliotek for automatisk sortering
av tabeller.
https://kryogenix.org/code/browser/sorttable/
