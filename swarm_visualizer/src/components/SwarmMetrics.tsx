import type { SwarmSnapshot } from "../hooks/useSwarmData";

type Props = {
  snapshot: SwarmSnapshot;
  status: string;
  lastUpdated: string;
};

export function SwarmMetrics({ snapshot, status, lastUpdated }: Props) {
  const cards = [
    { label: "Feed", value: status.toUpperCase() },
    { label: "Drones", value: String(snapshot.drone_count) },
    { label: "Avg Speed", value: `${snapshot.metrics.average_speed.toFixed(2)} m/s` },
    { label: "Mean Spacing", value: `${snapshot.metrics.mean_spacing.toFixed(2)} m` },
    { label: "Collision Pairs", value: String(snapshot.metrics.collision_pairs) },
    { label: "Obstacles", value: String(snapshot.metrics.obstacle_count) },
  ];

  return (
    <section className="metrics-grid">
      {cards.map((card) => (
        <article className="metric-card" key={card.label}>
          <span>{card.label}</span>
          <strong>{card.value}</strong>
        </article>
      ))}
      <article className="metric-card metric-card-wide">
        <span>Last Update</span>
        <strong>{lastUpdated}</strong>
      </article>
    </section>
  );
}
