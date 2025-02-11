package manifest

import (
	"github.com/osbuild/osbuild-composer/internal/artifact"
	"github.com/osbuild/osbuild-composer/internal/osbuild"
)

// A Tar represents the contents of another pipeline in a tar file
type Tar struct {
	Base
	Filename string

	inputPipeline *Base
}

// NewTar creates a new TarPipeline. The inputPipeline represents the
// filesystem tree which will be the contents of the tar file. The pipelinename
// is the name of the pipeline. The filename is the name of the output tar file.
func NewTar(m *Manifest,
	buildPipeline *Build,
	inputPipeline *Base,
	pipelinename string) *Tar {
	p := &Tar{
		Base:          NewBase(m, pipelinename, buildPipeline),
		inputPipeline: inputPipeline,
		Filename:      "image.tar",
	}
	if inputPipeline.manifest != m {
		panic("tree pipeline from different manifest")
	}
	buildPipeline.addDependent(p)
	m.addPipeline(p)
	return p
}

func (p *Tar) serialize() osbuild.Pipeline {
	pipeline := p.Base.serialize()

	tarStage := osbuild.NewTarStage(&osbuild.TarStageOptions{Filename: p.Filename}, p.inputPipeline.Name())
	pipeline.AddStage(tarStage)

	return pipeline
}

func (p *Tar) getBuildPackages() []string {
	return []string{"tar"}
}

func (p *Tar) Export() *artifact.Artifact {
	p.Base.export = true
	mimeType := "application/x-tar"
	return artifact.New(p.Name(), p.Filename, &mimeType)
}
