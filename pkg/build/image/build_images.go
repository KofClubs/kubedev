package image

import (
	"kubedev/pkg/env"
	imageGetter "kubedev/pkg/image"
	"log"
	"os/exec"
	"reflect"
)

var HardCodeString []string = []string{"KUBE_FASTBUILD", "KUBE_BUILD_HYPERKUBE", "KUBE_BUILD_CONFORMANCE", "KUBE_DOCKER_IMAGE_TAG", "KUBE_DOCKER_REGISTRY", "KUBE_GIT_VERSION_FILE", "KUBE_BUILD_PULL_LATEST_IMAGES", "GOFLAGS"}

type ImageConfig struct {
	KubeFastBuild             string
	KubeBuildHyperkube        string
	KubeBuildConformance      string
	KubeDockerImageTag        string
	KubeDockerRegistry        string
	KubeGitVersionFile        string
	KubeBuildPullLatestImages string
	GoFlags                   string
}

func NewDefaultImageConfig() *ImageConfig {
	return &ImageConfig{
		KubeFastBuild:             "true",
		KubeBuildHyperkube:        "n",
		KubeBuildConformance:      "n",
		KubeDockerImageTag:        "",
		KubeDockerRegistry:        "",
		KubeGitVersionFile:        env.KubeVersionFile,
		KubeBuildPullLatestImages: "n",
		GoFlags:                   "-tags=nokmem",
	}
}

func (i *ImageConfig) String() string {
	v := reflect.ValueOf(*i)
	count := v.NumField()
	imageConfigString := ""
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if f.String() != "" {
			imageConfigString = imageConfigString + " " + HardCodeString[i] + "=" + f.String()
		}

	}
	return imageConfigString
}

func (i *ImageConfig) SetKubeDockerImageTag(s string) {
	i.KubeDockerImageTag = s
}

func (i *ImageConfig) SetKubeDockerRegistry(s string) {
	i.KubeDockerRegistry = s
}

func (i *ImageConfig) SetKubeGitVersionFile(s string) {
	i.KubeGitVersionFile = s
}

func BuildImages() error {
	imageConfig := NewDefaultImageConfig()
	log.Printf("The image config is: %s", imageConfig.String())

	// Step 1: pull all images
	prePullImages()

	// Step 2: generate version file
	env.WriteVersionFile(env.KubeVersionFile)

	// Step 2: make release
	cmd := exec.Command("make", "release-images", imageConfig.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error building images: %s", err.Error())
		return err
	}
	log.Printf("Build images: %s", string(out))
	return nil
}

func prePullImages() {
	kubeImages := env.GetAllImages()
	imageGetter.PullImage(kubeImages.DebianBase)
	// imageGetter.PullImage(kubeImages.DebianHyperKubeBase)
	// imageGetter.PullImage(kubeImages.KubeCross)
	imageGetter.PullImage(kubeImages.KubePause)
	imageGetter.PullImage(kubeImages.DebianIptables)
}
