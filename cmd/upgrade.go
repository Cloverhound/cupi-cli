package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const upgradeRepo = "Cloverhound/cupi-cli"

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cupi to the latest release from GitHub",
	Long: `Check GitHub for the latest cupi release and upgrade the binary in-place if a newer version is available.

The upgrade replaces the currently running binary with the latest release for your
OS and architecture (darwin/linux/windows, amd64/arm64).

Examples:
  cupi upgrade
  cupi upgrade --check   # Print latest version without installing`,
	RunE: runUpgrade,
}

var upgradeCheckOnly bool

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeCheckOnly, "check", false, "Check for updates without installing")
	rootCmd.AddCommand(upgradeCmd)
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	fmt.Printf("Current version : cupi %s\n", Version)
	fmt.Printf("Checking         : https://github.com/%s/releases/latest\n", upgradeRepo)

	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}

	latestTag := release.TagName
	fmt.Printf("Latest version  : cupi %s\n", latestTag)

	if !isNewer(Version, latestTag) {
		fmt.Println("Already up to date.")
		return nil
	}

	if upgradeCheckOnly {
		fmt.Printf("Run 'cupi upgrade' to install %s.\n", latestTag)
		return nil
	}

	// Determine target asset name
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goarch == "aarch64" {
		goarch = "arm64"
	}

	version := strings.TrimPrefix(latestTag, "v")
	var assetName string
	if goos == "windows" {
		assetName = fmt.Sprintf("cupi_%s_%s_%s.zip", version, goos, goarch)
	} else {
		assetName = fmt.Sprintf("cupi_%s_%s_%s.tar.gz", version, goos, goarch)
	}

	downloadURL := ""
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no release asset found for %s/%s (looked for %q)\nCheck %s manually",
			goos, goarch, assetName, fmt.Sprintf("https://github.com/%s/releases/tag/%s", upgradeRepo, latestTag))
	}

	// Find the current binary path
	selfPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine current binary path: %w", err)
	}
	selfPath, err = filepath.EvalSymlinks(selfPath)
	if err != nil {
		return fmt.Errorf("could not resolve binary symlink: %w", err)
	}

	fmt.Printf("Downloading     : %s\n", assetName)

	tmpDir, err := os.MkdirTemp("", "cupi-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Extract the binary
	binaryName := "cupi"
	if goos == "windows" {
		binaryName = "cupi.exe"
	}
	newBinaryPath := filepath.Join(tmpDir, binaryName)

	if strings.HasSuffix(assetName, ".tar.gz") {
		if err := extractTarGz(archivePath, binaryName, newBinaryPath); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	} else {
		if err := extractZip(archivePath, binaryName, newBinaryPath); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	}

	// Atomic replace: write to a temp file next to the current binary, then rename
	tmpBinary := selfPath + ".new"
	if err := copyFile(newBinaryPath, tmpBinary); err != nil {
		return fmt.Errorf("failed to stage new binary: %w", err)
	}
	if err := os.Chmod(tmpBinary, 0755); err != nil {
		os.Remove(tmpBinary)
		return fmt.Errorf("failed to set binary permissions: %w", err)
	}
	if err := os.Rename(tmpBinary, selfPath); err != nil {
		os.Remove(tmpBinary)
		return fmt.Errorf("failed to replace binary (try with sudo?): %w", err)
	}

	fmt.Printf("Upgraded        : cupi %s → cupi %s\n", Version, latestTag)
	fmt.Printf("Binary          : %s\n", selfPath)
	return nil
}

func fetchLatestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", upgradeRepo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", fmt.Sprintf("cupi-cli/%s", Version))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release JSON: %w", err)
	}
	if release.TagName == "" {
		return nil, fmt.Errorf("no releases found for %s", upgradeRepo)
	}
	return &release, nil
}

// isNewer returns true if latestTag is a newer version than current.
// Handles "dev" builds (always considered outdated) and "vX.Y.Z" semver tags.
func isNewer(current, latest string) bool {
	if current == "dev" || current == "" {
		return true
	}
	// Normalize: strip leading "v"
	c := strings.TrimPrefix(current, "v")
	l := strings.TrimPrefix(latest, "v")
	return l != c && versionGreater(l, c)
}

// versionGreater returns true if a > b using simple semver comparison.
func versionGreater(a, b string) bool {
	aParts := strings.SplitN(a, ".", 3)
	bParts := strings.SplitN(b, ".", 3)
	for len(aParts) < 3 {
		aParts = append(aParts, "0")
	}
	for len(bParts) < 3 {
		bParts = append(bParts, "0")
	}
	for i := 0; i < 3; i++ {
		av := versionNum(aParts[i])
		bv := versionNum(bParts[i])
		if av > bv {
			return true
		}
		if av < bv {
			return false
		}
	}
	return false
}

func versionNum(s string) int {
	// Strip any pre-release suffix like "-rc1"
	if idx := strings.IndexAny(s, "-+"); idx >= 0 {
		s = s[:idx]
	}
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d downloading %s", resp.StatusCode, url)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func extractTarGz(archivePath, binaryName, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) == binaryName {
			out, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, tr) //nolint:gosec
			return err
		}
	}
	return fmt.Errorf("%q not found in archive", binaryName)
}

func extractZip(archivePath, binaryName, dest string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			out, err := os.Create(dest)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, rc) //nolint:gosec
			return err
		}
	}
	return fmt.Errorf("%q not found in zip archive", binaryName)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
