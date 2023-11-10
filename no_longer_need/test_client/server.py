import http.server
import socketserver

PORT = 3030

handler = http.server.SimpleHTTPRequestHandler

with socketserver.TCPServer(("", PORT), handler) as httpd:
    print("Server started at http://localhost:" + str(PORT))
    httpd.serve_forever()