// +build !remoteclient

package integration

import (
	"os"
	"path/filepath"
	"text/template"

	. "github.com/containers/libpod/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var yamlTemplate = `
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2019-07-17T14:44:08Z"
  labels:
    app: {{ .Name }}
  name: {{ .Name }}
spec:
  containers:
{{ with .Containers }}
  {{ range . }}
  - command:
    {{ range .Cmd }}
    - {{.}}
    {{ end }}
    env:
    - name: PATH
      value: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
    - name: TERM
      value: xterm
    - name: HOSTNAME
    - name: container
      value: podman
    image: {{ .Image }}
    name: {{ .Name }}
    resources: {}
    securityContext:
      allowPrivilegeEscalation: true
      {{ if .Caps }}
      capabilities:
        {{ with .CapAdd }}
        add:
          {{ range . }}
          - {{.}}
          {{ end }}
        {{ end }}
        {{ with .CapDrop }}
        drop:
          {{ range . }}
          - {{.}}
          {{ end }}
        {{ end }}
      {{ end }}
      privileged: false
      readOnlyRootFilesystem: false
    workingDir: /
  {{ end }}
{{ end }}
status: {}
`

type Pod struct {
	Name       string
	Containers []Container
}

type Container struct {
	Cmd     []string
	Image   string
	Name    string
	Caps    bool
	CapAdd  []string
	CapDrop []string
}

func generateKubeYaml(ctrs []Container, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	testPod := Pod{"test", ctrs}

	t, err := template.New("pod").Parse(yamlTemplate)
	if err != nil {
		return err
	}

	if err := t.Execute(f, testPod); err != nil {
		return err
	}

	return nil
}

var _ = Describe("Podman generate kube", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
		podmanTest.SeedImages()
	})

	AfterEach(func() {
		podmanTest.Cleanup()
		f := CurrentGinkgoTestDescription()
		processTestResult(f)
	})

	It("podman play kube test correct command", func() {
		ctrName := "testCtr"
		ctrCmd := []string{"top"}
		testContainer := Container{ctrCmd, ALPINE, ctrName, false, nil, nil}
		tempFile := filepath.Join(podmanTest.TempDir, "kube.yaml")

		err := generateKubeYaml([]Container{testContainer}, tempFile)
		Expect(err).To(BeNil())

		kube := podmanTest.Podman([]string{"play", "kube", tempFile})
		kube.WaitWithDefaultTimeout()
		Expect(kube.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"inspect", ctrName})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.OutputToString()).To(ContainSubstring(ctrCmd[0]))
	})

	It("podman play kube cap add", func() {
		ctrName := "testCtr"
		ctrCmd := []string{"cat", "/proc/self/status"}
		capAdd := "CAP_SYS_ADMIN"
		testContainer := Container{ctrCmd, ALPINE, ctrName, true, []string{capAdd}, nil}
		tempFile := filepath.Join(podmanTest.TempDir, "kube.yaml")

		err := generateKubeYaml([]Container{testContainer}, tempFile)
		Expect(err).To(BeNil())

		kube := podmanTest.Podman([]string{"play", "kube", tempFile})
		kube.WaitWithDefaultTimeout()
		Expect(kube.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"inspect", ctrName})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.OutputToString()).To(ContainSubstring(capAdd))
	})

	It("podman play kube cap add", func() {
		ctrName := "testCtr"
		ctrCmd := []string{"cat", "/proc/self/status"}
		capDrop := "CAP_SYS_ADMIN"
		testContainer := Container{ctrCmd, ALPINE, ctrName, true, []string{capDrop}, nil}
		tempFile := filepath.Join(podmanTest.TempDir, "kube.yaml")

		err := generateKubeYaml([]Container{testContainer}, tempFile)
		Expect(err).To(BeNil())

		kube := podmanTest.Podman([]string{"play", "kube", tempFile})
		kube.WaitWithDefaultTimeout()
		Expect(kube.ExitCode()).To(Equal(0))

		inspect := podmanTest.Podman([]string{"inspect", ctrName})
		inspect.WaitWithDefaultTimeout()
		Expect(inspect.ExitCode()).To(Equal(0))
		Expect(inspect.OutputToString()).To(ContainSubstring(capDrop))
	})
})
