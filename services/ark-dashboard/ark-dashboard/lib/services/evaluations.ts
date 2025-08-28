export interface Evaluation {
  id: string
  name: string
}

// Test data
const testEvaluations: Evaluation[] = [
  {
    id: "550e8400-e29b-41d4-a716-446655440001",
    name: "Response Quality Assessment"
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440002",
    name: "Latency Performance Test"
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440003",
    name: "Accuracy Benchmark"
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440004",
    name: "Safety Compliance Check"
  }
]

// Service with read-only operations
export const evaluationsService = {
  // Get all evaluations
  async getAll(): Promise<Evaluation[]> {
    return Promise.resolve([...testEvaluations])
  },

  // Get a single evaluation by ID
  async getById(id: string): Promise<Evaluation | null> {
    const evaluation = testEvaluations.find(e => e.id === id)
    return Promise.resolve(evaluation || null)
  }
}