 # Pub/Sub 📬
Imagine une boîte aux lettres intelligente :

Un **publisher** (expéditeur) envoie un message dans une "boîte" appelée topic
Plusieurs **subscribers** (abonnés) peuvent recevoir une copie de ce message
Les abonnés n'ont pas besoin d'être connectés au même moment

Exemple concret :
Tu as un site e-commerce. Quand un client passe commande :

Le service "Commandes" publie un message "Nouvelle commande #123"
Le service "Email" reçoit ce message → envoie un email de confirmation
Le service "Facturation" reçoit ce message → génère une facture
Le service "Stock" reçoit ce message → met à jour l'inventaire

Chacun fait son travail indépendamment, sans bloquer les autres.

# Eventarc 🎯
C'est un routeur d'événements intelligent :

Il écoute ce qui se passe sur GCP (un fichier uploadé, une base de données modifiée, etc.)
Il déclenche automatiquement une action (lance une fonction, une appli...)

Exemple concret :

Un utilisateur uploade une photo sur Cloud Storage
Eventarc détecte cet événement
Il déclenche automatiquement une Cloud Function qui redimensionne l'image

 **Différence avec Pub/Sub :** Eventarc est spécialisé pour connecter les services GCP entre eux, tandis que Pub/Sub est plus généraliste pour tout type de messaging.