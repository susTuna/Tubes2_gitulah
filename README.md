# Tubes2_gitulah

## Getting Started
These instructions are meant to setup and run Little Alchemy 2 Recipe Finder locally.

### Prerequisite
- Docker
- Docker Compose (comes built-in with Docker Desktop)

### Installing
#### Windows & macOS (via Docker Desktop)
1. **Download** [Docker Desktop](https://www.docker.com/products/docker-desktop)
2. **Install Docker Desktop**

    Follow the installer. Once done:
    - Verify Docker works by running `docker --version` in Command Prompt or Terminal.
    - Ensure WSL is installed (on Windows).
3. **Clone this repository**
    ```bash
    git clone https://github.com/filbertengyo/Tubes2_gitulah.git
    cd Tubes2_gitulah
    ```
4. **Run the app**
    ```bash
    docker compose up --build
    ```
#### Linux (Ubuntu)
1. **Install Docker Engine**
    ```bash
    sudo apt-get install docker.io
    sudo systemctl enable docker #enable the docker on service
    sudo systemctl start docker #start docker on service
    ```
2. **Install Docker Compose Plugin**
    ```bash
    sudo apt-get install docker-compose-plugin
    ```
3. **Verify Installation**
    ```bash
    docker --version
    docker compose version
    ```
4. **Clone the project & run**
    ```bash
    git clone https://github.com/filbertengyo/Tubes2_gitulah.git
    cd Tubes2_gitulah
    docker compose up --build
    ```
### Running the App
The recipe finder can be accessed [here](localhost:3000)

### Stopping the App
To stop the app
```bash
docker compose down
```