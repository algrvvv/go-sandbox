const editor = CodeMirror.fromTextArea(document.getElementById('codeEditor'), {
    mode: 'go',
    theme: 'material',
    lineNumbers: true,
    autoCloseBrackets: true,
    matchBrackets: true,
    extraKeys: {
        "utrl-Space": "autocomplete",
    },
    fontFamily: 'JetBrains Mono, monospace'
});
const searchParams = new URLSearchParams(window.location.search);
const session = searchParams.get("s");
const uid = searchParams.get("u");
const wsURL = `ws://${location.host}/ws?s=${session}&u=${uid}`
const socket = new WebSocket(wsURL)
let THISCHANGEFROMWS = false

socket.onopen = function (event) {
    console.log("connection established: " + wsURL)
}

socket.onmessage = function (event) {
    console.log(event.data)
    THISCHANGEFROMWS = true;
    editor.setValue(event.data);
}

socket.onclose = function (event) {
    console.log("connection closed")
}

socket.onerror = function (error) {
    console.error(error)
}

editor.on('change', function () {
    if (!THISCHANGEFROMWS) {
        var currentValue = editor.getValue();
        console.log(currentValue)
        socket.send(currentValue)
    } else {
        THISCHANGEFROMWS = false;
    }
})
