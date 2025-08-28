import { apiClient } from '@/lib/api/client';
import type { components } from '@/lib/api/generated/types';

export type HTTPRouteInfo = components['schemas']['HTTPRouteInfo'];
export type ArkService = components['schemas']['ArkService'];
export type ArkServiceListResponse = components['schemas']['ArkServiceListResponse'];

export const arkServicesService = {
  /**
   * Get all ARK services in a namespace
   */
  async getAll(namespace: string): Promise<ArkServiceListResponse> {
    return apiClient.get<ArkServiceListResponse>(`/api/v1/namespaces/${namespace}/ark-services`);
  }
};