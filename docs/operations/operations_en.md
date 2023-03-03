# Information about installing and using kubectl-ces

## Installation

1. Download the latest version of [`kubectl-ces`](https://github.com/cloudogu/kubectl-ces-plugin/releases) according to your operating system.
2. extract the archive (either like this, or use a UI tool of your preference):
   - Linux: `tar -xvzf kubectl-ces_linux_amd64.tar.gz`
   - Darwin: `tar -xvzf kubectl-ces_darwin_amd64.tar.gz`
   - Windows: `unzip kubectl-ces_darwin_amd64.tar.gz`
3. Copy the extracted binary to your executional path

After that, the plug-in is operational.

## Usage

### Modify dogu configuration values

List settable configuration keys for a given dogu:
`kubectl ces dogu-config ls <dogu-name>`

Interactively edit configuration keys (may also validate values when validation of values is supported):
`kubectl ces dogu-config edit <dogu-name>`

Fetch a configuration value for a given dogu and a configuration key:
`kubectl ces dogu-config get <dogu-name> <key>`

Set a configuration value for a given dogu and a configuration key:
`kubectl ces dogu-config set <dogu-name> <key> <value>`

Remove a configuration value for a given dogu and a configuration key (this removes the key as well):
`kubectl ces dogu-config delete <dogu-name> <key>`
