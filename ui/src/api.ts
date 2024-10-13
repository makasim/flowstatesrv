import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ServerService } from "./gen/flowstate/v1/server_connect";

export function createApiClient(baseUrl: string) {
  const transport = createConnectTransport({ baseUrl });

  return createClient(ServerService, transport);
}
