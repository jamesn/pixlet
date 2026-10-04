"""
Applet: Test App
Summary: For Testing
Description: It's an app for testing.
Author: Test Dev
"""

load("assert.star", "assert")
load("http.star", "http")
load("render.star", "render")

def main(config):
    url = config.get("url")

    resp = http.get(
        url = url,
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "MISS")

    resp = http.get(
        url = url,
        ttl_seconds = 3,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "HIT")

    resp = http.post(
        url = url,
        ttl_seconds = 0,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "MISS")

    resp = http.post(
        url = url,
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "HIT")

    resp = http.post(
        url = url,
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "HIT")

    return render.Root(
        child = render.Box(
            width = 64,
            height = 32,
        ),
    )
