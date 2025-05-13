import { BaseEdge, EdgeProps } from 'reactflow';

export default function StepEdge({ id , sourceX, sourceY, targetX, targetY, style } : EdgeProps) {
  const centerY = (targetY - sourceY) / 2 + sourceY;
 
  const edgePath = `M ${sourceX} ${sourceY} L ${sourceX} ${centerY} L ${targetX} ${centerY} L ${targetX} ${targetY}`;
 
  return (
    <>
        <BaseEdge id={id} path={edgePath} style={style} />
    </>
  );
}