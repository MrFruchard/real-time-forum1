/**
 * Crée un élément de conversation pour un utilisateur donné.
 * @param {Object} msg - Objet contenant les détails du message.
 * @param {string} msg.username - Nom du destinataire.
 * @param {string} msg.last_sender - Expéditeur du dernier message.
 * @param {string} msg.last_message - Contenu du dernier message.
 * @param {string} msg.last_message_date - Date du dernier message.
 * @param {Function} messageModal - Fonction pour ouvrir la modal du chat.
 * @returns {HTMLElement} - Élément de conversation créé.
 */
export const createMessageElement = (msg, messageModal) => {
  const msgElement = document.createElement("div");
  msgElement.classList.add("conversation");
  msgElement.id = `chat-${msg.username}`;

  // Nom du destinataire
  const msgTo = document.createElement("h3");
  msgTo.classList.add("conversation-title");
  msgTo.innerText = msg.username;

  // Dernier message (Expéditeur + Contenu)
  const msgContent = document.createElement("p");
  msgContent.classList.add("last-message");
  msgContent.innerHTML = `<strong>${msg.last_sender}:</strong> ${msg.last_message}`;

  // Date du dernier message (formatée)
  const msgDate = document.createElement("span");
  msgDate.classList.add("message-date");
  msgDate.innerText = `${formatDate(msg.last_message_date)}`;

  // Style du curseur pour interaction
  msgElement.style.cursor = "pointer";

  // Événement pour ouvrir la modal au clic
  msgElement.addEventListener("click", () => {
    messageModal(msg.username);
  });

  // Ajout des éléments dans la conversation
  msgElement.appendChild(msgTo);
  msgElement.appendChild(msgContent);
  msgElement.appendChild(msgDate);

  msgElement.addEventListener("click", () => {
    messageModal(msg.username);

    // On retire la pastille .notification-dot
    const notificationDot = msgElement.querySelector(".notification-dot");
    if (notificationDot) {
      notificationDot.remove();
    }
  });

  return msgElement;
};

export const formatDate = (dateString) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("fr-FR", {
    weekday: "long",
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
};
