import { router } from "../router.js";
import { login } from "./login.js";

export const register = () => {
    // Vider complètement le `body`
    document.body.innerHTML = "";

    // Créer un conteneur principal
    const container = document.createElement("div");
    container.id = "register-container";
    container.style.display = "flex";
    container.style.justifyContent = "center";
    container.style.alignItems = "center";
    container.style.height = "100vh";
    container.style.flexDirection = "column";

    // Ajouter un titre
    const title = document.createElement("h1");
    title.innerText = "Inscription au Forum";
    container.appendChild(title);

    // Créer un formulaire d'inscription
    const form = document.createElement("form");

    // Champs du formulaire
    const fields = [
        { label: "Email", type: "email", name: "email", placeholder: "Entrez votre email" },
        { label: "Mot de passe", type: "password", name: "password", placeholder: "Entrez votre mot de passe" },
        { label: "Nom d'utilisateur", type: "text", name: "username", placeholder: "Choisissez un pseudo" },
        { label: "Prénom", type: "text", name: "first_name", placeholder: "Votre prénom" },
        { label: "Nom", type: "text", name: "last_name", placeholder: "Votre nom" },
        { label: "Âge", type: "number", name: "age", placeholder: "Votre âge" },
    ];

    // Création des champs
    const inputs = {};
    fields.forEach(({ label, type, name, placeholder }) => {
        const fieldLabel = document.createElement("label");
        fieldLabel.innerText = label + " : ";
        const fieldInput = document.createElement("input");
        fieldInput.type = type;
        fieldInput.placeholder = placeholder;
        fieldInput.required = true;
        fieldInput.name = name;
        inputs[name] = fieldInput;

        form.appendChild(fieldLabel);
        form.appendChild(fieldInput);
        form.appendChild(document.createElement("br"));
    });

    // Sélecteur de genre
    const genreLabel = document.createElement("label");
    genreLabel.innerText = "Genre : ";
    const genreSelect = document.createElement("select");
    genreSelect.name = "genre";
    ["Homme", "Femme", "Autre"].forEach(optionText => {
        const option = document.createElement("option");
        option.value = optionText;
        option.innerText = optionText;
        genreSelect.appendChild(option);
    });

    inputs["genre"] = genreSelect;

    form.appendChild(genreLabel);
    form.appendChild(genreSelect);
    form.appendChild(document.createElement("br"));

    // Bouton d'inscription
    const submitButton = document.createElement("button");
    submitButton.innerText = "S'inscrire";
    submitButton.type = "submit";

    form.appendChild(submitButton);

    // Conteneur des boutons
    const buttonContainer = document.createElement("div");
    buttonContainer.style.display = "flex";
    buttonContainer.style.flexDirection = "column";
    buttonContainer.style.alignItems = "center";
    buttonContainer.style.marginTop = "10px";

    // Bouton de connexion
    const loginButton = document.createElement("button");
    loginButton.innerText = "Se connecter";
    loginButton.type = "button";
    loginButton.style.marginTop = "8px";
    loginButton.style.backgroundColor = "transparent";
    loginButton.style.color = "#050505"; // Bleu Twitter
    loginButton.style.border = "none";
    loginButton.style.cursor = "pointer";
    loginButton.style.fontWeight = "bold";

    loginButton.addEventListener("click", () => {
        login(); // Rediriger vers la page de connexion
    });

    buttonContainer.appendChild(loginButton);

    // Ajouter les boutons sous le formulaire
    container.appendChild(form);
    container.appendChild(buttonContainer);

    // Ajouter l'événement d'inscription
    form.addEventListener("submit", (event) => {
        event.preventDefault(); // Empêcher le rechargement de la page

        const userData = Object.fromEntries(
            Object.entries(inputs).map(([key, input]) => [key, input.value])
        );

        sendRegister(userData);
    });

    // Ajouter tout au body
    document.body.appendChild(container);
};

const sendRegister = async (userData) => {
    try {
        const response = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(userData),
        });

        const data = await response.json();

        if (response.status === 201) {
            router(); // Redirection après succès
        } else if (response.status === 409) {
            alert(`Erreur (${data.code}) : ${data.message}`);
        } else {
            throw new Error(`Erreur inattendue : ${response.status}`);
        }
    } catch (e) {
        console.error("Erreur lors de l'inscription :", e);
        alert("Une erreur est survenue lors de l'inscription.");
    }
};