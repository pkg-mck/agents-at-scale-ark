import { Command } from 'commander';

import { createGetIpCommand } from './get-ip.js';
import { createGetTypeCommand } from './get-type.js';

export function createClusterCommand(): Command {
  const cluster = new Command('cluster');
  cluster.description('Cluster management commands');

  cluster.addCommand(createGetTypeCommand());
  cluster.addCommand(createGetIpCommand());

  return cluster;
}
