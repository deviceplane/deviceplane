#### Build the binary

```
make build DEBUG=1
```

#### Run the binary

```
export DP_PRIVS_TOKEN="..."
export DP_DB_CONN="..."
export DP_SMTP_SERVER="..."
export DP_SMTP_PW="..."
export DP_AUTH0_AUD="..."
export DP_AUTH0_DOM="..."
script/run.bash
```

#### Run the binary as intercept to Kubernetes pod

Install telepresence

```
brew install datawire/blackbird/telepresence
```

Intercept

```
telepresence connect
telepresence intercept deviceplane -n deviceplane --port 80
```

Disconnect

```
telepresence leave deviceplane-deviceplane
telepresence quit
telepresence uninstall --everything
```
