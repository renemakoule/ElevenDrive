# Architecture de Synchronisation & Ergonomie Guichet
## Logiciel Hybride de Gestion d'Auto-École

---

## 1. Architecture de Synchronisation Détaillée (Offline-First)

Dans un environnement où les coupures réseau peuvent durer de quelques minutes à plusieurs jours, la synchronisation ne doit jamais être bloquante pour les utilisateurs en agence.

```mermaid
graph TD
    subgraph "Poste Desktop (Agence)"
        UI[Interface Secrétaire / Caisse] -->|Action locale| LDB[(SQLite Local)]
        UI -->|Génère un événement immuable| Q[File locale: journal_sync_events]
        SyncWorker[Background Sync Engine (Tauri/Rust)]
        SyncWorker <-->|Lit les événements 'EN_ATTENTE'| Q
        SyncWorker <-->|Met à jour les données| LDB
    end

    subgraph "Réseau Intermittent"
        NetworkDetector{Détection Réseau & Ping Cloud}
        SyncWorker --> NetworkDetector
    end

    subgraph "Serveur Central (Cloud)"
        Gateway[API Gateway / Ingestion Hub]
        ConflictResolver[Moteur d'Arbitrage & Déconfliction]
        CDB[(PostgreSQL Central)]
        
        Gateway --> ConflictResolver
        ConflictResolver <--> CDB
    end

    NetworkDetector -->|Connexion active : Envoi par lots / Batch| Gateway
    ConflictResolver -->|Accusé de réception (ACK) + Deltas distants| SyncWorker
```

### 1.1. Le Protocole de Synchronisation par Événements (Sync Queue)

