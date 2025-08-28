export interface ArkConfig {
  defaultModel?: string;
  defaultAgent?: string;
  defaultNamespace?: string;
  apiBaseUrl?: string;
  kubeconfig?: string;
  currentContext?: string;
  kubeNamespace?: string;
}

export interface KubernetesConfig {
  kubeconfig: string;
  currentContext?: string;
  namespace?: string;
  inCluster: boolean;
}

export interface ServiceStatus {
  name: string;
  status: 'healthy' | 'unhealthy' | 'not installed';
  url?: string;
  version?: string;
  details?: string;
}

export interface DependencyStatus {
  name: string;
  installed: boolean;
  version?: string;
  details?: string;
}

export interface StatusData {
  services: ServiceStatus[];
  dependencies: DependencyStatus[];
}

export interface CommandVersionConfig {
  command: string;
  versionArgs: string;
  versionExtract: (_output: string) => string;
}
