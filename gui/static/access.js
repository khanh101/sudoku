function interval_access() {
    httpPost("api/interval_access", "json", {
        key: key,
    }, function() {
        setTimeout(interval_access, 30000);
    });
}