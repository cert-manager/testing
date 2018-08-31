load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_bundle")

def bootstrap_image(name, base, runner = "//legacy/bootstrap:runner"):
    container_image(
        name = name,
        base = base,
        directory = "/workspace",
        # TODO: lookup name of 'runner' target
        entrypoint = ["/bin/bash", "/runner"],
        files = [runner],
    )

def bootstrap_image_bundle(name, images, stamp):
    i = 0
    new_images = {}
    for k, v in images.items():
        new_name = "%s.%d" % (name, i)
        bootstrap_image(new_name, v)
        new_images[k] = ":%s" % new_name
        i += 1

    container_bundle(
        name = name,
        images = new_images,
        stamp = stamp,
    )
