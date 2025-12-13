# Reactive Backend

## Introduction

Spring Boot Reactive application, used to send to the frontend application the data dispatched by the Hub.

## Development

### Requirements

- OpenJDK 21

### Build Application

```sh
./gradlew clean build
```

### Run Application

```sh
DATABASE_HOST=localhost DATABASE_NAME=busmap SPRING_R2DBC_USERNAME=postgres SPRING_R2DBC_PASSWORD=mysecretpassword java -jar build/libs/dispatch-0.0.1-SNAPSHOT.jar
```

### Format Code

```sh
./gradlew :spotlessApply
```
