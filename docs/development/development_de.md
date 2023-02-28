# `kubectl-ces` entwickeln

## Empfohlene Befehlsstruktur

Das verwendete CLI-Framework [github.com/spf13/cobra](https://github.com/spf13/cobra) bietet eine praktische Konvention für die Strukturierung von CLI-Befehlen:
> APPNAME VERB NOUN --ADJECTIVE

Allerdings wird dieses Plugin höchstwahrscheinlich mit einer ausreichenden Bandbreite an Anwendungsfällen zu tun haben, als dass die genannte Befehlsstruktur dafür ausreicht. Um eine passende Entsprechung für verschiedene Anwendungsfälle zu liefern, wird zusätzlich die Command-Struktur um einen GROUP-Befehl erweitert:
> kubectl ces GROUP VERB NOUN --ADJECTIVE

Dies ermöglicht es, eine leicht verständliche Befehlsstruktur beizubehalten, sodass Anwendungsfälle wie diese (einige sind möglicherweise noch nicht implementiert) in einem erlernbaren Rahmen bleibt:

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
