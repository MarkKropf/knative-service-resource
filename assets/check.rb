#!/usr/bin/env ruby

# frozen_string_literal: true

require_relative 'lib/common'

require 'json'

class KnativeServiceChecker
  def run(raw_input)
    source = Source.from_input(raw_input)
    client = KnativeClient.new(source)

    if first_time_checked?(raw_input)
      puts [client.observed_generation].to_json
    else
      latest_version_in_concourse = Version.from_input(raw_input)
      latest_version_in_knative = client.observed_generation

      if latest_version_in_knative > latest_version_in_concourse
        previously_unseen_generations = client.all_versions.select { |ver|
          ver <= latest_version_in_knative && ver > latest_version_in_concourse
        }

        puts previously_unseen_generations.to_json
      else
        puts [latest_version_in_concourse].to_json
      end
    end
  end

  private

  def first_time_checked?(raw_input)
    parsed_input = JSON.parse(raw_input)
    !parsed_input.has_key?('version') || parsed_input['version'].nil?
  end
end

KnativeServiceChecker.new.run(ARGF.read)
