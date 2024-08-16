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

function runHandler() {
    const btn = document.getElementById('running-button');
    btn.innerHTML = 'Running...';
    btn.disabled = true;

    const code = editor.getValue();
    const consoleOutput = document.getElementById('console');
    axios.post("/offline/run", {code: code}).then(res => {
        console.log(res)
        const err = res.data.error
        let output = ""

        if (res.data.output !== "") {
            // output += `Output:\n ${res.data.output}`
            output += `${res.data.output}`
        }

        if (err !== "") {
            const errMessage = err.split(" ").slice(1, err.length - 1).join(" ")
            // output += `\nError:\n ${err}`
            output += `${err}`
        }

        consoleOutput.value = output
    }).catch(err => {
        console.log(err)
    }).finally(() => {
        btn.disabled = false;
        btn.innerHTML = "Run";
    })
}
