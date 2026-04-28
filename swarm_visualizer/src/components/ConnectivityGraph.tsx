import type { SwarmSnapshot } from "../hooks/useSwarmData";

type Props = {
  snapshot: SwarmSnapshot;
};

export function ConnectivityGraph({ snapshot }: Props) {
  const ranked = [...snapshot.drones]
    .sort((left, right) => right.peer_ids.length - left.peer_ids.length)
    .slice(0, 5);

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="eyebrow">Mesh Graph</p>
        <h2>Strongest Peer Clusters</h2>
      </div>
      <div className="connectivity-list">
        {ranked.map((drone) => (
          <article className="connectivity-card" key={drone.agent_id}>
            <strong>{drone.agent_id}</strong>
            <span>{drone.peer_ids.length} active peer links</span>
            <p>{drone.peer_ids.slice(0, 4).join(", ") || "No current neighbors"}</p>
          </article>
        ))}
      </div>
    </section>
  );
}
