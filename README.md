# Automatic Scaling Load Balancer Documentation

## Project architecture
![Project architecture](https://github.com/asitsonawane/Building-an-Automatic-Scaling-Load-Balancer-with-Go/blob/main/Load-Balancer-with-Go.jpg)

## Table of Contents
1. [Introduction](#introduction)
2. [Load Balancer](#load-balancer)
3. [WebApp](#webapp)
4. [Configuration Manager](#configuration-manager)
5. [Statistics Collection](#statistics-collection)
6. [Multithreading](#multithreading)
7. [Conclusion](#conclusion)

---

## 1. Introduction <a name="introduction"></a>

This documentation provides a comprehensive overview of the Automatic Scaling Load Balancer system. The system is designed to distribute incoming requests to multiple WebApp workers, ensuring high availability and reliability.

## 2. Load Balancer <a name="load-balancer"></a>

### 2.1 Design

The Load Balancer implements round-robin functionality to evenly distribute incoming requests among the available WebApp workers. Additionally, it supports weighted round-robin as a bonus feature, allowing for finer control over request distribution based on worker capacity.

### 2.2 Automatic Scaling

The Load Balancer is capable of dynamically adjusting the number of WebApp workers based on system load and configuration parameters. This ensures optimal resource utilization and responsiveness during varying levels of traffic.

## 3. WebApp <a name="webapp"></a>

### 3.1 Endpoints

- **`/api/v1/hello`**: Responds with a JSON message after an average delay of `avg-delay` seconds.

- **`/worker/stats`**: Provides statistics for each worker, including successful requests, failed requests, total requests, and average request time.

### 3.2 Request Handling

- Some requests to `/api/v1/hello` may fail randomly, based on the defined failure percentage. This simulates real-world scenarios and allows for robustness testing.

### 3.3 Statistics Collection

The WebApp keeps track of request statistics for each worker and saves them in the specified stats directory. This information includes successful requests, failed requests, total requests, and average request time.

## 4. Configuration Manager <a name="configuration-manager"></a>

### 4.1 Input Handling

The Configuration Manager is responsible for reading the configuration file and extracting relevant parameters. It then uses this information to spawn the Load Balancer and the desired number of WebApp workers.

### 4.2 Configuration Parameters

- Number of workers
- Request pool size
- Stats directory
- Average delay
- Failure percentage

## 5. Statistics Collection <a name="statistics-collection"></a>

### 5.1 Thread Safety

To avoid race conditions, the system employs thread-safe mechanisms when writing statistics to the shared directory. This ensures accurate and reliable data collection across multiple workers.

## 6. Multithreading <a name="multithreading"></a>

The application server is built using a multithreading-capable programming language. This allows for parallel processing of incoming requests, significantly improving system performance under heavy loads.

### 6.1 Race Condition Handling

Special attention is given to handling race conditions, particularly when writing statistics. A serialized approach is implemented to ensure data integrity and prevent conflicts during simultaneous write operations.

## 7. Conclusion <a name="conclusion"></a>

The Automatic Scaling Load Balancer system provides a robust and scalable solution for distributing incoming requests across multiple WebApp workers. With features like automatic scaling, weighted round-robin, and thorough statistics collection, it ensures high availability and performance even under challenging conditions. The careful consideration of multithreading and race condition handling further contributes to the system's reliability.
