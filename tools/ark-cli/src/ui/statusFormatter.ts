import chalk from 'chalk';

import { StatusData, ServiceStatus, DependencyStatus } from '../lib/types.js';

export class StatusFormatter {
  /**
   * Print status check results to console
   */
  public static printStatus(statusData: StatusData): void {
    console.log(chalk.cyan.bold('\n🔍 ARK System Status Check'));
    console.log(chalk.gray('Checking ARK services and dependencies...\n'));

    // Print services status
    console.log(chalk.cyan.bold('📡 ARK Services:'));
    for (const service of statusData.services) {
      StatusFormatter.printService(service);
    }

    // Print dependencies status
    console.log(chalk.cyan.bold('\n🛠️  System Dependencies:'));
    for (const dep of statusData.dependencies) {
      StatusFormatter.printDependency(dep);
    }

    console.log();
  }

  private static printService(service: ServiceStatus): void {
    const statusColor =
      service.status === 'healthy'
        ? chalk.green('✓ healthy')
        : service.status === 'unhealthy'
          ? chalk.red('✗ unhealthy')
          : chalk.yellow('? not installed');

    console.log(`  • ${chalk.bold(service.name)}: ${statusColor}`);
    if (service.url) {
      console.log(`    ${chalk.gray(`URL: ${service.url}`)}`);
    }
    if (service.details) {
      console.log(`    ${chalk.gray(service.details)}`);
    }
  }

  private static printDependency(dep: DependencyStatus): void {
    const statusColor = dep.installed
      ? chalk.green('✓ installed')
      : chalk.red('✗ missing');

    console.log(`  • ${chalk.bold(dep.name)}: ${statusColor}`);
    if (dep.version) {
      console.log(`    ${chalk.gray(`Version: ${dep.version}`)}`);
    }
    if (dep.details) {
      console.log(`    ${chalk.gray(dep.details)}`);
    }
  }
}
