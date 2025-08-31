# Publishing the Firefly Terraform Provider

This guide walks through publishing the provider to the Terraform Registry.

## Prerequisites Completed ✅

- [x] Repository name: `terraform-provider-firefly` (correct format)
- [x] Repository is public on GitHub 
- [x] Provider code is complete and tested
- [x] Documentation structure created in `docs/`
- [x] Terraform registry manifest file created
- [x] GitHub Actions workflow configured

## Next Steps

### 1. Generate GPG Signing Key

```bash
# Install GPG if not available
brew install gnupg  # macOS
# or: apt-get install gnupg  # Ubuntu/Debian

# Generate key (use your real name and email)
gpg --full-generate-key

# Export public key (replace with your key ID)
gpg --armor --export YOUR_KEY_ID > public_key.asc

# Export private key for GitHub secrets
gpg --armor --export-secret-keys YOUR_KEY_ID > private_key.asc
```

### 2. Configure GitHub Secrets

Add these secrets to your GitHub repository (Settings > Secrets and variables > Actions):

- `GPG_PRIVATE_KEY`: Content of `private_key.asc`
- `PASSPHRASE`: Your GPG key passphrase (if any)

### 3. Create a Release

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow that:
- Builds binaries for all platforms
- Signs the checksums
- Creates a GitHub release with all assets

### 4. Add GPG Key to Terraform Registry

1. Go to [Terraform Registry](https://registry.terraform.io/)
2. Sign in with your GitHub account
3. Navigate to User Settings > Signing Keys
4. Upload the content of `public_key.asc`

### 5. Publish Provider

1. In Terraform Registry, go to "Publish > Provider"
2. Select the `gofireflyio/terraform-provider-firefly` repository
3. Follow the prompts to publish

### 6. Update Documentation

After publishing, update `docs/index.md` to use the official source:

```terraform
terraform {
  required_providers {
    firefly = {
      source = "gofireflyio/firefly"
      version = "~> 1.0"
    }
  }
}
```

## Verification

After publishing, test the provider:

```bash
# Remove local development setup
rm -rf .terraform*

# Use published provider
terraform init
terraform plan
```

## Files Created/Modified

- ✅ `terraform-registry-manifest.json` - Registry protocol specification
- ✅ `.github/workflows/release.yml` - Automated release workflow
- ✅ `PUBLISHING.md` - This publishing guide

## Ready for Publication

The provider is now ready for publication. Follow the steps above to complete the publishing process.