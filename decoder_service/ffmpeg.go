package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func startFfmpegEncrypt(url, filename, key, kid, path string) (<-chan container.ContainerWaitOKBody, <-chan error, io.Reader) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "jrottenberg/ffmpeg", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	// currPath, _ := os.Getwd()
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "jrottenberg/ffmpeg",

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
					Source: path,
					Target: path,
				},
			},
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	data, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		panic(err)
	}

	out, errors := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	return out, errors, data
}

func startFfmpegDecrypt(url, filename, key, path string) (<-chan container.ContainerWaitOKBody, <-chan error, io.Reader) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "jrottenberg/ffmpeg", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	// currPath, _ := os.Getwd()
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "jrottenberg/ffmpeg",

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
					Source: path,
					Target: path,
				},
			},
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	data, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		panic(err)
	}

	out, errors := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	return out, errors, data
}
