package helper

import (
	"math/rand"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/osde2e/pkg/config"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// New creates H, a helper used to expose common testing functions.
func New() *H {
	helper := &H{
		Config: config.Cfg,
	}
	ginkgo.BeforeEach(helper.Setup)
	ginkgo.AfterEach(helper.Cleanup)
	return helper
}

// H configures clients and sets up and destroys Projects for test isolation.
type H struct {
	// embed test configuration
	*config.Config

	// internal
	restConfig *rest.Config
	proj       string
}

// Setup configures a *rest.Config using the embedded kubeconfig then sets up a Project for tests to run in.
func (h *H) Setup() {
	var err error
	h.restConfig, err = clientcmd.RESTConfigFromKubeConfig(h.Kubeconfig)
	Expect(err).ShouldNot(HaveOccurred(), "failed to configure client")

	// setup project to run tests
	suffix := randomStr(5)
	proj, err := h.createProject(suffix)
	Expect(err).ShouldNot(HaveOccurred(), "failed to create project")
	Expect(proj).ShouldNot(BeNil())

	h.proj = proj.Name
}

// Cleanup deletes a Project after tests have been ran.
func (h *H) Cleanup() {
	err := h.cleanup(h.proj)
	Expect(err).ShouldNot(HaveOccurred(), "could not delete project '%s'", h.proj)

	h.restConfig = nil
	h.proj = ""
}

// CurrentProject returns the project being used for testing.
func (h *H) CurrentProject() string {
	return h.proj
}
