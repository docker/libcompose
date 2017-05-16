
wrappedNode(label: 'linux && x86_64') {
  deleteDir()
  checkout scm
  def image
  try {
    stage "build image"
    image = docker.build("dockerbuildbot/libcompose:${gitCommit()}")

    stage "validate"
    makeTask(image, "validate")

    stage "unit-test"
    makeTask(image, "test-unit")

    def integrations = [:]

    daemon_versions = ["1.9.1", "1.10.3", "1.11.2", "1.12.1"]
    for (int i=0; i < daemon_versions.size(); i++) {
        daemon_version = "${daemon_versions[i]}"
    	integrations["test-${daemon_version}"] = {
    	    makeTask(image, "test-integration", ["DAEMON_VERSION=${daemon_version}", "SHOWWARNING=false"])
    	}
    }

    stage "integrations tests"
    parallel integrations

    stage "build"
    makeTask(image, "cross-binary")
  } finally {
    try { archive "bundles" } catch (Exception exc) {}
    if (image) { sh "docker rmi ${image.id} ||:" }
  }
}

def makeTask(image, task, envVars=null) {
  // could send in the full list of envVars for each call or provide default env vars like this:
  withEnv((envVars ?: []) + ["LIBCOMPOSE_IMAGE=${image.id}"]) { // would need `def image` at top level of file instead of in the nested block
    withChownWorkspace {
      timeout(60) {
        sh "make -e ${task}"
      }
    }
  }
}
