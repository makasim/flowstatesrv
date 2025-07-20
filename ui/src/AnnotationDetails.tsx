import { useContext, useEffect, useState } from "react";
import { ApiContext } from "./ApiContext";
import { Command, GetDataCommand, Data, DataRef } from "./gen/flowstate/v1/messages_pb";

export const AnnotationDetails = ({ id, rev }: { id: string; rev: string }) => {
  const [info, setInfo] = useState<Command | null>(null);
  const client = useContext(ApiContext);

  useEffect(() => {
    if (!client) return;
    const command = new Command({
      datas: [
        new Data({ id, rev: BigInt(rev) })
      ],
      getData: new GetDataCommand({
        dataRef: new DataRef({ id, rev: BigInt(rev) })
      })
    });
    
    client.do(command).then(setInfo);
  }, [client]);

  if (!info) return "Loading...";

  return (
    <>
      {info.datas.map((d) => {
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
