FROM ubuntu
EXPOSE 8080
RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "ca-certificates"]
RUN ["update-ca-certificates"]
COPY vendor/leaflet/ /app/static/
COPY vendor/Leaflet.awesome-markers/ /app/static/
COPY vendor/reset/* /app/static/
COPY vendor/sorttable/* /app/static/
COPY vendor/jquery/* /app/static/
COPY vendor/ionicons/* /app/static/
COPY vendor/ionicons-fonts/ /app/static/fonts/
COPY bysykkel /app/bysykkel-server
COPY css/* /app/static/
COPY js/* /app/static/
WORKDIR /app
ENTRYPOINT ["/app/bysykkel-server"]

