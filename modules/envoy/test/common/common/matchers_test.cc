#include "envoy/common/exception.h"
#include "envoy/config/core/v3/base.pb.h"
#include "envoy/type/matcher/v3/metadata.pb.h"
#include "envoy/type/matcher/v3/string.pb.h"
#include "envoy/type/matcher/v3/value.pb.h"

#include "source/common/common/matchers.h"
#include "source/common/config/metadata.h"
#include "source/common/protobuf/protobuf.h"
#include "source/common/stream_info/filter_state_impl.h"

#include "test/mocks/server/server_factory_context.h"
#include "test/test_common/utility.h"

#include "gtest/gtest.h"

namespace Envoy {
namespace Matcher {
namespace {

class BaseTest : public ::testing::Test {
public:
  NiceMock<Server::Configuration::MockServerFactoryContext> context_;
};

class MetadataTest : public BaseTest {};

TEST_F(MetadataTest, MatchNullValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_null_value(ProtobufWkt::NullValue::NULL_VALUE);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_null_match();
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchDoubleValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_number_value(9);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_double_match()->set_exact(1);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_double_match()->set_exact(9);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  auto r = matcher.mutable_value()->mutable_double_match()->mutable_range();
  r->set_start(9.1);
  r->set_end(10);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  r = matcher.mutable_value()->mutable_double_match()->mutable_range();
  r->set_start(8.9);
  r->set_end(9);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  r = matcher.mutable_value()->mutable_double_match()->mutable_range();
  r->set_start(9);
  r->set_end(9.1);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchStringExactValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_string_value("prod");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_exact("prod");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchStringPrefixValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_string_value("prodabc");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_prefix("prodx");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_prefix("prod");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchStringSuffixValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_string_value("abcprod");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_suffix("prodx");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_suffix("prod");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchStringContainsValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_string_value("abcprodef");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_contains("pride");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_contains("prod");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchBoolValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_bool_value(true);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->set_bool_match(false);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->set_bool_match(true);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchPresentValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_number_value(1);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.b");
  matcher.add_path()->set_key("label");

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->set_present_match(false);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->set_present_match(true);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  matcher.clear_path();
  matcher.add_path()->set_key("unknown");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

TEST_F(MetadataTest, MatchStringOrBoolValue) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "label")
      .set_string_value("test");
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.b", "label")
      .set_bool_value(true);
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.c", "label")
      .set_bool_value(false);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.add_path()->set_key("label");

  auto* or_match = matcher.mutable_value()->mutable_or_match();
  or_match->add_value_matchers()->mutable_string_match()->set_exact("test");
  or_match->add_value_matchers()->set_bool_match(true);
  matcher.set_filter("envoy.filter.a");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.set_filter("envoy.filter.b");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.set_filter("envoy.filter.c");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

// Helper function to retrieve the reference of an entry in a ListMatcher from a MetadataMatcher.
envoy::type::matcher::v3::ValueMatcher*
listMatchEntry(envoy::type::matcher::v3::MetadataMatcher* matcher) {
  return matcher->mutable_value()->mutable_list_match()->mutable_one_of();
}

TEST_F(MetadataTest, MatchStringListValue) {
  envoy::config::core::v3::Metadata metadata;
  ProtobufWkt::Value& metadataValue =
      Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "groups");
  ProtobufWkt::ListValue* values = metadataValue.mutable_list_value();
  values->add_values()->set_string_value("first");
  values->add_values()->set_string_value("second");
  values->add_values()->set_string_value("third");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.a");
  matcher.add_path()->set_key("groups");

