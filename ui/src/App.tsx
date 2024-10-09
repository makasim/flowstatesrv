import "./App.css";
import { DataTable } from "./states/data-table";
import { columns } from "./states/columns";
import { useEffect, useState } from "react";
import { client } from "./api";
import { State } from "./gen/flowstate/v1/state_pb";

export default function App() {
  const [states, setStates] = useState<State[]>([]);

  useEffect(() => {
    const abortController = new AbortController();

    try {
      listenToStates(abortController.signal);
    } catch (error) {
      console.log("Game not found", error);
    }

    return () => abortController.abort();
  }, []);

  async function listenToStates(signal: AbortSignal) {
    for await (const res of client.watchStates({}, { signal })) {
      console.log(res);

      if (res.ping) continue;

      setStates((current) => (res.state ? [res.state, ...current] : current));
    }
  }

  function formatTransition({ from, to }: { from: string; to: string }) {
    const set = new Set();
    if (from) set.add(from);
    if (to) set.add(to);
    return [...set].join(" -> ");
  }

  const data = states.map((state) => ({
    id: state.id + "#" + state.rev,
    stateId: state.id,
    rev: state.rev,
    transition: state.transition ? formatTransition(state.transition) : "",
  }));

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data} />
    </div>
  );
}
