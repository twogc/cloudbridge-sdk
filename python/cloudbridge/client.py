import logging
from typing import Optional, Dict, Any
from .transport import Transport
from .connection import Connection

class Client:
    """
    CloudBridge SDK Client for Python.
    """

    def __init__(self, token: str, region: str = "eu-central", insecure_skip_verify: bool = False):
        self.token = token
        self.region = region
        self.logger = logging.getLogger("cloudbridge")
        self.transport = Transport(token, region, insecure_skip_verify)
        self.connections: Dict[str, Connection] = {}
        self.transport.on_message = self._handle_incoming_message
        self.logger.info(f"Initialized CloudBridge Client for region {region}")

    async def _handle_incoming_message(self, peer_id: str, data: bytes):
        if peer_id not in self.connections:
            connection = Connection(peer_id, self.transport)
            self.connections[peer_id] = connection
            await self._handle_new_connection(connection, data)
        else:
            # Existing connection, data handled by connection's internal handler if set
            # But wait, Connection registers itself with transport?
            # In transport.py: if peer_id in self._message_handlers: self._message_handlers[peer_id](payload)
            # So if connection is established, it should have registered a handler.
            # But here we are handling NEW connections.
            # We need to modify Transport to allow a global handler for unknown peers.
            pass

    async def connect(self, peer_id: str) -> Connection:
        """
        Connects to a peer in the CloudBridge network.
        """
        self.logger.info(f"Connecting to peer {peer_id}...")
        
        # Ensure transport is connected
        await self.transport.connect()
        
        if peer_id in self.connections:
            return self.connections[peer_id]

        connection = Connection(peer_id, self.transport)
        self.connections[peer_id] = connection
        return connection

    async def close(self):
        """
        Closes the client and all connections.
        """
        await self.transport.close()
        for conn in self.connections.values():
            await conn.close()
        self.connections.clear()

    def register_service(self, name: str, port: int, tags: Optional[list] = None):
        """
        Registers a service for discovery.
        """
        self.logger.info(f"Registering service {name} on port {port}")
        # Mock registration

    async def create_tunnel(self, config: Dict[str, Any]):
        """
        Creates a tunnel to a remote peer.
        """
        self.logger.info(f"Creating tunnel to {config['remote_peer']}:{config['remote_port']}")
        
        connection = await self.connect(config['remote_peer'])
        
        # Send handshake
        handshake = json.dumps({"type": "tunnel", "port": config['remote_port']})
        await connection.write(handshake.encode('utf-8'))
        
        # Start local listener
        server = await asyncio.start_server(
            lambda r, w: self._handle_local_connection(r, w, connection),
            '127.0.0.1', config['local_port']
        )
        
        self.logger.info(f"Tunnel listening on port {config['local_port']}")
        
        async with server:
            await server.serve_forever()

    async def _handle_local_connection(self, reader, writer, connection):
        try:
            # Forward local data to remote
            async def forward_to_remote():
                while True:
                    data = await reader.read(4096)
                    if not data:
                        break
                    await connection.write(data)
            
            # Forward remote data to local
            async def forward_to_local():
                while True:
                    data = await connection.read()
                    if not data:
                        break
                    writer.write(data)
                    await writer.drain()

            await asyncio.gather(forward_to_remote(), forward_to_local())
        except Exception as e:
            self.logger.error(f"Tunnel connection error: {e}")
        finally:
            writer.close()

    async def serve(self):
        """
        Waits for incoming connections.
        """
        self.logger.info("Waiting for incoming connections...")
        # Keep alive
        while True:
            await asyncio.sleep(1)

    async def _handle_new_connection(self, connection: Connection, initial_data: bytes):
        try:
            handshake = json.loads(initial_data.decode('utf-8'))
            if handshake.get("type") == "tunnel" and "port" in handshake:
                port = handshake["port"]
                self.logger.info(f"Incoming tunnel request for port {port}")
                
                try:
                    reader, writer = await asyncio.open_connection('127.0.0.1', port)
                    
                    # Forward remote data to local
                    async def forward_to_local():
                        # Initial data was the handshake, subsequent data is payload
                        while True:
                            data = await connection.read()
                            if not data:
                                break
                            writer.write(data)
                            await writer.drain()
                            
                    # Forward local data to remote
                    async def forward_to_remote():
                        while True:
                            data = await reader.read(4096)
                            if not data:
                                break
                            await connection.write(data)
                            
                    await asyncio.gather(forward_to_local(), forward_to_remote())
                    
                except Exception as e:
                    self.logger.error(f"Failed to connect to local service at port {port}: {e}")
                    await connection.close()
            else:
                # Not a tunnel or invalid handshake
                pass
        except Exception as e:
            self.logger.error(f"Error handling new connection: {e}")

