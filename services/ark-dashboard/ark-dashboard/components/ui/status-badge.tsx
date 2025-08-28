import React from 'react';

interface StatusBadgeProps {
  ready?: boolean;
  discovering?: boolean;
}

export function StatusBadge({ ready, discovering }: StatusBadgeProps) {
  if (ready === true) {
    return (
      <div className="px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
        Ready
      </div>
    );
  } else if (discovering === true) {
    return (
      <div className="px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
        Discovering
      </div>
    );
  } else {
    return (
      <div className="px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
        Not Ready
      </div>
    );
  }
}
