## Running these charts

### Create regcreds.yaml and auth0secrets.yaml files

These values should not be checked in to github. Look at the `regcreds-EXAMPLE.yaml` and `auth0secrets-EXAMPLE.yaml` files.

The auth0 value can be found in the Applications section of the Auth0 dashboard.

### Test values

Like this: 
```
helm install -f login-gravhl/auth0secrets.yaml -f login-gravhl/regcreds.yaml  --dry-run login-gravhl login-gravhl
```