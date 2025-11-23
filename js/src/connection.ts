import { Connection } from './types';
import { EventEmitter } from 'events';
import { Transport } from './transport';

export class P2PConnection extends EventEmitter implements Connection {
  public peerId: string;
  private transport: Transport;
  private connected: boolean = false;

  constructor(peerId: string, transport: Transport) {
    super();
    this.peerId = peerId;
    this.transport = transport;
    this.connected = true;

    // Listen for incoming messages from transport
    this.transport.on('message', (senderId: string, data: Uint8Array) => {
      if (senderId === this.peerId) {
        this.emit('data', data);
      }
    });
  }

  public async close(): Promise<void> {
    if (!this.connected) return;
    this.connected = false;
    this.emit('close');
  }

  public async read(): Promise<Uint8Array> {
    if (!this.connected) throw new Error('Connection closed');
    // For now, we rely on the 'data' event listener pattern which is more idiomatic in JS/Node
    // This method is kept for interface compatibility but might be deprecated or changed to AsyncIterator
    return new Promise((resolve) => {
      this.once('data', (data) => resolve(data));
    });
  }

  public async write(data: Uint8Array): Promise<void> {
    if (!this.connected) throw new Error('Connection closed');
    this.transport.send(this.peerId, data);
  }
}
