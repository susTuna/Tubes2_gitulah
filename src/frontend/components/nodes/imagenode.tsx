import React from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export default function ImageNode({ data } : NodeProps) {
  return (
    <div style={{ textAlign: 'center', width: '25px' }}>

      <img src={data.image_url} alt={data.label} style={{ width: '25px', height: '25px' }} />
      
      <div style={{ marginTop: '8px', fontSize: '8px', color: 'white'}}>{data.label}</div>
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ opacity: 0, width: 0, height: 0 }}
      />

      <Handle
        type="target"
        position={Position.Top}
        style={{ opacity: 0, width: 0, height: 0 }}
        />
    </div>
  );
};