#include "source/common/upstream/od_cds_api_impl.h"

#include "source/common/common/assert.h"
#include "source/common/grpc/common.h"

#include "absl/strings/str_join.h"

namespace Envoy {
namespace Upstream {

absl::StatusOr<OdCdsApiSharedPtr>
OdCdsApiImpl::create(const envoy::config::core::v3::ConfigSource& odcds_config,
                     OptRef<xds::core::v3::ResourceLocator> odcds_resources_locator,
                     Config::XdsManager& xds_manager, ClusterManager& cm,
                     MissingClusterNotifier& notifier, Stats::Scope& scope,
                     ProtobufMessage::ValidationVisitor& validation_visitor) {
  absl::Status creation_status = absl::OkStatus();
  auto ret =
      OdCdsApiSharedPtr(new OdCdsApiImpl(odcds_config, odcds_resources_locator, xds_manager, cm,
                                         notifier, scope, validation_visitor, creation_status));
  RETURN_IF_NOT_OK(creation_status);
  return ret;
}

OdCdsApiImpl::OdCdsApiImpl(const envoy::config::core::v3::ConfigSource& odcds_config,
                           OptRef<xds::core::v3::ResourceLocator> odcds_resources_locator,
                           Config::XdsManager& xds_manager, ClusterManager& cm,
                           MissingClusterNotifier& notifier, Stats::Scope& scope,
                           ProtobufMessage::ValidationVisitor& validation_visitor,
                           absl::Status& creation_status)
    : Envoy::Config::SubscriptionBase<envoy::config::cluster::v3::Cluster>(validation_visitor,
                                                                           "name"),
      helper_(cm, xds_manager, "odcds"), notifier_(notifier),
      scope_(scope.createScope("cluster_manager.odcds.")) {
  // TODO(krnowak): Move the subscription setup to CdsApiHelper. Maybe make CdsApiHelper a base
  // class for CDS and ODCDS.
  const auto resource_name = getResourceName();
  absl::StatusOr<Config::SubscriptionPtr> subscription_or_error;
  if (!odcds_resources_locator.has_value()) {
    subscription_or_error = cm.subscriptionFactory().subscriptionFromConfigSource(
        odcds_config, Grpc::Common::typeUrl(resource_name), *scope_, *this, resource_decoder_, {});
  } else {
    subscription_or_error = cm.subscriptionFactory().collectionSubscriptionFromUrl(
        *odcds_resources_locator, odcds_config, resource_name, *scope_, *this, resource_decoder_);
  }
  SET_AND_RETURN_IF_NOT_OK(subscription_or_error.status(), creation_status);
  subscription_ = std::move(*subscription_or_error);
}

absl::Status OdCdsApiImpl::onConfigUpdate(const std::vector<Config::DecodedResourceRef>& resources,
                                          const std::string& version_info) {
  UNREFERENCED_PARAMETER(resources);
  UNREFERENCED_PARAMETER(version_info);
  // On-demand cluster updates are only supported for delta, not sotw.
  PANIC("not supported");
}

absl::Status
OdCdsApiImpl::onConfigUpdate(const std::vector<Config::DecodedResourceRef>& added_resources,
                             const Protobuf::RepeatedPtrField<std::string>& removed_resources,
                             const std::string& system_version_info) {
  auto [_, exception_msgs] =
      helper_.onConfigUpdate(added_resources, removed_resources, system_version_info);
  sendAwaiting();
  status_ = StartStatus::InitialFetchDone;
  // According to the XDS specification, the server can send a reply with names in the
  // removed_resources field for requested resources that do not exist. That way we can notify the
  // interested parties about the missing resource immediately without waiting for some timeout to
  // be triggered.
  for (const auto& resource_name : removed_resources) {
    ENVOY_LOG(debug, "odcds: notifying about potential missing cluster {}", resource_name);
    notifier_.notifyMissingCluster(resource_name);
  }
  if (!exception_msgs.empty()) {
    return absl::InvalidArgumentError(
        fmt::format("Error adding/updating cluster(s) {}", absl::StrJoin(exception_msgs, ", ")));
  }
  return absl::OkStatus();
}

void OdCdsApiImpl::onConfigUpdateFailed(Envoy::Config::ConfigUpdateFailureReason reason,
                                        const EnvoyException*) {
  ASSERT(Envoy::Config::ConfigUpdateFailureReason::ConnectionFailure != reason);
  sendAwaiting();
  status_ = StartStatus::InitialFetchDone;
}

void OdCdsApiImpl::sendAwaiting() {
  if (awaiting_names_.empty()) {
    return;
  }
  // The awaiting names are sent only once. After the state transition from Starting to
  // InitialFetchDone (which happens on the first received response), the awaiting names list is not
  // used any more.
  ENVOY_LOG(debug, "odcds: sending request for awaiting cluster names {}",
            fmt::join(awaiting_names_, ", "));
  subscription_->requestOnDemandUpdate(awaiting_names_);
  awaiting_names_.clear();
}

void OdCdsApiImpl::updateOnDemand(std::string cluster_name) {
  switch (status_) {
  case StartStatus::NotStarted:
    ENVOY_LOG(trace, "odcds: starting a subscription with cluster name {}", cluster_name);
    status_ = StartStatus::Started;
    subscription_->start({std::move(cluster_name)});
    return;

  case StartStatus::Started:
    ENVOY_LOG(trace, "odcds: putting cluster name {} on awaiting list", cluster_name);
    awaiting_names_.insert(std::move(cluster_name));
    return;

  case StartStatus::InitialFetchDone:
    ENVOY_LOG(trace, "odcds: requesting for cluster name {}", cluster_name);
    subscription_->requestOnDemandUpdate({std::move(cluster_name)});
    return;
  }
  PANIC("corrupt enum");
}

} // namespace Upstream
} // namespace Envoy
