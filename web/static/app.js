let socket;
let username = "";

window.onload = () => {
    document.getElementById("saveUsername").onclick = () => {
        const name = document.getElementById("usernameInput").value.trim();
        if (name === "") return alert("Escribe un nombre");

        username = name;

        document.getElementById("usernameModal").style.display = "none";
        document.getElementById("msg").disabled = false;
        document.getElementById("sendBtn").disabled = false;

        iniciarWebSocket();
    };
};

function iniciarWebSocket() {
    socket = new WebSocket(`wss://${window.location.host}/ws`);


    socket.onopen = () => console.log("Conectado");

    socket.onmessage = (event) => {
        const msg = event.data;

        if (msg.startsWith(username + ":")) {
            agregarMensaje(msg, "me");
        } else {
            agregarMensaje(msg, "other");
        }
    };

    socket.onclose = () => console.log("Desconectado");
}

function agregarMensaje(text, type) {
    const messages = document.getElementById("messages");

    const div = document.createElement("div");
    div.classList.add("message", type);
    div.textContent = text;

    messages.appendChild(div);
    messages.scrollTop = messages.scrollHeight;
}

function sendMsg() {
    const input = document.getElementById("msg");
    const text = input.value.trim();
    if (!text) return;

    const msg = `${username}: ${text}`;
    socket.send(msg);

    input.value = "";
}
