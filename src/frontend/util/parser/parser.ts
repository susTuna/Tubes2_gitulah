// utils/recipeTree.ts
import dagre from 'dagre';
import StepEdge from '@/components/edges/stepedge';
import ImageNode from '@/components/nodes/imagenode';

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
  type?: string;
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

export const nodeTypes = {
  image: ImageNode,
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
    const image_url = element?.image_url.split("svg")[0] + "svg";

    nodeMap[index] = id;
    flowNodes.push({
      id,
      type: 'image',
      data: { 
        label,
        image_url 
      },
      position: { x: 0, y: 0 },
      width: 25,
      height: 25,
    });
    g.setNode(id, { width: 25, height: 25 });
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
    flowEdges.push({ id: `${target2}-${source}`, type: 'step', source: source, target: target2, style: { stroke: color } });

    g.setEdge(source, target1);
    g.setEdge(source, target2);
  });

  dagre.layout(g);

  // Get the Dagre-computed position of the first node
  const anchorId = flowNodes[0].id;
  const anchorNode = g.node(anchorId);

  // Calculate offset to move anchor node to (0, 0)
  const offsetX = anchorNode.x;
  const offsetY = anchorNode.y;

  flowNodes.forEach(node => {
    const nodeWithPosition = g.node(node.id);
    node.position = {
       x: nodeWithPosition.x - offsetX - node.width / 2,
       y: nodeWithPosition.y - offsetY - node.height / 2
};
  });

  return { flowNodes, flowEdges };
}