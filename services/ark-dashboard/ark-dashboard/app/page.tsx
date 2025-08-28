"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

export default function Page() {
  const router = useRouter()
  
  useEffect(() => {
    router.replace("/agents")
  }, [router])
  
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="text-center">
        <p className="text-lg">Redirecting to dashboard...</p>
      </div>
    </div>
  )
}
