"use client"

import { useState } from "react"
import { ChevronDown, ChevronRight } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"

export interface AdditionalField {
  key: string
  value: unknown
  label?: string
}

interface InfoDialogProps<T extends object> {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  data: T & {
    labels?: unknown
    annotations?: unknown
    namespace?: unknown
  }
  additionalFields?: AdditionalField[]
}

export function InfoDialog<T extends object>({ open, onOpenChange, title, data, additionalFields = [] }: InfoDialogProps<T>) {
  const [metadataExpanded, setMetadataExpanded] = useState(false)
  // Recursively parse JSON strings
  const recursiveJsonParse = (obj: unknown): unknown => {
    if (typeof obj === 'string') {
      // Try to parse as JSON
      try {
        const parsed = JSON.parse(obj)
        // Recursively parse the result
        return recursiveJsonParse(parsed)
      } catch {
        // Not valid JSON, return as is
        return obj
      }
    } else if (Array.isArray(obj)) {
      // Recursively parse array elements
      return obj.map(item => recursiveJsonParse(item))
    } else if (obj !== null && typeof obj === 'object') {
      // Recursively parse object values
      const result: Record<string, unknown> = {}
      for (const [key, value] of Object.entries(obj)) {
        result[key] = recursiveJsonParse(value)
      }
      return result
    }
    // Return primitives as is
    return obj
  }

  const formatValue = (value: unknown, key: string): React.ReactNode => {
    if (value === null || value === undefined) {
      return '-'
    }
    
    // Special handling for annotations field - parse JSON string recursively
    if (key === 'annotations') {
      const parsed = recursiveJsonParse(value)
      return (
        <pre className="text-xs bg-muted p-4 rounded overflow-auto max-h-[400px] min-h-[100px]">
          {JSON.stringify(parsed, null, 2)}
        </pre>
      )
    }
    
    if (typeof value === 'object') {
      return (
        <pre className="text-xs bg-muted p-2 rounded overflow-auto max-h-64">
          {JSON.stringify(value, null, 2)}
        </pre>
      )
    }
    
    return String(value)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-6xl">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>Detailed information</DialogDescription>
        </DialogHeader>
        <div className="mt-4 space-y-4 max-h-[70vh] overflow-y-auto">
          <>
            {additionalFields.map((field) => (
            <div key={field.key} className="grid grid-cols-3 gap-4 py-2 border-b">
              <div className="font-medium text-sm">
                {field.label || field.key.replace(/([A-Z])/g, ' $1').trim()}
              </div>
              <div className="col-span-2 text-sm text-muted-foreground">
                {formatValue(field.value, field.key)}
              </div>
            </div>
          ))}
          
          {Object.entries(data)
            .filter(([key]) => key !== 'namespace' && key !== 'labels' && key !== 'annotations')
            .map(([key, value]) => (
              <div key={key} className="grid grid-cols-3 gap-4 py-2 border-b">
                <div className="font-medium text-sm capitalize">
                  {key.replace(/([A-Z])/g, ' $1').trim()}
                </div>
                <div className="col-span-2 text-sm text-muted-foreground">
                  {formatValue(value, key)}
                </div>
              </div>
            ))}
          
          {!!(data.labels || data.annotations) && (
            <div className="border rounded-lg">
              <button
                onClick={() => setMetadataExpanded(!metadataExpanded)}
                className="w-full flex items-center justify-between p-3 hover:bg-muted/50 transition-colors"
              >
                <span className="font-medium text-sm">Metadata (Labels & Annotations)</span>
                {metadataExpanded ? (
                  <ChevronDown className="h-4 w-4" />
                ) : (
                  <ChevronRight className="h-4 w-4" />
                )}
              </button>
              
              {metadataExpanded && (
                <div className="p-3 pt-0 space-y-4">
                  {!!data.labels && (
                    <div>
                      <div className="font-medium text-sm mb-2">Labels</div>
                      <div className="text-sm text-muted-foreground">
                        {formatValue(data.labels, 'labels')}
                      </div>
                    </div>
                  )}
                  
                  {!!data.annotations && (
                    <div>
                      <div className="font-medium text-sm mb-2">Annotations</div>
                      <div className="text-sm text-muted-foreground">
                        {formatValue(data.annotations, 'annotations')}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          )}
          </>
        </div>
      </DialogContent>
    </Dialog>
  )
}