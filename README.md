# Service Watcher F5

Watches for services in kubernetes and creates corresponding F5 irules via a VIP and iRules.

## Settings

* `ALAMOAPI_SECRET` - The path in vault to the secret for accessing the region api. The secret in vault must contain the keys `location` (url to region api), `password` and `username`.
* `F5_SECRET` - The path in vault to the F5 secret.  The secret in vault must contain the keys `password`, `url` and `username` for the F5 interface.
* `CLUSTER` - The cluster name
* `DEFAULT_DOMAIN` - The default sub-domain for applications e.g., `myapp.domain.io`
* `F5_INSIDE_MONITOR` and `F5_MONITOR` - The `/PARTITION/monitor` to use for monitoring node ports for the INSIDE (private) apps, and public apps
* `F5_INSIDE_PARTITION` and `F5_PARTITION` - The inside and outside partition name (e.g., `DEV_INSIDE`, `DEV_OUTSIDE`)
* `F5_INSIDE_VIRTUAL` and `F5_VIRTUAL` - The inside and outside virtual IP name (without a partition). E.g. (`inside-apps`, `outside-apps`)
* `INSIDE_DOMAIN` - The default sub-domain for application inside e.g., `myapp-internal.domain.io`
* `KUBERNETES_API_SERVER` - The host name for the api server of kubernetes
* `KUBERNETES_CLIENT_TYPE` - The client type of kubernetes, must be set to token or cert.
* `KUBERNETES_TOKEN_SECRET` - The path in vault to the secret to access kubernetes. The vault secret must have a field called `token` if the client type is `token`. It must have fields `admin-crt`, `admin-key` and `ca-crt` if the client type is set to `cert`.
* `NAMESPACE_BLACKLIST` - A comma delimited list of namespaces to instruct the service watcher not to automatically create irules for.
* `PROFILE` - If set to `true` the environment variable `STACKIMPACT` must also be set to the API key for stack impact for profiling.
* `REGIONAPI_LOCATION` - The http path to the region api to use, this overrides the location set in `ALAMOAPI_SECRET`
* `UNIPOOL`: The name of the unipool to use. This pool must contain all of the kubernetes worker nodes. The port is automatically reset in the iRule.