"use client";

import type React from "react";
import { useState, useEffect } from "react";
import { toast } from "@/components/ui/use-toast";
import { A2AServersService, type A2AServer } from "@/lib/services";
import { A2AServerCard } from "@/components/cards";
import { useDelayedLoading } from "@/lib/hooks";
import { InfoDialog } from "@/components/dialogs/info-dialog";

interface A2AServersSectionProps {
  namespace: string;
}

export const A2AServersSection: React.FC<A2AServersSectionProps> = ({
  namespace
}) => {
  const [a2aServers, setA2AServers] = useState<A2AServer[]>([]);
  const [loading, setLoading] = useState(true);
  const showLoading = useDelayedLoading(loading);
  const [selectedServer, setSelectedServer] = useState<A2AServer | null>(null);
  const [infoDialogOpen, setInfoDialogOpen] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const data = await A2AServersService.getAll(namespace);
        setA2AServers(data);
      } catch (error) {
        console.error("Failed to load A2A servers:", error);
        toast({
          variant: "destructive",
          title: "Failed to Load A2A Servers",
          description:
            error instanceof Error
              ? error.message
              : "An unexpected error occurred"
        });
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [namespace]);

  const handleInfo = (server: A2AServer) => {
    setSelectedServer(server);
    setInfoDialogOpen(true);
  };
  if (showLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center py-8">Loading...</div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      <main className="flex-1 overflow-auto p-6">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 pb-6">
          {a2aServers.map((server) => (
            <A2AServerCard
              key={server.name || server.id}
              a2aServer={server}
              onInfo={handleInfo}
              namespace={namespace}
            />
          ))}
        </div>
      </main>

      {selectedServer && (
        <InfoDialog
          open={infoDialogOpen}
          onOpenChange={setInfoDialogOpen}
          title={`A2A Server: ${selectedServer.name || "Unnamed"}`}
          data={selectedServer}
        />
      )}
    </div>
  );
};
