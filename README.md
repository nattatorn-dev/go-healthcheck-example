# Example Go Healthcheck

A robust and flexible health check package for Go services, designed to follow clean architecture and SOLID principles. This package allows you to perform health checks on various components of your system, such as API servers, databases, Redis, Kafka, and external APIs.

## Features

- Perform health checks for different components.
- Support for individual health check intervals and timeouts.
- Graceful error handling and logging.
- Immediate and periodic health checks.
- Easy integration with HTTP servers.


## Health Checks Overview

### Supported Health Checks

1. **API Server Health Check**:
   - **Description**: Verifies if an API server is reachable and responding.
   - **Usage**: Typically used to check the status of internal or external RESTful APIs your service depends on.
   - **Example Configuration**:
     ```go
     apiServerConfig := entities.CheckerConfig{
         Timeout:  5 * time.Second,
         Interval: 1 * time.Minute,
     }
     apiServerHealth := api.NewAPIServerHealth("http://localhost:8080", apiServerConfig.Timeout)
     ```

2. **Database Health Check**:
   - **Description**: Pings a database to ensure it's accessible.
   - **Usage**: Used to check the availability of databases like MySQL, PostgreSQL, etc.
   - **Example Configuration**:
     ```go
     dbConfig := entities.CheckerConfig{
         Timeout:  5 * time.Second,
         Interval: 1 * time.Minute,
     }
     dsn := "user:password@tcp(localhost:3306)/"
     dbHealth, err := db.NewSQLDBHealth("mysql", dsn, dbConfig)
     if err != nil {
         log.Fatalf("Failed to initialize database health check: %v", err)
     }
     ```

3. **Redis Health Check**:
   - **Description**: Checks if the Redis server is reachable and operational.
   - **Usage**: Used to ensure Redis, often used for caching or as a key-value store, is available.
   - **Example Configuration**:
     ```go
     redisConfig := entities.CheckerConfig{
         Timeout:  5 * time.Second,
         Interval: 30 * time.Second,
     }
     redisClient := redispkg.NewClient(&redispkg.Options{Addr: "localhost:6379"})
     redisHealth := redis.NewRedisHealth(redisClient, redisConfig.Timeout)
     ```

4. **Kafka Health Check**:
   - **Description**: Verifies if the Kafka broker is reachable and operational.
   - **Usage**: Used to check the availability of Kafka brokers for message queuing.
   - **Example Configuration**:
     ```go
     kafkaConfig := entities.CheckerConfig{
         Timeout:  5 * time.Second,
         Interval: 45 * time.Second,
     }
     kafkaHealth := kafka.NewKafkaHealth("localhost:9092", kafkaConfig.Timeout)
     ```

5. **External API Health Check**:
   - **Description**: Sends a request to an external API to ensure it's reachable and responding.
   - **Usage**: Used to verify the availability and responsiveness of third-party APIs your service depends on.
   - **Example Configuration**:
     ```go
     externalAPIConfig := entities.CheckerConfig{
         Timeout:  5 * time.Second,
         Interval: 2 * time.Minute,
     }
     externalAPI := external.NewExternalAPIHealth("https://api.example.com/health", externalAPIConfig.Timeout)
     ```
