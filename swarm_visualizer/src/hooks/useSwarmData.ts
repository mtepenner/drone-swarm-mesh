import { useEffect, useMemo, useState } from "react";

export type SwarmDrone = {
  agent_id: string;
  x: number;
  y: number;
  z: number;
  vx: number;
  vy: number;
  vz: number;
  battery_level: number;
  peer_ids: string[];
};

export type MeshLink = {
  source: string;
  target: string;
};

export type SwarmSnapshot = {
  drone_count: number;
  drones: SwarmDrone[];
  links: MeshLink[];
  metrics: {
    average_speed: number;
    mean_spacing: number;
    collision_pairs: number;
    obstacle_count: number;
  };
};

const EMPTY_SNAPSHOT: SwarmSnapshot = {
  drone_count: 0,
  drones: [],
  links: [],
  metrics: {
    average_speed: 0,
    mean_spacing: 0,
    collision_pairs: 0,
    obstacle_count: 0,
  },
};

const SWARM_WS_URL = import.meta.env.VITE_SWARM_WS_URL ?? "ws://127.0.0.1:8000/ws";

export function useSwarmData() {
  const [snapshot, setSnapshot] = useState<SwarmSnapshot>(EMPTY_SNAPSHOT);
  const [status, setStatus] = useState<"connecting" | "live" | "simulated">("connecting");
  const [lastUpdated, setLastUpdated] = useState<string>("waiting for telemetry");

  useEffect(() => {
    let closed = false;
    let fallbackTimer: number | undefined;
    let socket: WebSocket | undefined;

    const startSimulation = () => {
      if (closed || fallbackTimer !== undefined) {
        return;
      }

      let tick = 0;
      setStatus("simulated");
      setSnapshot(buildSimulationSnapshot(tick));
      setLastUpdated("simulated feed");
      fallbackTimer = window.setInterval(() => {
        tick += 1;
        setSnapshot(buildSimulationSnapshot(tick));
        setLastUpdated(`simulated tick ${tick}`);
      }, 900);
    };

    try {
      socket = new WebSocket(SWARM_WS_URL);
      socket.addEventListener("open", () => {
        if (fallbackTimer !== undefined) {
          window.clearInterval(fallbackTimer);
          fallbackTimer = undefined;
        }
        setStatus("live");
      });

      socket.addEventListener("message", (event) => {
        const parsed = normalizeSnapshot(JSON.parse(event.data));
        setSnapshot(parsed);
        setStatus("live");
        setLastUpdated(new Date().toLocaleTimeString());
      });

      socket.addEventListener("error", startSimulation);
      socket.addEventListener("close", startSimulation);
    } catch (error) {
      console.warn("failed to connect to simulator websocket", error);
      startSimulation();
    }

    return () => {
      closed = true;
      if (fallbackTimer !== undefined) {
        window.clearInterval(fallbackTimer);
      }
      socket?.close();
    };
  }, []);

  return useMemo(
    () => ({ snapshot, status, lastUpdated }),
    [lastUpdated, snapshot, status]
  );
}

function normalizeSnapshot(raw: Partial<SwarmSnapshot>): SwarmSnapshot {
  return {
    drone_count: raw.drone_count ?? 0,
    drones: raw.drones ?? [],
    links: raw.links ?? [],
    metrics: {
      average_speed: raw.metrics?.average_speed ?? 0,
      mean_spacing: raw.metrics?.mean_spacing ?? 0,
      collision_pairs: raw.metrics?.collision_pairs ?? 0,
      obstacle_count: raw.metrics?.obstacle_count ?? 0,
    },
  };
}

function buildSimulationSnapshot(tick: number): SwarmSnapshot {
  const droneCount = 42;
  const drones = Array.from({ length: droneCount }, (_, index) => {
    const angle = (index / droneCount) * Math.PI * 2 + tick * 0.09;
    const radius = 14 + (index % 6) * 2.7;
    return {
      agent_id: `drone-${String(index + 1).padStart(3, "0")}`,
      x: Math.cos(angle) * radius,
      y: 8 + (index % 5),
      z: Math.sin(angle) * radius,
      vx: -Math.sin(angle) * 2.1,
      vy: Math.cos(angle * 0.5) * 0.4,
      vz: Math.cos(angle) * 2.1,
      battery_level: 74 + ((index + tick) % 22),
      peer_ids: [
        `drone-${String(((index + 1) % droneCount) + 1).padStart(3, "0")}`,
        `drone-${String(((index + 4) % droneCount) + 1).padStart(3, "0")}`,
      ],
    };
  });

  const links = drones.flatMap((drone) =>
    drone.peer_ids
      .filter((peerId) => drone.agent_id < peerId)
      .map((peerId) => ({ source: drone.agent_id, target: peerId }))
  );

  return {
    drone_count: droneCount,
    drones,
    links,
    metrics: {
      average_speed: 2.18,
      mean_spacing: 6.42,
      collision_pairs: tick % 6 === 0 ? 1 : 0,
      obstacle_count: 3,
    },
  };
}
