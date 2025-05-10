// utils/recipeTree.ts
import dagre from 'dagre';
import StepEdge from '@/components/edges/stepedge';

export interface RecipeJson {
  nodes: string[];
  dependencies: {
    node: number;
    dependency1: number;
    dependency2: number;
  }[];
}

export interface FlowNode {
  id: string;
  data: { label: string };
  position: { x: number; y: number };
}

export interface FlowEdge {
  id: string;
  type: string;
  source: string;
  target: string;
  style?: { [key: string]: any };
}

const generateRandomColor = () => {
  const r = Math.floor(Math.random() * 128) + 127;
  const g = Math.floor(Math.random() * 128) + 127;
  const b = Math.floor(Math.random() * 128) + 127;

  return `rgb(${r}, ${g}, ${b})`;
};  

export function parseRecipeJson(recipeJson: RecipeJson): { flowNodes: FlowNode[]; flowEdges: FlowEdge[]; edgeTypes: { [key: string]: any } } {
  const edgeTypes = {
    step: StepEdge,
  }

  const { nodes, dependencies } = recipeJson;
  const nodeMap: { [key: number]: string } = {};
  const flowNodes: FlowNode[] = [];
  const flowEdges: FlowEdge[] = [];

  const g = new dagre.graphlib.Graph();
  g.setGraph({});
  g.setDefaultEdgeLabel(() => ({}));
  
  nodes.forEach((name, index) => {
    const id = index.toString();
    nodeMap[index] = id;
    flowNodes.push({
      id,
      data: { label: name },
      position: { x: 0, y: 0 },
    });
    g.setNode(id, { width: 150, height: 50 });
  });

  let sauce = 'Yuuka'
  let color = generateRandomColor();
  dependencies.forEach(dep => {
    const source = nodeMap[dep.node];
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

  return { flowNodes, flowEdges, edgeTypes };
}