  listMatchEntry(&matcher)->mutable_string_match()->set_exact("second");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_string_match()->set_prefix("fi");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_string_match()->set_suffix("rd");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_string_match()->set_exact("fourth");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_string_match()->set_prefix("none");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  values->clear_values();
  metadataValue.Clear();
}

TEST_F(MetadataTest, MatchBoolListValue) {
  envoy::config::core::v3::Metadata metadata;
  ProtobufWkt::Value& metadataValue =
      Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "groups");
  ProtobufWkt::ListValue* values = metadataValue.mutable_list_value();
  values->add_values()->set_bool_value(false);
  values->add_values()->set_bool_value(false);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.a");
  matcher.add_path()->set_key("groups");

  listMatchEntry(&matcher)->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->set_bool_match(true);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->set_bool_match(false);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  values->clear_values();
  metadataValue.Clear();
}

TEST_F(MetadataTest, MatchDoubleListValue) {
  envoy::config::core::v3::Metadata metadata;
  ProtobufWkt::Value& metadataValue =
      Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.a", "groups");
  ProtobufWkt::ListValue* values = metadataValue.mutable_list_value();
  values->add_values()->set_number_value(10);
  values->add_values()->set_number_value(23);

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.a");
  matcher.add_path()->set_key("groups");

  listMatchEntry(&matcher)->mutable_string_match()->set_exact("test");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->set_bool_match(true);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_double_match()->set_exact(9);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  listMatchEntry(&matcher)->mutable_double_match()->set_exact(10);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  auto r = listMatchEntry(&matcher)->mutable_double_match()->mutable_range();
  r->set_start(10);
  r->set_end(15);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  r = listMatchEntry(&matcher)->mutable_double_match()->mutable_range();
  r->set_start(20);
  r->set_end(24);
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  r = listMatchEntry(&matcher)->mutable_double_match()->mutable_range();
  r->set_start(24);
  r->set_end(26);
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));

  values->clear_values();
  metadataValue.Clear();
}

TEST_F(MetadataTest, InvertMatch) {
  envoy::config::core::v3::Metadata metadata;
  Envoy::Config::Metadata::mutableMetadataValue(metadata, "envoy.filter.x", "label")
      .set_string_value("prod");

  envoy::type::matcher::v3::MetadataMatcher matcher;
  matcher.set_filter("envoy.filter.x");
  matcher.add_path()->set_key("label");
  matcher.set_invert(true);

  matcher.mutable_value()->mutable_string_match()->set_exact("test");
  EXPECT_TRUE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
  matcher.mutable_value()->mutable_string_match()->set_exact("prod");
  EXPECT_FALSE(Envoy::Matchers::MetadataMatcher(matcher, context_).match(metadata));
}

class StringMatcher : public BaseTest {};

TEST_F(StringMatcher, ExactMatchIgnoreCase) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_exact("exact");
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("exact"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("EXACT"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("exacz"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));

  matcher.set_ignore_case(true);
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("exact"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("EXACT"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("exacz"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));
}

TEST_F(StringMatcher, ExactMatchIgnoreCaseStringRepresentation) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_exact("eXaCt");
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "eXaCt");

  matcher.set_ignore_case(true);
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "eXaCt");
}

TEST_F(StringMatcher, PrefixMatchIgnoreCase) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_prefix("prefix");
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("prefix-abc"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("PREFIX-ABC"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("prefiz-abc"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));

  matcher.set_ignore_case(true);
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("prefix-abc"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("PREFIX-ABC"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("prefiz-abc"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));
}

TEST_F(StringMatcher, PrefixMatchIgnoreCaseStringRepresentation) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_prefix("pReFix");
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "pReFix");

  matcher.set_ignore_case(true);
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "pReFix");
}

TEST_F(StringMatcher, SuffixMatchIgnoreCase) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_suffix("suffix");
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("abc-suffix"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("ABC-SUFFIX"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("abc-suffiz"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));

  matcher.set_ignore_case(true);
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("abc-suffix"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("ABC-SUFFIX"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("abc-suffiz"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));
}

TEST_F(StringMatcher, SuffixMatchIgnoreCaseStringRepresentation) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_suffix("sUfFix");
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "sUfFix");

  matcher.set_ignore_case(true);
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "sUfFix");
}

