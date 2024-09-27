document.addEventListener('keydown', (e) => {
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.code === "KeyF") {
        gofmt()
    }
});

function gofmt() {
    const btn = document.getElementById('fmt-button');
    btn.innerHTML = 'Formating...';
    btn.disabled = true;

    const code = editor.getValue();
    axios.post("/gofmt", {code: code}).then(res => {
        console.log("response after gofmt", res)
        editor.setValue(res.data.code);
    }).catch(err => {
        console.error("failed to format file", err)
    }).finally(() => {
        btn.innerHTML = 'Format';
        btn.disabled = false;
    })
}
