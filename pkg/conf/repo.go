package conf

import (
	"fmt"
	"net/url"

	"arhat.dev/pkg/tlshelper"
	"go.uber.org/multierr"
)

type RepoSpec struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
	Auth struct {
		HTTPBasic struct {
			Username string `json:"username" yaml:"username"`
			Password string `json:"password" yaml:"password"`
		} `json:"httpBasic" yaml:"httpBasic"`
	} `json:"auth" yaml:"auth"`

	TLS tlshelper.TLSConfig `json:"tls" yaml:"tls"`
}

func (r RepoSpec) Validate() error {
	var err error
	if r.Name == "" {
		err = multierr.Append(err, fmt.Errorf("invalid helm repo with empty name"))
	}

	if r.URL == "" {
		err = multierr.Append(err, fmt.Errorf("invalid helm repo with no url"))
	} else {
		var u *url.URL
		u, err = url.Parse(r.URL)
		if err != nil {
			err = multierr.Append(err, fmt.Errorf("invalid repo url %q: %w", r.URL, err))
		}

		switch u.Scheme {
		case "http", "https":
		default:
			err = multierr.Append(err, fmt.Errorf("invalid url scheme %q, only http/https supported", u.Scheme))
		}
	}

	return err
}
