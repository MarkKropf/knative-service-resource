# knative-service-resource

For Concourse. This is a proof of concept, please don't use it in production.

This Resource works by tracking and updating Configurations embedded into a Service. Updating the Configuration causes Knative Serving to stamp out a new Revision.

## Source Configuration

* `name`: The name of your Knative Service.
* `kubernetes_uri`: The URI of your Kubernetes cluster controller. This can be found with `kubectl cluster-info`.
* `kubernetes_token`: Bearer token used to identify a Service Account to the Kubernetes controller.
* `kubernetes_ca`: The CA certificate from the Kubernetes controller, used to identify the cluster to this Resource.

## `check`

This resource watches the `serving.knative.dev/configurationGeneration` annotation on the Service. It relies on these numbers to be monotonic.

## `get` / `in`

Using `get` will produce two files under the Resource directory: `configuration.yaml` and `configuration.json`.

These are the same information -- the Configuration of the Service at the provided `configurationGeneration`.

## `put` / `out`

`put` updates an existing Service with a new image digest.

### Params

* `image_repository`: The Docker image repository string for the image to be used. For example, `ubuntu`, `gcr.io/your-project-name/image-name`.
* `image_digest_path`: Path to a file containing the image digest of the image. It is expected you will get this from the Docker image resource.

### Build URL

The Resource also sets a `concourseBuild` annotation on the Configuration. This is a full URL back to the Concourse build that included the `put`.

## Version metadata

The Resource provides four version metadata fields intended to help with more detailed debugging:

* `cluster_name`: The Kubernetes cluster name on the Service. This is often blank.
* `creation_timestamp`: The creation timestamp of the Service in Kubernetes.
* `resource_version`: The version number maintained by Kubernetes, as distinct from `configurationGeneration`.
* `uid`: The unique ID maintained by Kubernetes, as distinct from the Service name.

## Example

Refer to `pipeline.yml` for a pipeline I used during development.
