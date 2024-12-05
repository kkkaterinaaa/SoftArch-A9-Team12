# Event-Driven Messaging System with RabbitMQ and Pipes-and-Filters Architecture

This project implements an event-driven system using RabbitMQ as the message broker, consisting of four separately deployable services: a user-facing REST API server that receives user messages, a filter service that removes messages containing stop-words, a screaming service that converts messages to uppercase, and a publish service that sends an email with the processed message. Additionally, the project includes an alternative pipes-and-filters architecture where services are directly connected through pipes/queues without the message broker. 

## Event-Driven System vs. Pipes-and-Filters Architecture
We compared the performance of an Event-Driven System (EDS) using RabbitMQ and a Pipes-and-Filters architecture. The comparison was based on time behavior, resource utilization, and capacity.
### 1. Time Behavior
#### Event-Driven system (EDS):

    Mean Latency: 23.5 ms
    Effective Requests per Second (RPS): 831
    Percentile Latencies:
        50%: 22 ms
        90%: 37 ms
        95%: 43 ms
        99%: 60 ms
        100% (longest request): 111 ms

#### Pipes-and-Filters:

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


##### Memory Usage:

    EDS: 45 MB
    Pipes-and-Filters: 18 MB

The Pipes-and-Filters architecture has a slight advantage in CPU usage, consuming only 0.4% compared to the 0.6% of the EDS. However, the EDS performs better in terms of throughput and latency despite its higher memory usage (45 MB vs. 18 MB for Pipes-and-Filters). This suggests that while the EDS uses more resources, it provides better scalability and performance for handling large numbers of requests.

### 3. Capacity
#### EDS:

    Completed Requests: 16,632
    Total Time: 20.02 seconds
    Effective RPS: 831

#### Pipes-and-Filters:

    Completed Requests: 22
    Total Time: 20.016 seconds
    Effective RPS: 1

The EDS demonstrates much greater capacity, completing 16,632 requests in 20 seconds with an effective request rate of 831 RPS. On the other hand, the Pipes-and-Filters architecture only processes 22 requests in the same period, highlighting its limited scalability and efficiency in handling concurrent tasks.

### 4. Conclusion

The Event-Driven System provides a more robust and scalable solution for high-throughput, asynchronous workloads, with lower latency and better performance in terms of resource usage and capacity. The Pipes-and-Filters architecture, while potentially more lightweight in terms of CPU and memory, struggles to match the EDS in scalability and responsiveness, making it less suitable for handling large volumes of concurrent requests. Therefore, for systems requiring high performance and concurrency, an Event-Driven System is a more efficient choice.


### Metrics collection sources

The performance metrics were collected using the ```loadtest``` tool, which was configured to send HTTP POST requests to the system's `/send` endpoint. The load test was executed with a target time of 20 seconds, 20 concurrent clients, and a JSON payload containing an alias and content. The command used for the test was:

```bash
loadtest -t 20 -T application/json -P '{"alias": "Test Alias", "content": "Test Content"}' http://localhost:8080/send
```

During the test, the **mean latency**, **effective requests per second (RPS)**, and **latency percentiles** were recorded. CPU and memory usage statistics, which were averaged over the test duration, were gathered using the **[pidusage](https://github.com/struCoder/pidusage)** project to monitor system resources during the test. These metrics were used to compare the performance of the Event-Driven System and the Pipes-and-Filters architecture under load.

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
git clone https://github.com/kkkaterinaaa/SoftArch-A9-Team12.git
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
git clone https://github.com/kkkaterinaaa/SoftArch-A9-Team12.git
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
