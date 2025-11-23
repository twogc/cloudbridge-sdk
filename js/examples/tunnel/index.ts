import { Client } from '../../src';
import * as http from 'http';

async function main() {
  const args = process.argv.slice(2);
  const mode = args.find(arg => arg.startsWith('--mode='))?.split('=')[1] || 'server';
  const token = args.find(arg => arg.startsWith('--token='))?.split('=')[1];
  const peerId = args.find(arg => arg.startsWith('--peer='))?.split('=')[1];
  const localPort = parseInt(args.find(arg => arg.startsWith('--local-port='))?.split('=')[1] || '8080');
  const remotePort = parseInt(args.find(arg => arg.startsWith('--remote-port='))?.split('=')[1] || '8080');
  const region = args.find(arg => arg.startsWith('--region='))?.split('=')[1] || 'eu-central';

  if (!token) {
    console.error('Usage: npm start -- --token=<token> [--mode=server|client] ...');
    process.exit(1);
  }

  const client = new Client({
    token,
    region,
    logLevel: 'info'
  });

  try {
    // Ensure transport is connected
    // In server mode, serve() handles connection implicitly via transport events?
    // Actually serve() just waits. We need to ensure transport connects first.
    // Client.serve() doesn't explicitly call connect() in my implementation, 
    // but Transport connects lazily or we should call it.
    // Let's call a dummy connect or just expose connect() on Client?
    // Client.connect(peerId) connects. But for server mode we just want to be online.
    // I should probably expose client.connectTransport() or similar, or just rely on the fact 
    // that we need to be online.
    // For now, let's hack it by calling a private method or adding a public one.
    // Actually, let's add `await (client as any).transport.connect();` for now or better,
    // update Client to have a `start()` method.
    // Or just use `client.serve()` and update `serve` to ensure connection.
    
    // Let's assume serve() should ensure connection. I'll update client.ts in a bit if needed.
    // But wait, `serve` just waits for close.
    // I'll add `await (client as any).transport.connect();` here for safety as I can't easily change Client interface right now without another step.
    // Actually, I can just cast to any.
    await (client as any).transport.connect();

    if (mode === 'server') {
      // Start local service
      const server = http.createServer((req, res) => {
        res.writeHead(200);
        res.end(`Hello from CloudBridge Tunnel Service running on port ${localPort}!`);
      });
      server.listen(localPort, () => {
        console.log(`Starting local HTTP service on port ${localPort}...`);
      });

      console.log('Waiting for incoming tunnel connections...');
      await client.serve();
    } else {
      if (!peerId) {
        console.error('Peer ID is required for client mode');
        process.exit(1);
      }

      console.log(`Creating tunnel: localhost:${localPort} -> ${peerId}:${remotePort}`);
      const tunnel = await client.createTunnel({
        localPort,
        remotePeer: peerId,
        remotePort,
        protocol: 'tcp'
      });

      console.log(`Tunnel established! Access the service at http://localhost:${localPort}`);
      
      // Keep alive
      await new Promise(() => {});
    }
  } catch (err) {
    console.error('Error:', err);
    client.close();
    process.exit(1);
  }
}

main();
