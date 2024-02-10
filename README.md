# Project Title

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Go
- MongoDB
- Redis

### Installation

1. **Clone the project**

   Use git to clone the project onto your local machine.

   ```bash
   git clone <project-url>

   ```

2. **Build and run the project with Docker Navigate to the root directory of the project in your terminal and run the following command to build and run the project with Docker.**

- docker-compose up --build

3. **Running the tests**

- Make sure the Redis container is up and running in Docker.
- Navigate to the root directory of the project in your terminal.
- cd tests/e2e
- go test -v

4. **Running the project without Docker**

- If you need to run the project without Docker, don’t forget to change the MongoDB host in the database_config.yaml file to localhost.
