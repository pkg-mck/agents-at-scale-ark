import chalk from 'chalk';
import { Command } from 'commander';

import { detectClusterType } from '../../lib/cluster.js';

export function createGetTypeCommand(): Command {
  const getType = new Command('get-type');
  getType
    .description('Detect and display the current Kubernetes cluster type')
    .action(async () => {
      try {
        const clusterInfo = await detectClusterType();

        if (clusterInfo.error) {
          console.error(
            chalk.red('Error detecting cluster type:'),
            clusterInfo.error
          );
          process.exit(1);
        }

        console.log(clusterInfo.type);

        if (clusterInfo.context) {
          console.error(chalk.grey(`Context: ${clusterInfo.context}`));
        }
      } catch (error: any) {
        console.error(
          chalk.red('Failed to detect cluster type:'),
          error.message
        );
        process.exit(1);
      }
    });

  return getType;
}
