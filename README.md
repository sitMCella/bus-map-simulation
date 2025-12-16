# Bus Map Simulation

<img src="https://github.com/sitMCella/bus-map-simulation/wiki/images/bus_simulation.png">

## Table of contents

* [Introduction](#introduction)
* [Requirements](#requirements)
* [Run Application](#run-application)

## Introduction

The Bus Map Simulation application consists of a web application that displays a Leaflet map showing bus stops and bus line positions in near real time.

A Hub application provides bus configuration data to the frontend application and receives bus position data from a Bus application.

A Dispatch application is used to send the bus position data as a data stream to the frontend application.

Leaflet: https://leafletjs.com/
OpenStreetMap: https://www.openstreetmap.org/

## Requirements

- Docker (Docker compose)

## Run Application

```sh
docker-compose -f docker-compose.yml up
```

Open http://localhost:80 to view the application in your browser.
