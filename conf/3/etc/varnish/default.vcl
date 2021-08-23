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

# Default backend definition. Set this to point to your content server.
backend default {
    .host = "192.168.0.11";
    .port = "3000";
}

backend isucondition3 {
    .host = "127.0.0.1";
    .port = "3000";
}

# /api/trend は3で受ける
sub vcl_recv {
    # Happens before we check if we have this in cache already.
    #
    # Typically you clean up the request here, removing cookies you don't need,
    # rewriting the request, etc.
    if (req.url ~ "^/initialize") {
       ban("obj.http.url ~ ^/api/trend");
    }

    if (req.url ~ "^/api/trend") {
        unset req.http.cookie;
        set req.backend_hint = isucondition3;
    }
}

sub vcl_backend_response {
    # Happens after we have read the response headers from the backend.
    #
    # Here you clean the response headers, removing silly Set-Cookie headers
    # and other mistakes your backend does.
    if (bereq.url ~ "^/api/trend") {
        set beresp.grace = 0.2s;
        set beresp.ttl = 0.4s;
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
