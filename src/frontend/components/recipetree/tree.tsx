// components/RecipeTree.tsx
'use client';
import ReactFlow from 'reactflow';
import { Controls } from 'reactflow';
import 'reactflow/dist/style.css';
import { parseRecipeJson, RecipeJson, edgeTypes, nodeTypes } from '@/util/parser/parser';

export default function RecipeTree({ recipeJson } : { recipeJson : RecipeJson}) {
  const { flowNodes, flowEdges } = parseRecipeJson(recipeJson);

  return (
    <div style={{ width: '90%', height: '62.5vh' }}>
      <ReactFlow nodes={flowNodes} edges={flowEdges} edgeTypes={edgeTypes} nodeTypes={nodeTypes} fitView
        nodesDraggable={false}
        nodesConnectable={false}
        elementsSelectable={false}
        nodesFocusable={false}
        panOnScroll={true}
        zoomOnScroll={true}
        zoomOnPinch={true}
        panOnDrag={true}
        selectionOnDrag={false}
        preventScrolling={false}
      >
        <Controls />
      </ReactFlow>
    </div>
  );
};