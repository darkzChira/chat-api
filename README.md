# Chat API

This repository contains the backend service for a scalable and extensible chat application, designed to handle real-time messaging, user authentication, and chat history management. The application is built using Golang, ensuring high performance and reliability.

## Architecture

The architecture of this chat app is optimized for scalability, performance, and maintainability, with the following key components:

- **WebSocket Server:** Manages real-time, bidirectional communication between clients, including user online notifications.
- **REST API:** Provides endpoints for user authentication, registration, fetching chat history, and sending offline messages.
- **Database:** Stores user credentials and chat history.
- **Middleware:** Handles JWT authentication and error management.

## Technologies Used

- **Golang:** The primary programming language used for building the backend services.
- **Gorilla Mux:** A powerful HTTP router and dispatcher for building the REST API.
- **Gorilla WebSocket:** A Go library providing easy-to-use WebSocket support for real-time communication.
- **MongoDB:** A NoSQL database used for storing chat history and user information.
- **JWT (JSON Web Tokens):** Used for securing API endpoints with authentication.
- **Mongo Driver:** A Golang driver for interacting with the MongoDB database.

## Project Setup

Follow these steps to set up the project locally:

### Setup MongoDB

1. **Install MongoDB:**
    - Install MongoDB following the instructions for your operating system: [MongoDB Community Edition](https://www.mongodb.com/try/download/community).

2. **Connect to MongoDB:**
    - Use `mongosh` or MongoDB Compass to connect to your MongoDB instance.

3. **Create the Database and Collections:**
    - Using MongoDB Shell:
      ```bash
      use chatdb 
          
      db.createCollection("messages")  
      db.createCollection("users")    
      ```
    - Using MongoDB Compass:
        - Open MongoDB Compass and connect to your MongoDB instance.
        - Click on "Create Database."
        - Name your database (e.g., `chatdb`).
        - Create collections named `messages` and `users` under this database.

4. **Setup Indexes**

Indexes improve the performance of queries on large collections. For this chat application, create indexes on the `messages` collection:

```bash
db.messages.createIndex({ sender_id: 1, receiver_id: 1, timestamp: -1 }) 
 ```


### Project Setup
1. **Clone the repository**
```bash
git clone https://github.com/darkzChira/chat-api.git
cd chat-app
```

2. **Install dependencies*
   Ensure you have Go installed on your machine. Then, install the necessary Go packages
```bash
go mod tidy
```

3. **Configure environment variables**
   Create a .env file in the root directory and add the following variables
```bash
MONGO_URI = <MongoDB connection string>
DATABASE_NAME = <MongoDB database Name>
SERVER_ADDRESS = <Server Starting Port; Ex - :8080>
WS_ADDRESS = <Websocket Port; Ex - :8081>
JWT_SECRET = <Any string for secret>
```

4. **Run the application**
```bash
go run main.go
```

The server will start on the specified port, typically http://localhost:8080

## Deployment
The application has been deployed using the following services
* **MongoDB Atlas**: Managed cloud database service used to host MongoDB
* **DigitalOcean**: Cloud infrastructure provider used to deploy the backend service.


The production backend service is accessible at [https://goldfish-app-ewgy2.ondigitalocean.app/](https://goldfish-app-ewgy2.ondigitalocean.app/).


## Contact
For any inquiries or feedback, feel free to reach out:

* Email: [darakachiranjaya@gmail.com](darakachiranjaya@gmail.com)
* LinkedIn: [https://www.linkedin.com/in/daraka-chiranjaya/](https://www.linkedin.com/in/daraka-chiranjaya/)
