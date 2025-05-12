// utils/recipeTree.ts
import dagre from 'dagre';
import StepEdge from '@/components/edges/stepedge';

export interface RecipeJson {
  nodes: number[];
  dependencies: {
    result: number;
    dependency1: number;
    dependency2: number;
  }[];
  time_taken: number;
  nodes_searched: number;
  recipes_found: number;
  finished: boolean
  nameMap: Record<number, ElementInfo>;
}

export interface FlowNode {
  id: string;
  data: { label: string; image_url: string };
  position: { x: number; y: number };
  width: number;
  height: number;
}

export interface FlowEdge {
  id: string;
  type: string;
  source: string;
  target: string;
  style?: { [key: string]: any };
}

export const edgeTypes = {
  step: StepEdge,
}

const generateRandomColor = () => {
  const r = Math.floor(Math.random() * 128) + 127;
  const g = Math.floor(Math.random() * 128) + 127;
  const b = Math.floor(Math.random() * 128) + 127;

  return `rgb(${r}, ${g}, ${b})`;
};  

interface ElementInfo {
  id: number;
  name: string;
  tier: number;
  image_url: string;
}

export const fetchElementInfo = async (
  ids: number[]
): Promise<Record<number, ElementInfo>> => {
  const map: Record<number, ElementInfo> = {};
  await Promise.all(
    ids.map(async (id) => {
      try {
        const res = await fetch(`http://localhost:5761/elements/${id}?type=id`);
        const data: ElementInfo = await res.json();
        map[id] = data;
      } catch (err) {
        console.warn(`Failed to fetch element ${id}`, err);
      }
    })
  );
  return map;
};


export function parseRecipeJson(recipeJson: RecipeJson): { flowNodes: FlowNode[]; flowEdges: FlowEdge[] } {

  const { nodes, dependencies } = recipeJson;
  const nodeMap: { [key: number]: string } = {};
  const flowNodes: FlowNode[] = [];
  const flowEdges: FlowEdge[] = [];
  const nameMap = recipeJson.nameMap || {};

  const g = new dagre.graphlib.Graph();
  g.setGraph({});
  g.setDefaultEdgeLabel(() => ({}));
  
  nodes.forEach((name, index) => {
    const id = index.toString();
    const element = nameMap[name];
    const label = element?.name || `Element ${name}`;
    const image_url = element?.image_url;

    nodeMap[index] = id;
    flowNodes.push({
      id,
      data: { 
        label,
        image_url 
      },
      position: { x: 0, y: 0 },
      width: 50,
      height: 50,
    });
    g.setNode(id, { width: 50, height: 50 });
  });

  let sauce = 'Yuuka'
  let color = generateRandomColor();
  dependencies.forEach(dep => {
    const source = nodeMap[dep.result];
    if (source === sauce) color = generateRandomColor();
    sauce = source;
    const target1 = nodeMap[dep.dependency1];
    const target2 = nodeMap[dep.dependency2];

    flowEdges.push({ id: `${target1}-${source}`, type: 'step', source: source, target: target1, style: { stroke: color } });
    flowEdges.push({ id: `${target2}-${source}`, type: 'step', source, target: target2, style: { stroke: color } });

    g.setEdge(source, target1);
    g.setEdge(source, target2);
  });

  dagre.layout(g);

  flowNodes.forEach(node => {
    const nodeWithPosition = g.node(node.id);
    node.position = { x: nodeWithPosition.x, y: nodeWithPosition.y };
  });

  return { flowNodes, flowEdges };
}