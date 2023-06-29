package client

import (
	"github.com/auth0/go-auth0/management"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/auth0/terraform-provider-auth0/internal/value"
)

func expandClient(d *schema.ResourceData) *management.Client {
	config := d.GetRawConfig()

	client := &management.Client{
		Name:                           value.String(config.GetAttr("name")),
		Description:                    value.String(config.GetAttr("description")),
		AppType:                        value.String(config.GetAttr("app_type")),
		LogoURI:                        value.String(config.GetAttr("logo_uri")),
		IsFirstParty:                   value.Bool(config.GetAttr("is_first_party")),
		IsTokenEndpointIPHeaderTrusted: value.Bool(config.GetAttr("is_token_endpoint_ip_header_trusted")),
		OIDCConformant:                 value.Bool(config.GetAttr("oidc_conformant")),
		ClientAliases:                  value.Strings(config.GetAttr("client_aliases")),
		Callbacks:                      value.Strings(config.GetAttr("callbacks")),
		AllowedLogoutURLs:              value.Strings(config.GetAttr("allowed_logout_urls")),
		AllowedOrigins:                 value.Strings(config.GetAttr("allowed_origins")),
		AllowedClients:                 value.Strings(config.GetAttr("allowed_clients")),
		GrantTypes:                     value.Strings(config.GetAttr("grant_types")),
		OrganizationUsage:              value.String(config.GetAttr("organization_usage")),
		OrganizationRequireBehavior:    value.String(config.GetAttr("organization_require_behavior")),
		WebOrigins:                     value.Strings(config.GetAttr("web_origins")),
		SSO:                            value.Bool(config.GetAttr("sso")),
		SSODisabled:                    value.Bool(config.GetAttr("sso_disabled")),
		CrossOriginAuth:                value.Bool(config.GetAttr("cross_origin_auth")),
		CrossOriginLocation:            value.String(config.GetAttr("cross_origin_loc")),
		CustomLoginPageOn:              value.Bool(config.GetAttr("custom_login_page_on")),
		CustomLoginPage:                value.String(config.GetAttr("custom_login_page")),
		FormTemplate:                   value.String(config.GetAttr("form_template")),
		TokenEndpointAuthMethod:        value.String(config.GetAttr("token_endpoint_auth_method")),
		InitiateLoginURI:               value.String(config.GetAttr("initiate_login_uri")),
		EncryptionKey:                  value.MapOfStrings(config.GetAttr("encryption_key")),
		OIDCBackchannelLogout:          expandOIDCBackchannelLogout(d),
		ClientMetadata:                 expandClientMetadata(d),
		RefreshToken:                   expandClientRefreshToken(d),
		JWTConfiguration:               expandClientJWTConfiguration(d),
		Addons:                         expandClientAddons(d),
		NativeSocialLogin:              expandClientNativeSocialLogin(d),
		Mobile:                         expandClientMobile(d),
	}

	return client
}

func expandOIDCBackchannelLogout(d *schema.ResourceData) *management.OIDCBackchannelLogout {
	raw := d.GetRawConfig().GetAttr("oidc_backchannel_logout_urls")

	logoutUrls := value.Strings(raw)

	if logoutUrls == nil {
		return nil
	}

	return &management.OIDCBackchannelLogout{
		BackChannelLogoutURLs: logoutUrls,
	}
}

func expandClientRefreshToken(d *schema.ResourceData) *management.ClientRefreshToken {
	refreshTokenConfig := d.GetRawConfig().GetAttr("refresh_token")
	if refreshTokenConfig.IsNull() {
		return nil
	}

	var refreshToken management.ClientRefreshToken

	refreshTokenConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		refreshToken.RotationType = value.String(config.GetAttr("rotation_type"))
		refreshToken.ExpirationType = value.String(config.GetAttr("expiration_type"))
		refreshToken.Leeway = value.Int(config.GetAttr("leeway"))
		refreshToken.TokenLifetime = value.Int(config.GetAttr("token_lifetime"))
		refreshToken.InfiniteTokenLifetime = value.Bool(config.GetAttr("infinite_token_lifetime"))
		refreshToken.InfiniteIdleTokenLifetime = value.Bool(config.GetAttr("infinite_idle_token_lifetime"))
		refreshToken.IdleTokenLifetime = value.Int(config.GetAttr("idle_token_lifetime"))
		return stop
	})

	if refreshToken == (management.ClientRefreshToken{}) {
		return nil
	}

	return &refreshToken
}

