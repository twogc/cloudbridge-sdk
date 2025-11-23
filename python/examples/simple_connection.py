import asyncio
import argparse
import logging
from cloudbridge.client import Client

async def main():
    parser = argparse.ArgumentParser(description="CloudBridge Python SDK Simple Connection Example")
    parser.add_argument("--token", required=True, help="CloudBridge API Token")
    parser.add_argument("--peer", required=True, help="Peer ID to connect to")
    parser.add_argument("--region", default="eu-central", help="CloudBridge Region")
    args = parser.parse_args()

    logging.basicConfig(level=logging.INFO)

    client = Client(token=args.token, region=args.region)

    try:
        print(f"Connecting to peer {args.peer}...")
        conn = await client.connect(args.peer)
        print(f"Connected to {conn.peer_id}!")

        # Send a message
        message = b"Hello from CloudBridge Python SDK!"
        print(f"Sending: {message}")
        await conn.write(message)

        # Read response (with timeout)
        try:
            response = await asyncio.wait_for(conn.read(), timeout=5.0)
            print(f"Received: {response}")
        except asyncio.TimeoutError:
            print("Timeout waiting for response")

    except Exception as e:
        print(f"Error: {e}")
    finally:
        await client.close()

if __name__ == "__main__":
    asyncio.run(main())
