import { useContext, useEffect, useState } from "react";
import { ApiContext } from "./ApiContext";
import { DoCommandResponse } from "./gen/flowstate/v1/server_pb";

export const AnnotationDetails = ({ id, rev }: { id: string; rev: string }) => {
  const [info, setInfo] = useState<DoCommandResponse | null>(null);
  const client = useContext(ApiContext);

  useEffect(() => {
    if (!client) return;
    client
      .doCommand({
        data: [{ id, rev: BigInt(rev) }],
        commands: [
          {
            command: {
              case: "getData",
              value: { dataRef: { id, rev: BigInt(rev) } },
            },
          },
        ],
      })
      .then(setInfo);
  }, [client]);

  if (!info) return "Loading...";

  return (
    <>
      {info.data.map((d) => {
        if (d.binary) return <span>{d.b}</span>;

        try {
          return (
            <pre className="text-left">
              {JSON.stringify(JSON.parse(d.b), null, 2)}
            </pre>
          );
        } catch {
          return <span>{d.b}</span>;
        }
      })}
    </>
  );
};
