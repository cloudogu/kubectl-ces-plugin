# Informationen zur Installation und Verwendung von kubectl-ces

## Installation

1. Laden Sie die neueste Version von [`kubectl-ces`](https://github.com/cloudogu/kubectl-ces-plugin/releases) entsprechend Ihrem Betriebssystem herunter.
2. Entpacken Sie das Archiv (entweder wie folgt oder mit einem UI-Tool Ihrer Wahl):
   - Linux: `tar -xvzf kubectl-ces_linux_amd64.tar.gz`
   - Darwin: `tar -xvzf kubectl-ces_darwin_amd64.tar.gz`
   - Windows: `unzip kubectl-ces_darwin_amd64.tar.gz`
3. Kopieren Sie die entpackte Binärdatei in Ihren Ausführungspfad

Danach ist das Plug-in betriebsbereit.

## Verwendung

### Ändern der dogu-Konfigurationswerte

Auflisten der einstellbaren Konfigurationsschlüssel für einen bestimmten dogu:
`kubectl ces dogu-config ls <dogu-name>`

Interaktives Bearbeiten von Konfigurationsschlüsseln (kann auch Werte validieren, wenn die Validierung von Werten unterstützt wird):
`kubectl ces dogu-config edit <dogu-name>`

Abrufen eines Konfigurationswertes für einen gegebenen dogu und einen Konfigurationsschlüssel:
`kubectl ces dogu-config get <dogu-name> <key>`

Setzen eines Konfigurationswertes für eine gegebene dogu und einen Konfigurationsschlüssel:
`kubectl ces dogu-config set <dogu-name> <key> <value>`

Einen Konfigurationswert für eine bestimmte dogu und einen Konfigurationsschlüssel entfernen (dadurch wird auch der Schlüssel entfernt):
`kubectl ces dogu-config delete <dogu-name> <key>`
