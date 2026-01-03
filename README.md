# Battleship Commander

Un jeu de bataille navale en ligne de commande (TUI) écrit en Go. Ce projet permet de jouer contre des pairs en réseau via une architecture P2P HTTP.

## Fonctionnalités

- Jeu en réseau peer-to-peer (HTTP)
- Interface utilisateur en terminal (TUI)
- Système de chat intégré
- Événements aléatoires (Kraken)
- Bonus de tirs : Croix, Cercle, Ligne
- Notifications en temps réel (Touche, Rate, Coule)

## Installation

Prérequis : Go 1.20 ou supérieur.

Compiler le projet :

    go build -o battleship .

## Utilisation

Lancer le jeu en spécifiant un port :

    ./battleship -port 8080 -name "Joueur1"

Pour jouer contre un adversaire, lancez une seconde instance dans un autre terminal :

    ./battleship -port 8081 -name "Joueur2" -peers http://localhost:8080

**Note importante :** Si le Joueur 2 ajoute le Joueur 1 au démarrage, le Joueur 1 doit également ajouter le Joueur 2 manuellement pour pouvoir lui tirer dessus :

    add http://localhost:8081

Les pairs peuvent être ajoutés au démarrage avec le flag `-peers` ou dynamiquement en jeu avec la commande `add`.

## Commandes en jeu

- `x y` : Effectuer un tir simple aux coordonnées indiquées (ex: 5 5).
- `croix x y` : Effectuer un tir en croix (centre + 4 cases adjacentes).
- `cercle x y` : Effectuer un tir en cercle (8 cases autour du centre).
- `ligne x1 y1 x2 y2` : Effectuer un tir en ligne (longueur max 4 cases).
- `chat <message>` : Envoyer un message à l'adversaire.
- `add <url>` : Ajouter manuellement un adversaire en cours de partie.
- `target <url>` : Changer de cible active (si plusieurs adversaires).
- `q` : Quitter le jeu.

## Règles

1. Chaque joueur possède une grille de 10x10.
2. Les navires sont placés aléatoirement au démarrage.
3. Le but est de couler tous les navires de l'adversaire.
4. Des événements aléatoires (Attaque du Kraken) peuvent survenir périodiquement.

## Structure du Projet

- `main.go` : Point d'entrée de l'application.
- `client/` : Logique du client HTTP et gestion de l'interface utilisateur (TUI).
- `server/` : Serveur HTTP pour la réception des tirs et messages.
- `game/` : Modèles de données et logique métier (Plateau, Navires).
