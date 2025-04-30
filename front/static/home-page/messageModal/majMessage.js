import { messageModal } from "./messageModal.js";
import { createMessageElement } from "./messageGenerique.js";
import { formatDate } from "./messageGenerique.js";
import { refreshConversations } from "../home-components/body.js";

export const majMessage = (data) => {
  if (!data || !Array.isArray(data) || data.length < 3) return;

  const [msg, user, timestamp] = data;
  const sectionContent = document.querySelector("#messages-id");
  let chatUser = document.getElementById(`chat-${user}`);

  // Vérifie si la modal pour cet utilisateur est déjà ouverte (ID "chat-modal-{user}")
  const userModal = document.getElementById(`chat-modal-${user}`);

  // Fonction de mise à jour de la conversation existante
  const updateConversation = (conversationElem) => {
    // Mise à jour du dernier message
    conversationElem.querySelector(
      ".last-message"
    ).innerHTML = `<strong>${user}:</strong> ${msg}`;

    // Mise à jour de la date
    const formattedDate = formatDate(timestamp);
    conversationElem.querySelector(".message-date").textContent = formattedDate;

    // Ajoute la notification seulement si la modal n'est pas ouverte
    if (!userModal) {
      let notificationDot = conversationElem.querySelector(".notification-dot");
      if (!notificationDot) {
        notificationDot = createNotificationDot();
        conversationElem.appendChild(notificationDot);
      }
    }

    // Retirer complètement l'élément du DOM
    conversationElem.remove();

    // PUIS l'ajouter de nouveau en première position
    sectionContent.prepend(conversationElem);

    // Ajouter un log pour vérifier
    console.log("🔄 Conversation remontée:", user);
  };

  if (chatUser) {
    // Conversation déjà présente : mise à jour
    updateConversation(chatUser);
  } else {
    // Nouvelle conversation
    const msgElement = createMessageElement(
      {
        username: user,
        last_sender: user,
        last_message: msg,
        // ⚠️ Remplace la date brute par la date formatée,
        //    ou envoie seulement le timestamp et formate dans createMessageElement.
        last_message_date: timestamp,
      },
      messageModal
    );
    refreshConversations();
    // Si la modal pour cet utilisateur n'est pas déjà ouverte,
    // on y ajoute le dot de notification
    if (!userModal) {
      const notificationDot = createNotificationDot();
      msgElement.appendChild(notificationDot);
    }

    // Ajout du message en haut de la liste
    sectionContent.prepend(msgElement);
  }
};

const createNotificationDot = () => {
  const notificationDot = document.createElement("span");
  notificationDot.classList.add("notification-dot");
  return notificationDot;
};
