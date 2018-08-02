#!/usr/bin/env ruby

require_relative 'lib/knative_service/putter'

require 'json'

KnativeServicePutter.new.run(input_directory: ARGV[0], raw_input: STDIN.read)
