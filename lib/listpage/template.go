package listpage

import (
	"html/template"

	"kaldager.com/bysykkel/lib/oslobysykkel"
)

var (
	listpageTmpl = template.Must(template.New("").Parse(`
<html>

	<head>
		<title>oslobysykkel.app.kaldager.com</title>

		<meta name="description" content="Sanntidsoversikt over bysykler i Oslo.">
		<meta name="author" content="Steinar Kaldager">

		<script src="/static/jquery-3.5.1.min.js"></script>
		<script src="/static/sorttable.js"></script>
		<link rel="stylesheet" href="/static/reset.css">
		<link rel="stylesheet" href="/static/leaflet.css">
		<script src="/static/leaflet.js"></script>
		<link rel="stylesheet" href="/static/bysykkel.css">
		<link rel="stylesheet" href="/static/ionicons.min.css">
		<link rel="stylesheet" href="/static/leaflet.awesome-markers.css">
		<script src="/static/leaflet.awesome-markers.js"></script>
	</head>

	{{ if .Error }}
		Error: {{ .Error }}
	{{ else }}
		<div id="leaflet-map"></div>

		<div class="table-container">
		<table class="sortable">
			<thead>
				<tr>
					<td>Navn
					<td>Sykler 
					<td>Plasser
			</thead>
			
			<tbody>
				{{ range .Stations }}
					<tr class="station" data-station-id="{{ .Info.StationId }}" >
						<td>
						{{ with .Info }}
							{{ .Name }}
								<span class="tinylinks">
									<a href="https://www.google.com/maps/search/?api=1&query={{.Lat}},{{.Lon}}">kart</a>
									<a href="https://www.google.com/maps/@?api=1&map_action=pano&viewpoint={{.Lat}},{{.Lon}}">utsikt</a>
								</span>
								<p>
						    {{ if ne .Name .Address }}
									<span class="address">{{.Address}}</span>
								{{ end }}
								</p>

						{{ end }}

						<td class="avail-{{ if .Status.NumBikesAvailable }}ok{{else}}zero{{end}}"
						    sorttable_customkey="{{ .Status.NumBikesAvailable }}"
						>
							{{ .Status.NumBikesAvailable }}/{{ .Info.Capacity }}
						<td class="avail-{{ if .Status.NumDocksAvailable }}ok{{else}}zero{{end}}"
						    sorttable_customkey="{{ .Status.NumDocksAvailable }}"
						>
							{{ .Status.NumDocksAvailable }}/{{ .Info.Capacity }}
				{{ end }}
			</tbody>
		</table>
		</div>
	{{ end }}

	<script src="/static/bysykkel.js"></script>
</html>
`))
)

type pageParams struct {
	Error    error
	Stations []*oslobysykkel.Station
}