TEST_F(StringMatcher, ContainsMatchIgnoreCase) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_contains("contained-str");
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("abc-contained-str-def"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("contained-str"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("ABC-Contained-Str-DEF"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("abc-container-int-def"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));

  matcher.set_ignore_case(true);
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("abc-contained-str-def"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("abc-cOnTaInEd-str-def"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("abc-ContAineR-str-def"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("other"));
}

TEST_F(StringMatcher, ContainsMatchIgnoreCaseStringRepresentation) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_contains("ConTained-STR");
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "ConTained-STR");

  matcher.set_ignore_case(true);
  EXPECT_EQ(Matchers::StringMatcherImpl(matcher, context_).stringRepresentation(), "contained-str");
}

TEST_F(StringMatcher, SafeRegexValue) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.mutable_safe_regex()->mutable_google_re2();
  matcher.mutable_safe_regex()->set_regex("foo.*");
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("foo"));
  EXPECT_TRUE(Matchers::StringMatcherImpl(matcher, context_).match("foobar"));
  EXPECT_FALSE(Matchers::StringMatcherImpl(matcher, context_).match("bar"));
}

TEST_F(StringMatcher, SafeRegexValueIgnoreCase) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_ignore_case(true);
  matcher.mutable_safe_regex()->mutable_google_re2();
  matcher.mutable_safe_regex()->set_regex("foo");
  EXPECT_THROW_WITH_MESSAGE(Matchers::StringMatcherImpl(matcher, context_).match("foo"),
                            EnvoyException, "ignore_case has no effect for safe_regex.");
}

TEST_F(StringMatcher, NoMatcherRejected) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.set_ignore_case(true);
  EXPECT_THROW_WITH_MESSAGE(
      Matchers::StringMatcherImpl(matcher, context_).match("foo"), EnvoyException,
      fmt::format("Configuration must define a matcher: {}", matcher.DebugString()));
}

// Validates the amount of memory that is being used by the different string
// matchers. Requested as part of https://github.com/envoyproxy/envoy/pull/37782.
TEST_F(StringMatcher, Memory) {
  const uint32_t matchers_num = 1000;
  // Prefix matcher.
  {
    // Add 1000 Prefix-String Matchers of varying string lengths (1 to 1000).
    std::vector<Matchers::StringMatcherImpl> all_matchers;
    all_matchers.reserve(matchers_num);
    Memory::TestUtil::MemoryTest memory_test;
    for (uint32_t i = 0; i < matchers_num; ++i) {
      envoy::type::matcher::v3::StringMatcher matcher;
      matcher.set_prefix(std::string(i + 1, 'a'));
      all_matchers.emplace_back(Matchers::StringMatcherImpl(matcher, context_));
    }
    const size_t prefix_consumed_bytes = memory_test.consumedBytes();
    // The memory constraints were added to ensure that the amount of memory
    // used by matchers is carefully analyzed. These constraints can be relaxed
    // when additional features are added, but it should be done in a thoughtful manner.
    // Adding 3*8192 bytes because tcmalloc consumption estimation may return
    // different values depending on memory alignment.
    EXPECT_MEMORY_LE(prefix_consumed_bytes, 530176 + 3 * 8192);
  }
  // Regex matcher.
  {
    // Add 1000 Regex-String Matchers of varying string lengths (1 to 1000).
    std::vector<Matchers::StringMatcherImpl> all_matchers;
    all_matchers.reserve(matchers_num);
    Memory::TestUtil::MemoryTest memory_test;
    for (uint32_t i = 0; i < matchers_num; ++i) {
      envoy::type::matcher::v3::StringMatcher matcher;
      matcher.mutable_safe_regex()->mutable_google_re2();
      matcher.mutable_safe_regex()->set_regex(std::string(i + 1, 'a'));
      all_matchers.emplace_back(Matchers::StringMatcherImpl(matcher, context_));
    }
    const size_t regex_consumed_bytes = memory_test.consumedBytes();
    // The memory constraints were added to ensure that the amount of memory
    // used by matchers is carefully analyzed. These constraints can be relaxed
    // when additional features are added, but it should be done in a thoughtful  manner.
    // Adding 10*8192 bytes because tcmalloc consumption estimation may return
    // different values depending on memory alignment.
    EXPECT_MEMORY_LE(regex_consumed_bytes, 15038016 + 10 * 8192);
  }
}

