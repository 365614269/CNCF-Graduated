#pragma once

#include <atomic>

#include "envoy/common/time.h"

#include "source/common/buffer/buffer_impl.h"
#include "source/common/event/event_impl_base.h"
#include "source/common/event/file_event_impl.h"
#include "source/common/network/utility.h"

#include "base_listener_impl.h"

namespace Envoy {
namespace Network {

/**
 * Implementation of Network::Listener for UDP.
 */
class UdpListenerImpl : public BaseListenerImpl,
                        public virtual UdpListener,
                        public UdpPacketProcessor,
                        protected Logger::Loggable<Logger::Id::udp> {
public:
  UdpListenerImpl(Event::Dispatcher& dispatcher, SocketSharedPtr socket, UdpListenerCallbacks& cb,
                  TimeSource& time_source, const envoy::config::core::v3::UdpSocketConfig& config);
  ~UdpListenerImpl() override;
  uint32_t packetsDropped() { return packets_dropped_; }
  bool paused() const { return parent_drained_callback_registrar_ != absl::nullopt; }
  void unpause();

  // Network::Listener
  void disable() override;
  void enable() override;
  void setRejectFraction(UnitFloat) override {}
  void configureLoadShedPoints(Server::LoadShedPointProvider&) override {}
  bool shouldBypassOverloadManager() const override { return false; }

  // Network::UdpListener
  Event::Dispatcher& dispatcher() override;
  const Address::InstanceConstSharedPtr& localAddress() const override;
  Api::IoCallUint64Result send(const UdpSendData& data) override;
  Api::IoCallUint64Result flush() override;
  void activateRead() override;

  // Network::UdpPacketProcessor
  void processPacket(Address::InstanceConstSharedPtr local_address,
                     Address::InstanceConstSharedPtr peer_address, Buffer::InstancePtr buffer,
                     MonotonicTime receive_time, uint8_t tos,
                     Buffer::OwnedImpl saved_cmsg) override;
  uint64_t maxDatagramSize() const override { return config_.max_rx_datagram_size_; }
  void onDatagramsDropped(uint32_t dropped) override { cb_.onDatagramsDropped(dropped); }
  size_t numPacketsExpectedPerEventLoop() const override {
    return cb_.numPacketsExpectedPerEventLoop();
  }
  const IoHandle::UdpSaveCmsgConfig& saveCmsgConfig() const override {
    return cb_.udpSaveCmsgConfig();
  }

protected:
  void handleWriteCallback();
  void handleReadCallback();

  UdpListenerCallbacks& cb_;
  uint32_t packets_dropped_{0};

private:
  void onSocketEvent(short flags);
  void disableEvent();

  TimeSource& time_source_;
  const ResolvedUdpSocketConfig config_;
  OptRef<ParentDrainedCallbackRegistrar> parent_drained_callback_registrar_;
  // Taking a weak_ptr to this lets us detect if the listener has been destroyed.
  std::shared_ptr<bool> destruction_checker_ = std::make_shared<bool>(true);
  uint32_t events_when_unpaused_ = Event::FileReadyType::Read | Event::FileReadyType::Write;
};

class UdpListenerWorkerRouterImpl : public UdpListenerWorkerRouter {
public:
  UdpListenerWorkerRouterImpl(uint32_t concurrency);

  // UdpListenerWorkerRouter
  void registerWorkerForListener(UdpListenerCallbacks& listener) override;
  void unregisterWorkerForListener(UdpListenerCallbacks& listener) override;
  void deliver(uint32_t dest_worker_index, UdpRecvData&& data) override;

private:
  absl::Mutex mutex_;
  std::vector<UdpListenerCallbacks*> workers_ ABSL_GUARDED_BY(mutex_);
};

} // namespace Network
} // namespace Envoy
