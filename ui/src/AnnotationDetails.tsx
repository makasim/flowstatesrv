import { useContext, useEffect, useState } from "react";
import { ApiContext } from "./ApiContext";

export const AnnotationDetails = ({ id, rev }: { id: string; rev: string }) => {
  const [info, setInfo] = useState<object | null>(null);
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

  return <pre className="text-left">{JSON.stringify(info, null, 2)}</pre>;
};
