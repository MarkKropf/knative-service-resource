#!/usr/bin/env ruby

require_relative 'lib/common'

require 'json'
require 'yaml'

class KnativeServiceGetter
  def run(output_directory:, raw_input:)
    version = Version.from_input(raw_input)
    source = Source.from_input(raw_input)

    client = KnativeClient.new(source)

    generation = version.configuration_generation
    raw_configuration = client.get_configuration_at_generation(generation)
    version_metadata = VersionMetadata.from_input(raw_configuration)

    # both of these seem unnecessarily fiddly
    pretty_json = JSON.pretty_generate(JSON.parse(raw_configuration))
    pretty_yaml = JSON.parse(raw_configuration).to_hash.to_yaml

    json_out_path = "#{output_directory}/configuration.json"
    File.write(json_out_path, pretty_json)

    yaml_out_path = "#{output_directory}/configuration.yaml"
    File.write(yaml_out_path, pretty_yaml)

    get_output = {version: version, metadata: version_metadata}.to_json

    puts get_output
  end
end

KnativeServiceGetter.new.run(output_directory: ARGV[0], raw_input: STDIN.read)
