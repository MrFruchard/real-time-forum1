export const bubbleAnim = (username, isTyping) => {
    const chatModal = document.querySelector(`#chat-modal-${username}`);
    if (!chatModal) {
        const chatMessage = document.querySelector(`#chat-${username}`);
        if (!chatMessage) return;

        const lastMessageElement = chatMessage.querySelector('.last-message');
        if (!lastMessageElement) return;

        // Sauvegarde le dernier message si ce n'est pas encore fait
        if (!lastMessageElement.dataset.lastMessage) {
            lastMessageElement.dataset.lastMessage = lastMessageElement.textContent;
        }

        if (isTyping) {
            // Remplace par "est en train d'écrire..."
            lastMessageElement.textContent = `${username} est en train d'écrire ...`;
        } else {
            // Restaure le dernier message sauvegardé
            lastMessageElement.textContent = lastMessageElement.dataset.lastMessage;
            delete lastMessageElement.dataset.lastMessage; // Nettoyage après restauration
        }
    }

    const chatBody = chatModal?.querySelector(`.chat-body`);
    if (!chatBody) return;

    if (!isTyping) {
        const isTypingMessage = chatBody.querySelector(`.isTyping`);
        if (isTypingMessage) isTypingMessage.remove();
        return;
    }

    // Vérifie si un message "typing" est déjà affiché
    if (chatBody.querySelector(`.isTyping`)) return;

    const newMessage = document.createElement("p");
    newMessage.innerHTML = svgTyping();
    newMessage.classList.add("chat-message", "received", "isTyping");
    chatBody.appendChild(newMessage);

    // 🔽 Scroll automatique vers le bas
    chatBody.scrollTop = chatBody.scrollHeight;
}


const svgTyping = () => {
    return `<svg width="50" height="20" viewBox="0 0 50 20" xmlns="http://www.w3.org/2000/svg">
  <circle cx="6" cy="10" r="3" fill="#ffffff">
    <animate attributeName="cy" values="10;6;10" dur="1.2s" repeatCount="indefinite" begin="0s"/>
    <animate attributeName="opacity" values="0.6;1;0.6" dur="1.2s" repeatCount="indefinite" begin="0s"/>
  </circle>
  <circle cx="20" cy="10" r="3" fill="#ffffff">
    <animate attributeName="cy" values="10;6;10" dur="1.2s" repeatCount="indefinite" begin="0.2s"/>
    <animate attributeName="opacity" values="0.6;1;0.6" dur="1.2s" repeatCount="indefinite" begin="0.2s"/>
  </circle>
  <circle cx="34" cy="10" r="3" fill="#ffffff">
    <animate attributeName="cy" values="10;6;10" dur="1.2s" repeatCount="indefinite" begin="0.4s"/>
    <animate attributeName="opacity" values="0.6;1;0.6" dur="1.2s" repeatCount="indefinite" begin="0.4s"/>
  </circle>
</svg>`
}