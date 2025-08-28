export type ChatType = "model" | "team" | "agent"

export const openFloatingChat = (name: string, type: ChatType, namespace: string) => {
  window.dispatchEvent(new CustomEvent('open-floating-chat', { detail: { name, type, namespace } }))
}

export const toggleFloatingChat = (name: string, type: ChatType, namespace: string) => {
  window.dispatchEvent(new CustomEvent('toggle-floating-chat', { detail: { name, type, namespace } }))
}