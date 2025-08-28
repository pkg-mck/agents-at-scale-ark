import chalk from 'chalk';
import { Command } from 'commander';

import { getClusterIp } from '../../lib/cluster.js';

export function createGetIpCommand(): Command {
  const getIp = new Command('get-ip');
  getIp
    .description('Get the IP address of the current Kubernetes cluster')
    .option('-c, --context <context>', 'Kubernetes context to use')
    .action(async (options) => {
      try {
        const clusterInfo = await getClusterIp(options.context);

        if (clusterInfo.error) {
          console.error(
            chalk.red('Error getting cluster IP:'),
            clusterInfo.error
          );
          process.exit(1);
        }

        if (!clusterInfo.ip) {
          console.error(chalk.red('Could not determine cluster IP'));
          process.exit(1);
        }

        console.log(clusterInfo.ip);

        console.error(chalk.grey(`Cluster type: ${clusterInfo.type}`));
        if (clusterInfo.context) {
          console.error(chalk.grey(`Context: ${clusterInfo.context}`));
        }
      } catch (error: any) {
        console.error(chalk.red('Failed to get cluster IP:'), error.message);
        process.exit(1);
      }
    });

  return getIp;
}
