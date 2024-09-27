document.addEventListener('keydown', (e) => {
    if ((e.ctrlKey || e.metaKey) && e.code === "Enter") {
        runHandler()
    }
});
