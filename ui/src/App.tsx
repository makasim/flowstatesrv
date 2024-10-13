import { useState } from "react";
import "./App.css";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import { useLocalStorage } from "./hooks/useLocalStorage";
import { StatesPage } from "./StatesPage";

export default function App() {
  const [apiURL, setApiURL] = useState("");
  const [choosenServer, setChoosenServer] = useState("");
  const [servers, setServers] = useLocalStorage<string[]>("servers", []);

  if (apiURL) {
    return <StatesPage apiUrl={apiURL} />;
  }

  return (
    <div className="container mx-auto py-10">
      <form
        className="flex w-full max-w-sm items-center space-x-2"
        onSubmit={(e) => {
          e.preventDefault();
          if (!servers.includes(choosenServer)) {
            setServers([choosenServer, ...servers]);
          }
          setApiURL(choosenServer);
        }}
      >
        <Input
          placeholder="https://flowstate.makasim.com"
          value={choosenServer}
          onChange={(e) => setChoosenServer(e.target.value)}
        />
        <Button type="submit">Subscribe</Button>
      </form>

      <div className="mt-10">
        <h1 className="text-2xl font-bold">Servers</h1>

        <ul>
          {servers.map((server) => (
            <li
              key={server}
              className="flex w-full max-w-sm items-center space-x-2"
            >
              <Button
                onClick={() => setServers(servers.filter((s) => s !== server))}
              >
                x
              </Button>
              <a
                href={server}
                onClick={(e) => {
                  e.preventDefault();
                  setChoosenServer(server);
                }}
              >
                {server}
              </a>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
