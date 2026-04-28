# 🚁 Drone Swarm Mesh

## Description
Drone Swarm Mesh is a highly scalable, containerized simulation environment designed for developing and testing decentralized drone swarms. The architecture separates the global physics simulation, autonomous drone logic, and real-time visual telemetry. Each drone operates as its own independent Go-based container, simulating real-world decentralized mesh networking and autonomous decision-making. 

## 📑 Table of Contents
- [Features](#-features)
- [Technologies Used](#-technologies-used)
- [Installation](#-installation)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Features
* **Autonomous Go Agents:** Each drone acts as an autonomous brain in a lightweight (15-20MB) Go container, featuring local PID loops for 3D translation.
* **Decentralized Mesh Networking:** Implements UDP broadcasting for peer discovery and P2P gossip protocols to exchange position and velocity data.
* **Advanced Flight Behaviors:** Utilizes Boids algorithms for flocking (separation, alignment, cohesion) and ORCA/RVO logic for decentralized collision avoidance.
* **Global Airspace Physics:** A robust Python and FastAPI backend that tracks true positions and defines 3D boundaries/obstacles for swarm challenges.
* **Tactical 3D Visualizer:** A React and TypeScript frontend utilizing Three.js to render 100+ drone agents, visualizing P2P mesh links, average velocity, and inter-drone distances in real time.
* **Dynamic Scalability:** Orchestration scripts allow you to dynamically spin up 100+ drone containers on demand.

## 🛠️ Technologies Used
* **Drone Agent:** Go
* **Simulator Backend:** Python, FastAPI, WebSockets
* **Visualizer Frontend:** React, TypeScript, Three.js
* **Orchestration:** Docker, Docker Compose, Shell Scripting
* **CI/CD:** GitHub Actions

## ⚙️ Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mtepenner/drone-swarm-mesh.git
   cd drone-swarm-mesh
   ```

2. Build the Docker images (including the lightweight Go daemon):
   ```bash
   make build
   # Alternatively, refer to the .github/workflows for build instructions.
   ```

3. Launch the core infrastructure (simulator and basic swarm):
   ```bash
   cd orchestration
   docker-compose up -d
   ```

4. Local validation without Docker:
   ```bash
   python -m compileall swarm_simulator/app
   cd drone_agent && go test ./...
   cd ../swarm_visualizer && npm install && npm run build
   ```

## 💻 Usage

* **Access the Tactical HUD:** Once the containers are running, navigate to `http://localhost:3000` (or your configured port) to access the React-based visualizer.
* **Scale the Swarm:** To test load and collision logic with a larger swarm, use the included scaling script to spawn additional drone agents:
  ```bash
  ./orchestration/scale_swarm.sh 100
  ```
* **Simulator endpoints:** The FastAPI simulator exposes `GET /health`, `GET /snapshot`, and `WS /ws` on port `8000`.

## 📂 Project Structure
* `/swarm_simulator`: Contains the Python/FastAPI code for global airspace physics and WebSocket UDP relay.
* `/drone_agent`: Contains the Go application for the drone's autonomous brain, mesh networking, and flight controllers.
* `/swarm_visualizer`: Contains the React/TypeScript frontend for rendering the 3D airspace and swarm metrics.
* `/orchestration`: Contains the Docker Compose file and shell scripts for deploying and scaling the containerized swarm.
* `/.github/workflows`: Contains CI/CD pipelines for swarm logic verification and container building.

## 🤝 Contributing
Contributions are welcome! Please feel free to submit a Pull Request. To run the automated tests for P2P discovery and collision logic, please utilize the included GitHub Action workflows.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

