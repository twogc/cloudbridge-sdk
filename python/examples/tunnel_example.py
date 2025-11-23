import argparse
import asyncio
import logging
import sys
import os

# Add parent directory to path to import cloudbridge
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from cloudbridge import Client

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

async def run_server(token: str, region: str, port: int):
    client = Client(token, region)
    
    # Start local service
    async def handle_echo(reader, writer):
        writer.write(f"Hello from CloudBridge Tunnel Service running on port {port}!".encode())
        await writer.drain()
        writer.close()

    server = await asyncio.start_server(handle_echo, '127.0.0.1', port)
    logging.info(f"Starting local TCP service on port {port}...")
    
    # Connect to relay
    await client.transport.connect()
    
    logging.info("Waiting for incoming tunnel connections...")
    
    async with server:
        await asyncio.gather(
            server.serve_forever(),
            client.serve()
        )

async def run_client(token: str, region: str, peer_id: str, local_port: int, remote_port: int):
    client = Client(token, region)
    
    logging.info(f"Creating tunnel: localhost:{local_port} -> {peer_id}:{remote_port}")
    
    try:
        await client.create_tunnel({
            "remote_peer": peer_id,
            "remote_port": remote_port,
            "local_port": local_port
        })
    except KeyboardInterrupt:
        logging.info("Stopping tunnel...")
    finally:
        await client.close()

def main():
    parser = argparse.ArgumentParser(description="CloudBridge Tunnel Example")
    parser.add_argument("--token", required=True, help="API Token")
    parser.add_argument("--region", default="eu-central", help="Region")
    parser.add_argument("--mode", default="server", choices=["server", "client"], help="Mode")
    parser.add_argument("--peer", help="Peer ID (client mode)")
    parser.add_argument("--local-port", type=int, default=8080, help="Local port")
    parser.add_argument("--remote-port", type=int, default=8080, help="Remote port (client mode)")
    
    args = parser.parse_args()
    
    if args.mode == "client" and not args.peer:
        print("Error: --peer is required for client mode")
        sys.exit(1)
        
    if args.mode == "server":
        asyncio.run(run_server(args.token, args.region, args.local_port))
    else:
        asyncio.run(run_client(args.token, args.region, args.peer, args.local_port, args.remote_port))

if __name__ == "__main__":
    main()
