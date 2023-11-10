import numpy as np
import asyncio
import websockets
import queue
from vosk import Model, KaldiRecognizer

CHUNK_SIZE = 1024
FRAME_RATE = 16000
RECORD_TIME = 10

q = queue.Queue()
rec = KaldiRecognizer(Model(lang="en"), FRAME_RATE)

async def websocket_handler(websocket, path):
    async for message in websocket:
        audio_data = np.frombuffer(message, dtype=np.int16)
        q.put(bytes(audio_data.tobytes()))
        # https://github.com/openai/whisper/discussions/873
        # if q.qsize() >= FRAME_RATE * RECORD_TIME / CHUNK_SIZE:
        data = q.get()
        if rec.AcceptWaveform(data):
            print("Final Result: " + rec.Result())
        else:
            print(rec.PartialResult())

        #stream.write(message)

async def start_websocket_server():
    async with websockets.serve(websocket_handler, '', 8888):
        await asyncio.Future()  # Keep the server running


if __name__ == '__main__':
    asyncio.run(start_websocket_server())


# mic = whisper_mic.WhisperMic(model="base", device="cpu")
# mic.listen_loop(dictate=True)