import {router} from "../../router.js";

export const postModal = async () => {
    // Vérifier si la modal existe déjà
    let modal = document.getElementById("postModal");

    if (!modal) {
        modal = document.createElement("div");
        modal.id = "postModal";
        modal.style.position = "fixed";
        modal.style.top = "50%";
        modal.style.left = "50%";
        modal.style.transform = "translate(-50%, -50%)";
        modal.style.width = "400px";
        modal.style.background = "#ffffff";
        modal.style.padding = "20px";
        modal.style.borderRadius = "10px";
        modal.style.boxShadow = "0 4px 8px rgba(0, 0, 0, 0.2)";
        modal.style.zIndex = "9999";
        modal.style.display = "none";

        modal.innerHTML = `
            <h2 style="margin-bottom: 10px; font-size: 1.5rem; text-align: center;">Créer un post</h2>

            <label for="category">Catégorie :</label>
            <input type="text" id="category" style="width: 100%; padding: 8px; margin-bottom: 10px; border: 1px solid #ccc; border-radius: 5px;" placeholder="Nom de la catégorie (sans espace)" required>

            <label for="content">Contenu :</label>
            <textarea id="content" rows="4" style="width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 5px;" placeholder="Écrivez votre post ici..." required></textarea>

            <div style="display: flex; justify-content: space-between; margin-top: 15px;">
                <button id="closePostModal" style="background: #ccc; padding: 8px 12px; border: none; border-radius: 5px; cursor: pointer;">Annuler</button>
                <button id="submitPost" style="background: #000000; color: white; padding: 8px 12px; border: none; border-radius: 5px; cursor: pointer;">Publier</button>
            </div>
        `;

        document.body.appendChild(modal);

        // Événement fermeture de la modal
        document.getElementById("closePostModal").addEventListener("click", () => {
            modal.style.display = "none";
        });

        // Événement pour soumettre le post
        document.getElementById("submitPost").addEventListener("click", async () => {
            const categoryInput = document.getElementById("category");
            const contentInput = document.getElementById("content");

            // Supprimer les espaces dans la catégorie et trim() le contenu
            const category = categoryInput.value.replace(/\s+/g, "").trim();
            const content = contentInput.value.trim();

            if (!category || !content) {
                alert("Veuillez remplir tous les champs correctement.");
                return;
            }

            // Envoi des données au serveur
            try {
                const response = await fetch("http://localhost:8080/post", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ category: category, content: content })
                });

                if (!response.ok) {
                    throw new Error("Erreur lors de la création du post.");
                }

                console.log("Post ajouté avec succès !");
                modal.style.display = "none"; // Fermer après envoi
                router()
            } catch (error) {
                console.error("Erreur :", error);
                alert("Une erreur est survenue. Réessayez.");
            }
        });
    }

    modal.style.display = "block";
};