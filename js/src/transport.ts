import WebSocket from 'ws';
import { Config } from './types';
import { EventEmitter } from 'events';

export class Transport extends EventEmitter {
  private ws: WebSocket | null = null;
  private config: Config;
  private connected: boolean = false;
  private reconnectTimer: NodeJS.Timeout | null = null;

  constructor(config: Config) {
    super();
    this.config = config;
  }

  public async connect(): Promise<void> {
    if (this.connected) return;

    const url = `wss://relay.${this.config.region}.2gc.ru/v1/connect`;
    console.log(`Connecting to ${url}`);

    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(url, {
        headers: {
          'Authorization': `Bearer ${this.config.token}`
        },
        rejectUnauthorized: !this.config.insecureSkipVerify
      });

      this.ws.on('open', () => {
        console.log('Transport connected');
        this.connected = true;
        this.emit('open');
        resolve();
      });

      this.ws.on('message', (data: WebSocket.Data) => {
        this.handleMessage(data);
      });

      this.ws.on('close', () => {
        console.log('Transport closed');
        this.connected = false;
        this.emit('close');
        this.scheduleReconnect();
      });

      this.ws.on('error', (err) => {
        console.error('Transport error:', err);
        if (!this.connected) {
          reject(err);
        }
        this.emit('error', err);
      });
    });
  }

  public send(peerId: string, data: Uint8Array): void {
    if (!this.ws || !this.connected) {
      throw new Error('Transport not connected');
    }

    // Simple protocol: [peerId_len(1)][peerId][data]
    // In a real implementation, this would be a proper binary protocol or JSON
    // For Alpha, we'll assume a JSON envelope for simplicity in debugging
    const message = JSON.stringify({
      type: 'p2p',
      peer_id: peerId,
      payload: Buffer.from(data).toString('base64')
    });

    this.ws.send(message);
  }

  public close(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  private handleMessage(data: WebSocket.Data): void {
    try {
      const msg = JSON.parse(data.toString());
      if (msg.type === 'p2p') {
        const payload = Buffer.from(msg.payload, 'base64');
        this.emit('message', msg.peer_id, payload);
      }
    } catch (err) {
      console.error('Failed to parse message:', err);
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) return;
    console.log('Scheduling reconnect in 5s...');
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect().catch(err => console.error('Reconnect failed:', err));
    }, 5000);
  }
}