class PathMatcher : public BaseTest {};

TEST_F(PathMatcher, MatchExactPath) {
  const auto matcher = Envoy::Matchers::PathMatcher::createExact("/exact", false, context_);

  EXPECT_TRUE(matcher->match("/exact"));
  EXPECT_TRUE(matcher->match("/exact?param=val"));
  EXPECT_TRUE(matcher->match("/exact#fragment"));
  EXPECT_TRUE(matcher->match("/exact#fragment?param=val"));
  EXPECT_FALSE(matcher->match("/EXACT"));
  EXPECT_FALSE(matcher->match("/exacz"));
  EXPECT_FALSE(matcher->match("/exact-abc"));
  EXPECT_FALSE(matcher->match("/exacz?/exact"));
  EXPECT_FALSE(matcher->match("/exacz#/exact"));
}

TEST_F(PathMatcher, MatchExactPathIgnoreCase) {
  const auto matcher = Envoy::Matchers::PathMatcher::createExact("/exact", true, context_);

  EXPECT_TRUE(matcher->match("/exact"));
  EXPECT_TRUE(matcher->match("/EXACT"));
  EXPECT_TRUE(matcher->match("/exact?param=val"));
  EXPECT_TRUE(matcher->match("/Exact#fragment"));
  EXPECT_TRUE(matcher->match("/EXACT#fragment?param=val"));
  EXPECT_FALSE(matcher->match("/exacz"));
  EXPECT_FALSE(matcher->match("/exact-abc"));
  EXPECT_FALSE(matcher->match("/exacz?/exact"));
  EXPECT_FALSE(matcher->match("/exacz#/exact"));
}

TEST_F(PathMatcher, MatchPrefixPath) {
  const auto matcher = Envoy::Matchers::PathMatcher::createPrefix("/prefix", false, context_);

  EXPECT_TRUE(matcher->match("/prefix"));
  EXPECT_TRUE(matcher->match("/prefix-abc"));
  EXPECT_TRUE(matcher->match("/prefix?param=val"));
  EXPECT_TRUE(matcher->match("/prefix#fragment"));
  EXPECT_TRUE(matcher->match("/prefix#fragment?param=val"));
  EXPECT_FALSE(matcher->match("/PREFIX"));
  EXPECT_FALSE(matcher->match("/prefiz"));
  EXPECT_FALSE(matcher->match("/prefiz?/prefix"));
  EXPECT_FALSE(matcher->match("/prefiz#/prefix"));
}

TEST_F(PathMatcher, MatchPrefixPathIgnoreCase) {
  const auto matcher = Envoy::Matchers::PathMatcher::createPrefix("/prefix", true, context_);

  EXPECT_TRUE(matcher->match("/prefix"));
  EXPECT_TRUE(matcher->match("/prefix-abc"));
  EXPECT_TRUE(matcher->match("/Prefix?param=val"));
  EXPECT_TRUE(matcher->match("/Prefix#fragment"));
  EXPECT_TRUE(matcher->match("/PREFIX#fragment?param=val"));
  EXPECT_TRUE(matcher->match("/PREFIX"));
  EXPECT_FALSE(matcher->match("/prefiz"));
  EXPECT_FALSE(matcher->match("/prefiz?/prefix"));
  EXPECT_FALSE(matcher->match("/prefiz#/prefix"));
}

