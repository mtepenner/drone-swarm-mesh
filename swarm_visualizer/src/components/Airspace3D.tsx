import { useEffect, useRef } from "react";
import * as THREE from "three";

import type { SwarmSnapshot } from "../hooks/useSwarmData";

type SceneHandles = {
  renderer: THREE.WebGLRenderer;
  scene: THREE.Scene;
  camera: THREE.PerspectiveCamera;
  instanced: THREE.InstancedMesh;
  links: THREE.LineSegments;
  dummy: THREE.Object3D;
  frameId: number;
  onResize: () => void;
};

type Props = {
  snapshot: SwarmSnapshot;
};

export function Airspace3D({ snapshot }: Props) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const sceneRef = useRef<SceneHandles | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) {
      return;
    }

    const scene = new THREE.Scene();
    scene.background = new THREE.Color("#06131c");
    scene.fog = new THREE.Fog("#06131c", 40, 180);

    const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 300);
    camera.position.set(34, 24, 34);

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

    const ambient = new THREE.HemisphereLight("#f0fdfa", "#052e2b", 1.25);
    const key = new THREE.DirectionalLight("#ffffff", 1.6);
    key.position.set(18, 24, 12);
    const grid = new THREE.GridHelper(160, 24, "#0ea5e9", "#155e75");
    grid.position.y = -0.5;
    scene.add(ambient, key, grid);

    const instanced = new THREE.InstancedMesh(
      new THREE.ConeGeometry(0.42, 1.4, 7),
      new THREE.MeshStandardMaterial({ roughness: 0.32, metalness: 0.15 }),
      256
    );
    instanced.instanceMatrix.setUsage(THREE.DynamicDrawUsage);
    instanced.count = 0;
    scene.add(instanced);

    const links = new THREE.LineSegments(
      new THREE.BufferGeometry(),
      new THREE.LineBasicMaterial({ color: "#22d3ee", transparent: true, opacity: 0.45 })
    );
    scene.add(links);

    const dummy = new THREE.Object3D();

    const onResize = () => {
      const width = canvas.clientWidth;
      const height = canvas.clientHeight;
      if (width === 0 || height === 0) {
        return;
      }
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height, false);
    };

    const animate = () => {
      renderer.render(scene, camera);
      sceneRef.current!.frameId = window.requestAnimationFrame(animate);
    };

    sceneRef.current = {
      renderer,
      scene,
      camera,
      instanced,
      links,
      dummy,
      frameId: window.requestAnimationFrame(animate),
      onResize,
    };
    window.addEventListener("resize", onResize);
    onResize();

    return () => {
      if (!sceneRef.current) {
        return;
      }
      window.cancelAnimationFrame(sceneRef.current.frameId);
      window.removeEventListener("resize", sceneRef.current.onResize);
      sceneRef.current.renderer.dispose();
      sceneRef.current = null;
    };
  }, []);

  useEffect(() => {
    const handles = sceneRef.current;
    if (!handles) {
      return;
    }

    const { instanced, dummy, links } = handles;
    const drones = snapshot.drones.slice(0, 256);
    const droneMap = new Map(drones.map((drone) => [drone.agent_id, drone]));
    instanced.count = drones.length;

    drones.forEach((drone, index) => {
      dummy.position.set(drone.x, drone.y, drone.z);
      dummy.rotation.set(Math.PI, Math.atan2(drone.vx, drone.vz), 0);
      const speed = Math.hypot(drone.vx, drone.vy, drone.vz);
      const scale = THREE.MathUtils.mapLinear(speed, 0, 6, 0.8, 1.8);
      dummy.scale.set(scale, scale, scale);
      dummy.updateMatrix();
      instanced.setMatrixAt(index, dummy.matrix);
      instanced.setColorAt(
        index,
        new THREE.Color().setHSL(0.36 - Math.min(speed / 12, 0.18), 0.82, 0.55)
      );
    });
    instanced.instanceMatrix.needsUpdate = true;
    instanced.instanceColor!.needsUpdate = true;

    const positions: number[] = [];
    snapshot.links.forEach((link) => {
      const left = droneMap.get(link.source);
      const right = droneMap.get(link.target);
      if (!left || !right) {
        return;
      }
      positions.push(left.x, left.y, left.z, right.x, right.y, right.z);
    });

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute("position", new THREE.Float32BufferAttribute(positions, 3));
    links.geometry.dispose();
    links.geometry = geometry;
  }, [snapshot]);

  return <canvas className="airspace-canvas" ref={canvasRef} />;
}
