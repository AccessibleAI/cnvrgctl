package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cnvrgctl/pkg"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var ImagesParams = []Param{
	{Name: "registry", Value: "", Usage: "destination registry, example: docker.io"},
	{Name: "registry-repo", Value: "", Usage: "destination repository in registry, example: docker.io/<MY-REPO>"},
	{Name: "registry-user", Value: "", Usage: "registry user"},
	{Name: "registry-pass", Value: "", Usage: "registry password"},
	{Name: "dir", Value: ".", Usage: "destination/source directory for saving/loading docker images"},
	{Name: "image", Value: "", Usage: "override default images list with explicit image"},
}

var ImagesDumpParams = []Param{
	{Name: "list", Value: false, Usage: "print raw images list"},
	{Name: "pull", Value: false, Usage: "print images pull commands"},
	{Name: "save", Value: false, Usage: "print images save command"},
	{Name: "load", Value: false, Usage: "print images load command"},
	{Name: "tag", Value: false, Usage: "print images tag command"},
	{Name: "push", Value: false, Usage: "print images push command"},
}

var ImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "manage images",
}

var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump images commands",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("list") {
			DumpImagesList()
		}
		if viper.GetBool("pull") {
			DumpImagesPull()
		}
		if viper.GetBool("save") {
			DumpImagesSave()
		}
		if viper.GetBool("load") {
			dumpImagesLoad()
		}
		if viper.GetBool("tag") {
			dumpImagesTag()
		}
		if viper.GetBool("push") {
			dumpImagesPush()
		}
	},
}

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull cnvrg images",
	Run: func(cmd *cobra.Command, args []string) {
		pullImages()
	},
}

var SaveCmd = &cobra.Command{
	Use:   "save",
	Short: "save cnvrg images",
	Run: func(cmd *cobra.Command, args []string) {
		saveImages()
	},
}

var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load cnvrg images",
	Run: func(cmd *cobra.Command, args []string) {
		loadImages()
	},
}

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag cnvrg images",
	Run: func(cmd *cobra.Command, args []string) {
		tagImages()
	},
}

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "push cnvrg images",
	Run: func(cmd *cobra.Command, args []string) {
		pushImages()
	},
}

func DumpImagesList() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(image)
	}
}

func DumpImagesPull() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(fmt.Sprintf("docker pull %v", image))
	}
}

func DumpImagesSave() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(fmt.Sprintf("docker save --output %v.tar %v", imageArchiveName(image), image))
	}
}

func dumpImagesLoad() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(fmt.Sprintf("docker load < %v.tar", imageArchiveName(image)))
	}
}

func dumpImagesTag() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(fmt.Sprintf("docker tag %v %v", image, newImageTag(image)))
	}
}

func dumpImagesPush() {
	for _, image := range pkg.LoadCnvrgImages() {
		fmt.Println(fmt.Sprintf("docker push %v", newImageTag(image)))
	}
}

func imageArchiveName(image string) string {
	archiveName := strings.Replace(image, ":", "=", 1)
	archiveName = strings.ReplaceAll(archiveName, "/", "~")
	return archiveName
}

func newImageTag(image string) string {
	if viper.GetString("registry") == "" {
		logrus.Fatal("destination registry not set, please set --registry to your private registry")
	}
	registry := viper.GetString("registry")
	if viper.GetString("registry-repo") != "" {
		registry = fmt.Sprintf("%v/%v", registry, viper.GetString("registry-repo"))
	}
	newImage := strings.Split(image, "/")
	newImage[0] = registry
	return strings.Join(newImage, "/")
}

func archiveFullPath(image string) string {
	dir := viper.GetString("dir")
	if dir == "." {
		dir, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		return fmt.Sprintf("%v/%v", dir, imageArchiveName(image))
	}
	return fmt.Sprintf("%v/%v", dir, imageArchiveName(image))
}

func registryAuth() string {
	authConfig := types.AuthConfig{
		Username: viper.GetString("registry-user"),
		Password: viper.GetString("registry-pass"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)

}

func pullImages() {
	for _, image := range pkg.LoadCnvrgImages() {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		imagePullOptions := types.ImagePullOptions{}
		// add pull secret only for cnvrg app image
		if strings.Contains(image, "cnvrg/app") {
			imagePullOptions = types.ImagePullOptions{RegistryAuth: registryAuth()}
		}
		out, err := cli.ImagePull(ctx, image, imagePullOptions)
		if err != nil {
			logrus.Error(err)
		}
		buf := make([]byte, 1024)
		for {
			n, err := out.Read(buf)
			logrus.Info(strings.ReplaceAll(strings.TrimSpace(string(buf[:n])), "\\u003e", ">"))
			if err == io.EOF {
				break
			}
		}
		err = out.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func saveImages() {
	for _, image := range pkg.LoadCnvrgImages() {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		inspectData, _, err := cli.ImageInspectWithRaw(ctx, image)
		if err != nil {
			logrus.Fatal(err, " (make sure [cnvrgctl images pull] command finished without errors) ")
		}
		filePath := archiveFullPath(image)
		logrus.Infof("saving: %v path: %v id: %v", image, filePath, inspectData.ID)

		out, err := cli.ImageSave(ctx, []string{inspectData.ID})
		if err != nil {
			logrus.Error(err)
		}

		outFile, err := os.Create(filePath)
		_, err = io.Copy(outFile, out)
		err = outFile.Close()
		if err != nil {
			logrus.Fatal(err)
		}

	}
}

func loadImages() {
	for _, image := range pkg.LoadCnvrgImages() {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		archiveFile := imageArchiveName(image)
		logrus.Infof("loading image from: %v", archiveFile)
		f, err := os.Open(archiveFile)
		out, err := cli.ImageLoad(ctx, f, false)
		buf := make([]byte, 1024)
		for {
			n, err := out.Body.Read(buf)
			logrus.Info(strings.ReplaceAll(strings.TrimSpace(string(buf[:n])), "\\u003e", ">"))
			if err == io.EOF {
				break
			}
		}
	}
}

func tagImages() {
	for _, image := range pkg.LoadCnvrgImages() {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		oldImage := image
		newImage := newImageTag(image)
		logrus.Infof("tagging %v -> %v", oldImage, newImage)
		err = cli.ImageTag(ctx, oldImage, newImage)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func pushImages() {
	for _, image := range pkg.LoadCnvrgImages() {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		newImage := newImageTag(image)
		logrus.Infof("pushing image: %v", newImage)
		out, err := cli.ImagePush(ctx, newImage, types.ImagePushOptions{RegistryAuth: registryAuth()})
		if err != nil {
			logrus.Error(err)
		}
		buf := make([]byte, 1024)
		for {
			n, err := out.Read(buf)
			logrus.Info(strings.ReplaceAll(strings.TrimSpace(string(buf[:n])), "\\u003e", ">"))
			if err == io.EOF {
				break
			}
		}
		err = out.Close()
		if err != nil {
			logrus.Fatal(err)
		}

	}
}
