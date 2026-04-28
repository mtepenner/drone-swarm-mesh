from __future__ import annotations

from collections import defaultdict
from dataclasses import dataclass
from math import dist
from typing import Iterable


Vec3 = tuple[float, float, float]


@dataclass(slots=True)
class DronePoint:
    agent_id: str
    position: Vec3
    velocity: Vec3


def clamp_point(point: Vec3, bounds: tuple[Vec3, Vec3]) -> Vec3:
    lower, upper = bounds
    return (
        min(max(point[0], lower[0]), upper[0]),
        min(max(point[1], lower[1]), upper[1]),
        min(max(point[2], lower[2]), upper[2]),
    )


class SpatialGrid:
    def __init__(self, cell_size: float = 6.0) -> None:
        self.cell_size = cell_size
        self._cells: dict[tuple[int, int, int], list[DronePoint]] = defaultdict(list)

    def rebuild(self, drones: Iterable[DronePoint]) -> None:
        self._cells.clear()
        for drone in drones:
            self._cells[self._cell_for(drone.position)].append(drone)

    def nearest_distance_mean(self) -> float:
        distances: list[float] = []
        drones = [drone for bucket in self._cells.values() for drone in bucket]
        for drone in drones:
            candidates = [
                other
                for other in drones
                if other.agent_id != drone.agent_id
            ]
            if not candidates:
                continue
            nearest = min(dist(drone.position, other.position) for other in candidates)
            distances.append(nearest)
        return sum(distances) / len(distances) if distances else 0.0

    def collision_pairs(self, radius: float = 1.4) -> list[tuple[str, str]]:
        drones = [drone for bucket in self._cells.values() for drone in bucket]
        pairs: list[tuple[str, str]] = []
        for index, left in enumerate(drones):
            for right in drones[index + 1 :]:
                if dist(left.position, right.position) <= radius:
                    pairs.append((left.agent_id, right.agent_id))
        return pairs

    def _cell_for(self, point: Vec3) -> tuple[int, int, int]:
        return (
            int(point[0] // self.cell_size),
            int(point[1] // self.cell_size),
            int(point[2] // self.cell_size),
        )
