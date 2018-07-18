# frozen_string_literal: true

require 'kubeclient'
require 'json'

class KnativeClient
  KNATIVE_SERVING_API_GROUP = 'serving.knative.dev'
  KNATIVE_SERVING_GROUP_VERSION = 'v1alpha1'

  def initialize(source)
    @source = source

    Tempfile.create('kubernetes-ca.crt') do |cert_file|
      cert_file.write(@source.kubernetes_ca) # kubeclient gem assumes it's a file on-disk

      @client = Kubeclient::Client.new(
          "#{@source.kubernetes_uri}/apis/#{KNATIVE_SERVING_API_GROUP}",
          KNATIVE_SERVING_GROUP_VERSION,
          auth_options: {bearer_token: @source.kubernetes_token},
          ssl_options: { ca_file: cert_file.path }
      )
    end

    @client.discover

    @client
  end

  def observed_generation
    service_configuration = @client.get_configuration(@source.name, 'default')
    observed = service_configuration[:status][:observedGeneration]

    Version.new(observed)
  end

  def all_versions
    @client.get_revisions(label_selector: "serving.knative.dev/configuration=#{@source.name}").map do |data|
      Version.new(data[:metadata][:annotations]['serving.knative.dev/configurationGeneration'].to_i)
    end
  end

  def get_configuration_at_generation(generation)
    @client.get_configuration(@source.name, 'default', {field_selector: "generation=#{generation}", as: :raw})
  end

  def patch_service(name:, patch:)
    with_patch_header do
      @client.patch_service(name, patch, 'default')
    end
  end

  private

  def with_patch_header # See https://github.com/abonas/kubeclient/issues/268
    old = @client.headers['Content-Type']
    @client.headers['Content-Type'] = 'application/json-patch+json'
    yield
  ensure
    @client.headers['Content-Type'] = old
  end
end

class Version
  include Comparable

  attr_reader :configuration_generation

  def initialize(configuration_generation)
    @configuration_generation = configuration_generation
  end

  def self.from_input(raw_input)
    parsed_input = JSON.parse(raw_input)
    configuration_generation = parsed_input['version']['configuration_generation'].to_i

    Version.new(configuration_generation)
  end

  def to_json(*args_ignored)
    to_hash.to_json
  end

  def to_hash
    { configuration_generation: @configuration_generation.to_s }
  end

  def <=>(other)
    @configuration_generation <=> other.configuration_generation
  end
end

class Source
  attr_reader :name, :kubernetes_uri, :kubernetes_token, :kubernetes_ca

  def initialize(name, kubernetes_uri, kubernetes_token, kubernetes_ca)
    @name = name
    @kubernetes_uri = kubernetes_uri
    @kubernetes_token = kubernetes_token
    @kubernetes_ca = kubernetes_ca
  end

  def self.from_input(raw_input)
    parsed_input = JSON.parse(raw_input)
    source = parsed_input['source']

    Source.new(source['name'], source['kubernetes_uri'], source['kubernetes_token'], source['kubernetes_ca'])
  end
end

class Params
  attr_reader :image_repository, :image_digest

  def initialize(input_directory:, image_repository:, image_digest_path:)
    @image_repository = image_repository
    @image_digest = File.read("#{input_directory}/#{image_digest_path}").chomp
  end

  def self.from_input(input_directory, raw_input)
    parsed_input = JSON.parse(raw_input)
    params = parsed_input['params']

    Params.new(input_directory: input_directory, image_repository: params['image_repository'], image_digest_path: params['image_digest_path'])
  end
end

class VersionMetadata
  def initialize(cluster_name:, creation_timestamp:, resource_version:, uid:)
    @cluster_name = cluster_name
    @creation_timestamp = creation_timestamp
    @resource_version = resource_version
    @uid = uid
  end

  def self.from_input(raw_input)
    parsed_input = JSON.parse(raw_input)
    metadata = parsed_input['metadata']

    VersionMetadata.new(
        cluster_name: metadata['clusterName'],
        creation_timestamp: metadata['creationTimestamp'],
        resource_version: metadata['resourceVersion'],
        uid: metadata['uid']
    )
  end


  def to_json(*args_ignored)
    [
        {name: 'cluster_name', value: @cluster_name},
        {name: 'creation_timestamp', value: @creation_timestamp},
        {name: 'resource_version', value: @resource_version},
        {name: 'uid', value: @uid}
    ].to_json
  end
end

class ConcourseBuildUrlHelper
  def self.url
    atc_external_url = ENV['ATC_EXTERNAL_URL']
    team_name = ENV['BUILD_TEAM_NAME']

    if ENV.include?'BUILD_PIPELINE_NAME'
      pipeline_name = ENV['BUILD_PIPELINE_NAME']
      job_name = ENV['BUILD_JOB_NAME']
      build_name = ENV['BUILD_NAME']

      "#{atc_external_url}/teams/#{team_name}/pipelines/#{pipeline_name}/jobs/#{job_name}/builds/#{build_name}"
    else
      "#{atc_external_url}/teams/#{team_name} (one-off build)"
    end
  end
end

