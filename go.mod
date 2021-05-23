module github.com/noyako/Audio-Decoder

go 1.15

replace github.com/noyako/swolf => /home/aoyako/projects/swolf

require (
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/docker/docker v20.10.6+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/lib/pq v1.10.2
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/noyako/swolf v0.0.0-20210523170027-b57e908efb52
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)
