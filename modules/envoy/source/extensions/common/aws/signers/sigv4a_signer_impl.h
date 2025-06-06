#pragma once

#include <memory>

#include "source/extensions/common/aws/credentials_provider.h"
#include "source/extensions/common/aws/signer_base_impl.h"
#include "source/extensions/common/aws/signers/sigv4a_key_derivation.h"

namespace Envoy {
namespace Extensions {
namespace Common {
namespace Aws {

/**
 * Implementation of the Signature V4A signing process.
 * See https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html
 *
 * Query parameter support is implemented as per:
 * https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
 */

class SigV4ASignerImpl : public SignerBaseImpl {

  // Allow friend access for signer corpus testing
  friend class SigV4ASignerImplFriend;

public:
  SigV4ASignerImpl(
      absl::string_view service_name, absl::string_view region,
      const CredentialsProviderChainSharedPtr& credentials_provider,
      Server::Configuration::CommonFactoryContext& context,
      const AwsSigningHeaderExclusionVector& matcher_config, const bool query_string = false,
      const uint16_t expiration_time = SignatureQueryParameterValues::DefaultExpiration,
      std::unique_ptr<SigV4AKeyDerivationBase> key_derivation_ptr =
          std::make_unique<SigV4AKeyDerivation>())
      : SignerBaseImpl(service_name, region, credentials_provider, context, matcher_config,
                       query_string, expiration_time),
        key_derivation_ptr_(std::move(key_derivation_ptr)) {}

private:
  void addRegionHeader(Http::RequestHeaderMap& headers,
                       const absl::string_view override_region) const override;

  void addRegionQueryParam(Envoy::Http::Utility::QueryParamsMulti& query_params,
                           const absl::string_view override_region) const override;

  std::string createCredentialScope(const absl::string_view short_date,
                                    const absl::string_view override_region) const override;

  std::string createStringToSign(const absl::string_view canonical_request,
                                 const absl::string_view long_date,
                                 const absl::string_view credential_scope) const override;

  std::string
  createSignature(const absl::string_view access_key_id, const absl::string_view secret_access_key,
                  ABSL_ATTRIBUTE_UNUSED const absl::string_view short_date,
                  const absl::string_view string_to_sign,
                  ABSL_ATTRIBUTE_UNUSED const absl::string_view override_region) const override;

  std::string createAuthorizationHeader(const absl::string_view access_key_id,
                                        const absl::string_view credential_scope,
                                        const std::map<std::string, std::string>& canonical_headers,
                                        const absl::string_view signature) const override;

  absl::string_view getAlgorithmString() const override;
  std::unique_ptr<SigV4AKeyDerivationBase> key_derivation_ptr_;
};

} // namespace Aws
} // namespace Common
} // namespace Extensions
} // namespace Envoy
