# Real-Time Forum

Un forum moderne avec fonctionnalités de communication en temps réel, permettant aux utilisateurs de publier des posts, interagir via des commentaires et échanger des messages privés instantanés.

## 📋 Table des matières

- [Architecture](#architecture)
- [Fonctionnalités](#fonctionnalités)
- [Structure du projet](#structure-du-projet)
- [Prérequis](#prérequis)
- [Installation et démarrage](#installation-et-démarrage)
- [API REST](#api-rest)
- [WebSockets](#websockets)
- [Modèle de données](#modèle-de-données)
- [Sécurité](#sécurité)
- [Fonctionnement technique](#fonctionnement-technique)
- [Dépendances](#dépendances)

## 🏗️ Architecture

Le projet est construit selon une architecture client-serveur classique:

### Backend (Go)

- **Serveur HTTP**: Gère les requêtes API REST et sert les fichiers statiques
- **Base de données SQLite**: Stockage persistant des données
- **WebSockets**: Communication bidirectionnelle en temps réel

### Frontend (JavaScript)

- **Vanilla JavaScript**: Code modulaire sans framework
- **HTML/CSS**: Interface utilisateur responsive
- **WebSockets**: Communication en temps réel avec le serveur

## ✨ Fonctionnalités

### 🔐 Authentification

- **Inscription**: Création de compte avec validation des données
- **Connexion**: Authentification par email/username et mot de passe
- **Sessions**: Gestion par cookies HTTP-only
- **Déconnexion**: Suppression de session sécurisée

### 📝 Publications

- **Création de posts**: Interface pour publier du contenu avec catégorie
- **Affichage dynamique**: Liste des posts triés par date
- **Interactions sociales**: Système de likes/dislikes
- **Catégorisation**: Organisation des posts par thème

### 💬 Commentaires

- **Ajout de commentaires**: Possibilité de commenter chaque post
- **Interactions**: Likes/dislikes sur les commentaires
- **Notifications**: Alerte à l'auteur du post lors d'un nouveau commentaire

### 📨 Messagerie privée

- **Conversations en temps réel**: Échange instantané de messages
- **Historique des conversations**: Chargement paginé des messages
- **Indicateur de frappe**: Affichage "est en train d'écrire..."
- **Indicateurs de statut**: Visualisation des utilisateurs connectés

### 🔔 Système de notifications

- **Alertes en temps réel**: Notification pour les likes, dislikes, commentaires
- **Interface dédiée**: Panel de visualisation des notifications
- **Statut de lecture**: Distinction entre notifications lues/non lues

### 👥 Utilisateurs

- **Profils utilisateurs**: Informations de base sur les utilisateurs
- **Indicateurs de présence**: Liste des utilisateurs connectés
- **États de connexion**: Mise à jour en temps réel des statuts

## 📁 Structure du projet

```
real-time-forum/
├── back/                     # Code source backend (Go)
│   ├── handlers/             # Gestionnaires de requêtes HTTP
│   ├── models/               # Structures de données
│   ├── services/             # Logique métier
│   ├── utils/                # Fonctions utilitaires
│   ├── websocketFile/        # Implémentation WebSocket
│   ├── go.mod                # Dépendances Go
│   └── main.go               # Point d'entrée du serveur
├── front/                    # Code source frontend (JavaScript)
│   ├── index.html            # Page principale
│   └── static/               # Ressources statiques
│       ├── home-page/        # Composants page principale
│       ├── login-page/       # Pages d'authentification
│       ├── websocket-integration/ # Intégration WebSocket
│       ├── router.js         # Routeur client
│       └── styles.css        # Styles globaux
└── database/                 # Fichiers de base de données
    └── real-time-forum.db    # Base de données SQLite
```

## 🔧 Prérequis

- **Go** (version 1.20 ou supérieure)
- **SQLite** (version 3)
- **Navigateur web moderne** (Chrome, Firefox, Edge, Safari)

## 🚀 Installation et démarrage

### Configuration de l'environnement

1. Cloner le dépôt
   ```bash
   git clone https://github.com/username/real-time-forum.git
   cd real-time-forum
   ```

2. Installation des dépendances Go
   ```bash
   cd back
   go mod download
   ```

### Démarrage du serveur

```bash
cd back
go run main.go
```

Le serveur démarre par défaut sur le port 8080. Une fois lancé, vous pouvez accéder à l'application via:
```
http://localhost:8080
```

## 📡 API REST

L'application expose les endpoints suivants:

### Authentification

| Méthode | Endpoint    | Description                       | Corps de requête                                       | Réponse                                                 |
|---------|-------------|-----------------------------------|--------------------------------------------------------|---------------------------------------------------------|
| POST    | /register   | Inscription d'un nouvel utilisateur | `{email, password, username, first_name, last_name, age, genre}` | `{message: "Inscription réussie"}` + Cookie session     |
| POST    | /login      | Connexion utilisateur             | `{email, password}`                                   | `{message: "Connexion réussie"}` + Cookie session       |
| POST    | /logout     | Déconnexion                       | -                                                      | Suppression du cookie de session                        |

### Posts

| Méthode | Endpoint    | Description                       | Corps de requête                | Réponse                              |
|---------|-------------|-----------------------------------|----------------------------------|--------------------------------------|
| POST    | /post       | Création d'un post                | `{content, category, image}`    | `{status: "Created", message: "..."}`|
| GET     | /post       | Récupération de tous les posts    | -                                | `[{id, username, content, ...}]`     |
| GET     | /post/{id}  | Récupération d'un post spécifique | -                                | `{id, username, content, ...}`       |
| PUT     | /post/{id}  | Modification d'un post            | `{content, category}`           | `{message: "Post mis à jour"}`       |
| DELETE  | /post/{id}  | Suppression d'un post             | -                                | `{message: "Post supprimé"}`         |

### Commentaires

| Méthode | Endpoint      | Description                        | Corps de requête               | Réponse                             |
|---------|---------------|------------------------------------|---------------------------------|-------------------------------------|
| POST    | /comment      | Ajout d'un commentaire             | `{id_post, content}`           | Status 200 OK                       |
| GET     | /comment/{id} | Récupération des commentaires d'un post | -                         | `[{id, username, content, ...}]`    |

### Messages privés

| Méthode | Endpoint       | Description                        | Corps de requête                | Réponse                            |
|---------|----------------|------------------------------------|---------------------------------|------------------------------------|
| POST    | /message       | Envoi d'un message privé           | `{receiver, message}`           | `{message_id, conversation_id, ...}`|
| GET     | /message?user={username} | Récupération des messages d'une conversation | -            | `[{sender, message, date}]`        |
| GET     | /conversation  | Liste des conversations            | -                                | `[{username, last_message, ...}]`  |

### Interactions

| Méthode | Endpoint      | Description                   | Corps de requête                      | Réponse       |
|---------|---------------|-------------------------------|--------------------------------------|---------------|
| POST    | /event        | Like/dislike post ou commentaire | `{type, content_type, id}`         | Status 200 OK |
| GET     | /notification | Récupération des notifications | -                                    | `[{sender, type, ...}]` |
| POST    | /notification | Marquer comme lu              | `{id}`                               | Status 200 OK |

## 🔌 WebSockets

Le serveur expose un endpoint WebSocket à `/ws` pour la communication en temps réel:

### Types de messages WebSocket

| Type              | Direction      | Contenu                            | Description                        |
|-------------------|----------------|------------------------------------|------------------------------------|
| `get_user`        | Client → Serveur | -                                | Demande la liste des utilisateurs connectés |
| `connected_users` | Serveur → Client | `{content: [usernames]}`         | Liste des utilisateurs connectés   |
| `new_user`        | Serveur → Client | `{content: [username]}`          | Notification d'un nouvel utilisateur connecté |
| `user_disconnected` | Serveur → Client | `{content: [username]}`        | Notification d'un utilisateur déconnecté |
| `private`         | Serveur → Client | `{content: [message, sender, timestamp]}` | Réception d'un message privé |
| `is_typing`       | Client → Serveur | `{content: username}`            | Notification "est en train d'écrire" |
| `is_not_typing`   | Client → Serveur | `{content: username}`            | Arrêt de la notification d'écriture |
| `notification`    | Serveur → Client | -                                | Alerte de nouvelle notification    |

## 🗄️ Modèle de données

### Tables principales

#### USER
- **ID**: UUID unique (clé primaire)
- **EMAIL**: Email de l'utilisateur (unique)
- **PASSWORD**: Mot de passe hashé
- **USERNAME**: Nom d'utilisateur (unique)
- **FIRSTNAME**: Prénom
- **LASTNAME**: Nom de famille
- **AGE**: Âge
- **GENRE**: Genre
- **CREATED_AT**: Date de création du compte

#### SESSION
- **ID**: UUID unique (clé primaire)
- **USERID**: ID de l'utilisateur (clé étrangère)
- **CREATED_AT**: Date de création
- **LAST_ACTIVE_AT**: Dernière activité
- **EXPIRES_AT**: Date d'expiration

#### POST
- **ID**: ID auto-incrémenté (clé primaire)
- **USER_ID**: ID de l'auteur (clé étrangère)
- **CONTENT**: Contenu du post
- **CATEGORY**: Catégorie
- **IMAGE**: Chemin de l'image (optionnel)
- **CREATED_AT**: Date de création

#### COMMENT
- **ID**: ID auto-incrémenté (clé primaire)
- **USERID**: ID de l'auteur (clé étrangère)
- **POST_ID**: ID du post (clé étrangère)
- **CONTENT**: Contenu du commentaire
- **CREATED_AT**: Date de création

#### CONVERSATION
- **ID**: UUID unique (clé primaire)
- **USER1ID**: Premier participant (clé étrangère)
- **USER2ID**: Second participant (clé étrangère)
- **STARTED_AT**: Date de début

#### MESSAGE
- **ID**: ID auto-incrémenté (clé primaire)
- **CONVERSATION_ID**: ID de la conversation (clé étrangère)
- **SENDER_ID**: ID de l'expéditeur (clé étrangère)
- **CONTENT**: Contenu du message
- **SENT_AT**: Date d'envoi

#### NOTIFICATION
- **ID**: ID auto-incrémenté (clé primaire)
- **RECEIVER_ID**: ID du destinataire (clé étrangère)
- **SENDER_ID**: ID de l'expéditeur (clé étrangère)
- **TYPE**: Type de notification (like, comment, etc.)
- **RELATED_ID**: ID de l'objet concerné
- **STATUS**: État (read/unread)
- **CREATED_AT**: Date de création

### Tables d'interaction

- **LIKES**: Stocke les likes sur les posts
- **DISLIKE**: Stocke les dislikes sur les posts
- **LIKE_COMMENT**: Stocke les likes sur les commentaires
- **DISLIKE_COMMENT**: Stocke les dislikes sur les commentaires

## 🔒 Sécurité

### Authentification

- Mots de passe hashés avec **bcrypt** (cost=14)
- Sessions stockées côté serveur avec expiration
- Cookies HTTP-only pour prévenir le vol de session via XSS
- Option SameSite=Strict pour les cookies

### Validation des données

- Vérification des données utilisateur à l'inscription
- Validation des inputs côté serveur avant insertion en base
- Protection contre les injections SQL via requêtes préparées

### Contrôle d'accès

- Vérification des droits pour la modification/suppression de contenu
- Sessions invalidées à la déconnexion
- Vérification systématique de l'authentification pour les routes protégées

## 🔄 Fonctionnement technique

### Flux d'authentification

1. L'utilisateur soumet ses identifiants via le formulaire de connexion
2. Le serveur valide les informations et génère un token de session UUID
3. Le token est stocké en base avec une date d'expiration (24h par défaut)
4. Un cookie de session est envoyé au client
5. Les requêtes suivantes incluent automatiquement ce cookie pour l'authentification

### Communication en temps réel

1. Après authentification, une connexion WebSocket est établie
2. Le serveur maintient un `Hub` central qui gère toutes les connexions
3. Chaque client est associé à son username dans le Hub
4. Les messages et notifications sont routés via ce Hub vers les destinataires appropriés
5. Les états d'utilisateur (connexion, déconnexion, "est en train d'écrire") sont diffusés en temps réel

### Cycle de vie des posts et commentaires

1. Un utilisateur crée un post qui est immédiatement ajouté à la base de données
2. L'interface principale est actualisée pour afficher le nouveau post
3. Les commentaires sont chargés à la demande lors de l'expansion de la section commentaires
4. Les interactions (likes, dislikes) mettent à jour la base et sont reflétées immédiatement dans l'interface
5. Les auteurs reçoivent des notifications en temps réel pour les nouvelles interactions

### Messagerie privée

1. Un utilisateur clique sur un autre utilisateur dans la liste des connectés
2. Une modal de chat s'ouvre avec l'historique des messages chargé depuis le serveur
3. Les nouveaux messages sont envoyés via HTTP POST puis diffusés en temps réel via WebSocket
4. L'indicateur "est en train d'écrire" est activé pendant la saisie et désactivé après un délai d'inactivité
5. Les conversations sont automatiquement triées par date du dernier message

## 📦 Dépendances

### Backend (Go)

- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - Driver SQLite
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket) - Implémentation WebSocket
- [github.com/gofrs/uuid](https://github.com/gofrs/uuid) - Génération d'UUIDs
- [golang.org/x/crypto](https://golang.org/x/crypto) - Fonctions cryptographiques (bcrypt)

### Frontend (JavaScript)

- Vanilla JavaScript (ES6+)
- WebSocket API native
- Fetch API pour les requêtes HTTP
