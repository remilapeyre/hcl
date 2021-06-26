package gohcl

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
)

type Service struct {
	Name string   `hcl:"name,label"`
	Exe  []string `hcl:"executable"`
}
type Constraints struct {
	OS   string `hcl:"os"`
	Arch string `hcl:"arch"`
}
type App struct {
	Name        string                 `hcl:"name"`
	Desc        string                 `hcl:"description"`
	Constraints *Constraints           `hcl:"constraints,block"`
	Services    []Service              `hcl:"service,block"`
	Meta        map[string]string      `hcl:"meta,block"`
	Config      map[string]interface{} `hcl:"config,block"`
}

var testCases = map[string]struct {
	arg  interface{}
	want string
}{
	"map-blocks": {
		arg: App{
			Name: "awesome-app",
			Desc: "Such an awesome application",
			Constraints: &Constraints{
				OS:   "linux",
				Arch: "amd64",
			},
			Meta: map[string]string{
				"hello": "world",
			},
			Config: map[string]interface{}{
				"command": "/bin/sleep",
				"args":    []interface{}{"1"},
			},
			Services: []Service{
				{
					Name: "web",
					Exe:  []string{"./web", "--listen=:8080"},
				},
				{
					Name: "worker",
					Exe:  []string{"./worker"},
				},
			},
		},
		want: `name="awesome-app"
description="Such an awesome application"

constraints{
os="linux"
arch="amd64"
}

service"web"{
executable=["./web", "--listen=:8080"]
}
service"worker"{
executable=["./worker"]
}

meta{
hello="world"
}

config{
args=["1"]
command="/bin/sleep"
}
`,
	},
}

func TestEncodeAsBlock(t *testing.T) {
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			got := EncodeAsBlock(tt.arg, "app")
			require.Equal(t, tt.want, string(got.Body().BuildTokens(nil).Bytes()))
		})
	}
}

func TestEncodeIntoBody(t *testing.T) {
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			f := hclwrite.NewEmptyFile()
			EncodeIntoBody(tt.arg, f.Body())
			expected, diags := hclwrite.ParseConfig([]byte(tt.want), "test", hcl.InitialPos)
			require.Nil(t, diags)
			require.Equal(t, string(expected.Bytes()), string(f.Bytes()))
		})
	}
}
