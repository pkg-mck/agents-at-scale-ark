import { describe, it, expect } from 'vitest'
import { isValidKubernetesName, getKubernetesNameError } from '@/lib/utils/kubernetes-validation'

describe('isValidKubernetesName', () => {
  describe('valid names', () => {
    it('should accept simple lowercase names', () => {
      expect(isValidKubernetesName('myapp')).toBe(true)
      expect(isValidKubernetesName('app123')).toBe(true)
      expect(isValidKubernetesName('123app')).toBe(true)
    })

    it('should accept names with hyphens', () => {
      expect(isValidKubernetesName('my-app')).toBe(true)
      expect(isValidKubernetesName('app-123')).toBe(true)
      expect(isValidKubernetesName('my-awesome-app')).toBe(true)
    })

    it('should accept single character names', () => {
      expect(isValidKubernetesName('a')).toBe(true)
      expect(isValidKubernetesName('1')).toBe(true)
    })

    it('should accept maximum length names', () => {
      const maxLengthName = 'a'.repeat(253)
      expect(isValidKubernetesName(maxLengthName)).toBe(true)
    })
  })

  describe('invalid names', () => {
    it('should reject empty or null names', () => {
      expect(isValidKubernetesName('')).toBe(false)
      expect(isValidKubernetesName(null as any)).toBe(false)
      expect(isValidKubernetesName(undefined as any)).toBe(false)
    })

    it('should reject names with uppercase letters', () => {
      expect(isValidKubernetesName('MyApp')).toBe(false)
      expect(isValidKubernetesName('MYAPP')).toBe(false)
      expect(isValidKubernetesName('myApp')).toBe(false)
    })

    it('should reject names with special characters', () => {
      expect(isValidKubernetesName('my_app')).toBe(false)
      expect(isValidKubernetesName('my.app')).toBe(false)
      expect(isValidKubernetesName('my@app')).toBe(false)
      expect(isValidKubernetesName('my app')).toBe(false)
      expect(isValidKubernetesName('my/app')).toBe(false)
    })

    it('should reject names starting with hyphen', () => {
      expect(isValidKubernetesName('-myapp')).toBe(false)
      expect(isValidKubernetesName('-123')).toBe(false)
    })

    it('should reject names ending with hyphen', () => {
      expect(isValidKubernetesName('myapp-')).toBe(false)
      expect(isValidKubernetesName('123-')).toBe(false)
    })

    it('should reject names exceeding maximum length', () => {
      const tooLongName = 'a'.repeat(254)
      expect(isValidKubernetesName(tooLongName)).toBe(false)
    })
  })
})

describe('getKubernetesNameError', () => {
  it('should return null for valid names', () => {
    expect(getKubernetesNameError('myapp')).toBe(null)
    expect(getKubernetesNameError('my-app')).toBe(null)
    expect(getKubernetesNameError('app123')).toBe(null)
  })

  it('should return error for empty names', () => {
    expect(getKubernetesNameError('')).toBe('Name is required')
    expect(getKubernetesNameError(null as any)).toBe('Name is required')
    expect(getKubernetesNameError(undefined as any)).toBe('Name is required')
  })

  it('should return error for names exceeding maximum length', () => {
    const tooLongName = 'a'.repeat(254)
    expect(getKubernetesNameError(tooLongName)).toBe('Name must be 253 characters or less')
  })

  it('should return error for names not starting with alphanumeric', () => {
    expect(getKubernetesNameError('-myapp')).toBe('Name must start with a lowercase letter or number')
    expect(getKubernetesNameError('_myapp')).toBe('Name must start with a lowercase letter or number')
    expect(getKubernetesNameError('.myapp')).toBe('Name must start with a lowercase letter or number')
  })

  it('should return error for names not ending with alphanumeric', () => {
    expect(getKubernetesNameError('myapp-')).toBe('Name must end with a lowercase letter or number')
    expect(getKubernetesNameError('myapp_')).toBe('Name must end with a lowercase letter or number')
    expect(getKubernetesNameError('myapp.')).toBe('Name must end with a lowercase letter or number')
  })

  it('should return error for names with invalid characters', () => {
    expect(getKubernetesNameError('my_app')).toBe('Name can only contain lowercase letters, numbers, and hyphens')
    expect(getKubernetesNameError('my.app')).toBe('Name can only contain lowercase letters, numbers, and hyphens')
    expect(getKubernetesNameError('MyApp')).toBe('Name can only contain lowercase letters, numbers, and hyphens')
    expect(getKubernetesNameError('my app')).toBe('Name can only contain lowercase letters, numbers, and hyphens')
  })

  it('should return the first applicable error', () => {
    // This name has multiple issues: starts with hyphen and contains uppercase
    expect(getKubernetesNameError('-MyApp')).toBe('Name must start with a lowercase letter or number')
    
    // This name starts valid but has uppercase
    expect(getKubernetesNameError('myApp')).toBe('Name can only contain lowercase letters, numbers, and hyphens')
  })
})