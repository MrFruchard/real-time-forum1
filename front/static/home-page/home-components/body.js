import { createMessageElement } from "../messageModal/messageGenerique.js";
import { messageModal } from "../messageModal/messageModal.js";
import { socket } from "../../router.js";

/**
 * Construit l'interface principale de l'application.
 */
export const bodyHtml = () => {
  // Création de l'élément <main> principal
  const main = document.createElement("main");
  main.classList.add("main-container");

  /**
   * Crée une section sticky avec un titre et un contenu scrollable.
   * @param {string} titleText - Le titre de la section.
   * @param {string} [id=""] - Un identifiant optionnel pour le contenu.
   * @returns {object} - Un objet contenant la section et le contenu.
   */
  const createStickySection = (titleText, id = "") => {
    const section = document.createElement("div");
    section.classList.add("section");

    // Création du titre de la section
    const title = document.createElement("h3");
    title.classList.add("section-title");
    title.innerText = titleText;

    // Création du conteneur scrollable
    const content = document.createElement("div");
    content.classList.add("section-content");
    content.id = id;

    section.appendChild(title);
    section.appendChild(content);
    return { section, content };
  };

  // Colonne gauche : Utilisateurs connectés
  const { section: usersContainer } = createStickySection("Utilisateurs connectés");
  usersContainer.classList.add("users-container");

  // Section centrale : Publications (posts)
  const { section: postsContainer, content: postsContent } = createStickySection("Publications");
  postsContainer.classList.add("posts-container");

  /**
   * Récupère et affiche les posts.
   */
  const fetchPosts = async () => {
    try {
      const posts = await getAllPost();

      posts.forEach((post) => {
        // Création d'un conteneur pour le post avec un identifiant unique
        const postElement = document.createElement("div");
        postElement.classList.add("post");
        postElement.setAttribute("data-post-id", post.id);

        // En-tête du post : affiche l'auteur, la catégorie et la date
        const postHeader = document.createElement("div");
        postHeader.classList.add("post-header");

        const postAuthor = document.createElement("strong");
        postAuthor.innerText = `${post.username} - ${post.category}`;

        const postDate = document.createElement("small");
        postDate.innerText = new Date(post.created_at).toLocaleString();
        postDate.classList.add("post-date");

        postHeader.appendChild(postAuthor);
        postHeader.appendChild(postDate);

        // Contenu principal du post
        const postContent = document.createElement("p");
        postContent.innerText = post.content;
        postContent.classList.add("post-content");

        // Section des interactions (likes/dislikes)
        const postInteractions = document.createElement("div");
        postInteractions.classList.add("post-interactions");

        // Bouton Like
        const likeBtn = document.createElement("button");
        likeBtn.classList.add("post-button", "like");
        if (post.liked) {
          likeBtn.classList.add("liked");
        }
        likeBtn.innerHTML = svgLike + `<span>${post.likes}</span>`;

        // Bouton Dislike
        const dislikeBtn = document.createElement("button");
        dislikeBtn.classList.add("post-button", "dislike");
        if (post.disliked) {
          dislikeBtn.classList.add("disliked");
        }
        dislikeBtn.innerHTML = svgDislike + `<span>${post.dislikes}</span>`;

        // Gestion du clic sur le bouton Like
        likeBtn.addEventListener("click", async () => {
          try {
            const isLiked = likeBtn.classList.contains("liked");
            const isDisliked = dislikeBtn.classList.contains("disliked");

            // Envoi de l'événement like au serveur
            await fetch("http://localhost:8080/event", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                type: "post",
                content_type: "like",
                id: post.id,
              }),
            });

            let likeCount = parseInt(likeBtn.querySelector("span").innerText, 10);
            let dislikeCount = parseInt(dislikeBtn.querySelector("span").innerText, 10);

            if (isLiked) {
              // Annule le like si déjà liké
              likeCount -= 1;
              likeBtn.classList.remove("liked");
            } else {
              // Ajoute un like
              likeCount += 1;
              likeBtn.classList.add("liked");
              if (isDisliked) {
                // Annule le dislike s'il existe
                dislikeCount -= 1;
                dislikeBtn.classList.remove("disliked");
              }
            }

            // Sélectionne le nom de l'auteur du post courant et envoie une notification
            const nameElement = postHeader.querySelector("strong");
            if (nameElement) {
              const namePost = nameElement.textContent.split(" ")[0].trim();
              console.log("Notification Like:", namePost);
              socket.send(JSON.stringify({ type: "notify", content: namePost, action: "like" }));
            }

            // Mise à jour des compteurs affichés
            likeBtn.querySelector("span").innerText = likeCount;
            dislikeBtn.querySelector("span").innerText = dislikeCount;
          } catch (e) {
            console.error("Erreur lors de l'envoi du like:", e.message);
          }
        });

        // Gestion du clic sur le bouton Dislike
        dislikeBtn.addEventListener("click", async () => {
          try {
            const isLiked = likeBtn.classList.contains("liked");
            const isDisliked = dislikeBtn.classList.contains("disliked");

            // Envoi de l'événement dislike au serveur
            await fetch("http://localhost:8080/event", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                type: "post",
                content_type: "dislike",
                id: post.id,
              }),
            });

            let likeCount = parseInt(likeBtn.querySelector("span").innerText, 10);
            let dislikeCount = parseInt(dislikeBtn.querySelector("span").innerText, 10);

            if (isDisliked) {
              // Annule le dislike si déjà disliké
              dislikeCount -= 1;
              dislikeBtn.classList.remove("disliked");
            } else {
              // Ajoute un dislike
              dislikeCount += 1;
              dislikeBtn.classList.add("disliked");
              if (isLiked) {
                // Annule le like s'il existe
                likeCount -= 1;
                likeBtn.classList.remove("liked");
              }
            }

            // Sélectionne le nom de l'auteur du post courant et envoie une notification
            const nameElement = postHeader.querySelector("strong");
            if (nameElement) {
              const namePost = nameElement.textContent.split(" ")[0].trim();
              console.log("Notification Dislike:", namePost);
              socket.send(JSON.stringify({ type: "notify", content: namePost, action: "dislike" }));
            }

            // Mise à jour des compteurs affichés
            likeBtn.querySelector("span").innerText = likeCount;
            dislikeBtn.querySelector("span").innerText = dislikeCount;
          } catch (e) {
            console.error("Erreur lors de l'envoi du dislike:", e.message);
          }
        });

        postInteractions.appendChild(likeBtn);
        postInteractions.appendChild(dislikeBtn);

        // Section des commentaires
        const commentsSection = document.createElement("div");
        commentsSection.classList.add("comment-section");

        // Titre pour afficher/masquer les commentaires
        const commentTitle = document.createElement("strong");
        commentTitle.classList.add("comment-title");
        commentTitle.innerHTML = svgComment + ` Commentaires`;
        commentTitle.style.cursor = "pointer";
        commentTitle.style.color = "#000000";

        // Conteneur pour afficher les commentaires
        const commentsContainer = document.createElement("div");
        commentsContainer.classList.add("comments-container");
        commentsContainer.style.display = "none";

        // Formulaire pour ajouter un commentaire
        const commentForm = document.createElement("form");
        commentForm.classList.add("comment-form");
        commentForm.style.display = "none";

        const commentInput = document.createElement("input");
        commentInput.type = "text";
        commentInput.placeholder = "Écrire un commentaire...";
        commentInput.required = true;

        const commentSubmitBtn = document.createElement("button");
        commentSubmitBtn.type = "submit";
        commentSubmitBtn.innerText = "Envoyer";

        commentForm.appendChild(commentInput);
        commentForm.appendChild(commentSubmitBtn);

        // Gestion de la soumission du commentaire
        commentForm.addEventListener("submit", async (e) => {
          e.preventDefault();
          const content = commentInput.value.trim();
          if (!content) return;

          // Envoie du commentaire au serveur
          await sendComment(post.id, content);

          // Rechargement des commentaires pour ce post
          commentsContainer.innerHTML = "";
          const newComments = await fetchComments(post.id);
          if (newComments.length > 0) {
            newComments.forEach((comment) => {
              const commentElement = makeCommentElement(comment);
              commentsContainer.appendChild(commentElement);
            });
          } else {
            commentsContainer.innerText = "Aucun commentaire.";
          }

          // Envoi d'une notification pour le commentaire
          const nameElement = postHeader.querySelector("strong");
          if (nameElement) {
            const namePost = nameElement.textContent.split(" ")[0].trim();
            console.log("Notification Comment:", namePost);
            socket.send(JSON.stringify({ type: "notify", content: namePost, action: "comment" }));
          }

          // Vider le champ de saisie
          commentInput.value = "";
        });

        // Permet d'afficher ou masquer les commentaires lors du clic sur le titre
        commentTitle.addEventListener("click", async () => {
          if (commentsContainer.childNodes.length === 0) {
            const comments = await fetchComments(post.id);
            if (comments.length > 0) {
              comments.forEach((comment) => {
                const commentElement = makeCommentElement(comment);
                commentsContainer.appendChild(commentElement);
              });
            } else {
              commentsContainer.innerText = "Aucun commentaire.";
            }
          }
          const currentDisplay = commentsContainer.style.display;
          commentsContainer.style.display = currentDisplay === "none" ? "block" : "none";
          commentForm.style.display = commentsContainer.style.display;
        });

        // Assemblage final de la section commentaires
        commentsSection.appendChild(commentTitle);
        commentsSection.appendChild(commentForm);
        commentsSection.appendChild(commentsContainer);

        // Assemblage final du post
        postElement.appendChild(postHeader);
        postElement.appendChild(postContent);
        postElement.appendChild(postInteractions);
        postElement.appendChild(commentsSection);
        postsContent.appendChild(postElement);
      });
    } catch (error) {
      console.error("Erreur lors du chargement des posts:", error.message);
    }
  };

  fetchPosts();

  // Colonne droite : Messages envoyés
  const { section: messagesContainer } = createStickySection("Messages", "messages-id");
  messagesContainer.classList.add("messages-container");

  // Appliquer les styles pour les conversations
  const initMessageStyles = () => {
    const style = document.createElement("style");
    style.innerHTML = `
      .conversation {
          display: flex;
          flex-direction: column;
          padding: 12px;
          border-bottom: 1px solid #444;
          background: #000000FF;
          color: white;
          transition: background 0.3s;
          border-radius: 10px;
          margin-bottom: 10px;
          cursor: pointer;
      }

      .conversation:hover {
          background: #3a3b3c;
      }

      .conversation-title {
          font-size: 16px;
          font-weight: bold;
          margin-bottom: 4px;
      }

      .last-message {
          font-size: 14px;
          color: #b0b3b8;
          margin-bottom: 4px;
      }

      .message-date {
          font-size: 12px;
          color: #b0b3b8;
          align-self: flex-end;
      }
      .empty-conversations, .error-message {
          padding: 20px;
          text-align: center;
          color: #b0b3b8;
          background: #000000FF;
          border-radius: 10px;
          margin: 10px;
      }

      .empty-conversations p {
          margin: 5px 0;
      }

      .error-message {
          color: #ff6b6b;
      }
    `;
    document.head.appendChild(style);
  };

  initMessageStyles();
  // Appel de refreshConversations pour charger les conversations (cette fonction est désormais exportée)
  refreshConversations();

  // Assemblage final des colonnes et ajout dans le conteneur #app
  main.appendChild(usersContainer);
  main.appendChild(postsContainer);
  main.appendChild(messagesContainer);
  document.querySelector("#app").appendChild(main);
};

