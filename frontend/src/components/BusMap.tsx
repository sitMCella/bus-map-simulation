import { useEffect, useState, useRef } from 'react';
import { MapContainer, TileLayer, Marker, Popup, CircleMarker } from 'react-leaflet';
import L, { type LatLngTuple } from 'leaflet';
import iconUrl from "leaflet/dist/images/marker-icon.png";
import iconRetinaUrl from "leaflet/dist/images/marker-icon-2x.png";
import shadowUrl from "leaflet/dist/images/marker-shadow.png";

function BusMap() {
  interface BusPosition {
    id: number;
    creationtime: string;
    busId: string;
    latitude: number;
    longitude: number;
    nextBusStopId: string;
    isBusStop: boolean;
  }

  interface BusPositionDisplay extends BusPosition {
    position: LatLngTuple;
    time: string;
  }

  interface BusStop {
    id: string;
    name: string;
    latitude: number;
    longitude: number;
    position: LatLngTuple;
  }

  interface Bus {
    id: string;
    latitude: number;
    longitude: number;
  }

  interface BusTimeTable {
    bus_id: string;
    bus_stop_id: string;
    time_seconds: number;
    timestamp: string;
  }

  interface BusTimeTableDisplay extends BusTimeTable {
    time: string;
  }

  interface BusStopDisplay extends BusStop {
    position: LatLngTuple;
    busTimeTables: BusTimeTableDisplay[];
  }

  interface BusTimeDelay {
    bus_id: string;
    bus_stop_id: string;
    delay_seconds: number;
  }

  function useStateRef(value: any) {
    const ref = useRef(value);

    useEffect(() => {
      ref.current = value;
    }, [value]);

    return [ref];
  }

  const [busPositions, setBusPositions] = useState<BusPositionDisplay[]>([]);

  const [centerMao, _setCenterMap] = useState<LatLngTuple>([41.9096, 12.52975]);

  const [busStopEntries, setBusStopEntries] = useState<BusStopDisplay[]>([]);

  const [_, setBusEntries] = useState<Bus[]>([]);

  const [busTimeTableEntries, setBusTimeTableEntries] = useState<BusTimeTableDisplay[]>([]);

  const [busTimeTableEntriesRef] = useStateRef(busTimeTableEntries);

  const [busTimeDelays, setBusTimeDelays] = useState<BusTimeDelay[]>([]);

  const [busTimeDelaysRef] = useStateRef(busTimeDelays);

  const getBusStopEntries = async (signal: AbortSignal) => {
    const headers = {
      Accepted: 'application/json',
    };
    // Dev environment prefix with "http://localhost:9090"
    const response = await fetch('/hub/bus_stop', {
      method: 'GET',
      headers,
      signal,
    });
    if (!response.ok) {
      throw new Error(JSON.stringify(response));
    }
    const responseData: BusStop[] = await response.json();
    const busStopsWithPosition: BusStopDisplay[] = responseData.map((b) => ({
      ...b,
      position: [b.latitude, b.longitude] as LatLngTuple,
      busTimeTables: [],
    }));
    return busStopsWithPosition;
  };

  const getBusEntries = async (signal: AbortSignal) => {
    const headers = {
      Accepted: 'application/json',
    };
    // Dev environment prefix with "http://localhost:9090"
    const response = await fetch('/hub/bus', {
      method: 'GET',
      headers,
      signal,
    });
    if (!response.ok) {
      throw new Error(JSON.stringify(response));
    }
    const responseData: Bus[] = await response.json();
    setBusEntries(responseData);
    return responseData;
  };

  const getBusTimeTableEntries = async (busId: string, signal: AbortSignal) => {
    const headers = {
      Accepted: 'application/json',
    };
    // Dev environment prefix with "http://localhost:9090"
    const response = await fetch('/hub/bus/' + busId + '/time_table', {
      method: 'GET',
      headers,
      signal,
    });
    if (!response.ok) {
      throw new Error(JSON.stringify(response));
    }
    const responseData: BusTimeTable[] = await response.json();
    const busTimeTables: BusTimeTableDisplay[] = responseData.map((r: BusTimeTable) => ({
      ...r,
      time: new Date(r.timestamp).toLocaleTimeString(),
    }));
    setBusTimeTableEntries(busTimeTables);
    return busTimeTables;
  };

  const diffInSeconds = (a: Date, b: Date): number => {
    return Math.floor((b.getTime() - a.getTime()) / 1000);
  };

  const calculateBusTimeDelay = (busPositionDisplay: BusPositionDisplay) => {
    const currentBusStopTimeTableEntry = busTimeTableEntriesRef.current.filter(
      (btt: BusTimeTableDisplay) =>
        btt.bus_id === busPositionDisplay.busId &&
        btt.bus_stop_id === busPositionDisplay.nextBusStopId,
    );
    var delay = 0;
    if (currentBusStopTimeTableEntry.length > 0) {
      delay = diffInSeconds(
        new Date(currentBusStopTimeTableEntry[0].timestamp),
        new Date(busPositionDisplay.creationtime),
      );
    }
    const cleanedBusTimeDelay = busTimeDelaysRef.current.filter(
      (item: BusTimeDelay) => !(item.bus_id === busPositionDisplay.busId),
    );
    var newBusTimeDelays: BusTimeDelay[] = [];
    busTimeTableEntriesRef.current
      .filter((btt: BusTimeTableDisplay) => btt.bus_id === busPositionDisplay.busId)
      .forEach((btt: BusTimeTableDisplay) => {
        const busTimeDelay: BusTimeDelay = {
          bus_id: btt.bus_id,
          bus_stop_id: btt.bus_stop_id,
          delay_seconds: delay,
        };
        newBusTimeDelays = [...newBusTimeDelays, busTimeDelay];
      });
    setBusTimeDelays([...cleanedBusTimeDelay, ...newBusTimeDelays]);
  };

  const getDelay = (busTimeTable: BusTimeTable): number => {
    const busStopTimeTable = busTimeDelaysRef.current.filter(
      (bttd: BusTimeDelay) =>
        bttd.bus_stop_id === busTimeTable.bus_stop_id && bttd.bus_id === busTimeTable.bus_id,
    );
    if (busStopTimeTable.length > 0) {
      return busStopTimeTable[0].delay_seconds;
    }
    return 0;
  };

  const getBuStopName = (busPosition: BusPositionDisplay): string => {
    const busStopEntry = busStopEntries.filter((be: BusStopDisplay) => (busPosition.nextBusStopId === be.id))
    if(busStopEntry.length > 0) {
      return busStopEntry[0].name;
    }
    return busPosition.nextBusStopId;
  }

  useEffect(() => {
    L.Icon.Default.mergeOptions({
      iconRetinaUrl,
      iconUrl,
      shadowUrl,
    });

    const controller = new AbortController();
    const signal = controller.signal;

    getBusStopEntries(signal)
      .then((busStopEntries: BusStopDisplay[]) => {
        getBusEntries(signal).then((busEntries: Bus[]) => {
          busEntries.forEach((bus: Bus) => {
            getBusTimeTableEntries(bus.id, signal).then((btte: BusTimeTableDisplay[]) => {
              const updatedStopEntries = busStopEntries.map((bs: BusStopDisplay) => ({
                ...bs,
                busTimeTables: [
                  ...bs.busTimeTables,
                  ...btte.filter((btt: BusTimeTableDisplay) => btt.bus_stop_id === bs.id),
                ],
              }));
              setBusStopEntries(updatedStopEntries);
            });
          });
        });
      })
      .catch((err) => {
        console.log('Cannot retrieve the Bus entries ' + err.message);
      });

    // Dev environment: prefix with "http://localhost:8080"
    const busPositionEventSource = new EventSource('/api/bus/position');
    busPositionEventSource.onmessage = (event) => {
      if (event.data) {
        const busPosition: BusPosition = JSON.parse(event.data);
        const busPositionDisplay: BusPositionDisplay = {
          ...busPosition,
          time: new Date(busPosition.creationtime).toLocaleTimeString(),
          position: [busPosition.latitude, busPosition.longitude] as LatLngTuple,
        };
        if (!busPositions.some((item: BusPositionDisplay) => busPosition.busId === item.busId)) {
          setBusPositions([...busPositions, busPositionDisplay]);
        } else {
          const cleanedBusPositions = busPositions.filter(
            (item: BusPositionDisplay) => busPosition.busId !== item.busId,
          );
          setBusPositions([...cleanedBusPositions, busPositionDisplay]);
        }
        if (busPositionDisplay.isBusStop) {
          calculateBusTimeDelay(busPositionDisplay);
        }
      }
    };
    return () => controller.abort();
  }, []);

  return (
    <MapContainer
      center={centerMao}
      zoom={13}
      scrollWheelZoom={true}
      style={{ minHeight: '100vh', minWidth: '100vw' }}
    >
      <TileLayer
        attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {busStopEntries.map((b) => (
        <Marker key={b.id} position={b.position}>
          <Popup>
            <div className="popup-card">
              <div className="popup-title">
                <b>Bus Stop:</b> {b.name}
              </div>
              <div className="popup-row">
                <span className="label">Latitude:</span>
                <span className="value">{b.latitude}</span>
              </div>
              <div className="popup-row">
                <span className="label">Longitude:</span>
                <span className="value"> {b.longitude}</span>
              </div>
              <div className="popup-row">
                <table className="popup-table">
                  <thead>
                    <tr>
                      <th>Bus Line</th>
                      <th>Time</th>
                      <th>Delay (sec)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {b.busTimeTables.map((t) => (
                      <tr key={t.bus_id}>
                        <td>{t.bus_id}</td>
                        <td>{t.time}</td>
                        <td>{getDelay(t)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </Popup>
        </Marker>
      ))}
      {busPositions.map((b) => (
        <CircleMarker key={b.busId} center={b.position} radius={20}>
          <Popup>
            <div className="popup-card">
              <div className="popup-title">
                <b>Bus Line:</b> {b.busId}
                <br />
              </div>
              <div className="popup-row">
                <span className="label">Time:</span>
                <span className="value">{b.time}</span>
              </div>
              <div className="popup-row">
                <span className="label">Latitude:</span>
                <span className="value">{b.latitude}</span>
              </div>
              <div className="popup-row">
                <span className="label">Longitude:</span>
                <span className="value">{b.longitude}</span>
              </div>
              <div className="popup-row">
                <span className="label">Next Stop:</span>
                <span className="value">{getBuStopName(b)}</span>
              </div>
              <div className="popup-row">
                <span className="label">Status:</span>
                {b.isBusStop ? <span className="value">Onboarding</span> : <span className="value">Moving</span>}
              </div>
            </div>
          </Popup>
        </CircleMarker>
      ))}
    </MapContainer>
  );
}

export default BusMap;