func expandClientJWTConfiguration(d *schema.ResourceData) *management.ClientJWTConfiguration {
	jwtConfig := d.GetRawConfig().GetAttr("jwt_configuration")
	if jwtConfig.IsNull() {
		return nil
	}

	var jwt management.ClientJWTConfiguration

	jwtConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		jwt.LifetimeInSeconds = value.Int(config.GetAttr("lifetime_in_seconds"))
		jwt.Algorithm = value.String(config.GetAttr("alg"))
		jwt.Scopes = value.MapOfStrings(config.GetAttr("scopes"))

		if d.IsNewResource() {
			jwt.SecretEncoded = value.Bool(config.GetAttr("secret_encoded"))
		}

		return stop
	})

	if jwt == (management.ClientJWTConfiguration{}) {
		return nil
	}

	return &jwt
}

func expandClientNativeSocialLogin(d *schema.ResourceData) *management.ClientNativeSocialLogin {
	nativeSocialLoginConfig := d.GetRawConfig().GetAttr("native_social_login")
	if nativeSocialLoginConfig.IsNull() {
		return nil
	}

	var nativeSocialLogin management.ClientNativeSocialLogin

	nativeSocialLoginConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		nativeSocialLogin.Apple = expandClientNativeSocialLoginSupportEnabled(config.GetAttr("apple"))
		nativeSocialLogin.Facebook = expandClientNativeSocialLoginSupportEnabled(config.GetAttr("facebook"))
		return stop
	})

	if nativeSocialLogin == (management.ClientNativeSocialLogin{}) {
		return nil
	}

	return &nativeSocialLogin
}

func expandClientNativeSocialLoginSupportEnabled(config cty.Value) *management.ClientNativeSocialLoginSupportEnabled {
	if config.IsNull() {
		return nil
	}

	var support management.ClientNativeSocialLoginSupportEnabled

	config.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		support.Enabled = value.Bool(config.GetAttr("enabled"))
		return stop
	})

	if support == (management.ClientNativeSocialLoginSupportEnabled{}) {
		return nil
	}

	return &support
}

func expandClientMobile(d *schema.ResourceData) *management.ClientMobile {
	mobileConfig := d.GetRawConfig().GetAttr("mobile")
	if mobileConfig.IsNull() {
		return nil
	}

	var mobile management.ClientMobile

	mobileConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		mobile.Android = expandClientMobileAndroid(config.GetAttr("android"))
		mobile.IOS = expandClientMobileIOS(config.GetAttr("ios"))
		return stop
	})

	if mobile == (management.ClientMobile{}) {
		return nil
	}

	return &mobile
}

func expandClientMobileAndroid(androidConfig cty.Value) *management.ClientMobileAndroid {
	if androidConfig.IsNull() {
		return nil
	}

	var android management.ClientMobileAndroid

	androidConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		android.AppPackageName = value.String(config.GetAttr("app_package_name"))
		android.KeyHashes = value.Strings(config.GetAttr("sha256_cert_fingerprints"))
		return stop
	})

	if android == (management.ClientMobileAndroid{}) {
		return nil
	}

	return &android
}

func expandClientMobileIOS(iosConfig cty.Value) *management.ClientMobileIOS {
	if iosConfig.IsNull() {
		return nil
	}

	var ios management.ClientMobileIOS

	iosConfig.ForEachElement(func(_ cty.Value, config cty.Value) (stop bool) {
		ios.TeamID = value.String(config.GetAttr("team_id"))
		ios.AppID = value.String(config.GetAttr("app_bundle_identifier"))
		return stop
	})

	if ios == (management.ClientMobileIOS{}) {
		return nil
	}

	return &ios
}

func expandClientMetadata(d *schema.ResourceData) *map[string]interface{} {
	if !d.HasChange("client_metadata") {
		return nil
	}

	oldMetadata, newMetadata := d.GetChange("client_metadata")
	oldMetadataMap := oldMetadata.(map[string]interface{})
	newMetadataMap := newMetadata.(map[string]interface{})

	for key := range oldMetadataMap {
		if _, ok := newMetadataMap[key]; !ok {
			newMetadataMap[key] = nil
		}
	}

	return &newMetadataMap
}