/**
 * Récupère l'ensemble des posts depuis l'API.
 * @returns {Promise<Array>} - Liste des posts.
 */
const getAllPost = async () => {
  try {
    const response = await fetch("http://localhost:8080/post", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });
    if (!response.ok) throw new Error(`Erreur HTTP: ${response.status}`);
    return await response.json();
  } catch (e) {
    throw new Error(`Erreur lors de la récupération des posts: ${e.message}`);
  }
};

/**
 * Récupère les commentaires d'un post donné.
 * @param {number|string} postId - L'identifiant du post.
 * @returns {Promise<Array>} - Liste des commentaires.
 */
const fetchComments = async (postId) => {
  try {
    const response = await fetch(`http://localhost:8080/comment/${postId}`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });
    if (!response.ok) throw new Error(`Erreur HTTP: ${response.status}`);
    let comments = await response.json();
    // Trier les commentaires du plus récent au plus ancien
    comments.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    return comments;
  } catch (e) {
    console.error(`Erreur lors de la récupération des commentaires: ${e.message}`);
    return [];
  }
};

/**
 * Envoie un nouveau commentaire pour un post donné.
 * @param {number|string} postId - L'identifiant du post.
 * @param {string} content - Le contenu du commentaire.
 */
const sendComment = async (postId, content) => {
  try {
    const response = await fetch("http://localhost:8080/comment", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id_post: postId, content }),
    });
    if (!response.ok) throw new Error(`Erreur HTTP: ${response.status}`);
  } catch (e) {
    console.error(`Erreur lors de l'envoi du commentaire: ${e.message}`);
  }
};

