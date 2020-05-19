"use strict";

// Frontend Javascript for the oslobysykkel app.

console.log("oslobysykkel.app.kaldager.com loading");

function humanQuantity(n, sing, plural) {
  if (n === 1) {
    return '' + n + ' ' + sing;
  }
  return '' + n + ' ' + plural;
}

const INITIAL_POSITION = [59.913246, 10.739387];  // Stortinget
const INITIAL_ZOOM = 13;  // Overview of the city center. Higher level is further zoomed in.
const PAN_ZOOM = 16;  // Overview of a neighbourhood.

// Use the Leaflet library to create a map based on OpenStreetMap data.
const map = L.map('leaflet-map').setView(INITIAL_POSITION, INITIAL_ZOOM);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

// Differently-coloured markers to be used depending on whether there are free bikes/docks.
const markers = {
  bikes_docks: L.AwesomeMarkers.icon({icon: 'checkmark-circled', prefix: 'ion', markerColor: 'green'}),
  nobikes_nodocks: L.AwesomeMarkers.icon({icon: 'close-circled', prefix: 'ion', markerColor: 'red'}),
  bikes_nodocks: L.AwesomeMarkers.icon({icon: 'plus-circled', prefix: 'ion', markerColor: 'orange'}),
  nobikes_docks: L.AwesomeMarkers.icon({icon: 'minus-circled', prefix: 'ion', markerColor: 'beige'}),
};

const stationMap = {};
const markerMap = {};
let onNextMoveEnd = null;

// Certain actions are delayed to the end of a map panning action, to avoid choppy animations.
function onMoveEnd() {
  if (onNextMoveEnd) {
    onNextMoveEnd();
  }
  onNextMoveEnd = null;
}
map.on("moveend", onMoveEnd);

// focusOn, triggered if a user clicks a station in the table or on the map.
// The effect is panning towards the station on the map, activating its popup, and
// sorting the table by distance from the station (i.e. showing other stations nearby).
function focusOn(stationId) {
  const station = stationMap[stationId];
  if (!station) {
    return;
  }
  const pos = [station.info.lat, station.info.lon];
  map.flyTo(pos, PAN_ZOOM);
  markerMap[stationId].openPopup();

  // Capture rows in order to sort them later.
  const tb = $("tbody").toArray()[0];
  const sortableRows = [];
  $.each(tb.rows, (index, value) => { sortableRows.push(value); });

  // Blank out the rows for now.
  $("tbody").empty();

  const sortByDistance = () => {
    $.getJSON("/api/get-distances?origin=" + stationId, null, response => {
      // Sort rows by distance from origin.
      sortableRows.sort((a, b) => { return response[$(a).data("station-id")] - response[$(b).data("station-id")]; });

      $("tbody").empty();
      $.each(sortableRows, (index, value) => { tb.appendChild(value); });
      ensureTableRowsClickable();

      // Remove sorttable sorting indicators, since we've set a custom sort.
      $(".sorttable_sorted").removeClass("sorttable_sorted");
      $(".sorttable_sorted_reverse").removeClass("sorttable-sorted_reverse");
      $("#sorttable_sortfwdind").remove();
      $("#sorttable_sortrevind").remove();
    });
  };

  onNextMoveEnd = sortByDistance;
}

// Make an API request to get the data to set up the map.
$.getJSON("/api/get-all-stations", null, response => {
  $.each(response, (stationId, station) => {
    stationMap[stationId] = station;
    const pos = [station.info.lat, station.info.lon];
    let popupText = station.info.name;
    popupText += '<br/> ' + humanQuantity(station.status.num_bikes_available, 'ledig sykkel', 'ledige sykler');
    popupText += '<br/> ' + humanQuantity(station.status.num_docks_available, 'ledig plass', 'ledige plasser');

    const markerClass = (((station.status.num_bikes_available > 0) ? "bikes" : "nobikes") + "_" +
                         ((station.status.num_docks_available > 0) ? "docks" : "nodocks"));
    const marker = L.marker(pos, {icon: markers[markerClass]}).addTo(map);
    marker.bindPopup(popupText, {opacity: 0.75});
    marker.on("click", e => { focusOn(stationId); });
    marker.bindTooltip(station.info.name, {direction: "auto", opacity: 0.75});
    markerMap[stationId] = marker;
  });
});

function ensureTableRowsClickable() {
  // Make table rows clickable.
  // Unbind the handler first to make it safe to call this multiple times.
  $("tr.station").unbind("click").click(e => {
    const target = e.target;
    const stationId = $(target).closest("tr").data("station-id");
    focusOn(stationId);
  });
}

ensureTableRowsClickable();
