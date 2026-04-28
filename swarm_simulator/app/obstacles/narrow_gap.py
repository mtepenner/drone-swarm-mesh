from __future__ import annotations


AIRSPACE_BOUNDS = ((-60.0, 0.0, -60.0), (60.0, 35.0, 60.0))

NARROW_GAP_OBSTACLES = [
    {
        "id": "canyon-wall-west",
        "min": (-8.0, 0.0, -12.0),
        "max": (-2.5, 28.0, 12.0),
    },
    {
        "id": "canyon-wall-east",
        "min": (2.5, 0.0, -12.0),
        "max": (8.0, 28.0, 12.0),
    },
    {
        "id": "tower",
        "min": (-20.0, 0.0, 22.0),
        "max": (-12.0, 26.0, 30.0),
    },
]
