import asyncio
import json
import logging
import websockets
from typing import Optional, Callable, Dict

class Transport:
    def __init__(self, token: str, region: str, insecure_skip_verify: bool = False):
        self.token = token
        self.region = region
        self.insecure_skip_verify = insecure_skip_verify
        self.ws: Optional[websockets.WebSocketClientProtocol] = None
        self.logger = logging.getLogger("cloudbridge.transport")
        self.connected = False
        self._message_handlers: Dict[str, Callable[[bytes], None]] = {}
        self.on_message: Optional[Callable[[str, bytes], None]] = None
        self._connect_lock = asyncio.Lock()

    async def connect(self):
        async with self._connect_lock:
            if self.connected:
                return

            url = f"wss://relay.{self.region}.2gc.ru/v1/connect"
            self.logger.info(f"Connecting to {url}")

            import ssl
            if self.insecure_skip_verify:
                ssl_context = ssl.create_default_context()
                ssl_context.check_hostname = False
                ssl_context.verify_mode = ssl.CERT_NONE
            else:
                ssl_context = ssl.create_default_context()

            connect_args = {"additional_headers": {"Authorization": f"Bearer {self.token}"}}
            connect_args["ssl"] = ssl_context
            # If no ssl_context is provided, websockets library handles SSL for wss:// automatically
            # However, if we passed ssl=None explicitly it might complain? 
            # In my code: if ssl_context: connect_args["ssl"] = ssl_context. 
            # So if ssl_context is None, "ssl" key is NOT in connect_args.
            # So why did it fail?
            # "ssl=None is incompatible with a wss:// URI"
            # Maybe `websockets` version behavior?
            # Let's look at the error again.
            # It seems I am NOT passing ssl=None.
            # Wait, I see the logs: `Connecting to wss://relay.eu-central.cloudbridge.global/v1/connect`
            # Oh, I see. I updated the domain to 2gc.ru but the logs show cloudbridge.global?
            # Ah, I might have missed updating the example or the transport file wasn't saved?
            # Let's check the file content again.
            # And also fix the SSL issue.
            # If I don't pass ssl arg, it should work.
            # Maybe I need to pass `ssl=True` or a context?
            # Let's try creating a default context if None.

            try:
                self.ws = await websockets.connect(url, **connect_args)
                self.connected = True
                self.logger.info("Transport connected")
                
                # Start listening loop
                asyncio.create_task(self._listen())
                
            except Exception as e:
                self.logger.error(f"Failed to connect: {e}")
                raise

    async def _listen(self):
        try:
            while self.connected and self.ws:
                message = await self.ws.recv()
                await self._handle_message(message)
        except websockets.exceptions.ConnectionClosed:
            self.logger.info("Transport connection closed")
            self.connected = False
        except Exception as e:
            self.logger.error(f"Error in listen loop: {e}")
            self.connected = False

    async def _handle_message(self, message: str):
        try:
            data = json.loads(message)
            if data.get("type") == "p2p":
                peer_id = data.get("peer_id")
                payload_b64 = data.get("payload")
                import base64
                payload = base64.b64decode(payload_b64)
                
                if peer_id in self._message_handlers:
                    self._message_handlers[peer_id](payload)
                elif self.on_message:
                    # Run in task to avoid blocking
                    asyncio.create_task(self.on_message(peer_id, payload))
        except Exception as e:
            self.logger.error(f"Failed to handle message: {e}")

    async def send(self, peer_id: str, data: bytes):
        if not self.connected or not self.ws:
            raise Exception("Transport not connected")

        import base64
        payload_b64 = base64.b64encode(data).decode('utf-8')
        
        message = {
            "type": "p2p",
            "peer_id": peer_id,
            "payload": payload_b64
        }
        
        await self.ws.send(json.dumps(message))

    def register_handler(self, peer_id: str, handler: Callable[[bytes], None]):
        self._message_handlers[peer_id] = handler

    def unregister_handler(self, peer_id: str):
        if peer_id in self._message_handlers:
            del self._message_handlers[peer_id]

    async def close(self):
        self.connected = False
        if self.ws:
            await self.ws.close()
