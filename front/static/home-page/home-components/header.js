import { closeWebSocket, router } from "../../router.js";
import { modalNotification } from "../notification/modalNotification.js";
import {postModal} from "../post-detail/post.js";

// Fonction d'affichage du header
export const header = () => {
    const header = document.createElement("header");
    header.style.display = "flex";
    header.style.justifyContent = "space-between";
    header.style.alignItems = "center";
    header.style.padding = "15px";
    header.style.backgroundColor = "rgb(0,0,0)";
    header.style.color = "#fff";
    header.style.boxShadow = "0 2px 5px rgba(0, 0, 0, 0.2)";

    const appName = document.createElement("h1");
    appName.textContent = "Real Time Forum";
    appName.style.margin = "0";
    appName.style.fontSize = "1.5rem";
    appName.style.textAlign = "center";

    const actionsContainer = document.createElement("div");
    actionsContainer.style.display = "flex";
    actionsContainer.style.alignItems = "center";
    actionsContainer.style.gap = "15px";

    // Bouton "+"
    const addPostButton = document.createElement("button");
    addPostButton.innerHTML = "+";
    addPostButton.style.fontSize = "20px";
    addPostButton.style.fontWeight = "bold";
    addPostButton.style.display = "flex";
    addPostButton.style.alignItems = "center";
    addPostButton.style.justifyContent = "center";
    addPostButton.style.width = "35px";
    addPostButton.style.height = "35px";
    addPostButton.style.borderRadius = "50%";
    addPostButton.style.backgroundColor = "#ffffff";
    addPostButton.style.color = "#000000";
    addPostButton.style.border = "none";
    addPostButton.style.cursor = "pointer";
    addPostButton.style.transition = "0.3s";

    addPostButton.addEventListener("click", () => {
        postModal()
    });

    // Bouton de notification
    const notifButton = document.createElement("button");
    notifButton.innerHTML = `
      <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20">
          <path d="M12 22c1.1 0 2-.9 2-2h-4a2 2 0 0 0 2 2zm6-6v-5c0-3.07-1.64-5.64-4.5-6.32V4a1.5 1.5 0 0 0-3 0v.68C7.64 5.36 6 7.92 6 11v5l-1 1v1h16v-1l-1-1z"/>
      </svg>`;
    notifButton.classList.add("notif");
    notifButton.style.display = "flex";
    notifButton.style.alignItems = "center";
    notifButton.style.justifyContent = "center";
    notifButton.style.background = "none";
    notifButton.style.border = "none";
    notifButton.style.cursor = "pointer";
    notifButton.style.color = "#fff";
    notifButton.style.position = "relative";

    // Création de la modal de notifications
    const notifModal = document.createElement("div");
    notifModal.style.position = "absolute";
    notifModal.style.top = "40px";
    notifModal.style.right = "0";
    notifModal.style.width = "300px";
    notifModal.style.background = "#000000";
    notifModal.style.color = "#fff";
    notifModal.style.borderRadius = "10px";
    notifModal.style.boxShadow = "0 4px 6px rgba(0, 0, 0, 0.3)";
    notifModal.style.padding = "15px";
    notifModal.style.display = "none";
    notifModal.style.zIndex = "9999";
    notifModal.style.maxHeight = "400px";
    notifModal.style.overflowY = "auto";

    // Bouton de fermeture en SVG
    const closeButton = document.createElement("button");
    closeButton.innerHTML = `
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="18" height="18">
          <path fill="white" d="M18.3 5.71a1 1 0 00-1.41 0L12 10.59 7.11 5.71a1 1 0 00-1.41 1.41L10.59 12l-4.89 4.88a1 1 0 101.41 1.41L12 13.41l4.88 4.89a1 1 0 001.41-1.41L13.41 12l4.89-4.88a1 1 0 000-1.41z"/>
      </svg>`;
    closeButton.style.position = "absolute";
    closeButton.style.top = "10px";
    closeButton.style.right = "10px";
    closeButton.style.background = "none";
    closeButton.style.border = "none";
    closeButton.style.cursor = "pointer";

    notifModal.appendChild(closeButton);

    notifButton.addEventListener("click", async () => {
        if (notifModal.style.display === "none") {
            await modalNotification(notifModal);
            notifModal.style.display = "block";
        } else {
            notifModal.style.display = "none";
        }
    });

    closeButton.addEventListener("click", () => {
        notifModal.style.display = "none";
    });

    document.addEventListener("click", (event) => {
        if (!notifButton.contains(event.target) && !notifModal.contains(event.target)) {
            notifModal.style.display = "none";
        }
    });

    // Bouton de déconnexion
    const logoutButton = document.createElement("button");
    logoutButton.innerHTML = `
      <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20" style="margin-right: 5px;">
          <path d="M16 13v-2H7V8l-5 4 5 4v-3zM20 3h-8v2h8v14h-8v2h8a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2z"/>
      </svg> Déconnexion`;
    logoutButton.style.display = "flex";
    logoutButton.style.alignItems = "center";
    logoutButton.style.backgroundColor = "#e31414";
    logoutButton.style.color = "#fff";
    logoutButton.style.border = "none";
    logoutButton.style.cursor = "pointer";
    logoutButton.style.borderRadius = "5px";
    logoutButton.style.padding = "8px 12px";

    logoutButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/logout", {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json"
                }
            });

            if (!response.ok) {
                throw new Error("Erreur lors de la déconnexion");
            }

            closeWebSocket();
            router();
        } catch (error) {
            console.error("Erreur lors de la requête de déconnexion :", error);
            alert("Erreur lors de la déconnexion. Veuillez réessayer.");
        }
    });

    actionsContainer.appendChild(addPostButton);
    actionsContainer.appendChild(notifButton);
    actionsContainer.appendChild(notifModal);
    actionsContainer.appendChild(logoutButton);

    header.appendChild(appName);
    header.appendChild(actionsContainer);

    const app = document.querySelector("#app");
    if (app) {
        app.prepend(header);
    }
};