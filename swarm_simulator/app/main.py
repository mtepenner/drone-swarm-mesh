from __future__ import annotations

import asyncio
import json
import os
import time
from dataclasses import dataclass, field
from math import sqrt
from typing import Any

from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware

from app.obstacles.narrow_gap import AIRSPACE_BOUNDS, NARROW_GAP_OBSTACLES
from app.physics.spatial_grid import DronePoint, SpatialGrid, clamp_point


@dataclass(slots=True)
class DroneState:
    agent_id: str
    position: tuple[float, float, float]
    velocity: tuple[float, float, float]
    battery_level: float
    peer_ids: list[str] = field(default_factory=list)
    updated_at: float = field(default_factory=lambda: time.time())

    @property
    def speed(self) -> float:
        vx, vy, vz = self.velocity
        return sqrt(vx * vx + vy * vy + vz * vz)


class SwarmRuntime:
    def __init__(self) -> None:
        self.bounds = AIRSPACE_BOUNDS
        self.grid = SpatialGrid(cell_size=8.0)
        self.drones: dict[str, DroneState] = {}

    def ingest(self, payload: dict[str, Any]) -> None:
        position = clamp_point(
            (
                float(payload.get("x", 0.0)),
                float(payload.get("y", 0.0)),
                float(payload.get("z", 0.0)),
            ),
            self.bounds,
        )
        velocity = (
            float(payload.get("vx", 0.0)),
            float(payload.get("vy", 0.0)),
            float(payload.get("vz", 0.0)),
        )
        peer_ids = [str(peer) for peer in payload.get("peer_ids", [])]
        state = DroneState(
            agent_id=str(payload.get("agent_id", "unknown")),
            position=position,
            velocity=velocity,
            battery_level=float(payload.get("battery_level", 100.0)),
            peer_ids=peer_ids,
        )
        self.drones[state.agent_id] = state
        self._refresh_grid()

    def snapshot(self) -> dict[str, Any]:
        drones = sorted(self.drones.values(), key=lambda drone: drone.agent_id)
        links: list[dict[str, Any]] = []
        for drone in drones:
            for peer_id in drone.peer_ids:
                if drone.agent_id < peer_id:
                    links.append(
                        {
                            "source": drone.agent_id,
                            "target": peer_id,
                        }
                    )

        average_speed = sum(drone.speed for drone in drones) / len(drones) if drones else 0.0
        collisions = self.grid.collision_pairs()
        return {
            "drone_count": len(drones),
            "drones": [
                {
                    "agent_id": drone.agent_id,
                    "x": drone.position[0],
                    "y": drone.position[1],
                    "z": drone.position[2],
                    "vx": drone.velocity[0],
                    "vy": drone.velocity[1],
                    "vz": drone.velocity[2],
                    "battery_level": drone.battery_level,
                    "peer_ids": drone.peer_ids,
                }
                for drone in drones
            ],
            "links": links,
            "metrics": {
                "average_speed": average_speed,
                "mean_spacing": self.grid.nearest_distance_mean(),
                "collision_pairs": len(collisions),
                "obstacle_count": len(NARROW_GAP_OBSTACLES),
            },
            "obstacles": NARROW_GAP_OBSTACLES,
        }

    def _refresh_grid(self) -> None:
        now = time.time()
        self.drones = {
            agent_id: drone
            for agent_id, drone in self.drones.items()
            if now - drone.updated_at <= 5.0
        }
        self.grid.rebuild(
            DronePoint(drone.agent_id, drone.position, drone.velocity)
            for drone in self.drones.values()
        )


class TelemetryProtocol(asyncio.DatagramProtocol):
    def __init__(self, runtime: SwarmRuntime) -> None:
        self.runtime = runtime

    def datagram_received(self, data: bytes, addr: tuple[str, int]) -> None:
        try:
            payload = json.loads(data.decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError):
            return

        if "agent_id" not in payload:
            payload["agent_id"] = f"agent-{addr[0]}:{addr[1]}"
        self.runtime.ingest(payload)


runtime = SwarmRuntime()
transport: asyncio.DatagramTransport | None = None

app = FastAPI(title="Drone Swarm Mesh Simulator")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.on_event("startup")
async def startup() -> None:
    global transport
    udp_host = os.getenv("SWARM_SIM_UDP_HOST", "0.0.0.0")
    udp_port = int(os.getenv("SWARM_SIM_UDP_PORT", "9010"))
    loop = asyncio.get_running_loop()
    transport, _ = await loop.create_datagram_endpoint(
        lambda: TelemetryProtocol(runtime),
        local_addr=(udp_host, udp_port),
    )


@app.on_event("shutdown")
async def shutdown() -> None:
    if transport is not None:
        transport.close()


@app.get("/health")
async def health() -> dict[str, Any]:
    snapshot = runtime.snapshot()
    return {
        "status": "ok",
        "drone_count": snapshot["drone_count"],
        "collision_pairs": snapshot["metrics"]["collision_pairs"],
    }


@app.get("/snapshot")
async def snapshot() -> dict[str, Any]:
    return runtime.snapshot()


@app.websocket("/ws")
async def websocket_feed(websocket: WebSocket) -> None:
    await websocket.accept()
    try:
        while True:
            await websocket.send_json(runtime.snapshot())
            await asyncio.sleep(0.4)
    except WebSocketDisconnect:
        return
