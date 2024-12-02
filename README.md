## Prerequisites

Before running the project, make sure you have the following installed on your local machine:

- **Go** (v1.18+)
- **RabbitMQ** (running locally or remotely)
- **Golang dependencies** (`github.com/joho/godotenv`, `github.com/rabbitmq/amqp091-go`)
- **SMTP credentials** (for sending emails)

### 1. Clone the Repository

Start by cloning the project repository to your local machine:

```bash
git clone https://github.com/iucd2/A9_EDS.git
cd A9_EDS
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
