// Package sftp implements an SFTP server to serve an rclone VFS

// +build !plan9

package sftp

import (
	"github.com/ncw/rclone/cmd"
	"github.com/ncw/rclone/fs/config/flags"
	"github.com/ncw/rclone/fs/rc"
	"github.com/ncw/rclone/vfs"
	"github.com/ncw/rclone/vfs/vfsflags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Options contains options for the http Server
type Options struct {
	ListenAddr     string // Port to listen on
	Key            string // Path to private key
	AuthorizedKeys string // Path to authorized keys file
	User           string // single username
	Pass           string // password for user
	NoAuth         bool   // allow no authentication on connections
}

// DefaultOpt is the default values used for Options
var DefaultOpt = Options{
	ListenAddr:     "localhost:2022",
	AuthorizedKeys: "~/.ssh/authorized_keys",
}

// Opt is options set by command line flags
var Opt = DefaultOpt

// AddFlags adds flags for the sftp
func AddFlags(flagSet *pflag.FlagSet, Opt *Options) {
	rc.AddOption("sftp", &Opt)
	flags.StringVarP(flagSet, &Opt.ListenAddr, "addr", "", Opt.ListenAddr, "IPaddress:Port or :Port to bind server to.")
	flags.StringVarP(flagSet, &Opt.Key, "key", "", Opt.Key, "SSH private key file (leave blank to auto generate)")
	flags.StringVarP(flagSet, &Opt.AuthorizedKeys, "authorized-keys", "", Opt.AuthorizedKeys, "Authorized keys file")
	flags.StringVarP(flagSet, &Opt.User, "user", "", Opt.User, "User name for authentication.")
	flags.StringVarP(flagSet, &Opt.Pass, "pass", "", Opt.Pass, "Password for authentication.")
	flags.BoolVarP(flagSet, &Opt.NoAuth, "no-auth", "", Opt.NoAuth, "Allow connections with no authentication if set.")
}

func init() {
	vfsflags.AddFlags(Command.Flags())
	AddFlags(Command.Flags(), &Opt)
}

// Command definition for cobra
var Command = &cobra.Command{
	Use:   "sftp remote:path",
	Short: `Serve the remote over SFTP.`,
	Long: `rclone serve sftp implements an SFTP server to serve the remote
over SFTP.  This can be used with an SFTP client or you can make a
remote of type sftp to use with it.

You can use the filter flags (eg --include, --exclude) to control what
is served.

The server will log errors.  Use -v to see access logs.

--bwlimit will be respected for file transfers.  Use --stats to
control the stats printing.

You must provide some means of authentication, either with --user/--pass,
an authorized keys file (specify location with --authorized-keys - the
default is the same as ssh) or set the --no-auth flag for no
authentication when logging in.

Note that this also implements a small number of shell commands so
that it can provide md5sum/sha1sum/df information for the sftp
backend.

If you don't supply a --key then rclone will generate one and cache it
for later use.

By default the server binds to localhost:2022 - if you want it to be
reachable externally then supply "--addr :2022" for example.

Note that the default of "--vfs-cache-mode off" is fine for the rclone
sftp backend, but it may not be with other SFTP clients.

` + vfs.Help,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(1, 1, command, args)
		f := cmd.NewFsSrc(args)
		cmd.Run(false, true, command, func() error {
			s := newServer(f, &Opt)
			err := s.Serve()
			if err != nil {
				return err
			}
			s.Wait()
			return nil
		})
	},
}
