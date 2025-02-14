# Realtime AlgoTrading App

## Overview

This is a high-performance, low-latency algorithmic trading application built with Golang. Designed for high-frequency trading, it leverages an event-driven architecture and efficient data processing techniques to ensure real-time execution.

## Why Golang?

Golang was chosen over Python due to its superior performance, lower latency, and built-in concurrency support, making it an ideal choice for real-time algorithmic trading.

## Database Choice: PostgreSQL + TimescaleDB

PostgreSQL was selected over MongoDB due to the need for TimescaleDB, which is optimized for time-series data. This choice ensures efficient order book storage and retrieval while also enabling future capabilities such as portfolio management and report generation using SQL’s powerful querying capabilities.

## No ORM - Raw SQL Only

To minimize latency and eliminate unnecessary abstraction layers, the application avoids ORMs and instead interacts with the database using raw SQL for maximum performance.

## Database Migrations

Currently, no migration tools are used for database schema management. However, `goose` may be integrated in the future if needed.

## Architecture

- **Data Source:** The application retrieves USDT/BTC order book data from Binance via WebSocket.
- **Event-Driven Design:** Order book data is processed using an event-driven architecture facilitated by Redis Pub/Sub, ensuring modularity and scalability.
- **Sliding Window Technique:** Implements Simple Moving Average (SMA) calculations for 50 and 200 records with an **O(1) complexity** approach.
- **Signal Generation:** Detects trend changes to generate buy/sell signals, automatically closing existing orders and creating new ones as needed.
- **Loosely Coupled Monolithic Structure:** While monolithic, the application adheres to event-driven principles to ensure flexibility and scalability.

## Deployment

The entire project is fully **Dockerized** and can be deployed using:

```sh
docker-compose up --build
```

## Database Initialization

The database schema, including TimescaleDB tables, is created using `init.sql`. This script is executed automatically when running the project via Docker.

## Metrics & Monitoring

- **Prometheus Metrics:**
  - Order book latency
  - Order creation latency
  - Trade signal latency
  - CPU & memory usage
  - Errors & data loss counts
  
  Available at: `http://localhost:8080/metrics` and on Prometheus: `http://localhost:9090`

- **Health & Readiness Checks:**
  - `http://localhost:8080/healthz`
  - `http://localhost:8080/readiness`

- **Grafana Dashboard:**
  - Preloaded dashboards available at: `http://localhost:3000`
  - Default credentials: `admin / admin`
  - If metrics are not visible, edit and save a metric to refresh the dashboard.

## Performance Metrics

After extensive testing with high-frequency order book data from Binance:

- **Average Order Book Processing Latency:** ~5ms
- **Memory Usage:** 20-25 MB
- **CPU Usage:** 5-6%

## Future Enhancements

- Develop separate data ingestion and processing layers.
- Implement batch processing for order book data to further optimize CPU usage.
- Add support for multiple trading pairs.
- Integrate `goose` for database migrations.


## Scalability, Fault Tolerance, and Security

To ensure the application remains reliable and scalable:

- Retry Mechanism & Connection Handling: The application implements retry logic for handling WebSocket disconnections, Database and Redis failures, ensuring continuous data flow.
- Low Resource Consumption: With minimal CPU and memory usage, the app can efficiently scale without excessive hardware demands.
- Raw SQL with Security Measures: Database queries are executed using raw SQL while adhering to best practices to prevent SQL injection.
- Event-Driven Architecture: Redis Pub/Sub ensures that services remain loosely coupled, improving fault tolerance and scalability.
- Monitoring & Alerting: Prometheus collects critical performance metrics, allowing early detection of issues like errors, data loss, connection failures.
- Health & Readiness Endpoints for Kubernetes: Dedicated endpoints monitor WebSocket, Redis, and database on `/readiness` and `/healthz` for Kubernetes.