import { Tool } from "../services/tools";

export const groupToolsByLabel = (tools: Tool[]) => {
  const groups: Record<string, { tools: Tool[]; isMcp: boolean }> = {};

  tools.forEach((tool) => {
    let groupName = "Built In";
    let isMcp = false;

    if (tool.labels && typeof tool.labels === "object") {
      const labels = tool.labels as Record<string, string>;
      if (labels["mcp/server"]) {
        groupName = labels["mcp/server"];
        isMcp = true;
      }
    }

    if (!groups[groupName]) {
      groups[groupName] = { tools: [], isMcp };
    }
    groups[groupName].tools.push(tool);
  });

  return Object.entries(groups).map(([groupName, { tools, isMcp }]) => ({
    groupName,
    tools,
    isMcp
  }));
};
