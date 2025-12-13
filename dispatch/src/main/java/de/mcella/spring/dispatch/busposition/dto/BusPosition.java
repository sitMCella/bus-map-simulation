package de.mcella.spring.dispatch.busposition.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;
import java.time.Instant;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonInclude(JsonInclude.Include.NON_NULL)
public record BusPosition(
    Long id,
    Instant creationtime,
    String busId,
    Double latitude,
    Double longitude,
    String nextBusStopId,
    Boolean isBusStop) {
  @Override
  public String toString() {
    return "BusPosition{"
        + "id='"
        + id
        + '\''
        + ", creationtime="
        + creationtime
        + ", busId="
        + busId
        + ", latitude="
        + latitude
        + ", longitude="
        + longitude
        + ", nextBusStopId="
        + nextBusStopId
        + ", isBusStop="
        + isBusStop
        + '}';
  }
}
