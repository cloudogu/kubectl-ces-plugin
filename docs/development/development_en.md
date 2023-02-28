# Developing `kubectl-ces`

## Encouraged Command Structure

```
kubectl ces "group" "verb" "noun" "adjective"
kubectl ces dogu-config edit redmine
kubectl ces dogu-config delete redmine <key>
kubectl ces global-config edit
kubectl ces dogu install official/redmine
kubectl ces dogu detail official/redmine 1.2.3-4
kubectl ces dogu delete redmine
kubectl ces dogu stop redmine
kubectl ces dogu start redmine
kubectl ces dogu restart redmine
kubectl ces dogu upgrade redmine <version>
kubectl ces blueprint create "SI 6.2.1"
kubectl ces blueprint upgrade "SI 6.2.1"
kubectl ces blueprint delete <blueprint>
kubectl ces backup create --all
kubectl ces backup delete "15d8eb48" --force -y --dont-ask-again
kubectl ces backup list
kubectl ces backup restore "15d8eb48" -y
kubectl ces instance register
```