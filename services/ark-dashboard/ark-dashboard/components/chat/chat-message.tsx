import type { ChatMessageData } from "@/lib/types/chat"
import { useMarkdownProcessor } from "@/lib/hooks/use-markdown-processor"

interface ChatMessageProps extends Pick<ChatMessageData, "role" | "content" | "status"> {
  className?: string
  viewMode?: 'text' | 'markdown'
}

export function ChatMessage({ role, content, status, className, viewMode = 'text' }: ChatMessageProps) {
  const isUser = role === "user"
  const isFailed = status === "failed"
  const markdownContent = useMarkdownProcessor(content)
  
  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"} ${className || ""}`}>
      <div
        className={`max-w-[80%] rounded-lg px-3 py-2 ${
          isUser
            ? "bg-primary text-primary-foreground" 
            : isFailed
            ? "bg-destructive/10 text-destructive"
            : "bg-muted"
        }`}
      >
        {viewMode === 'markdown' ? (
          <div className="text-sm">
            {markdownContent}
          </div>
        ) : (
          <pre className="text-sm whitespace-pre-wrap font-mono bg-transparent p-0 m-0 border-0">{content}</pre>
        )}
      </div>
    </div>
  )
}