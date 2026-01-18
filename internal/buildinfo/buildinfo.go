package buildinfo

// Ces variables sont typiquement injectées à la compilation via -ldflags.
// Exemple :
//
//	-X github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo.Version=v0.0.0
//	-X github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo.Commit=abcdef
//	-X github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/buildinfo.Date=2026-01-18
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
	Date    string `json:"date,omitempty"`
}

func Current() Info {
	return Info{Version: Version, Commit: Commit, Date: Date}
}
