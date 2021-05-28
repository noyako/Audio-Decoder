package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

func startFfmpegEncrypt(url, filename, key, kid, path string, source string) (<-chan container.ContainerWaitOKBody, <-chan error, io.Reader) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("client:", err)
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	_, err = cli.ImagePull(ctx, viper.GetString("converter.image"), types.ImagePullOptions{})
	if err != nil {
		log.Println(err)
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}
	// io.Copy(os.Stdout, reader)
	fmt.Println(path, source)
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: viper.GetString("converter.image"),

			Entrypoint: []string{"ffmpeg", "-i", url,
				"-f", "mp4",
				"-encryption_scheme", "cenc-aes-ctr",
				"-encryption_key", key,
				"-encryption_kid", kid,
				filename,
			},
			WorkingDir: path,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: source,
					Target: path,
				},
			},
			AutoRemove: true,
		}, nil, nil, "")
	if err != nil {
		log.Println("create:", err)
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println("start:", err)
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	data, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		log.Println(err)
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	out, errors := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	return out, errors, data
}

func startFfmpegDecrypt(url, filename, key, path string, source string) (<-chan container.ContainerWaitOKBody, <-chan error, io.Reader) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	_, err = cli.ImagePull(ctx, viper.GetString("converter.image"), types.ImagePullOptions{})
	if err != nil {
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}
	// io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: viper.GetString("converter.image"),

			Entrypoint: []string{"ffmpeg",
				"-decryption_key", key,
				"-i", url,
				"-f", "mp3",
				filename,
			},
			WorkingDir: path,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: source,
					Target: path,
				},
			},
			Privileged: true,
			AutoRemove: true,
		}, nil, nil, "")
	if err != nil {
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	data, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		errorc := make(chan error)
		errorc <- err
		return make(chan container.ContainerWaitOKBody, 0),
			make(chan error), strings.NewReader("")
	}

	out, errors := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	return out, errors, data
}

func ffmpegTest() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, viper.GetString("converter.image"), types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	_, err = cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: viper.GetString("converter.image"),

			Entrypoint: []string{"ffmpeg",
				"-L",
			},
		},
		&container.HostConfig{
			AutoRemove: true,
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}
}
