import { Client } from '../../src';

async function main() {
  const args = process.argv.slice(2);
  const token = args.find(arg => arg.startsWith('--token='))?.split('=')[1];
  const peerId = args.find(arg => arg.startsWith('--peer='))?.split('=')[1];
  const region = args.find(arg => arg.startsWith('--region='))?.split('=')[1] || 'eu-central';

  if (!token || !peerId) {
    console.error('Usage: npm start -- --token=<token> --peer=<peer_id> [--region=<region>]');
    process.exit(1);
  }

  console.log(`Initializing client for region ${region}...`);
  const client = new Client({
    token,
    region,
    logLevel: 'info'
  });

  try {
    console.log(`Connecting to peer ${peerId}...`);
    const connection = await client.connect(peerId);
    
    console.log(`Connected to ${connection.peerId}!`);

    // Handle incoming data
    connection.on('data', (data) => {
      console.log(`Received: ${Buffer.from(data).toString()}`);
    });

    connection.on('close', () => {
      console.log('Connection closed');
      process.exit(0);
    });

    // Send a message
    const message = 'Hello from CloudBridge JS SDK!';
    console.log(`Sending: ${message}`);
    await connection.write(Buffer.from(message));

    // Keep alive for a bit to receive response
    setTimeout(() => {
      console.log('Closing connection...');
      connection.close();
    }, 5000);

  } catch (err) {
    console.error('Error:', err);
    client.close();
    process.exit(1);
  }
}

main();
