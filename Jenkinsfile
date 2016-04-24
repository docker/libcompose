wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  checkout scm
  def image
  try {
    // TODO: split up into phases, create real test reports, etc.
    stage "build image"
    image = docker.build("dockerbuildbot/libcompose:${gitCommit()}")

    stage "validate, test, build"
    withEnv(["TESTVERBOSE=1", "LIBCOMPOSE_IMAGE=${image.id}"]) {
      withChownWorkspace {
        timeout(60) {
          sh "make -e all"
        }
      }
    }
  } finally {
    try { archive "bundles" } catch (Exception exc) {}
    if (image) { sh "docker rmi ${image.id} ||:" }
  }
}
