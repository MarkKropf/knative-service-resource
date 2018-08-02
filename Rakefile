PROJECT_NAME='cf-elafros-dog'
IMAGE_NAME='knative-service-resource'
IMAGE_REPOSITORY="gcr.io/#{PROJECT_NAME}/#{IMAGE_NAME}"

desc 'Runs image:update'
task default: 'image:update'

namespace :image do
	desc 'Update the resource image'
	task update: [:build, :push]

	desc 'builds and tags the resource docker image'
	task :build do
		tag_opt = "--tag #{IMAGE_REPOSITORY}"

		puts `docker build . #{tag_opt}`
	end

	desc 'pushes the resource docker image to GCR'
	task :push do
		puts `docker push #{IMAGE_REPOSITORY}`
	end
end