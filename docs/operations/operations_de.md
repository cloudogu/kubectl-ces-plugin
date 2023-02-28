# Informationen zur Installation und Verwendung von kubectl-ces

## Installation

1. [Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) installieren
1. Das Plugin mit `kubectl krew install ces` installieren

Danach ist das Plug-in betriebsbereit.

## Verwendung

### Ändern der dogu-Konfigurationswerte

Auflisten der einstellbaren Konfigurationsschlüssel für einen bestimmten dogu:
`kubectl ces dogu-config ls <dogu-name>`

Interaktives Bearbeiten von Konfigurationsschlüsseln (kann auch Werte validieren, wenn die Validierung von Werten unterstützt wird):
`kubectl ces dogu-config edit <dogu-name>`

Abrufen eines Konfigurationswertes für einen gegebenen dogu und einen Konfigurationsschlüssel:
kubectl ces dogu-config get <dogu-name> <key>`

Setzen eines Konfigurationswertes für eine gegebene dogu und einen Konfigurationsschlüssel:
`kubectl ces dogu-config set <dogu-name> <key> <value>`

Einen Konfigurationswert für eine bestimmte dogu und einen Konfigurationsschlüssel entfernen (dadurch wird auch der Schlüssel entfernt):
`kubectl ces dogu-config delete <dogu-name> <key>`
