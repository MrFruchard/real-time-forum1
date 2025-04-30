import { router } from "../router.js";
import { register } from "./register.js"; // Importer la fonction register

export const login = () => {
  // Vider complètement le `body`
  document.body.innerHTML = "";

  // Créer un conteneur principal
  const container = document.createElement("div");
  container.id = "login-container";
  container.style.display = "flex";
  container.style.justifyContent = "center";
  container.style.alignItems = "center";
  container.style.height = "100vh";
  container.style.flexDirection = "column";

  // Ajouter un titre
  const title = document.createElement("h1");
  title.innerText = "Connexion au Forum";
  container.appendChild(title);

  // Créer un formulaire de connexion
  const form = document.createElement("form");

  // Champ email
  const emailLabel = document.createElement("label");
  emailLabel.innerText = "Email : ";
  const emailInput = document.createElement("input");
  emailInput.placeholder = "Entrez votre email ou username";
  emailInput.required = true;

  // Champ mot de passe
  const passLabel = document.createElement("label");
  passLabel.innerText = "Mot de passe : ";
  const passInput = document.createElement("input");
  passInput.type = "password";
  passInput.placeholder = "Entrez votre mot de passe";
  passInput.required = true;

  // Bouton de connexion
  const submitButton = document.createElement("button");
  submitButton.innerText = "Se connecter";
  submitButton.type = "submit";

  // Bouton d'inscription
  const registerButton = document.createElement("button");
  registerButton.innerText = "Créer un compte";
  registerButton.type = "button";
  registerButton.style.marginTop = "10px";
  registerButton.style.backgroundColor = "transparent";
  registerButton.style.color = "#000000"; // Bleu Twitter
  registerButton.style.border = "none";
  registerButton.style.cursor = "pointer";
  registerButton.style.fontWeight = "bold";

  registerButton.addEventListener("click", () => {
    register(); // Rediriger vers register()
  });

  // Ajout des éléments au formulaire
  form.appendChild(emailLabel);
  form.appendChild(emailInput);
  form.appendChild(document.createElement("br"));
  form.appendChild(passLabel);
  form.appendChild(passInput);
  form.appendChild(document.createElement("br"));
  form.appendChild(submitButton);
  form.appendChild(registerButton); // Ajout du bouton d'inscription

  // Ajout du formulaire au conteneur
  container.appendChild(form);

  // Ajouter l'événement de connexion
  form.addEventListener("submit", (event) => {
    event.preventDefault(); // Empêcher le rechargement de la page
    sendLog(emailInput.value, passInput.value);
  });

  // Ajouter tout au body
  document.body.appendChild(container);
};

const sendLog = async (email, mdp) => {
  try {
    const response = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email, password: mdp }),
    });

    if (!response.ok) {
      throw new Error(`Erreur : ${response.status} ${response.statusText}`);
    }
    router();
  } catch (e) {
    console.error("Erreur lors de la connexion :", e);
    alert("Une erreur est survenue lors de la connexion.");
  }
};