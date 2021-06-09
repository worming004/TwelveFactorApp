# 12 factors app, une méthode pour ranger ses micro-services

## Des principes cloud friendly

**Naïvement**, voici les composants et la famille de bonne pratique associée
![zones](../naive-best-practice-zones.png)

### Cloud over simplifié

![over-simplifiés](../over-simplest-cloud.png)

## Les principes

### I. Base de code unique

* Une seule base de code
* Un build, plusieurs déploiements

### II. Des dépendances explicites

* téléchargeables, ne dépendent pas de l'ordinateur (éviter GAC, utilisation d'un VS spécifique)

### III. Configuration dans l'environnement

* Configurations et secrets sont enregistrés en tant que variable d'environnement
* Avoir un fichier de configuration est acceptable, mais non recommandé

### IV. Les services externes sont à traiter comme des ressources locales

* Changer une chaîne de connexion, de type de base de données (mysql -> sql server) ne nécessite pas de changement de code

### V. Séparez les étapes d'assemblage, publication et d'éxécution

* Un build
* L'assemblage des configurations se fait en amont
* Le runtime s'exécute plus tard

### VI. Des processus sans états

* Share Nothing architecture. Les instances sont indépendantes en utilisation de CPU, mémoire et disque
* Pas de variable de session, ou de 'sticky session'

### VII. Un service est associé à un port

* Identifiez vos services via le numéro de port

### VIII. La concurrence se scale horizontalement

* Ayez une architecture qui scale horizontalement. Un service est plus facile à scale en ajoutant des machines qu'en les grossissant

### IX. Des instances jetables

* Les services ont pour intention d'être constamment déplacé. Ils peuvent s'allumer et s'éteindre a tout moment
* Les services doivent démarrer le plus vite possible, et s'éteindre gracieusement. Ils devront terminer leur requête, et ne plus accepter les nouvelles

### X. Parité dev prod

* Devops friendly, les applications doivent être codées et déployées par la même équipe
* Les différents environnements doivent utiliser la même infrastructure, et la même pile technologique

### XI. Loggez en tant que stream

* Les logs doivent être traités en tant que stream dans stdout
* Pas de tampon ni de gestion de fichier local sur le serveur applicatif
* La persistance des logs se fait par l'environnement d'éxécution

### XII. Les tâches d'administration sont gérées comme tout autre process

* Les mises à jour des schémas sql ou n'importe quelle tâche 'one shot' doivent suivre le cycle de vie des 12 factors app. Ils doivent suivre les 11 premiers facteurs