TEST_F(PathMatcher, SlashPrefixMatcherShared) {
  // Create 3 matchers and verify that the same instance is being reused for them.
  const auto matcher1 = Envoy::Matchers::PathMatcher::createPrefix("/", false, context_);
  const auto matcher2 = Envoy::Matchers::PathMatcher::createPrefix("/", false, context_);
  const auto matcher3 = Envoy::Matchers::PathMatcher::createPrefix("/", true, context_);

  EXPECT_EQ(matcher1, matcher2);
  EXPECT_EQ(matcher1, matcher3);

  // Sanity check that the matcher works as expected.
  EXPECT_TRUE(matcher1->match("/bla"));
  EXPECT_FALSE(matcher1->match("bla"));
}

TEST_F(PathMatcher, EmptyPrefixMatcherShared) {
  // Create 3 matchers and verify that the same instance is being reused for them.
  const auto matcher1 = Envoy::Matchers::PathMatcher::createPrefix("", false, context_);
  const auto matcher2 = Envoy::Matchers::PathMatcher::createPrefix("", false, context_);
  const auto matcher3 = Envoy::Matchers::PathMatcher::createPrefix("", true, context_);

  EXPECT_EQ(matcher1, matcher2);
  EXPECT_EQ(matcher1, matcher3);

  // Sanity check that the matcher works as expected.
  EXPECT_TRUE(matcher1->match("/bla"));
  EXPECT_TRUE(matcher1->match("bla"));
}

TEST_F(PathMatcher, MatchSuffixPath) {
  envoy::type::matcher::v3::PathMatcher matcher;
  matcher.mutable_path()->set_suffix("suffix");

  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/suffix"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/abc-suffix"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/suffix?param=val"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/suffix#fragment"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/suffix#fragment?param=val"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/suffiz"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/suffiz?param=suffix"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/suffiz#suffix"));
}

