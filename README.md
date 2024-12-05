# Event-Driven Messaging System with RabbitMQ and Pipes-and-Filters Architecture

This project implements an event-driven system using RabbitMQ as the message broker, consisting of four separately deployable services: a user-facing REST API server that receives user messages, a filter service that removes messages containing stop-words, a screaming service that converts messages to uppercase, and a publish service that sends an email with the processed message. Additionally, the project includes an alternative pipes-and-filters architecture where services are directly connected through pipes/queues without the message broker. 

## Performance Report: Event-Driven System vs. Pipes-and-Filters Architecture
This report compares the performance of an Event-Driven System (EDS) using RabbitMQ as the message broker and a Pipes-and-Filters architecture where services are connected directly through pipes/queues. The comparison is based on metrics such as time behavior, resource utilization, and capacity.
### 1. Time Behavior
#### Event-Driven System (EDS):

    Mean Latency: 23.5 ms
    Effective Requests per Second (RPS): 831
    Percentile Latencies:
        50%: 22 ms
        90%: 37 ms
        95%: 43 ms
        99%: 60 ms
        100% (longest request): 111 ms

#### Pipes-and-Filters Architecture:

    Mean Latency: 8650.6 ms
    Effective Requests per Second (RPS): 1
    Percentile Latencies:
        50%: 9498 ms
        90%: 15681 ms
        95%: 16405 ms
        99%: 19744 ms
        100% (longest request): 19744 ms

The EDS outperforms the Pipes-and-Filters architecture in terms of response times, with 50% of requests processed in under 22 ms and 99% within 60 ms, demonstrating RabbitMQ's efficiency in managing high-throughput asynchronous tasks. In contrast, the Pipes-and-Filters system experiences significantly higher latencies, with a mean latency of over 8 seconds and the 99th percentile reaching nearly 20 seconds, indicating difficulties in handling concurrent, asynchronous tasks when using direct pipes between services.

### 2. Resource Utilization
##### CPU Usage:

    EDS: 0.6%
    Pipes-and-Filters: 0.4%

The slightly lower CPU usage in the Pipes-and-Filters architecture (0.4%) compared to the Event-Driven System (0.6%) suggests that, despite its higher latency, the Pipes-and-Filters system may be more lightweight in terms of resource consumption. This could be due to its synchronous processing model, which avoids the overhead of message passing and queuing involved in the RabbitMQ-based EDS, but at the cost of performance and scalability.

##### Memory Usage:

    EDS: 45 MB
    Pipes-and-Filters: 18 MB

The EDS consumes more memory (45 MB) than the Pipes-and-Filters architecture (18 MB), suggesting that RabbitMQ or the additional services in the event-driven system require more memory. Despite the higher memory usage, the event-driven system performs better in terms of response time and throughput. This suggests that the benefits of using RabbitMQ in handling large numbers of asynchronous requests outweigh the additional memory overhead.

### 3. Capacity
#### EDS:

    Completed Requests: 16,632
    Total Time: 20.02 seconds
    Effective RPS: 831

#### Pipes-and-Filters:

    Completed Requests: 22
    Total Time: 20.016 seconds
    Effective RPS: 1

The Event-Driven System (EDS) efficiently processes 16,632 requests in 20 seconds, achieving an effective requests per second rate of 831, demonstrating its scalability and ability to handle high concurrency with low latency. In contrast, the Pipes-and-Filters architecture handles only 22 requests in the same period, with an effective RPS of just 1, highlighting its limited capacity to manage concurrent requests, likely due to its more synchronous and linear processing approach.

### How the Metrics Were Obtained

The performance metrics were collected using the ```loadtest``` tool, which was configured to send HTTP POST requests to the system's `/send` endpoint. The load test was executed with a target time of 20 seconds, 20 concurrent clients, and a JSON payload containing an alias and content. The command used for the test was:

```bash
loadtest -t 20 -T application/json -P '{"alias": "Test Alias", "content": "Test Content"}' http://localhost:8080/send
```

