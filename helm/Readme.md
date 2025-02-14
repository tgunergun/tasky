
## Pushing Helm Chart to Azure Container Registry (ACR)

1. **Login to Azure Container Registry:**
   ```sh
   az acr login --name wizdemoacr000897
   ```

2. **Package the Helm Chart:**
   ```sh
   helm package /home/dixon/wiz/tasky/helm
   ```

3. **Push the Helm Chart to ACR:**
   ```sh
   helm push tasky-0.2.0.tgz oci://wizdemoacr000897.azurecr.io/helm
   ```
