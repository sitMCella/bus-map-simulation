package de.mcella.spring.dispatch.busposition;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import de.mcella.spring.dispatch.busposition.dto.BusPosition;
import io.r2dbc.postgresql.PostgresqlConnectionFactory;
import io.r2dbc.postgresql.api.PostgresqlConnection;
import io.r2dbc.postgresql.api.PostgresqlResult;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.data.r2dbc.repository.config.EnableR2dbcRepositories;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Service
@EnableR2dbcRepositories(considerNestedRepositories = true)
public class BusPositionService {

  private final PostgresqlConnection connection;

  private final ObjectMapper objectMapper;

  private final Logger logger;

  private final Set<Long> seenIds;

  BusPositionService(PostgresqlConnectionFactory connectionFactory, ObjectMapper objectMapper) {
    this.connection =
        Mono.from(connectionFactory.create()).cast(PostgresqlConnection.class).block();
    this.objectMapper = objectMapper;
    this.logger = LoggerFactory.getLogger(BusPositionService.class);
    this.seenIds = ConcurrentHashMap.newKeySet();
  }

  @PostConstruct
  private void postConstruct() {
    connection
        .createStatement("LISTEN bus_position_notification")
        .execute()
        .flatMap(PostgresqlResult::getRowsUpdated)
        .subscribe();
  }

  @PreDestroy
  private void preDestroy() {
    connection.close().subscribe();
  }

  @Scheduled(fixedRate = 1000)
  private void cleanSeenIds() {
    this.seenIds.clear();
  }

  public Flux<BusPosition> streamBusPositions() {
    return connection
        .getNotifications()
        .onBackpressureBuffer(1000, dropped -> logger.warn("Dropped notification: " + dropped))
        .map(
            notification -> {
              try {
                var busPosition =
                    objectMapper.readValue(notification.getParameter(), BusPosition.class);
                logger.warn(busPosition.toString());
                return busPosition;
              } catch (JsonProcessingException e) {
                logger.error("Cannot send position", e);
              }
              return null;
            })
        .filter(busPosition -> busPosition != null && seenIds.add(busPosition.id()));
  }
}