During the test, the **mean latency**, **effective requests per second (RPS)**, and **latency percentiles** were recorded. CPU and memory usage statistics were gathered using the **[pidusage](https://github.com/struCoder/pidusage)** project to monitor system resources during the test. These metrics were used to compare the performance of the Event-Driven System and the Pipes-and-Filters architecture under load.

# How to Run the System

## Event-Driven System (EDS)

### Prerequisites

Before running the project, make sure you have the following installed on your local machine:

- **Go** (v1.18+)
- **RabbitMQ** (running locally or remotely)
- **Golang dependencies** (`github.com/joho/godotenv`, `github.com/rabbitmq/amqp091-go`)
- **SMTP credentials** (for sending emails)

### 1. Clone the Repository

Start by cloning the project repository to your local machine:

```bash
git clone https://github.com/iucd2/A9_EDS.git
cd SoftArch-A9-Team12
cd EDS
```

### 2. Set Up `.env` File

The `.env` file contains environment variables necessary for the application to run. Specifically, you will need to provide your SMTP credentials in this file.

Create a `.env` file in the root of the project directory with the following content:

```env
EMAIL_MAIL=your_email@example.com
EMAIL_PASSWORD=your_email_password
```

> **Note**: These environment variables are required for authenticating with the SMTP server. If they are not set correctly, the email sending feature will not work.

### 3. Install Dependencies

Ensure that Go modules are enabled, and install the necessary dependencies:

```bash
go mod tidy
```

This will download and install all the required Go packages for the project, including the RabbitMQ client (`github.com/rabbitmq/amqp091-go`) and the environment variable loader (`github.com/joho/godotenv`).

### 4. Start RabbitMQ (if not running)

If you don't already have RabbitMQ running locally, you can start it using Docker:

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

This will run RabbitMQ with the default guest credentials (`guest` / `guest`). You can access the RabbitMQ management console at `http://localhost:15672`.

### 5. Run the Project

Once you have RabbitMQ running and the `.env` file configured, you can start the project (it`s services):

```bash
go run api/main.go
go run filter-service/main.go
go run publish-service/main.go
go run screaming-service/main.go
```

This will start the services, which will:

- Connect to RabbitMQ
- Listen for messages from the related queues
- Process the messages and send emails using the SMTP server

If everything is set up correctly, you should see logs indicating that the services is successfully connected to RabbitMQ.

### 6. Testing the Application

To test the email functionality:

1. Send a message via the API endpoint. Example: POST http://localhost:8080/send body: content/json: {
    "alias": "Test Alias",
    "content": "Test Content"
}
2. The application will process the message and send an email to the recipients specified in the `sendEmail` function.

You can modify the recipient list or the message contents in the code.

### Troubleshooting

- **Missing `.env` Variables**: If you receive an error about missing `EMAIL_MAIL` or `EMAIL_PASSWORD`, ensure that these variables are correctly defined in your `.env` file.
- **SMTP Authentication Issues**: If the SMTP authentication fails, double-check your credentials and ensure that the email server allows access for less secure apps.
- **RabbitMQ Connection Issues**: Make sure RabbitMQ is running locally.


## Pipes-and-Filters System

### Prerequisites

Before running the project, make sure you have the following installed on your local machine:

- **Go** (v1.18+)
- **SMTP credentials** (for sending emails)

### 1. Clone the Repository

Start by cloning the project repository to your local machine:

```bash
git clone https://github.com/iucd2/A9_EDS.git
cd SoftArch-A9-Team12
cd pipes_and_filters
```

### 2. Set Up `.env` File

The `.env` file contains environment variables necessary for the application to run. Specifically, you will need to provide your SMTP credentials in this file.

Create a `.env` file in the root of the project directory with the following content:

```env
EMAIL_MAIL=your_email@example.com
EMAIL_PASSWORD=your_email_password
```

> **Note**: These environment variables are required for authenticating with the SMTP server. If they are not set correctly, the email sending feature will not work.

### 3. Install Dependencies

Ensure that Go modules are enabled, and install the necessary dependencies:

```bash
go mod tidy
```

This will download and install all the required Go packages for the project.

### 4. Run the Project

Once the dependencies are installed and the `.env` file is configured, you can start the project by running:

```bash
go run api.go
```

This will start the API server, which will:

- Process incoming messages
- Apply filtering, transformation, and email publishing through the pipeline

### 5. Testing the Application

To test the email functionality:

1. Send a message via the API endpoint. Example: POST http://localhost:8080/send body: content/json: 
    ```json
    {
      "alias": "Test Alias",
      "content": "Test Content"
    }
    ```
2. The application will process the message through the filter, screaming, and publishing services, ultimately sending an email to the recipients specified in the `sendEmail` function.

You can modify the recipient list or the message contents in the code.

### Troubleshooting

- **Missing `.env` Variables**: If you receive an error about missing `EMAIL_MAIL` or `EMAIL_PASSWORD`, ensure that these variables are correctly defined in your `.env` file.
- **SMTP Authentication Issues**: If the SMTP authentication fails, double-check your credentials and ensure that the email server allows access for less secure apps.

# Link to demonstration:
**[Demo](https://youtu.be/K6GaoTu3ynQ?si=vEQ3nbYfBsbxYBg1)**
