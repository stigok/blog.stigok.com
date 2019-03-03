import datetime
import hashlib
import http.server
import json
import os
import os.path as path
import re
import socketserver
import urllib.parse


OUT_DIR        = "collections/_comments"
JSON_FEED_PATH = "_site/comments_subject_ids.json"
HTTP_PORT      = 8000
MAX_REQUEST_BODY_BYTES = 10_000


# TODO: Memoize
def get_post_ids():
    post_ids = []
    with open(JSON_FEED_PATH, mode='r') as f:
        post_ids += json.load(f)
    return post_ids


def write_liquid_file(filename, body, *, metadata=dict):
    now = datetime.datetime.utcnow()
    unixnow = str(round(now.timestamp() * 1000))

    out_file = path.join(OUT_DIR, filename + "-" + unixnow + ".md")

    # Set comment creation date if not already present
    metadata['date'] = metadata.get('date', now.isoformat())

    # Save comment as markdown file to disk
    with open(out_file, mode='x') as f:
        # Liquid template front matter
        f.write("---\n")
        for k, v in metadata.items():
            f.write("%s: %s\n" % (k, v))
        f.write("---\n\n")

        # File body
        f.write(body + "\n")

    print("Saved file to %s" % out_file)


class CommentRequestHandler(http.server.BaseHTTPRequestHandler):
    error_message_format = "%(code)d %(message)s\n"
    server_version = "jecose"
    sys_version = "0.1"
    protocol_version = "HTTP/1.1"

    def do_POST(self):
        # Only accept GET requests to /comments/<id>
        # with an optional trailing slash
        match = re.match(r"/comments/(?P<id>[-\w]+)/?$", self.path.split("?")[0])
        post_id = match.group("id") if match else None

        if not post_id or post_id not in get_post_ids():
            self.send_error(404)
            return

        # Validate input
        if self.headers.get("content-type", "") != "application/x-www-form-urlencoded":
            self.send_error(400, message="Invalid content type")
            return

        try:
            req_len = int(self.headers.get("content-length", 0))
        except ValueError:
            self.send_error(400)
            return

        if req_len == 0:
            self.send_error(411)
            return

        if req_len > MAX_REQUEST_BODY_BYTES:
            self.send_error(413)
            return

        # Read request body
        req_body = self.rfile.read(req_len)
        comment = urllib.parse.parse_qs(req_body.decode())

        # urllib.parse returns values in arrays
        body   = comment.get('body',   []).pop()
        author = comment.get('author', []).pop()
        email  = comment.get('email',  []).pop()

        if not body:
            self.send_error(400, "Field 'body' missing or too short")
            return

        if not author:
            self.send_error(400, "Field 'author' missing or too short")
            return

        if email:
            if not re.match(r"[^@]+@[^@\.]+\.[^@]+", email):
                self.send_error(400, "Failed to parse field 'email': invalid format")
                return

        try:
            # Create the comment
            metadata = dict(author=author, email=email)
            write_liquid_file(post_id, body, metadata=metadata)

            data = b"201 Created"
            self.send_response(201)
            self.send_header("Content-Length", len(data))
            self.send_header("Content-Type", "text/plain")
            self.end_headers()
            self.wfile.write(data)
            return
        except ValueError as e: # TODO: Determine what errors to catch here
            self.send_error(500)
            self.log_error("Failed to save comment to file", e)
            return


server = http.server.HTTPServer(("", HTTP_PORT), CommentRequestHandler)

print("serving at port", HTTP_PORT)
server.serve_forever()
