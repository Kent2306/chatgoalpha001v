const ws = new WebSocket(`ws://${window.location.host}/ws`);
const elements = {
    input: document.getElementById("message"),
    btn: document.getElementById("send-button"),
    files: document.getElementById("file-input"),
    container: document.getElementById("messages")
};

function displayMessage(msg) {
    const message = document.createElement("div");
    message.className = msg.is_system ? "system-message" :
        (msg.username === username ? "message outgoing" : "message incoming");

    if (msg.is_file) {
        message.innerHTML = `
            <div class="file-message">
                ${msg.file_type === "image" ?
            `<div class="file-message-image">
                    <img src="${msg.text}" class="image-preview">
                    <a href="${msg.text}" download class="download-btn">
                        <i class="fas fa-download"></i>
                    </a>
                </div>` :
            `<a href="${msg.text}" download>${msg.file_name}</a>`}
            </div>
        `;
    } else {
        message.innerHTML = `
            <div class="message-header">
                <span class="message-username">${msg.username}</span>
                <span class="message-time">${msg.time}</span>
            </div>
            <div class="message-content">${msg.text}</div>
        `;
    }

    elements.container.appendChild(message);
    scrollToBottom();
}

function scrollToBottom() {
    elements.container.scrollTop = elements.container.scrollHeight;
}

function sendMessage() {
    const text = elements.input.value.trim();
    if (text) {
        ws.send(JSON.stringify({
            text: text,
            time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        }));
        elements.input.value = "";
    }
}

elements.btn.addEventListener("click", sendMessage);
elements.input.addEventListener("keypress", (e) => e.key === "Enter" && sendMessage());

elements.files.addEventListener("change", (e) => {
    if (!e.target.files[0]) return;
    const formData = new FormData();
    formData.append("file", e.target.files[0]);
    fetch("/upload", { method: "POST", body: formData });
    e.target.value = "";
});

ws.onmessage = (e) => {
    const data = JSON.parse(e.data);

    if (data.action === "clear_chat") {
        document.getElementById('messages').innerHTML = '';
        return;
    }

    displayMessage(data);
};

function updateOnlineUsers() {
    fetch('/online-users')
        .then(response => response.json())
        .then(users => {
            const onlineList = document.getElementById('online-list');
            const onlineCount = document.getElementById('online-count');

            onlineList.innerHTML = '';
            onlineCount.textContent = users.length;

            users.forEach(username => {
                const userElement = document.createElement('div');
                userElement.className = 'online-user';
                userElement.innerHTML = `
                    <div class="online-user-avatar">${username.charAt(0)}</div>
                    <div class="online-user-name">${username}</div>
                `;
                onlineList.appendChild(userElement);
            });
        });
}

setInterval(updateOnlineUsers, 5000);
updateOnlineUsers();

document.getElementById('clear-chat')?.addEventListener('click', function() {
    if (confirm('Вы уверены, что хотите очистить весь чат для всех пользователей?')) {
        fetch('/clear-chat', { method: 'POST' })
            .then(response => {
                if (!response.ok) {
                    alert('Ошибка при очистке чата');
                }
            });
    }
});

document.addEventListener('DOMContentLoaded', scrollToBottom);
// Обработка модального окна
document.addEventListener('DOMContentLoaded', function() {
    // Проверяем, показывали ли уже предупреждение
    if (!localStorage.getItem('warningAccepted')) {
        const modal = document.getElementById('welcome-modal');
        if (modal) {
            modal.style.display = 'block';

            document.getElementById('accept-btn').addEventListener('click', function() {
                modal.style.display = 'none';
                // Запоминаем, что пользователь принял условия
                localStorage.setItem('warningAccepted', 'true');
            });
        }
    }
});