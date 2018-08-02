#!/usr/bin/env ruby

require_relative 'lib/knative_service/getter'

KnativeServiceGetter.new.run(output_directory: ARGV[0], raw_input: STDIN.read)
