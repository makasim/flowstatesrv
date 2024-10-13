import "./App.css";
import { DataTable } from "./components/data-table";
import React, { useEffect, useState } from "react";
import { State } from "./gen/flowstate/v1/state_pb";

import { ColumnDef } from "@tanstack/react-table";
import { createApiClient } from "./api";

export type StateData = {
  id: string;
  stateId: string;
  rev: bigint;
  transition: string;
};

export const columns: ColumnDef<StateData>[] = [
  { accessorKey: "stateId", header: "ID" },
  { accessorKey: "rev", header: "REV" },
  { accessorKey: "transition", header: "Transtion" },
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

  async function listenToStates(
    client: ApiClient,
    signal: AbortSignal
  ) {
    for await (const res of client.watchStates({}, { signal })) {
      console.log(res);
      if (res.ping) continue;
      setStates((v) => (res.state ? [res.state, ...v] : v));
    }
  }

  function formatTransition({ from, to }: { from: string; to: string }) {
    return from && from !== to ? `${from} -> ${to}` : to;
  }

  const data = states.map(({ id, rev, transition }) => ({
    id: `${id}#${rev}`,
    stateId: id,
    rev,
    transition: transition ? formatTransition(transition) : "",
  }));

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data} />
    </div>
  );
};
