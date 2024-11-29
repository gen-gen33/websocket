document.addEventListener('DOMContentLoaded', () => {
    let loc = window.location;
    let uri = (loc.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + loc.host + '/ws';

    const ws = new WebSocket(uri);

    ws.onopen = function () {
        console.log('Connected to WebSocket server');
    };

    ws.onmessage = function (evt) {
        let out = document.getElementById('output');
        out.innerHTML += evt.data + '<br>';
    };

    ws.onclose = function () {
        console.log('Connection closed. Attempting to reconnect...');
        setTimeout(() => {
            location.reload(); // ページを再読込
        }, 1000);
    };

    ws.onerror = function (err) {
        console.error('WebSocket error:', err);
    };

    const btn = document.querySelector('.btn');
    btn.addEventListener('click', () => {
        const input = document.getElementById('input').value;
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(input);
        } else {
            console.warn('WebSocket is not open. Message not sent.');
        }
    });
});
