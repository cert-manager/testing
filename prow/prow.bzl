# +skip_license_check

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_k8s//k8s:object.bzl", "k8s_object")
load("@io_bazel_rules_k8s//k8s:objects.bzl", "k8s_objects")

MULTI_KIND = None

CORE_CONTEXT = "{STABLE_PROW_CONTEXT}"  # For components like hook
BUILD_CONTEXT = "{STABLE_BUILD_CONTEXT}"  # For untrusted test code

def object(name, context=CORE_CONTEXT, **kwargs):
  k8s_object(
      name = name,
      context = context,
      **kwargs
  )

# component generates k8s_object rules and returns a {kind: [targets]} map.
#
# This will generate a k8s_object rule for each specified kind.
# Use MULTI_KIND for a multi-document yaml (this returns nothing).
# Assumes files exist at <cmd>_<kind>.yaml
#
# Concretely, component("hook", "service", "deployment") will create the following:
#   object("hook_service", kind="service", template=":hook_service.yaml")
#   object("hook_deployment", kind="deployment", template=":hook_deployment.yaml")
# And return the following:
#   {
#     "hook": [":hook_service", ":hook_deployment",
#     "service": [":hook_service"],
#     "deployment": [":hook_deployment"],
#   }
def component(cmd, *kinds, images = {}, **kwargs):
  targets = {}
  for k in kinds:
      if k == MULTI_KIND:
        n = cmd
      else:
        n = "%s_%s" % (cmd, k)
      kwargs["name"] = n
      kwargs["kind"] = k
      kwargs["template"] = ":%s.yaml" % n
      args = dict(kwargs)
      if k == "deployment" or k == "cronjob":
        args["images"] = images
      object(**args)
      if k != MULTI_KIND:
        targets.setdefault(cmd,[]).append(":%s" % n)
        targets.setdefault(k,[]).append(":%s" % n)
  return targets

# release packages multiple components into a release.
#
# Generates a k8s_objects() rule for each component and kind, as well as an
# target which includes everything.
#
# Thus you can do things like:
#   bazel run //prow/cluster:hook.apply  # Update all hook resources
#   bazel run //prow/cluster:deployment.apply  # Update all deployments in prow
#
# Concretely, the following:
#   release(
#     "fancy",
#     component("hook", "deployment", "service"),
#     compoennt("plank", "deployment"),
#   )
# Generates the five following rules:
#   k8s_objects(name = "hook", objects=[":hook_deployment", ":hook_service"])
#   k8s_objects(name = "plank", objects=[":plank_deployment"])
#   k8s_objects(name = "deployment", objects=[":hook_deployment", ":plank_deployment"])
#   k8s_objects(name = "service", objects=[":hook_service"])
#   k8s_objects(name = "fancy", objects=[":hook", ":plank", ":deployment", ":service"])
def release(name, *components):
  targets = {}
  objs = []
  for cs in components:
    for (n, ts) in cs.items():
      targets.setdefault(n, []).extend(ts)
  for (piece, ts) in targets.items():
    k8s_objects(
        name = piece,
        objects = ts,
    )
    objs.append(":%s" % piece)
  k8s_objects(
      name = name,
      objects=objs,
  )