/**
 * Crée un élément DOM pour un commentaire.
 * @param {object} comment - L'objet commentaire.
 * @returns {HTMLElement} - L'élément DOM représentant le commentaire.
 */
const makeCommentElement = (comment) => {
  const commentWrapper = document.createElement("div");
  commentWrapper.classList.add("comment");

  const commentHeader = document.createElement("div");
  commentHeader.classList.add("comment-header");

  const userInfo = document.createElement("strong");
  userInfo.classList.add("comment-username");
  userInfo.innerText = comment.username;

  const commentDate = document.createElement("small");
  commentDate.classList.add("comment-date");
  commentDate.innerText = new Date(comment.created_at).toLocaleString();

  commentHeader.appendChild(userInfo);
  commentHeader.appendChild(commentDate);

  const contentInfo = document.createElement("p");
  contentInfo.classList.add("comment-content");
  contentInfo.innerText = comment.content;

  const commentActions = document.createElement("div");
  commentActions.classList.add("comment-actions");

  commentWrapper.appendChild(commentHeader);
  commentWrapper.appendChild(contentInfo);
  commentWrapper.appendChild(commentActions);

  return commentWrapper;
};

// Définition des SVG pour les icônes
const svgLike = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
  <path d="M1 21h4V9H1v12zM23 10c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32c0-.41-.17-.79-.44-1.06L14.17 2 7.59 8.59C7.22 8.95 7 9.45 7 10v9c0 1.1.9 2 2 2h9c.78 0 1.45-.45 1.75-1.11l3.58-7.16c.08-.14.12-.3.12-.44v-2z"/>