Chaque modification effectuée en local (inscription, versement d'acompte, présence au code) génère un objet d'événement dans `journal_sync_events` :
```json
{
  "event_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
  "agency_id": "ag_cocody_01",
  "entity_name": "paiements",
  "operation": "INSERT",
  "entity_id": "c8a4a58b-9e4a-4a0e-bc21-7e8e19c3b8a1",
  "payload": {
    "dossier_id": "dos_2026_042",
    "montant": 50000.00,
    "mode_reglement": "ESPECES",
    "reference_recu": "REC-2026-10-0089"
  },
  "local_timestamp_utc": "2026-10-19T14:45:10Z",
  "sync_status": "EN_ATTENTE"
}
```

### 1.2. Matrice de Résolution des Conflits

Comment le système réagit-il lorsque deux personnes modifient la même donnée ou partagent la même ressource en mode déconnecté ?

| Type de Donnée | Exemple de Conflit | Règle d'Arbitrage Automatique | Résultat & Comportement UX |
| :--- | :--- | :--- | :--- |
| **A. Données Financières (Paiements & Reçus)** | La secrétaire encaisse un acompte en espèces au guichet hors-ligne pendant que l'élève paie par Mobile Money sur le Web. | **Append-Only (Addition pure)** : Aucun conflit possible. | Les deux versements sont acceptés. Le Cloud additionne les deux montants. Le solde restant diminue de la somme des deux. Zéro litige. |
| **B. Ressources Uniques (Véhicule ou Créneau Conduite)** | La secrétaire réserve la Peugeot 208 pour Amadou mardi à 10h hors-ligne. En même temps sur le Web, le moniteur réserve la même voiture pour Fatou. | **Premier validé par le Cloud gagne (First-Committed Wins) + Alerte** | La première transaction arrivée au Cloud est validée. La seconde est refusée avec le statut `CONFLIT_CRENEAU_OCCUPE`. Une alerte visuelle s'affiche au guichet pour réassigner le créneau. |
| **C. Fiche Élève (État Civil / Contact)** | La secrétaire corrige l'adresse de l'élève en local, tandis que le manager met à jour son numéro de téléphone sur son laptop. | **Fusion au niveau des champs (Field-Level Merge) + Timestamp UTC** | Les deux modifications sont fusionnées (la nouvelle adresse ET le nouveau numéro sont conservés). Si le même champ est touché, la modification la plus récente en UTC l'emporte avec traçabilité dans l'historique. |

---

## 2. Conception UX : L'Écran Guichet de la Secrétaire (100% Clavier)

Le guichet d'une auto-école en Afrique est un lieu à fort trafic. La secrétaire doit pouvoir enrôler un élève ou encaisser une tranche **sans jamais lever les mains de son clavier**.

### 2.1. Wireframe Conceptuel de l'Écran Principal Guichet (Vue Desktop 1920x1080)

```text
+---------------------------------------------------------------------------------------------------------------------------------------+
|  [LOGO AUTO-ECOLE] AUTO-ÉCOLE LA RÉFÉRENCE | Agence Cocody | Session Caissière: Aminata K.       ● CONNECTÉ (Synchro: OK)   [14:32]   |
+---------------------------------------------------------------------------------------------------------------------------------------+
|  RACCOURCIS RAPIDES :  [F1] Nouvelle Inscription   [F2] Encaissement Rapide   [F3] Recherche Dossier   [F4] Point Caisse Z   [F5] Code  |
+---------------------------------------------------------------------------------------------------------------------------------------+
|                                                                   |                                                                   |
|   PANEL GAUCHE : RECHERCHE & DOSSIER ACTIF                        |   PANEL DROIT : OPÉRATION ACTIVE & APERÇU REÇU                    |
|                                                                   |                                                                   |
|   Recherche Rapide [F3 ou Alt+S] :                                |   [MODE ACTIF : ENCAISSEMENT D'UNE TRANCHE]                       |
|   +-----------------------------------------------------------+   |                                                                   |
|   | 07089921.. ou "KOUASSI" ou Scanner QR Code du Reçu        |   |   Dossier : REF-2026-B-0042                                       |
|   +-----------------------------------------------------------+   |   Candidat: KOUASSI Yao Jean-Luc                                  |
|                                                                   |   Forfait : Permis B Standard (Total: 150 000 FCFA)               |
|   Fiche Candidat Sélectionné :                                    |                                                                   |
|   +-----------------------------------------------------------+   |   SITUATION FINANCIÈRE :                                          |
|   | [PHOTO]   KOUASSI Yao Jean-Luc           Permis B         |   |   +--------------------+--------------------+-----------------+    |
|   | 4x4       Né le 12/05/1998 à Bouaké      Dossier #0042    |   |   | Total Forfait      | Déjà Encaissé      | RESTE À PAYER   |    |
|   |           Tél WhatsApp: +225 07 08 99 21 34               |   |   | 150 000 FCFA       | 50 000 FCFA        | 100 000 FCFA    |    |
|   |           Pièce: CNI CI00389218 (Valide)                  |   |   +--------------------+--------------------+-----------------+    |
|   |           Médical: Apte sans réserve                      |   |                                                                   |
|   +-----------------------------------------------------------+   |   FORMULAIRE D'ENCAISSEMENT [Tab pour naviguer] :                 |
|                                                                   |                                                                   |
|   Progression Pédagogique :                                       |   1. Montant versé aujourd'hui :                                  |
|   - Code de la route : 14/20 heures effectuées (Assidu)           |      +------------------------------------+                       |
|   - Conduite pratique: 0/20 heures (Verrou: Payer Tranche 2)      |      | 50 000 FCFA                        |                       |
|                                                                   |      +------------------------------------+                       |
|   Historique des Versements :                                     |   2. Mode de Règlement :                                          |
|   • 10/10/2026 : 50 000 FCFA (Espèces) - Reçu REC-00389           |      (•) 1. Espèces   ( ) 2. Wave   ( ) 3. Orange MoMo            |
|                                                                   |   3. Espèces Remises par le Client :                              |
|                                                                   |      +-----------------+   -> RENDU DE MONNAIE CALCULÉ :          |
|                                                                   |      | 50 000 FCFA     |      [ 0 FCFA ]                          |
|                                                                   |      +-----------------+                                          |
|                                                                   |                                                                   |
|                                                                   |   [APERÇU TICKET THERMIQUE 80mm - PRÊT POUR IMPRESSION]           |
|                                                                   |   +-----------------------------------------------------------+   |
|                                                                   |   | AUTO-ÉCOLE LA RÉFÉRENCE - REÇU N° REC-2026-10-0412        |   |
|                                                                   |   | Date: 19/10/2026 14:32:15  | Caissier: Aminata K.         |   |
|                                                                   |   | Élève: KOUASSI Yao Jean-Luc (REF-2026-B-0042)             |   |
|                                                                   |   | Montant Versé : 50 000 FCFA (Cinquante mille)             |   |
|                                                                   |   | Reste Dû      : 50 000 FCFA                               |   |
|                                                                   |   | [QR CODE SÉCURISÉ D'AUTHENTICITÉ]                         |   |
|                                                                   |   +-----------------------------------------------------------+   |
|                                                                   |                                                                   |
|                                                                   |   ACTION FINALE :                                                 |
|                                                                   |   [ ENTRÉE ou F10 ] : VALIDER L'ENCAISSEMENT & IMPRIMER TICKET    |
+---------------------------------------------------------------------------------------------------------------------------------------+
|  BARRE D'ÉTAT BASSE : F1: Inscrire | F2: Encaisser | F4: Fermer Caisse Z | Échap: Annuler | Espace: Prendre Photo Webcam               |
+---------------------------------------------------------------------------------------------------------------------------------------+
```

---

### 2.2. Cinématique 100% Clavier d'un Encaissement (15 Secondes Chrono)

Voici la séquence de touches exécutée par la secrétaire sans jamais toucher à la souris :

1. **`F2`** : Ouvre instantanément le mode *Encaissement Rapide*. Le curseur se positionne automatiquement dans la barre de recherche.
2. **Taper `0708` ou scanner le QR code** du ticket précédent :
   - La liste d'autocomplétion filtre en temps réel.
   - **`Touche Bas` + `Entrée`** : Le dossier de l'élève est chargé. Son solde restant (`100 000 FCFA`) clignote en orange.
3. Le curseur saute immédiatement dans le champ **"Montant versé"** :
   - La secrétaire tape `50000` puis presse **`Tab`**.
4. Mode de règlement :
   - Par défaut sélectionné sur `1. Espèces`. (Taper `2` pour Wave ou `3` pour Orange Money).
   - Touche **`Tab`** vers "Espèces remises par le client".
5. Si le client donne un billet de 10 000 pour un montant de 8 000 :
   - Elle tape `10000` -> Le système affiche en grand et en vert : **RENDU : 2 000 FCFA**.
6. **`Entrée` (ou `F10`)** :
   - L'enregistrement SQL est exécuté en 5 millisecondes.
   - L'imprimante thermique USB reçoit la commande directe ESC/POS et imprime le ticket en 1 seconde.
   - Une impulsion électrique ouvre automatiquement le tiroir-caisse physique.
   - L'écran émet un bip discret de confirmation et se réinitialise pour le client suivant.

---

### 2.3. Cinématique de Prise de Photo Webcam à l'Enrôlement (Écran F1)

Lorsqu'un nouveau candidat s'inscrit :
1. La secrétaire remplit l'état civil au clavier avec des tabulations automatiques.
2. Arrivée à l'étape "Photo d'identité" :
   - La fenêtre active la webcam connectée en direct.
   - Une mire ovale guide le candidat pour centrer son visage.
   - La secrétaire appuie sur la **`Barre d'Espace`** :
     - La photo est capturée instantanément.
     - L'algorithme recadre automatiquement le visage au ratio 4x4 officiel.
   - La secrétaire tape **`Entrée`** pour valider et passer aux pièces justificatives.
   - **Zéro logiciel externe, zéro souris nécessaire.**

---

### 2.4. Indicateurs Visuels de Résilience Réseau

Dans le bandeau supérieur de l'application :
* 🟢 **Vert (`Connecté`)** : Données synchronisées en temps réel avec le Cloud.
* 🟠 **Orange clignotant (`Hors-Ligne - 4 actions en attente`)** :
  - L'application continue à fonctionner à 100% de sa vitesse.
  - La secrétaire voit le nombre d'opérations locales en attente.
  - Dès que le signal Internet revient, l'icône tourne pendant 2 secondes et redevient verte : *"Synchronisation terminée (4 éléments transmis)"*.
* 🔴 **Rouge (`Alerte Périphérique`)** :
  - Ex : *"Imprimante ticket hors-ligne ou plus de papier thermique"* avec invite immédiate : *"Réessayer [Entrée] ou imprimer en PDF A4 [F8]"*.
