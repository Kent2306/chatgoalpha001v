const ws = new WebSocket("ws://" + window.location.host + "/ws");
const messageInput = document.getElementById("message");
const sendButton = document.getElementById("send-button");
const fileInput = document.getElementById("file-input");

// Theme toggle
document.getElementById('theme-toggle')?.addEventListener('click', function() {
    fetch('/toggle-theme', { method: 'POST' })
        .then(() => window.location.reload());
});

// File upload
fileInput?.addEventListener('change', function(e) {
    const file = e.target.files[0];
    if (file) {
        const formData = new FormData();
        formData.append('file', file);
        
        fetch('/upload', {
            method: 'POST',
            body: formData
        });
        
        // Reset input to allow uploading same file again
        e.target.value = '';
    }
});

// Message sending
function sendMessage() {
    const message = messageInput.value.trim();
    if (message) {
        ws.send(JSON.stringify({
            text: message,
            time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        }));
        messageInput.value = "";
    }
}

messageInput?.addEventListener("keypress", function(e) {
    if (e.key === "Enter") {
        sendMessage();
    }
});

sendButton?.addEventListener("click", sendMessage);

// Message receiving
ws.onmessage = function(event) {
    const msg = JSON.parse(event.data);
    const messages = document.getElementById("messages");
    
    const messageElement = document.createElement("div");
    messageElement.className = `message ${msg.username === username ? 'outgoing' : 'incoming'}`;
    
    if (msg.is_file) {
        messageElement.innerHTML = `
            <div class="message-header">
                <span class="message-username">${msg.username}</span>
                <span class="message-time">${msg.time}</span>
            </div>
            <div class="message-content">
                <a href="${msg.text}" download="${msg.fileName}" class="file-message">
                    <i class="fas fa-file"></i> ${msg.fileName}
                </a>
            </div>
        `;
    } else {
        messageElement.innerHTML = `
            <div class="message-header">
                <span class="message-username">${msg.username}</span>
                <span class="message-time">${msg.time}</span>
            </div>
            <div class="message-content">
                <div class="message-text">${msg.text}</div>
            </div>
        `;
    }
    
    messages.appendChild(messageElement);
    messages.scrollTop = messages.scrollHeight;
};