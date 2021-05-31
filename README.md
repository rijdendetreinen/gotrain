GoTrain
=======

![Go](https://github.com/rijdendetreinen/gotrain/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/rijdendetreinen/gotrain/branch/master/graph/badge.svg)](https://codecov.io/gh/rijdendetreinen/gotrain)
[![Maintainability](https://api.codeclimate.com/v1/badges/e35ed750fa4facc99e92/maintainability)](https://codeclimate.com/github/rijdendetreinen/gotrain/maintainability)

GoTrain is a server application for receiving, processing and distributing
real-time data about train services in the Netherlands.

GoTrain is designed to continuously receive data streams offered as open data
by the [Dutch Railways (NS)](http://www.ns.nl/). The data is processed and
saved in-memory in order to offer a very fast REST API, which can be used by
numerous clients to show live information about train departures, arrivals etc.

GoTrain is currently able to process the following data streams:

* Arrivals (arriving service at a station)
* Departures (departing service from a station)
* Services (data about a complete trip for a single train)

And it offers the received data through a number of REST APIs, which allow you
to:

* Request a summary of all departing trains for a single station
* Detailed information for a departing train
* All upcoming arrivals for a single station
* Detailed information about a single train journey

You can also use GoTrain to store all services to a Redis queue for further
processing (archive function).

It is easy to extend GoTrain's functionality or to build your own applications
using the REST API. For example: create your live departures board, or analyze
which trains are currently delayed or cancelled on the Dutch rail network.

REST API
--------

GoTrain includes a convenient REST API. 

An example request: Request all departures from Utrecht Centraal: 
`/v2/departures/station/UT`

Response (shortened):

```json
{
    "departures": [
        {
            "cancelled": false,
            "company": "NS",
            "delay": 0,
            "departure_time": "2019-07-14T00:49:00+02:00",
            "destination_actual": "'t Harde",
            "destination_actual_codes": [
                "HDE"
            ],
            "destination_planned": "'t Harde",
            "name": null,
            "platform_actual": "12",
            "platform_changed": false,
            "platform_planned": "12",
            "remarks": [],
            "service_date": "2019-07-13",
            "service_id": "591",
            "service_number": "591",
            "station": "UT",
            "status": 2,
            "timestamp": "2019-07-13T22:49:00Z",
            "tips": [
                "Stopt ook in Harderwijk"
            ],
            "type": "Intercity",
            "type_code": "IC",
            "via": "Amersfoort",
            "wings": []
        }
   ],
    "status": "UP"
}
```

The following endpoints are available:

* `/` - API version
* `/v2/status` - System status
* `/v2/arrivals/stats` - Arrival statistics
* `/v2/arrivals/station/{station}` - Arrivals for `{station}` (e.g. `UT`)
* `/v2/arrivals/arrival/{id}/{station}/{date}` - Specific arrival details
* `/v2/departures/stats` - Departures statistics
* `/v2/departures/station/{station}` - Departures for `{station}` (e.g. `UT`)
* `/v2/departures/departure/{id}/{station}/{date}` - Specific departure details
* `/v2/services/stats` - Services statistics
* `/v2/services/service/{service_number}/{date}` - Specific service details

The full API documentation, including parameters and response formats, is included
in the [GoTrain OpenAPI specification](openapi.yaml). Or check out the nicely
formatted [GoTrain API on Apiary](https://rijdendetreinen.docs.apiary.io/).

Archiver
--------

The archive function allows you to store all data to an archive for further
processing or analysis at a later time. When running `gotrain archiver`,
GoTrain simply pushes all received services to a Redis queue (in JSON format).

Installation
------------

Binary packages will be provided in a future update. For now, the best way to install GoTrain is by
[downloading the source code](https://github.com/rijdendetreinen/gotrain/releases) and manually compile
the application.

1. Download source code
2. Install dependencies: `go get`
3. Compile: `go build`

Now, configure the application:

4. Go the the `config/` directory
5. Copy `example.yaml` to `config.yaml`
6. Modify the configuration parameters. Change at least the 'server' line (which should point to the NDOV ZeroMQ server - see below), and the API port (unless you want to run on port 8080).

Usage
-----

7. Start your server by running `./gotrain server`.
8. GoTrain keeps all information in memory, so you should aim to keep the server running for a long time.
   Process monitors like [supervisord](http://supervisord.org/) may help with that.
9. The REST API should be available on the address you have specified in your configuration file.

The CLI interface of GoTrain also allows you to:

* Request the current status of your server: 
  `./gotrain status -u http://localhost:8080/` 
  Initially, your server will have the status UNKNOWN and then RECOVERING, as it slowly starts to build up a complete dataset.
  Arrivals and departures should be UP after approximately 80 minutes.
* Inspect a single XML message, for example: 
  `./gotrain inspect departure parsers/testdata/departure.xml` 
* Run `./gotrain help` to show all commands.

Data feed
---------

You need access to a data feed from NS called InfoPlus. This data feed is distributed freely by the [NDOVloket](https://ndovloket.nl/).

Tip: you can find a public best-effort data feed from the NDOV on this page: [NDOVloket realtime](http://data.ndovloket.nl/REALTIME.TXT).
You need the address listed for 'NS InfoPlus', likely in the format `tcp://pubsub.besteffort.ndovloket.nl:...` 
Consider [signing up](https://ndovloket.nl/aanmelden/) for free if you are planning to use this application for production purposes.
The data feed is the same, but it's covered by a SLA, and it helps NDOVloket to keep track of their users.

Development roadmap
-------------------

The main objectives for GoTrain have now been developed, but there is a roadmap
for further development. The main planned improvements are:

* Increase test coverage - the API is currently not tested
* Better monitoring tools - analyze the data streams, monitor for errors, etc.
* Packaging - make it easier to install gotrain on a server by just downloading
  and installing the binaries

Status
------

GoTrain is currently being used in production by
[Rijden de Treinen](https://www.rijdendetreinen.nl/) as a source for realtime
departure times, arrivals and for trip details, both on the website and in the
mobile app. Please let it known when you use GoTrain for a cool project, a big
application or for some other purpose!

License
-------

Copyright (c) 2019-2021 Geert Wirken, Rijden de Treinen

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
