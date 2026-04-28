import { Airspace3D } from "./components/Airspace3D";
import { ConnectivityGraph } from "./components/ConnectivityGraph";
import { SwarmMetrics } from "./components/SwarmMetrics";
import { useSwarmData } from "./hooks/useSwarmData";

export default function App() {
  const { snapshot, status, lastUpdated } = useSwarmData();

  return (
    <main className="app-shell">
      <section className="hero-panel">
        <div>
          <p className="eyebrow">Decentralized Airspace</p>
          <h1>Drone Swarm Mesh HUD</h1>
          <p className="hero-copy">
            Track decentralized drone agents, watch peer-to-peer links form in real time,
            and monitor flocking stability as the swarm squeezes through constrained airspace.
          </p>
        </div>
        <SwarmMetrics snapshot={snapshot} status={status} lastUpdated={lastUpdated} />
      </section>

      <section className="content-grid">
        <div className="visual-panel panel">
          <div className="panel-header">
            <p className="eyebrow">Three-Dimensional Telemetry</p>
            <h2>Airspace and Mesh Links</h2>
          </div>
          <Airspace3D snapshot={snapshot} />
        </div>

        <ConnectivityGraph snapshot={snapshot} />
      </section>
    </main>
  );
}
