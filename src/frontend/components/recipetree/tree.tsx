// components/RecipeTree.tsx
'use client';
import ReactFlow from 'reactflow';
import 'reactflow/dist/style.css';
import { parseRecipeJson, RecipeJson, edgeTypes, nodeTypes } from '@/util/parser/parser';

interface RecipeTreeProps {
  recipeJson: RecipeJson;
}

const RecipeTree: React.FC<RecipeTreeProps> = ({ recipeJson }) => {
  const { flowNodes, flowEdges } = parseRecipeJson(recipeJson);

  return (
    <div style={{ width: '90%', height: '62.5vh' }}>
      <ReactFlow nodes={flowNodes} edges={flowEdges} edgeTypes={edgeTypes} nodeTypes={nodeTypes} fitView>
      </ReactFlow>
    </div>
  );
};

export default RecipeTree;