# Information about installing and using kubectl-ces

## Installation

1. Install [Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)
1. Run `kubectl krew install ces` to install the plug-in

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
