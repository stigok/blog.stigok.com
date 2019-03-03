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
HTTP_HOST      = ""
HTTP_PORT      = 8000
TRUST_PROXY    = False
MAX_REQUEST_BODY_BYTES = 10000


# TODO: Memoize
def get_post_ids():
    post_ids = []
    with open(JSON_FEED_PATH, mode='r') as f:
        post_ids += json.load(f)
    return post_ids


def write_liquid_comment_file(post_id, body, *, metadata=dict):
    now = datetime.datetime.utcnow()

    out_file = path.join(OUT_DIR, post_id, now.strftime("%Y%m%d-%H%M%S-%f.md"))
    os.makedirs(os.path.dirname(out_file), exist_ok=True)

    metadata['date'] = metadata.get('date', now.isoformat())
    metadata['subject_id'] = post_id

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

    def address_string(self):
        """Helpful when running behind an nginx proxy"""
        orig = super(CommentRequestHandler, self).address_string()
        if TRUST_PROXY:
            addresses = self.headers.get('X-Forwarded-For', orig).split(', ')
            return addresses[-1]
        else:
            return orig

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

        try:
            comment = urllib.parse.parse_qs(req_body.decode(), strict_parsing=True)
        except ValueError:
            # This may indicate invalid content-length in the request
            self.send_error(400, "Errors while parsing request body. Invalid Content-Length?")
            return

        # urllib.parse returns values in arrays
        body   = comment.get('body',   [None]).pop()
        author = comment.get('author', [None]).pop()
        email  = comment.get('email',  [None]).pop()

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
            write_liquid_comment_file(post_id, body, metadata=metadata)

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

if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('--hostname',    default=HTTP_HOST)
    parser.add_argument('--port',        type=int, default=HTTP_PORT)
    parser.add_argument('--trust-proxy', action="store_true",
                        help="Trust X-Forwarded-For header to determine remote address",
                        default=TRUST_PROXY)
    args = parser.parse_args()

    TRUST_PROXY = args.trust_proxy
    server = http.server.HTTPServer((args.hostname, args.port), CommentRequestHandler)

    print("Listening on %s:%d" % (args.hostname, args.port))
    server.serve_forever()
