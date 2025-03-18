// @generated by protoc-gen-es v1.10.0 with parameter "target=ts"
// @generated from file flowstate/v1/server.proto (package flowstate.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Data, StateContext } from "./state_pb.js";
import { AnyCommand, AnyResult } from "./commands_pb.js";

/**
 * @generated from message flowstate.v1.ErrorConflict
 */
export class ErrorConflict extends Message<ErrorConflict> {
  /**
   * @generated from field: repeated string committable_state_ids = 2;
   */
  committableStateIds: string[] = [];

  constructor(data?: PartialMessage<ErrorConflict>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "flowstate.v1.ErrorConflict";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 2, name: "committable_state_ids", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ErrorConflict {
    return new ErrorConflict().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ErrorConflict {
    return new ErrorConflict().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ErrorConflict {
    return new ErrorConflict().fromJsonString(jsonString, options);
  }

  static equals(a: ErrorConflict | PlainMessage<ErrorConflict> | undefined, b: ErrorConflict | PlainMessage<ErrorConflict> | undefined): boolean {
    return proto3.util.equals(ErrorConflict, a, b);
  }
}

/**
 * @generated from message flowstate.v1.DoCommandRequest
 */
export class DoCommandRequest extends Message<DoCommandRequest> {
  /**
   * @generated from field: repeated flowstate.v1.StateContext state_contexts = 1;
   */
  stateContexts: StateContext[] = [];

  /**
   * @generated from field: repeated flowstate.v1.Data data = 2;
   */
  data: Data[] = [];

  /**
   * @generated from field: repeated flowstate.v1.AnyCommand commands = 3;
   */
  commands: AnyCommand[] = [];

  constructor(data?: PartialMessage<DoCommandRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "flowstate.v1.DoCommandRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "state_contexts", kind: "message", T: StateContext, repeated: true },
    { no: 2, name: "data", kind: "message", T: Data, repeated: true },
    { no: 3, name: "commands", kind: "message", T: AnyCommand, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DoCommandRequest {
    return new DoCommandRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DoCommandRequest {
    return new DoCommandRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DoCommandRequest {
    return new DoCommandRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DoCommandRequest | PlainMessage<DoCommandRequest> | undefined, b: DoCommandRequest | PlainMessage<DoCommandRequest> | undefined): boolean {
    return proto3.util.equals(DoCommandRequest, a, b);
  }
}

/**
 * @generated from message flowstate.v1.DoCommandResponse
 */
export class DoCommandResponse extends Message<DoCommandResponse> {
  /**
   * @generated from field: repeated flowstate.v1.StateContext state_contexts = 1;
   */
  stateContexts: StateContext[] = [];

  /**
   * @generated from field: repeated flowstate.v1.Data data = 2;
   */
  data: Data[] = [];

  /**
   * @generated from field: repeated flowstate.v1.AnyResult results = 3;
   */
  results: AnyResult[] = [];

  constructor(data?: PartialMessage<DoCommandResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "flowstate.v1.DoCommandResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "state_contexts", kind: "message", T: StateContext, repeated: true },
    { no: 2, name: "data", kind: "message", T: Data, repeated: true },
    { no: 3, name: "results", kind: "message", T: AnyResult, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DoCommandResponse {
    return new DoCommandResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DoCommandResponse {
    return new DoCommandResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DoCommandResponse {
    return new DoCommandResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DoCommandResponse | PlainMessage<DoCommandResponse> | undefined, b: DoCommandResponse | PlainMessage<DoCommandResponse> | undefined): boolean {
    return proto3.util.equals(DoCommandResponse, a, b);
  }
}

/**
 * @generated from message flowstate.v1.RegisterFlowRequest
 */
export class RegisterFlowRequest extends Message<RegisterFlowRequest> {
  /**
   * @generated from field: string flow_id = 1;
   */
  flowId = "";

  /**
   * @generated from field: string http_host = 2;
   */
  httpHost = "";

  constructor(data?: PartialMessage<RegisterFlowRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "flowstate.v1.RegisterFlowRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "flow_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "http_host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RegisterFlowRequest {
    return new RegisterFlowRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RegisterFlowRequest {
    return new RegisterFlowRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RegisterFlowRequest {
    return new RegisterFlowRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RegisterFlowRequest | PlainMessage<RegisterFlowRequest> | undefined, b: RegisterFlowRequest | PlainMessage<RegisterFlowRequest> | undefined): boolean {
    return proto3.util.equals(RegisterFlowRequest, a, b);
  }
}

/**
 * @generated from message flowstate.v1.RegisterFlowResponse
 */
export class RegisterFlowResponse extends Message<RegisterFlowResponse> {
  constructor(data?: PartialMessage<RegisterFlowResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "flowstate.v1.RegisterFlowResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RegisterFlowResponse {
    return new RegisterFlowResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RegisterFlowResponse {
    return new RegisterFlowResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RegisterFlowResponse {
    return new RegisterFlowResponse().fromJsonString(jsonString, options);
  }

  static equals(a: RegisterFlowResponse | PlainMessage<RegisterFlowResponse> | undefined, b: RegisterFlowResponse | PlainMessage<RegisterFlowResponse> | undefined): boolean {
    return proto3.util.equals(RegisterFlowResponse, a, b);
  }
}

