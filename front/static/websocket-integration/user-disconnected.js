export const disconnectedUser = (data) => {
    const usersContent = document.querySelector(".section-content");

    if (!usersContent) {
        console.error("❌ L'élément usersContent n'existe pas !");
        return;
    }

    console.log("👀 Utilisateurs avant suppression :", Array.from(usersContent.children).map(child => child.id));

    data.forEach((userId) => {
        const userElement = document.querySelector(`#user-${userId}`);

        if (userElement) {
            userElement.remove();
            console.log(`✅ Utilisateur ${userId} supprimé.`);
        } else {
            console.warn(`⚠️ L'utilisateur ${userId} n'a pas été trouvé.`);
        }
    });

    console.log("🧹 Mise à jour de la liste des utilisateurs !");
};