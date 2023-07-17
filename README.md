# About the project

 A comprehensive project that explores the various aspects of building web servers in Go. Through this endeavour, I aimed to gain a deeper insight into the different types of web servers, their functionalities, advantages, and disadvantages. In order to feel more real and practical the project simulates a simplified banking application and covers topics such as:
 -  **server architecture**, **authentication mechanisms**, **gRPC implementation**, **gRPC gateway**, **asynchronous task execution**, **database interactions**, **CI/CD**, **Kubernetes orchestration**

## Table of Contents

- [Features](#features)
- [Installation](#installation)

## Features
This is a list of topics that were the main focus of development, comparison and analysis:


- **HTTP Server Development with Gin Framework**: A fully-functional HTTP server using the popular Gin Framework. To demonstrate real-world applicability, this server simulates a simplified bank application, complete with user management, transaction handling, and bank account management.
    - code location -> **/api folder**
    

- **Flexible Authentication Layer**: A modular authentication layer that leverages middleware and tokens. By employing interfaces, it can switch between alternative authentication methods. Notably, this project showcases the implementation of two exemplary options: JSON Web Tokens (JWT) and Platform-Agnostic Security Tokens (PASETO). Also token IDs are persistently stored in the database, enabling efficient blacklisting if required.
    - code location -> **/auth folder**

- **gRPC Server Implementation**: A gRPC version of the Gin HTTP server to compare and understand the differences between the two approaches. Used protobuf files to define the server structure and interceptors for implementing the middleware functionality.
    - code location -> **/grpc + /proto folder**

- **gRPC Gateway and Swagger Documentation**: Exposed the gRPC server through HTTP using a gRPC gateway. Additionally, a Swagger Documentation UI is generated automatically when running the gateway server -> gateway_url/swagger/
    - code location -> **/grpc + /proto folder**

- **Asynchronous Task Execution**: Integrated asynchronous email sending using Redis Queue and the asynq library. Set up a separate server for processing queued tasks and utilize mailhog as development SMTP server.
    - code location -> **/async folder**

- **Database Interactions with Postgres**: Using Postgres as the database layer and leverage tools like sqlc for auto-generating Go code from SQL definitions. Handle database migrations using golang-migrate and implement database transactions for safe fund transfers and user creation.
    - code location -> **/db folder**

- **Docker Containers for Easy Deployment**: Simplify deployment by leveraging Docker containers. The project utilizes three containers to streamline the setup: one for the API, another for the Redis Workers responsible for asynchronous task processing, and a third for handling database migrations.

- **CI/CD with GitHub Actions**: Upon triggering a workflow (**push to main**), automated tests are executed to validate the system's integrity. If successful, Docker images are built and pushed to the GitHub Container Registry. After that another worflow could be added that deploys them to Cloud Run or any other cloud container service.

- **Kubernetes Orchestration**:  Utilizes the Docker images generated through the CI/CD process. It contains definitions for:
    - Deployments - for the API and Redis Workers
    - Stateful Set - for the DB + executes migrations / Redis
    - Ingress - for routing http to the Gin & Gateway server / For routing to the gRPC server
    

## Installation

To install and run the project locally, please follow these steps:

1. Clone the repository and install dependencies:

   ```bash
   git clone https://github.com/BogoCvetkov/go_masterclass
   ```

   ```bash
   go mod download
   ```

2. Run the necessary services (requires Docker Compose):

    ```bash
    docker-compose up redis db mailhog  --build 
    ```

3. Run db migrations (requires make to be installed):

    ```bash
    make db-migrate-up DB_URL=your-postgres-url
    ```

4. Start the API's - HTTP + gRPC + Gateway:

    ```bash
    make start-api
    ```

5. Start Redis workers - HTTP + gRPC + Gateway:

    ```bash
    make start-workers
    ```
6. Additional:
    
    generate GO code from sql definitions:
        
    ```bash
        make generate-models
    ```

    create a new db-migration:
    ```bash
        make db-new-migration MIGRATION_NAME=my_new_migration
    ```

    run tests:
    ```bash
        make test
    ```

    generate grpc + gateway server from proto definitions:
    ```bash
        make proto-generate
    ```


