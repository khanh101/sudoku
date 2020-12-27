function access() {
    httpPost("api/access", "json", {
        key: key,
    }, function() {
        setTimeout(access, 30000);
    });
}