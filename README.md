# Bus Map Simulation Application

<img src="https://github.com/sitMCella/bus-map-simulation/wiki/images/bus_simulation.png">

## Table of contents

* [Introduction](#introduction)
* [Requirements](#requirements)
* [Run Application](#run-application)

## Introduction

The Bus Map Simulation application consists of a web application displaying a Leaflet map showing bus stops and the bus line positions in near real time.

A Hub application provides the bus configurations to the frontend application, and receives the bus positions from a Bus application.

A dispatch application is used to send the bus positions as a data flow to the frontend application.

Leaflet: https://leafletjs.com/
OpenStreetMap: https://www.openstreetmap.org/

## Requirements

- Docker (Docker compose)

## Run Application

```sh
docker-compose -f docker-compose.yml up
```

Open http://localhost:80 to view the application in your browser.
