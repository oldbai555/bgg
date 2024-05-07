function download(sUrl) {
    doGet("/api/download/" + sUrl)
}

function doGet(url) {
    window.location.href = url;
}
