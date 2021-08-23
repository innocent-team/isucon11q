#
# This is an example VCL file for Varnish.
#
# It does not do anything by default, delegating control to the
# builtin VCL. The builtin VCL is called when there is no explicit
# return statement.
#
# See the VCL chapters in the Users Guide at https://www.varnish-cache.org/docs/
# and https://www.varnish-cache.org/trac/wiki/VCLExamples for more examples.

# Marker to tell the VCL compiler that this VCL has been adapted to the
# new 4.0 format.
vcl 4.0;

import directors;

# Default backend definition. Set this to point to your content server.
backend isucondition1 {
    .host = "192.168.0.11";
    .port = "3000";
}

backend isucondition3 {
    .host = "127.0.0.1";
    .port = "3000";
}

sub vcl_init {
    new bar = directors.round_robin();
    bar.add_backend(isucondition1);
    bar.add_backend(isucondition3);
}

sub vcl_recv {
    # Happens before we check if we have this in cache already.
    #
    # Typically you clean up the request here, removing cookies you don't need,
    # rewriting the request, etc.
    # アイコン関係ありそうなものは1に寄せる
    if (req.url ~ "/icon$") {
        set req.backend_hint = isucondition1;
    } else {
        set req.backend_hint = bar.backend();
    }

    if (req.url ~ "^/initialize") {
       ban("obj.http.url ~ ^/api/trend");
    }

    if (req.url ~ "^/api/trend") {
        unset req.http.cookie;
    }
}

sub vcl_backend_response {
    # Happens after we have read the response headers from the backend.
    #
    # Here you clean the response headers, removing silly Set-Cookie headers
    # and other mistakes your backend does.
    if (bereq.url ~ "^/api/trend") {
        set beresp.grace = 0.2s;
        set beresp.ttl = 1.8s;
        set beresp.http.Cache-Control = "public, max-age=1";
    }
    set beresp.do_gzip = true;
}

sub vcl_deliver {
    # Happens when we have all the pieces we need, and are about to send the
    # response to the client.
    #
    # You can do accounting or modifying the final object here.
    if (obj.hits != 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }
}
