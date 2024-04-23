document.addEventListener("DOMContentLoaded", function() {
  const chatForm = document.getElementById("chat-form");
  const chatInput = document.getElementById("chat-input");
  const chatMessages = document.getElementById("chat-messages");
  const clearChatButton = document.getElementById("clear-chat");

  let messages = JSON.parse(localStorage.getItem("messages")) || [];

  function renderMessages() {
    chatMessages.innerHTML = "";
    messages.forEach(message => {
      const messageElement = document.createElement("div");
      messageElement.classList.add("message");
      
      messageElement.innerHTML = `
        <div class="message-sender">${message.sender}</div>
        <div class="message-text">${message.text}</div>
        <div class="message-timestamp">${message.timestamp}</div>
      `;
      chatMessages.appendChild(messageElement);
    });
    // Scroll to the bottom
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  chatForm.addEventListener("submit", function(event) {
    event.preventDefault();
    const messageText = chatInput.value.trim();

    if (messageText !== "") {
      const timestamp = new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
      const message = { sender: "You", text: messageText, timestamp: timestamp };
      messages.push(message);
      localStorage.setItem("messages", JSON.stringify(messages));

      renderMessages();

      chatInput.value = "";
    }
  });

  clearChatButton.addEventListener("click", function() {
    localStorage.removeItem("messages");
    messages = [];
    renderMessages();
  });

  renderMessages();
});
