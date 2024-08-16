function connect() {
    const sid = document.getElementById("session_id").value;
    window.location.href = `/connect?s=${sid}`;
}