func expandClientAddons(d *schema.ResourceData) *management.ClientAddons {
	if !d.HasChange("addons") {
		return nil
	}

	var addons management.ClientAddons

	d.GetRawConfig().GetAttr("addons").ForEachElement(func(_ cty.Value, addonsCfg cty.Value) (stop bool) {
		addons.AWS = expandClientAddonAWS(addonsCfg.GetAttr("aws"))
		addons.AzureBlob = expandClientAddonAzureBlob(addonsCfg.GetAttr("azure_blob"))
		addons.AzureSB = expandClientAddonAzureSB(addonsCfg.GetAttr("azure_sb"))
		addons.RMS = expandClientAddonRMS(addonsCfg.GetAttr("rms"))
		addons.MSCRM = expandClientAddonMSCRM(addonsCfg.GetAttr("mscrm"))
		addons.Slack = expandClientAddonSlack(addonsCfg.GetAttr("slack"))
		addons.Sentry = expandClientAddonSentry(addonsCfg.GetAttr("sentry"))
		addons.EchoSign = expandClientAddonEchoSign(addonsCfg.GetAttr("echosign"))
		addons.Egnyte = expandClientAddonEgnyte(addonsCfg.GetAttr("egnyte"))
		addons.Firebase = expandClientAddonFirebase(addonsCfg.GetAttr("firebase"))
		addons.NewRelic = expandClientAddonNewRelic(addonsCfg.GetAttr("newrelic"))
		addons.Office365 = expandClientAddonOffice365(addonsCfg.GetAttr("office365"))
		return stop
	})

	if addons == (management.ClientAddons{}) {
		return nil
	}

	return &addons
}

func expandClientAddonAWS(awsCfg cty.Value) *management.AWSClientAddon {
	var awsAddon management.AWSClientAddon

	awsCfg.ForEachElement(func(_ cty.Value, awsCfg cty.Value) (stop bool) {
		awsAddon = management.AWSClientAddon{
			Principal:         value.String(awsCfg.GetAttr("principal")),
			Role:              value.String(awsCfg.GetAttr("role")),
			LifetimeInSeconds: value.Int(awsCfg.GetAttr("lifetime_in_seconds")),
		}

		return stop
	})

	return &awsAddon
}

func expandClientAddonAzureBlob(azureCfg cty.Value) *management.AzureBlobClientAddon {
	var azureAddon management.AzureBlobClientAddon

	azureCfg.ForEachElement(func(_ cty.Value, azureCfg cty.Value) (stop bool) {
		azureAddon = management.AzureBlobClientAddon{
			AccountName:      value.String(azureCfg.GetAttr("account_name")),
			StorageAccessKey: value.String(azureCfg.GetAttr("storage_access_key")),
			ContainerName:    value.String(azureCfg.GetAttr("container_name")),
			BlobName:         value.String(azureCfg.GetAttr("blob_name")),
			Expiration:       value.Int(azureCfg.GetAttr("expiration")),
			SignedIdentifier: value.String(azureCfg.GetAttr("signed_identifier")),
			BlobRead:         value.Bool(azureCfg.GetAttr("blob_read")),
			BlobWrite:        value.Bool(azureCfg.GetAttr("blob_write")),
			BlobDelete:       value.Bool(azureCfg.GetAttr("blob_delete")),
			ContainerRead:    value.Bool(azureCfg.GetAttr("container_read")),
			ContainerWrite:   value.Bool(azureCfg.GetAttr("container_write")),
			ContainerDelete:  value.Bool(azureCfg.GetAttr("container_delete")),
			ContainerList:    value.Bool(azureCfg.GetAttr("container_list")),
		}

		return stop
	})

	return &azureAddon
}

func expandClientAddonAzureSB(azureCfg cty.Value) *management.AzureSBClientAddon {
	var azureAddon management.AzureSBClientAddon

	azureCfg.ForEachElement(func(_ cty.Value, azureCfg cty.Value) (stop bool) {
		azureAddon = management.AzureSBClientAddon{
			Namespace:  value.String(azureCfg.GetAttr("namespace")),
			SASKeyName: value.String(azureCfg.GetAttr("sas_key_name")),
			SASKey:     value.String(azureCfg.GetAttr("sas_key")),
			EntityPath: value.String(azureCfg.GetAttr("entity_path")),
			Expiration: value.Int(azureCfg.GetAttr("expiration")),
		}

		return stop
	})

	return &azureAddon
}

func expandClientAddonRMS(rmsCfg cty.Value) *management.RMSClientAddon {
	var rmsAddon management.RMSClientAddon

	rmsCfg.ForEachElement(func(_ cty.Value, rmsCfg cty.Value) (stop bool) {
		rmsAddon = management.RMSClientAddon{
			URL: value.String(rmsCfg.GetAttr("url")),
		}

		return stop
	})

	if rmsAddon == (management.RMSClientAddon{}) {
		return nil
	}

	return &rmsAddon
}