</svg>`;

const svgDislike = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
  <path d="M23 3h-4v12h4V3zM1 14c0 1.1.9 2 2 2h6.31l-.95 4.57-.03.32c0 .41.17-.79.44 1.06L9.83 22l6.58-6.59C16.78 14.05 17 13.55 17 13V4c0-1.1-.9-2-2-2H3c-.78 0-1.45.45-1.75 1.11L.67 8.27c-.08.14-.12.3-.12.44v2z"/>
</svg>`;

const svgComment = `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
  <path d="M21 6h-2v9H5v2c0 .55.45 1 1 1h11l4 4V7c0-.55-.45-1-1-1zM17 2H3c-.55 0-1 .45-1 1v14l4-4h11c.55 0 1-.45 1-1V3c0-.55-.45-1-1-1z"/>
</svg>`;

/**
 * Récupère les conversations et met à jour la section Messages.
 * Cette fonction est exportée pour pouvoir être utilisée par d'autres modules.
 */
export const refreshConversations = async () => {
  const sectionContent = document.querySelector("#messages-id");
  if (!sectionContent) return;

  try {
    // Récupération de toutes les conversations via la fonction messages (définie ailleurs ou à ajouter)
    const conversations = await messages();
    sectionContent.innerHTML = "";

    if (!conversations || !Array.isArray(conversations) || conversations.length === 0) {
      const emptyMessage = document.createElement("div");
      emptyMessage.classList.add("empty-conversations");
      emptyMessage.innerHTML = `
          <p>Aucune conversation pour le moment.</p>
          <p>Cliquez sur un utilisateur connecté pour démarrer une discussion.</p>
      `;
      sectionContent.appendChild(emptyMessage);
      return;
    }

    // Tri des conversations par date du dernier message (du plus récent au plus ancien)
    const conversationsWithDates = conversations.map((conv) => ({
      ...conv,
      timestamp: new Date(conv.last_message_date).getTime(),
    }));
    conversationsWithDates.sort((a, b) => b.timestamp - a.timestamp);

    // Création et ajout des éléments pour chaque conversation
    conversationsWithDates.forEach((conv) => {
      const msgElement = createMessageElement(conv, messageModal);
      sectionContent.appendChild(msgElement);
    });

    console.log("📊 Conversations triées par date du dernier message");
  } catch (error) {
    console.error("❌ Erreur lors du tri des conversations:", error);
    sectionContent.innerHTML = `
        <div class="error-message">
          <p>Impossible de charger les conversations.</p>
        </div>
      `;
  }
};

/**
 * Exemple de fonction pour récupérer les conversations depuis l'API.
 * Si vous avez déjà une fonction 'messages' ailleurs, assurez-vous qu'elle est importée ou définie ici.
 */
export const messages = async () => {
  try {
    const response = await fetch("http://localhost:8080/conversation");
    if (!response.ok) throw new Error(`Erreur HTTP: ${response.status} - ${response.statusText}`);
    return await response.json();
  } catch (error) {
    console.error("Erreur lors du chargement des conversations:", error.message);
    return []; // Retourne un tableau vide en cas d'erreur pour éviter un crash
  }
};