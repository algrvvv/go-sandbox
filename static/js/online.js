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

document.getElementById("session-id").innerHTML = `Session id: ${session}`

socket.onopen = function (event) {
    console.log("connection established: " + wsURL)
}

socket.onmessage = function (event) {
    console.log(event.data)
    THISCHANGEFROMWS = true;

    let msg = JSON.parse(event.data);
    let type = msg.type;
    let data = msg.data;

    if (type === "code") {
        editor.setValue(data);
    } else if (type === "console") {
        const output = JSON.parse(event.data)
        let finalString = "";

        const out = output.data.output;
        const err = output.data.error;
        if (out !== "") finalString += out;
        if (err !== "") finalString += err;

        document.getElementById('console').value = finalString;
    } else if (type === "userCount") {
        console.log("active user count: ", data)
        document.getElementById('active-user-count').innerHTML = `Active user count: ${data}`
    }
}

socket.onclose = function (event) {
    console.log("connection closed")
}

socket.onerror = function (error) {
    console.error(error)
}

editor.on('change', function () {
    if (!THISCHANGEFROMWS) {
        let currentValue = editor.getValue();
        const message = {
            type: "code",
            data: currentValue
        }
        console.log(message)
        socket.send(JSON.stringify(message))
    } else {
        THISCHANGEFROMWS = false;
    }
})

function runHandler() {
    const btn = document.getElementById('running-button');
    btn.innerHTML = 'Running...';
    btn.disabled = true;

    const code = editor.getValue();

    axios.post("/online/run", {
        code: code, session: session, uid: uid
    }).then(res => {
        console.log(res)
    }).catch(err => {
        console.error(err)
    }).finally(() => {
        btn.disabled = false;
        btn.innerHTML = "Run";
    })
}
