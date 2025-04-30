export const modalNotification = async (notifModal) => {
    try {
        const notifications = await fetchNotif();
        console.log(notifications);

        if (!notifications || notifications.length === 0) {
            notifModal.innerHTML = "<p style='text-align:center;'>Aucune notification.</p>";
            return;
        }

        // Trier les notifications par date (du plus récent au plus ancien)
        notifications.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        // Nettoyage de la modal avant d'afficher les nouvelles notifications
        notifModal.innerHTML = "<h3 style='margin-bottom:10px; text-align:center;'>Notifications</h3>";

        // Générer la liste des notifications
        notifications.forEach(notif => {
            const notifItem = document.createElement("div");
            notifItem.style.padding = "8px";
            notifItem.style.borderBottom = "1px solid rgba(255,255,255,0.2)";
            notifItem.style.fontSize = "14px";

            // Formatage de la date
            const date = new Date(notif.created_at);
            const formattedDate = date.toLocaleString("fr-FR", {
                day: "2-digit",
                month: "short",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit"
            });

            // Déterminer l'icône et le texte en fonction du type
            let actionText = "";
            let icon = "";
            switch (notif.type) {
                case "like":
                    icon = "<svg class=\"icon\" xmlns=\"http://www.w3.org/2000/svg\"\n      viewBox=\"0 0 24 24\" width=\"16\" height=\"16\">\n      <path d=\"M1 21h4V9H1v12zM23 10c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32\n      c0-.41-.17-.79-.44-1.06L14.17 2 7.59 8.59C7.22 8.95 7 9.45 7 10v9\n      c0 1.1.9 2 2 2h9c.78 0 1.45-.45 1.75-1.11l3.58-7.16\n      c.08-.14.12-.3.12-.44v-2z\"/>\n    </svg>";
                    actionText = "a aimé";
                    break;
                case "dislike":
                    icon = "<svg class=\"icon\" xmlns=\"http://www.w3.org/2000/svg\"\n      viewBox=\"0 0 24 24\" width=\"16\" height=\"16\">\n      <path d=\"M23 3h-4v12h4V3zM1 14c0 1.1.9 2 2 2h6.31l-.95 4.57-.03.32\n      c0 .41.17.79.44 1.06L9.83 22l6.58-6.59C16.78 14.05 17 13.55\n      17 13V4c0-1.1-.9-2-2-2H3c-.78 0-1.45.45-1.75 1.11L.67 8.27\n      c-.08.14-.12.3-.12.44v2z\"/>\n    </svg>";
                    actionText = "n'a pas aimé";
                    break;
                default:
                    icon = "ℹ️";
                    actionText = "a effectué une action";
            }

            // Affichage final de la notification
            notifItem.innerHTML = `
                <div style="display: flex; align-items: center; gap: 10px;">
                    <span style="font-size: 18px;">${icon}</span>
                    <div>
                        <p style="margin: 0;">
                            <strong>${notif.sender}</strong> ${actionText} le post ${notif.related_id}
                        </p>
                        <small style="color: rgba(255,255,255,0.6);">${formattedDate}</small>
                    </div>
                </div>
            `;

            notifModal.appendChild(notifItem);
        });

        // Affichage de la modal
        notifModal.style.display = "block";

    } catch (error) {
        console.error("Erreur lors de la récupération des notifications :", error);
        notifModal.innerHTML = "<p style='color:red; text-align:center;'>Erreur de chargement</p>";
    }
};

// Fonction pour récupérer les notifications depuis l'API
const fetchNotif = async () => {
    try {
        const response = await fetch("http://localhost:8080/notification");
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (e) {
        console.error("Erreur dans fetchNotif :", e);
        return [];
    }
};