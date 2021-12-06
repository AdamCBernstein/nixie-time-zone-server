# nixie-time-zone-server
Go implementation of isparkes/time-zone-server

This is a Go server which gives back the local time anywhere in the world, given a Unix Style location as input, e.g. 'Europe/Zurich', or 'America/Los_Angeles'.

This is a 100% project, and has no external dependency on a web server. This program IS a very tiny web server

A full list of supported time zones is available at: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones

#API:

##GET /area/location ##GET /area/location/city

Gets the local time right now for the given TZ style input area and city n the format "yyyymmdd HHMMSS"

curl http://localhost:8081/Europe/Berlin
2021,12,06,03,31,44

curl http://localhost:8081/America/Los_Angeles
2021,12,05,18,31,44
