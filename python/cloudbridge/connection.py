import asyncio
from typing import Optional
from .transport import Transport

class Connection:
    def __init__(self, peer_id: str, transport: Transport):
        self.peer_id = peer_id
        self.transport = transport
        self.connected = True
        self._read_queue = asyncio.Queue()
        
        # Register handler for incoming messages
        self.transport.register_handler(self.peer_id, self._handle_data)

    def _handle_data(self, data: bytes):
        self._read_queue.put_nowait(data)

    async def read(self) -> bytes:
        if not self.connected:
            raise Exception("Connection closed")
        return await self._read_queue.get()

    async def write(self, data: bytes):
        if not self.connected:
            raise Exception("Connection closed")
        await self.transport.send(self.peer_id, data)

    async def close(self):
        self.connected = False
        self.transport.unregister_handler(self.peer_id)
