# Developing `kubectl-ces`

## Encouraged Command Structure

The used CLI framework [github.com/spf13/cobra](https://github.com/spf13/cobra) provides a handy convention on how to structure a CLI commands:
> APPNAME VERB NOUN --ADJECTIVE

Anyhow, this plug-in will most likely deal with a large enough bandwidth of use cases for this plugin to enhance this command structure. To provide a matching narrative for different use cases the structure adds a GROUP command:
> kubectl ces GROUP VERB NOUN --ADJECTIVE

This still keeps the command structure easily understandable and enables use cases like these (some may not be implemented yet):

```
kubectl ces dogu-config edit redmine
kubectl ces dogu-config delete redmine logging/root
kubectl ces global-config edit
kubectl ces dogu detail-remote official/redmine 1.2.3-4
kubectl ces dogu install official/redmine 1.2.3-4
kubectl ces dogu delete redmine
kubectl ces dogu stop redmine
kubectl ces dogu start redmine
kubectl ces dogu restart redmine
kubectl ces dogu upgrade redmine 2.3.4-5
kubectl ces blueprint create "Blueprint 6.2.1"
kubectl ces blueprint upgrade "Blueprint 6.2.1"
kubectl ces blueprint delete "Blueprint 6.2.1"
kubectl ces backup list
kubectl ces backup create --all
kubectl ces backup delete "15d8eb48" --force -y --dont-ask-again
kubectl ces backup restore "15d8eb48" -y
kubectl ces instance register
```
