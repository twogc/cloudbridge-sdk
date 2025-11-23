import { Tunnel, TunnelConfig, Connection } from './types';
import * as net from 'net';

export class TCPTunnel implements Tunnel {
  public id: string;
  public config: TunnelConfig;
  private server: net.Server | null = null;
  private connectToPeer: (peerId: string) => Promise<Connection>;

  constructor(config: TunnelConfig, connectToPeer: (peerId: string) => Promise<Connection>) {
    this.config = config;
    this.connectToPeer = connectToPeer;
    this.id = `tunnel-${Math.random().toString(36).substr(2, 9)}`;
  }

  public async start(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.server = net.createServer(async (socket) => {
        console.log(`[Tunnel ${this.id}] New connection`);
        
        try {
          const connection = await this.connectToPeer(this.config.remotePeer);
          
          // Send handshake
          const handshake = JSON.stringify({ type: 'tunnel', port: this.config.remotePort });
          await connection.write(Buffer.from(handshake));
          
          // Forward data to remote peer via P2P connection
          socket.on('data', async (data) => {
            try {
              await connection.write(data);
            } catch (err) {
              console.error(`[Tunnel ${this.id}] Write error:`, err);
              socket.destroy();
            }
          });

          // Receive data from P2P connection
          connection.on('data', (data) => {
            socket.write(data);
          });

          connection.on('close', () => {
             socket.end();
          });

          socket.on('end', () => {
            console.log(`[Tunnel ${this.id}] Connection closed`);
            connection.close();
          });
          
          socket.on('error', (err) => {
             console.error(`[Tunnel ${this.id}] Socket error:`, err);
             connection.close();
          });

        } catch (err) {
          console.error(`[Tunnel ${this.id}] Failed to establish tunnel:`, err);
          socket.destroy();
        }
      });

      this.server.listen(this.config.localPort, () => {
        console.log(`[Tunnel ${this.id}] Listening on port ${this.config.localPort}`);
        resolve();
      });

      this.server.on('error', (err) => {
        reject(err);
      });
    });
  }

  public async stop(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (!this.server) {
        resolve();
        return;
      }

      this.server.close((err) => {
        if (err) reject(err);
        else resolve();
      });
    });
  }
}
