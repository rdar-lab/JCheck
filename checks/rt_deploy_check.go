package checks

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/rdar-lab/JCheck/common"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

func GetRTDeployCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "RTDeployCheck",
		Group:       "Artifactory",
		Description: "Deploy a large file to Artifactory, download and verify checksum",
		IsReadOnly:  false,
		CheckFunc: func(ctx context.Context) (string, error) {
			rtDetails, err := config.GetDefaultServerConf()
			if err != nil {
				return "", err
			}
			serviceManager, err := utils.CreateServiceManager(rtDetails, -1, false)
			if err != nil {
				return "", err
			}
			createRepoParams := services.NewGenericLocalRepositoryParams()
			randomSource := rand.NewSource(time.Now().UnixNano())
			randomGenerator := rand.New(randomSource)
			createRepoParams.Key = fmt.Sprintf("jcheck-%d", randomGenerator.Intn(9999))
			err = serviceManager.CreateLocalRepository().Generic(createRepoParams)
			if err != nil {
				return "", err
			}
			stateMap, err := common.GetStateMapFromContext(ctx)
			if err != nil {
				return "", err
			}
			stateMap["repo"] = createRepoParams.Key
			f, err := ioutil.TempFile("", "randomfile")
			if err != nil {
				return "", err
			}
			defer f.Close()
			defer os.Remove(f.Name())
			hasher := sha256.New()
			buf := make([]byte, 1024*1024)
			for i := 0; i < 100; i++ {
				_, err := randomGenerator.Read(buf)
				if err != nil {
					return "", err
				}
				_, err = f.Write(buf)
				if err != nil {
					return "", err
				}
				_, err = hasher.Write(buf)
				if err != nil {
					return "", err
				}
			}
			sha256Uploaded := hex.EncodeToString(hasher.Sum(nil))
			up := services.NewUploadParams()
			up.Pattern = f.Name()
			up.Target = createRepoParams.Key + f.Name()
			totalUploaded, totalFailed, err := serviceManager.UploadFiles(up)
			if err != nil {
				return "", err
			}
			if totalFailed != 0 {
				return "", errors.New("failure on upload")
			}
			if totalUploaded != 1 {
				return "", errors.New("failure on upload")
			}
			downParams := services.NewDownloadParams()
			downParams.Pattern = up.Target
			workDir, err := coreutils.GetWorkingDirectory()
			if err != nil {
				return "", err
			}
			tempFile, err := ioutil.TempFile(workDir, "download")
			if err != nil {
				return "", err
			}
			defer os.Remove(tempFile.Name())
			err = tempFile.Close()
			if err != nil {
				return "", err
			}

			downParams.Target = tempFile.Name()
			downloadSummary, err := serviceManager.DownloadFilesWithSummary(downParams)
			if err != nil {
				return "", err
			}
			if downloadSummary.TotalSucceeded != 1 {
				return "", errors.New("failure on download")
			}
			if downloadSummary.TotalFailed > 0 {
				return "", errors.New("failure on download")
			}
			hasher = sha256.New()
			_, file := path.Split(tempFile.Name())
			downloadedPath := workDir + os.TempDir() + "/" + file
			downloadedFile, err := os.Open(downloadedPath)
			defer downloadedFile.Close()
			defer os.Remove(downloadedFile.Name())
			if err != nil {
				return "", err
			}
			defer downloadedFile.Close()
			_, err = io.Copy(hasher, downloadedFile)
			if err != nil {
				return "", err
			}
			sha256Downloaded := hex.EncodeToString(hasher.Sum(nil))
			if sha256Downloaded != sha256Uploaded {
				return "", errors.New("checksums mismatch")
			}
			return "", nil
		},
		CleanupFunc: func(ctx context.Context) error {
			stateMap, err := common.GetStateMapFromContext(ctx)
			if err != nil {
				return err
			}
			repo := stateMap["repo"]
			if repo != nil {
				repoStrings, ok := repo.(string)
				if !ok {
					return errors.New("failed to cleanup repository")
				}
				rtDetails, err := config.GetDefaultServerConf()
				if err != nil {
					return err
				}
				serviceManager, err := utils.CreateServiceManager(rtDetails, -1, false)
				if err != nil {
					return err
				}
				err = serviceManager.DeleteRepository(repoStrings)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