TEST_F(PathMatcher, MatchContainsPath) {
  envoy::type::matcher::v3::PathMatcher matcher;
  matcher.mutable_path()->set_contains("contains");

  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/contains"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/abc-contains"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/contains-abc"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/abc-contains-def"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/abc-contains-def?param=val"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/abc-contains-def#fragment"));
  EXPECT_FALSE(
      Matchers::PathMatcher(matcher, context_).match("/abc-def#containsfragment?param=contains"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/abc-curtains-def"));
}

TEST_F(PathMatcher, MatchRegexPath) {
  envoy::type::matcher::v3::StringMatcher matcher;
  matcher.mutable_safe_regex()->mutable_google_re2();
  matcher.mutable_safe_regex()->set_regex(".*regex.*");

  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/regex"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/regex/xyz"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/xyz/regex"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/regex?param=val"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/regex#fragment"));
  EXPECT_TRUE(Matchers::PathMatcher(matcher, context_).match("/regex#fragment?param=val"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/regez"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/regez?param=regex"));
  EXPECT_FALSE(Matchers::PathMatcher(matcher, context_).match("/regez#regex"));
}

class FilterStateMatcher : public BaseTest {};

TEST_F(FilterStateMatcher, MatchAbsentFilterState) {
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key("test.key");
  matcher.mutable_string_match()->set_exact("exact");
  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_FALSE((*filter_state_matcher)->match(filter_state));
}

class TestObject : public StreamInfo::FilterState::Object {
public:
  TestObject(absl::optional<std::string> value) : value_(value) {}
  absl::optional<std::string> serializeAsString() const override { return value_; }

private:
  absl::optional<std::string> value_;
};

TEST_F(FilterStateMatcher, MatchFilterStateWithoutString) {
  const std::string key = "test.key";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  matcher.mutable_string_match()->set_exact("exact");
  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(key, std::make_shared<TestObject>(absl::nullopt),
                       StreamInfo::FilterState::StateType::ReadOnly);
  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_FALSE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, MatchFilterStateDifferentString) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  matcher.mutable_string_match()->set_exact(value);
  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(key,
                       std::make_shared<TestObject>(absl::make_optional<std::string>("different")),
                       StreamInfo::FilterState::StateType::ReadOnly);
  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_FALSE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, MatchFilterState) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  matcher.mutable_string_match()->set_exact(value);
  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(key, std::make_shared<TestObject>(absl::make_optional<std::string>(value)),
                       StreamInfo::FilterState::StateType::ReadOnly);
  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_TRUE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, MatchFilterStateAddressMatchIpv4) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  auto* cidrv4 = matcher.mutable_address_match()->add_ranges();
  cidrv4->set_address_prefix("4.5.6.7");
  cidrv4->mutable_prefix_len()->set_value(32);
  auto* cidrv6 = matcher.mutable_address_match()->add_ranges();
  cidrv6->set_address_prefix("2001:db8::");
  cidrv6->mutable_prefix_len()->set_value(32);

  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(
      key,
      std::make_shared<Network::Address::InstanceAccessor>(
          Envoy::Network::Utility::parseInternetAddressNoThrow("4.5.6.7", 456, false)),
      StreamInfo::FilterState::StateType::Mutable);

  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_TRUE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, NoMatchFilterStateAddressMatchIpv4) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  auto* cidrv4 = matcher.mutable_address_match()->add_ranges();
  cidrv4->set_address_prefix("4.5.6.7");
  cidrv4->mutable_prefix_len()->set_value(32);
  auto* cidrv6 = matcher.mutable_address_match()->add_ranges();
  cidrv6->set_address_prefix("2001:db8::");
  cidrv6->mutable_prefix_len()->set_value(32);

  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(
      key,
      std::make_shared<Network::Address::InstanceAccessor>(
          Envoy::Network::Utility::parseInternetAddressNoThrow("4.5.6.8", 456, false)),
      StreamInfo::FilterState::StateType::Mutable);

  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_FALSE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, MatchFilterStateAddressMatchIpv6) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  auto* cidrv4 = matcher.mutable_address_match()->add_ranges();
  cidrv4->set_address_prefix("4.5.6.7");
  cidrv4->mutable_prefix_len()->set_value(32);
  auto* cidrv6 = matcher.mutable_address_match()->add_ranges();
  cidrv6->set_address_prefix("2001:db8::");
  cidrv6->mutable_prefix_len()->set_value(32);

  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(
      key,
      std::make_shared<Network::Address::InstanceAccessor>(
          Envoy::Network::Utility::parseInternetAddressNoThrow("2001:db8::1", 8080, false)),
      StreamInfo::FilterState::StateType::Mutable);

  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_TRUE((*filter_state_matcher)->match(filter_state));
}

TEST_F(FilterStateMatcher, NoMatchFilterStateAddressMatchIpv6) {
  const std::string key = "test.key";
  const std::string value = "exact_value";
  envoy::type::matcher::v3::FilterStateMatcher matcher;
  matcher.set_key(key);
  auto* cidrv4 = matcher.mutable_address_match()->add_ranges();
  cidrv4->set_address_prefix("4.5.6.7");
  cidrv4->mutable_prefix_len()->set_value(32);
  auto* cidrv6 = matcher.mutable_address_match()->add_ranges();
  cidrv6->set_address_prefix("2001:db8::");
  cidrv6->mutable_prefix_len()->set_value(32);

  StreamInfo::FilterStateImpl filter_state(StreamInfo::FilterState::LifeSpan::Connection);
  filter_state.setData(
      key,
      std::make_shared<Network::Address::InstanceAccessor>(
          Envoy::Network::Utility::parseInternetAddressNoThrow("2001:db7::1", 8080, false)),
      StreamInfo::FilterState::StateType::Mutable);

  auto filter_state_matcher = Matchers::FilterStateMatcher::create(matcher, context_);
  ASSERT_TRUE(filter_state_matcher.ok());
  EXPECT_FALSE((*filter_state_matcher)->match(filter_state));
}

} // namespace
} // namespace Matcher
} // namespace Envoy
