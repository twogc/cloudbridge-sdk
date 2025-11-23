import { Config, Mesh, Peer, Service, Connection, Tunnel, TunnelConfig } from './types';
import { mergeConfig, validateConfig } from './config';
import { P2PConnection } from './connection';
import { TCPTunnel } from './tunnel';
import { Transport } from './transport';

export class Client {
  private config: Config;
  private transport: Transport;
  private connections: Map<string, Connection> = new Map();
  private tunnels: Map<string, Tunnel> = new Map();

  constructor(config: Config) {
    validateConfig(config);
    this.config = mergeConfig(config);
    this.transport = new Transport(this.config);

    this.transport.on('message', (peerId, data) => {
      if (!this.connections.has(peerId)) {
        const connection = new P2PConnection(peerId, this.transport);
        this.connections.set(peerId, connection);
        
        connection.on('close', () => {
          this.connections.delete(peerId);
        });

        this.handleIncomingConnection(connection, data);
      }
    });
  }

  public async connect(peerId: string): Promise<Connection> {
    console.log(`Connecting to peer ${peerId} in region ${this.config.region}`);
    
    // Ensure transport is connected
    await this.transport.connect();

    // Check if connection already exists
    if (this.connections.has(peerId)) {
      return this.connections.get(peerId)!;
    }

    const connection = new P2PConnection(peerId, this.transport);
    this.connections.set(peerId, connection);
    
    connection.on('close', () => {
      this.connections.delete(peerId);
    });

    return connection;
  }

  public async createTunnel(config: TunnelConfig): Promise<Tunnel> {
    console.log(`Creating tunnel to ${config.remotePeer}:${config.remotePort}`);
    
    const tunnel = new TCPTunnel(config, (peerId) => this.connect(peerId));
    await tunnel.start();
    
    this.tunnels.set(tunnel.id, tunnel);
    return tunnel;
  }

  public async joinMesh(networkName: string): Promise<Mesh> {
    console.log(`Joining mesh network ${networkName}`);
    
    return {
      join: async () => console.log('Joined mesh'),
      leave: async () => console.log('Left mesh'),
      broadcast: async (data: any) => console.log('Broadcasting:', data),
      send: async (peerId: string, data: any) => console.log(`Sending to ${peerId}:`, data),
      peers: async () => [],
    };
  }

  public async registerService(service: Service): Promise<void> {
    console.log(`Registering service ${service.name}`);
    // TODO: Implement service registration
  }

  public async discoverServices(serviceName: string): Promise<Service[]> {
    console.log(`Discovering services with name ${serviceName}`);
    return [];
  }

  public async close(): Promise<void> {
    this.transport.close();
    // Close all connections and tunnels
    for (const conn of this.connections.values()) {
      await conn.close();
    }
    for (const tunnel of this.tunnels.values()) {
      await tunnel.stop();
    }
  }
  public async serve(): Promise<void> {
    return new Promise((resolve) => {
      this.transport.on('close', resolve);
    });
  }

  private async handleIncomingConnection(connection: Connection, initialData: Uint8Array) {
    try {
      const jsonString = Buffer.from(initialData).toString();
      const handshake = JSON.parse(jsonString);
      
      if (handshake.type === 'tunnel' && handshake.port) {
        console.log(`Incoming tunnel request for port ${handshake.port}`);
        const net = await import('net');
        const localSocket = net.connect(handshake.port, 'localhost');
        
        localSocket.on('connect', () => {
          connection.on('data', (data) => localSocket.write(data));
          localSocket.on('data', (data) => connection.write(data));
        });
        
        localSocket.on('error', (err) => {
          console.error('Local socket error:', err);
          connection.close();
        });
        
        connection.on('close', () => localSocket.end());
        localSocket.on('close', () => connection.close());
      } else {
        (connection as P2PConnection).emit('data', initialData);
      }
    } catch (err) {
      (connection as P2PConnection).emit('data', initialData);
    }
  }
}
