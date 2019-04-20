GoTrain
=======

[![Build Status](https://travis-ci.org/rijdendetreinen/gotrain.svg?branch=master)](https://travis-ci.org/rijdendetreinen/gotrain)
[![codecov](https://codecov.io/gh/rijdendetreinen/gotrain/branch/master/graph/badge.svg)](https://codecov.io/gh/rijdendetreinen/gotrain)

GoTrain is a server application for receiving, processing and distributing
real-time data about train services in the Netherlands.

GoTrain is designed to continuously receive data streams offered as open data
by the Dutch Railways ([NS](http://www.ns.nl/)). The data is processed and
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

It is easy to extend GoTrain's functionality or to build your own applications
using the REST API. For example; create your live departures board.

Status
------

GoTrain is currently **under development** and is not yet used for production
workloads.

License
-------

Copyright (c) 2019 Geert Wirken, Rijden de Treinen

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
