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
import { ApiContext } from "./ApiContext";
import { AnnotationDetails } from "./AnnotationDetails";

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
    accessorKey: "annotations",
    header: "Data",
    cell: ({ row }) =>
      Object.values(row.original.annotations)
        .filter((x) => x.startsWith("data:"))
        .map((x) => x.slice(5).split(":"))
        .map(([id, rev]) => (
          <Dialog modal key={`${id}:${rev}`}>
            <DialogTrigger className="text-slate">
              <svg
                className="w-6 h-6 text-white"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  stroke-width="2"
                  d="M21 12c0 1.2-4.03 6-9 6s-9-4.8-9-6c0-1.2 4.03-6 9-6s9 4.8 9 6Z"
                />
                <path
                  stroke="currentColor"
                  stroke-width="2"
                  d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                />
              </svg>
            </DialogTrigger>

            <DialogContent>
              <DialogTitle className="pb-4 sticky top-0 bg-background">
                {id}:{rev}
              </DialogTitle>
              <DialogDescription>
                <AnnotationDetails id={id} rev={rev} />
              </DialogDescription>
            </DialogContent>
          </Dialog>
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

type ApiClient = ReturnType<typeof createApiClient>;

export const StatesPage = () => {
  const [states, setStates] = useState<State[]>([]);
  const client = React.useContext(ApiContext);

  useEffect(() => {
    if (!client) return;

    const abortController = new AbortController();

    listenToStates(client, abortController.signal).catch((error) =>
      console.log("Listening error", error)
    );

    return () => abortController.abort();
  }, [client]);

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
