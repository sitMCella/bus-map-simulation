package de.mcella.spring.dispatch;

import io.r2dbc.postgresql.PostgresqlConnectionConfiguration;
import io.r2dbc.postgresql.PostgresqlConnectionFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class PostgresqlConnectionFactoryConfiguration {

  @Value("${DATABASE_HOST}")
  private String databaseHost;

  @Value("${DATABASE_NAME}")
  private String databaseName;

  @Value("${SPRING_R2DBC_USERNAME}")
  private String databaseUsername;

  @Value("${SPRING_R2DBC_PASSWORD}")
  private String databasePassword;

  @Bean
  public PostgresqlConnectionFactory postgresqlConnectionFactory() {
    return new PostgresqlConnectionFactory(
        PostgresqlConnectionConfiguration.builder()
            .host(databaseHost)
            .database(databaseName)
            .username(databaseUsername)
            .password(databasePassword)
            .build());
  }
}
