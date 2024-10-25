import "./App.css";
import { DataTable } from "./components/data-table";
import React, { useEffect, useState } from "react";
import { State } from "./gen/flowstate/v1/state_pb";

import { ColumnDef } from "@tanstack/react-table";
import { createApiClient } from "./api";
import { Badge } from "./components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "./components/ui/dialog";

type StateData = {
  id: string;
  stateId: string;
  rev: bigint;
  transition: string;
  annotations: Record<string, string>;
  labels: Record<string, string>;
  state: State;
};

const columns: ColumnDef<StateData>[] = [
  { accessorKey: "stateId", header: "ID" },
  { accessorKey: "rev", header: "REV" },
  { accessorKey: "transition", header: "Transtion" },
  {
    accessorKey: "annotations",
    header: "Annotations",
    cell: ({ row }) =>
      Object.entries(row.original.annotations).map(([key, value]) => (
        <div key={key} className="text-left">
          <Badge variant="outline">
            <span className="text-green-700">{key}:&nbsp;</span>
            <span className="text-purple-700">{String(value)}</span>
          </Badge>
        </div>
      )),
  },
  {
    accessorKey: "labels",
    header: "Labels",
    cell: ({ row }) =>
      Object.entries(row.original.labels).map(([key, value]) => (
        <div key={key} className="text-left">
          <Badge variant="outline">
            <span className="text-green-700">{key}:&nbsp;</span>
            <span className="text-purple-700">{String(value)}</span>
          </Badge>
        </div>
      )),
  },
  {
    accessorKey: "state",
    header: "State",
    cell: ({ row }) => (
      <Dialog modal>
        <DialogTrigger className="text-slate-100">Show State</DialogTrigger>
        <DialogContent>
          <DialogTitle>State: {row.original.state.id}</DialogTitle>
          <DialogDescription>
            <pre className="text-left">
              {JSON.stringify(row.original.state.toJson(), null, 2)}
            </pre>
          </DialogDescription>
        </DialogContent>
      </Dialog>
    ),
  },
];

type Props = { apiUrl: string };
type ApiClient = ReturnType<typeof createApiClient>;

export const StatesPage: React.FC<Props> = ({ apiUrl }) => {
  const [states, setStates] = useState<State[]>([]);

  useEffect(() => {
    if (!apiUrl) return;

    const client = createApiClient(apiUrl);
    const abortController = new AbortController();

    listenToStates(client, abortController.signal).catch((error) =>
      console.log("Listening error", error)
    );

    return () => abortController.abort();
  }, [apiUrl]);

  async function listenToStates(client: ApiClient, signal: AbortSignal) {
    for await (const res of client.watchStates({}, { signal })) {
      console.log(res);
      if (res.ping) continue;
      setStates((v) => (res.state ? [res.state, ...v] : v));
    }
  }

  function formatTransition({ from, to }: { from: string; to: string }) {
    return from && from !== to ? `${from} -> ${to}` : to;
  }

  const data = states.map((state) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const { id, rev, transition, annotations, labels } = state.toJson() as any;
    return {
      id: `${id}#${rev}`,
      stateId: id,
      rev,
      transition: transition ? formatTransition(transition) : "",
      annotations,
      labels,
      state,
    };
  });

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data} />
    </div>
  );
};
