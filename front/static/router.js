import { login } from "./login-page/login.js";
import { home } from "./home-page/home.js";
import { connectedUser } from "./websocket-integration/user-connected.js";
import { disconnectedUser } from "./websocket-integration/user-disconnected.js";
import { privateMessage } from "./home-page/privateMessage/privateMsg.js";
import { majMessage } from "./home-page/messageModal/majMessage.js";
import { refreshConversations } from "./home-page/home-components/body.js";
import {bubbleAnim} from "./home-page/messageModal/bubbleAnim.js";
import {notify} from "./websocket-integration/user-notification.js";

export let socket = null; // Déclare la variable mais ne l'initialise pas immédiatement

export const router = () => {
  // 🔌 Si un WebSocket est déjà ouvert, on le ferme avant d'en créer un nouveau
  if (socket && socket.readyState === WebSocket.OPEN) {
    closeWebSocket();
  }

  // 📡 Création d'une nouvelle connexion WebSocket
  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    console.log("✅ WebSocket connecté !");
    home(); // Charge la page d'accueil
    socket.send(JSON.stringify({ type: "get_user" })); // Demande la liste des utilisateurs

    refreshConversations(); // Met à jour les conversations
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("📩 Message WebSocket reçu :", data);

      if (data.type === "connected_users" || data.type === "new_user") {
        if (
          data.content.length > 0 &&
          data.content[0] !== "No users connected"
        ) {
          console.log("👥 Liste des utilisateurs connectés :", data.content);
          connectedUser(data.content);
        }
      }

      if (data.type === "user_disconnected") {
        if (
          data.content.length > 0 &&
          data.content[0] !== "No users connected"
        ) {
          console.log("🚪 Utilisateur déconnecté :", data.content);
          disconnectedUser(data.content);
        }
      }

      if (data.type === "private") {
        privateMessage(data.content);
        majMessage(data.content);
      }

      if (data.type === "is_typing" || data.type === "is_not_typing") {
        bubbleAnim(data.content, data.is_typing);
      }

      if (data.type === "notification") {
        notify()
      }

    } catch (error) {
      console.error(
        "❌ Erreur lors de la réception du message WebSocket :",
        error
      );
    }
  };

  socket.onerror = (error) => {
    console.error("⚠️ Erreur WebSocket :", error);
    login();
  };

  socket.onclose = (event) => {
    console.warn("🔌 WebSocket fermé :", event.reason);
    login();
  };
};

// 🔥 Fonction pour fermer proprement le WebSocket
export const closeWebSocket = () => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("🛑 Fermeture du WebSocket...");
    socket.close(1000, "Déconnexion de l'utilisateur");
    socket = null; // Réinitialisation pour éviter d'appeler des événements sur un socket fermé
  }
};
