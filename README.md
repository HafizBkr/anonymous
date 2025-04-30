
# Anonymous - Backend API

## Description

**Anonymous** est une API REST construite en Go avec le framework `chi` qui alimente une plateforme sociale anonyme. Elle permet aux utilisateurs de :
- S’inscrire, se connecter et gérer leur compte (l'email et les mots de passe sont crypter)
- Créer, liker, commenter et réagir à des publications
- Rejoindre et gérer des communautés
- Envoyer des messages en privé ou dans des chats communautaires
- Donner des points à des profils

Le backend est modulaire et utilise PostgreSQL comme base de données.

---

## Structure du projet

| Package          | Description |
|------------------|-------------|
| `auth`           | Gestion de l'inscription, connexion et vérification d'email |
| `users`          | Gestion des utilisateurs (profil, mot de passe, statut) |
| `posts`          | Gestion des publications et interactions (likes, réactions) |
| `comments`       | Commentaires sur les publications |
| `replies`        | Réponses aux commentaires |
| `comunauter`     | Création, consultation et adhésion aux communautés |
| `chat`           | Messagerie privée entre utilisateurs |
| `communitychats` | Discussions dans les communautés |
| `points`         | Système de points pour les profils |
| `search_algorithm` | Fonction de recherche d’utilisateurs |
| `middlewares`    | Middlewares d’authentification |
| `provider`       | JWT & transactions |
| `postgres`       | Connexion à PostgreSQL |

---

## Endpoints principaux

### Authentification

- `POST /auth/register` – Inscription
- `POST /auth/login` – Connexion
- `GET /auth/verify-email` – Vérification d'email
- `POST /auth//forgot-password` – Mot de passe oublier
- `POST /aut/reset-password` – Renitialiser le mot de passe 

### Utilisateurs

- `GET /users/` – Récupérer tous les utilisateurs
- `PATCH /users/status` – Activer/désactiver un compte
- `PATCH /users/password` – Changer le mot de passe
- `GET /users/{userID}` – Détails d’un utilisateur

### Publications

- `POST /posts/` – Créer une publication
- `GET /posts/` – Liste de toutes les publications
- `PATCH /posts/{postID}` – Modifier une publication
- `DELETE /posts/{postID}` – Supprimer une publication
- `POST /posts/{postID}/like` – Liker une publication
- `DELETE /posts/{postID}/like` – Retirer le like
- `POST /posts/{postID}/reaction` – Ajouter une réaction
- `DELETE /posts/{postID}/reaction` – Retirer une réaction

### Commentaires & Réponses

- `POST /{postID}/comments/` – Commenter un post
- `GET /{postID}/comments/` – Liste des commentaires
- `PUT /{commentID}/reactions` – Réagir à un commentaire
- `POST /{commentID}/replies/` – Répondre à un commentaire

### Communautés

- `POST /comunity/` – Créer une communauté
- `POST /comunity/{communityID}` – Rejoindre une communauté
- `GET /comunity/` – Lister les communautés
- `GET /comunity/{communityID}` – Détails d’une communauté

### Chat

- `GET /chat/conversations` – Obtenir ses conversations
- `POST /chat/` – Envoyer un message privé
- `GET /chat/messages/{user1ID}/{user2ID}` – Messages entre deux utilisateurs
- `GET /chat/messages/owner` – Tous les messages de l’utilisateur connecté

### Chats communautaires

- `POST /community_chats/{communityID}/messages` – Envoyer un message
- `GET /community_chats/{communityID}/messages` – Lire les messages

### Points

- `POST /points/` – Donner un point à un utilisateur
- `GET /points/{userID}` – Nombre de points reçus

### Recherche

- `GET /search` – Rechercher un utilisateur

---

## Variables d’environnement

- `PORT` : Port du serveur HTTP (ex: `8080`)
- `DB_URL` : URL de connexion à PostgreSQL
- Autres : variables SMTP si envoi d’emails activé

---

## Lancer le serveur

1. Crée un fichier `.env` à la racine :
   ```
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
   ```

2. Lance le backend :
   ```bash
   go run main.go
   ```

---

## Statique

Les fichiers statiques sont servis depuis le dossier `./static` via le chemin `/static/*`.

---

## Surveillance

- Endpoint `GET /health` : vérifie si le serveur est en ligne.
