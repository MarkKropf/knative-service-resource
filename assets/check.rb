#!/usr/bin/env ruby

require_relative 'lib/knative_service/checker'

KnativeServiceChecker.new.run(ARGF.read)
