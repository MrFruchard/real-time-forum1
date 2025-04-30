import { header } from "./home-components/header.js";
import { bodyHtml } from "./home-components/body.js";

export const home = () => {
    // Supprime tout le contenu de <body>
    document.body.innerHTML = '';

    // Crée un nouveau div #app et l'ajoute à <body>
    const app = document.createElement("div");
    app.id = "app";
    document.body.appendChild(app);

    // Ajoute le header et le contenu du forum
    header();
    bodyHtml();
};