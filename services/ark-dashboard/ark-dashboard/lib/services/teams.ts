import { apiClient } from '@/lib/api/client'
import type { components } from '@/lib/api/generated/types'

// Helper type for axios errors
interface AxiosError extends Error {
  response?: {
    status: number
  }
}

// Use the generated types from OpenAPI
export type TeamResponse = components['schemas']['TeamResponse']
export type TeamDetailResponse = components['schemas']['TeamDetailResponse']
export type TeamListResponse = components['schemas']['TeamListResponse']
export type TeamCreateRequest = components['schemas']['TeamCreateRequest']
export type TeamUpdateRequest = components['schemas']['TeamUpdateRequest']
export type TeamMember = components['schemas']['TeamMember']

// For UI compatibility, we'll map the API response to include an id field
export type Team = TeamDetailResponse & { id: string }

// CRUD Operations
export const teamsService = {
  // Get all teams
  async getAll(namespace: string): Promise<Team[]> {
    const response = await apiClient.get<TeamListResponse>(`/api/v1/namespaces/${namespace}/teams`)
    
    // Map the response items to include id for UI compatibility
    const teams = await Promise.all(
      response.items.map(async (item) => {
        // Fetch detailed info for each team to get full data
        const detailed = await teamsService.getByName(namespace, item.name)
        return detailed!
      })
    )
    
    return teams
  },

  // Get a single team by name
  async getByName(namespace: string, name: string): Promise<Team | null> {
    try {
      const response = await apiClient.get<TeamDetailResponse>(`/api/v1/namespaces/${namespace}/teams/${name}`)
      return {
        ...response,
        id: response.name // Use name as id for UI compatibility
      }
    } catch (error) {
      if ((error as AxiosError).response?.status === 404) {
        return null
      }
      throw error
    }
  },

  // Get a single team by ID (for UI compatibility - ID is actually the name)
  async getById(namespace: string, id: number | string): Promise<Team | null> {
    // Convert numeric ID to string name
    const name = String(id)
    return teamsService.getByName(namespace, name)
  },

  // Create a new team
  async create(namespace: string, team: TeamCreateRequest): Promise<Team> {
    const response = await apiClient.post<TeamDetailResponse>(`/api/v1/namespaces/${namespace}/teams`, team)
    return {
      ...response,
      id: response.name
    }
  },

  // Update an existing team
  async update(namespace: string, name: string, updates: TeamUpdateRequest): Promise<Team | null> {
    try {
      const response = await apiClient.put<TeamDetailResponse>(`/api/v1/namespaces/${namespace}/teams/${name}`, updates)
      return {
        ...response,
        id: response.name
      }
    } catch (error) {
      if ((error as AxiosError).response?.status === 404) {
        return null
      }
      throw error
    }
  },

  // Update by ID (for UI compatibility)
  async updateById(namespace: string, id: number | string, updates: TeamUpdateRequest): Promise<Team | null> {
    const name = String(id)
    return teamsService.update(namespace, name, updates)
  },

  // Delete a team
  async delete(namespace: string, name: string): Promise<boolean> {
    try {
      await apiClient.delete(`/api/v1/namespaces/${namespace}/teams/${name}`)
      return true
    } catch (error) {
      if ((error as AxiosError).response?.status === 404) {
        return false
      }
      throw error
    }
  },

  // Delete by ID (for UI compatibility)
  async deleteById(namespace: string, id: number | string): Promise<boolean> {
    const name = String(id)
    return teamsService.delete(namespace, name)
  }
}