func expandClientAddonMSCRM(mscrmCfg cty.Value) *management.MSCRMClientAddon {
	var mscrmAddon management.MSCRMClientAddon

	mscrmCfg.ForEachElement(func(_ cty.Value, mscrmCfg cty.Value) (stop bool) {
		mscrmAddon = management.MSCRMClientAddon{
			URL: value.String(mscrmCfg.GetAttr("url")),
		}

		return stop
	})

	if mscrmAddon == (management.MSCRMClientAddon{}) {
		return nil
	}

	return &mscrmAddon
}

func expandClientAddonSlack(slackCfg cty.Value) *management.SlackClientAddon {
	var slackAddon management.SlackClientAddon

	slackCfg.ForEachElement(func(_ cty.Value, slackCfg cty.Value) (stop bool) {
		slackAddon = management.SlackClientAddon{
			Team: value.String(slackCfg.GetAttr("team")),
		}

		return stop
	})

	if slackAddon == (management.SlackClientAddon{}) {
		return nil
	}

	return &slackAddon
}

func expandClientAddonSentry(sentryCfg cty.Value) *management.SentryClientAddon {
	var sentryAddon management.SentryClientAddon

	sentryCfg.ForEachElement(func(_ cty.Value, sentryCfg cty.Value) (stop bool) {
		sentryAddon = management.SentryClientAddon{
			OrgSlug: value.String(sentryCfg.GetAttr("org_slug")),
			BaseURL: value.String(sentryCfg.GetAttr("base_url")),
		}

		return stop
	})

	return &sentryAddon
}

func expandClientAddonEchoSign(echoSignCfg cty.Value) *management.EchoSignClientAddon {
	var echoSignAddon management.EchoSignClientAddon

	echoSignCfg.ForEachElement(func(_ cty.Value, echoSignCfg cty.Value) (stop bool) {
		echoSignAddon = management.EchoSignClientAddon{
			Domain: value.String(echoSignCfg.GetAttr("domain")),
		}

		return stop
	})

	return &echoSignAddon
}

func expandClientAddonEgnyte(egnyteCfg cty.Value) *management.EgnyteClientAddon {
	var egnyteAddon management.EgnyteClientAddon

	egnyteCfg.ForEachElement(func(_ cty.Value, egnyteCfg cty.Value) (stop bool) {
		egnyteAddon = management.EgnyteClientAddon{
			Domain: value.String(egnyteCfg.GetAttr("domain")),
		}

		return stop
	})

	return &egnyteAddon
}

func expandClientAddonFirebase(firebaseCfg cty.Value) *management.FirebaseClientAddon {
	var firebaseAddon management.FirebaseClientAddon

	firebaseCfg.ForEachElement(func(_ cty.Value, firebaseCfg cty.Value) (stop bool) {
		firebaseAddon = management.FirebaseClientAddon{
			Secret:            value.String(firebaseCfg.GetAttr("secret")),
			PrivateKeyID:      value.String(firebaseCfg.GetAttr("private_key_id")),
			PrivateKey:        value.String(firebaseCfg.GetAttr("private_key")),
			ClientEmail:       value.String(firebaseCfg.GetAttr("client_email")),
			LifetimeInSeconds: value.Int(firebaseCfg.GetAttr("lifetime_in_seconds")),
		}

		return stop
	})

	return &firebaseAddon
}

func expandClientAddonNewRelic(newRelicCfg cty.Value) *management.NewRelicClientAddon {
	var newRelicAddon management.NewRelicClientAddon

	newRelicCfg.ForEachElement(func(_ cty.Value, newRelicCfg cty.Value) (stop bool) {
		newRelicAddon = management.NewRelicClientAddon{
			Account: value.String(newRelicCfg.GetAttr("account")),
		}

		return stop
	})

	return &newRelicAddon
}

func expandClientAddonOffice365(office365Cfg cty.Value) *management.Office365ClientAddon {
	var office365Addon management.Office365ClientAddon

	office365Cfg.ForEachElement(func(_ cty.Value, office365Cfg cty.Value) (stop bool) {
		office365Addon = management.Office365ClientAddon{
			Domain:     value.String(office365Cfg.GetAttr("domain")),
			Connection: value.String(office365Cfg.GetAttr("connection")),
		}

		return stop
	})

	return &office365Addon
}

func clientHasChange(c *management.Client) bool {
	return c.String() != "{}"
}
