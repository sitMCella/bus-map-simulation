FROM gradle:8.13.0-jdk21-alpine AS builder
COPY ./dispatch /code/
WORKDIR /code
RUN ./gradlew clean build

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /code/build/libs/dispatch-0.0.1-SNAPSHOT.jar ./dispatch.jar
CMD ["java", "-jar", "dispatch.jar"]