import websockets
import inspect

print(f"Websockets version: {websockets.version.version}")
print(f"Connect signature: {inspect.signature(websockets.connect)}")
