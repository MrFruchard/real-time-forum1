import { messageModal } from "../home-page/messageModal/messageModal.js";

export const connectedUser = (data) => {
    const usersContent = document.querySelector(".section-content"); // Récupère l'élément parent

    if (!usersContent) {
        console.error("L'élément usersContent n'existe pas !");
        return;
    }

    data.forEach((user) => {
        // Vérifier si l'utilisateur est déjà affiché
        const existingUser = document.getElementById(`user-${user}`);
        if (!existingUser) {
            // Création du conteneur de l'utilisateur
            const userContainer = document.createElement("div");
            userContainer.classList.add("user-container");
            userContainer.id = `user-${user}`;

            // Ajout de l'icône verte pour l'état "connecté"
            const statusIndicator = document.createElement("span");
            statusIndicator.classList.add("status-indicator");

            // Création du nom d'utilisateur cliquable
            const userElement = document.createElement("p");
            userElement.textContent = user;
            userElement.classList.add("user-name");

            // Ajout de l'écouteur d'événements pour ouvrir le chat
            userElement.addEventListener("click", () => {
                messageModal(user);
            });

            // Ajout des éléments au conteneur
            userContainer.appendChild(statusIndicator);
            userContainer.appendChild(userElement);

            // Insertion dans l'ordre alphabétique
            insertSorted(usersContent, userContainer);
        }
    });
};

// Fonction pour insérer un élément dans le parent de façon triée
function insertSorted(parent, newElement) {
    const newUserName = newElement.querySelector('.user-name').textContent.toLowerCase();
    const children = Array.from(parent.getElementsByClassName('user-container'));

    let inserted = false;
    for (let child of children) {
        const childUserName = child.querySelector('.user-name').textContent.toLowerCase();
        if (newUserName < childUserName) {
            parent.insertBefore(newElement, child);
            inserted = true;
            break;
        }
    }
    if (!inserted) {
        parent.appendChild(newElement);
    }
}

// Ajout du style CSS dynamique
const style = document.createElement("style");
style.innerHTML = `
    .user-container {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px;
        background: #000000FF;
        color: white;
        border-radius: 5px;
        margin-bottom: 5px;
        cursor: pointer;
        transition: background 0.3s;
    }

    .user-container:hover {
        background: #444;
    }

    .status-indicator {
        width: 10px;
        height: 10px;
        background: #28a745;
        border-radius: 50%;
    }

    .user-name {
        font-size: 14px;
    }
`;

// Ajout du style au <head>
document.head.appendChild(style);