export interface Config {
  token: string;
  region?: string;
  timeout?: number;
  logLevel?: 'debug' | 'info' | 'warn' | 'error';
  insecureSkipVerify?: boolean;
}

export interface Peer {
  id: string;
  status: 'connected' | 'disconnected';
  latency?: number;
}

export interface Service {
  id: string;
  name: string;
  port: number;
  tags?: string[];
  metadata?: Record<string, string>;
}

export interface Connection {
  peerId: string;
  close(): Promise<void>;
  read(): Promise<Uint8Array>;
  write(data: Uint8Array): Promise<void>;
  on(event: 'data', listener: (data: Uint8Array) => void): void;
  on(event: 'close', listener: () => void): void;
  on(event: 'error', listener: (err: Error) => void): void;
}

export interface TunnelConfig {
  localPort: number;
  remotePeer: string;
  remotePort: number;
  protocol?: 'tcp' | 'udp' | 'quic';
}

export interface Tunnel {
  id: string;
  config: TunnelConfig;
  start(): Promise<void>;
  stop(): Promise<void>;
}

export interface Mesh {
  join(): Promise<void>;
  leave(): Promise<void>;
  broadcast(data: any): Promise<void>;
  send(peerId: string, data: any): Promise<void>;
  peers(): Promise<Peer[]>;
}
