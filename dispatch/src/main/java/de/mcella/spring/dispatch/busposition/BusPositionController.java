package de.mcella.spring.dispatch.busposition;

import de.mcella.spring.dispatch.busposition.dto.BusPosition;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/api/bus/position")
public class BusPositionController {

  private final BusPositionService busPositionService;

  BusPositionController(BusPositionService busPositionService) {
    this.busPositionService = busPositionService;
  }

  @GetMapping(produces = MediaType.TEXT_EVENT_STREAM_VALUE)
  @ResponseStatus(HttpStatus.OK)
  @CrossOrigin(origins = {"http://localhost", "http://localhost:3000", "http://localhost:5173"})
  public Flux<BusPosition> streamBusPositions() {
    return this.busPositionService.streamBusPositions();
  }
}
