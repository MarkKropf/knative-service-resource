
require_relative 'common'

require 'json'

class KnativeServicePutter
  def run(input_directory:, raw_input:)
    source = Source.from_input(raw_input)
    params = Params.from_input(input_directory, raw_input)
    client = KnativeClient.new(source)

    patch = [
        {
            op: 'replace',
            path: '/spec/runLatest/configuration/revisionTemplate/spec/container/image',
            value: "#{params.image_repository}@#{params.image_digest}"
        },
        {
            op: 'replace',
            path: '/spec/runLatest/configuration/revisionTemplate/metadata/annotations/concourseBuild',
            value: ConcourseBuildUrlHelper::url
        }
    ]

    client.patch_service(name: source.name, patch: patch)

    version = client.observed_generation
    raw_configuration = client.get_configuration_at_generation(version.configuration_generation)
    metadata = VersionMetadata.from_input(raw_configuration)

    put_output = {version: version, metadata: metadata}.to_json

    puts put_output
  end
end
