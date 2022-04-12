package cache

import (
	"context"
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	defaultOnce  sync.Once
	defaultCache *diskCacheManager
)

type diskCacheManager struct {
	cacheDir string
}

func NewDiskCacheManager(cacheDir string) (*diskCacheManager, error) {
	defaultOnce.Do(func() {
		initDefaultCache(cacheDir)
	})
	return defaultCache, defaultDirErr
}

// initDefaultCache does the work of finding the default cache
// the first time Default is called.
func initDefaultCache(cacheDir string) {
	dir := DefaultDir(cacheDir)
	if err := os.MkdirAll(dir, 0744); err != nil {
		log.Fatalf("failed to initialize build cache at %s: %s\n", dir, err)
	}

	c, err := Open(dir)
	if err != nil {
		log.Fatalf("failed to initialize build cache at %s: %s\n", dir, err)
	}
	defaultCache = c
}

var (
	defaultDirOnce sync.Once
	defaultDir     string
	defaultDirErr  error
)

// DefaultDir returns the effective GOLANGCI_LINT_CACHE setting.
func DefaultDir(cacheDir string) string {
	defaultDirOnce.Do(func() {
		if cacheDir != "" {
			if filepath.IsAbs(cacheDir) {
				return
			}

			defaultDirErr = fmt.Errorf("%s is not an absolute path", cacheDir)
			return
		}

		// Compute default location.
		dir, err := os.UserCacheDir()
		if err != nil {
			defaultDirErr = fmt.Errorf("cacheDir is not defined and %v", err)
			return
		}
		defaultDir = filepath.Join(dir, "bazel-differ")
	})

	return defaultDir
}

// Open opens and returns the cache in the given directory.
//
// It is safe for multiple processes on a single machine to use the
// same cache directory in a local file system simultaneously.
// They will coordinate using operating system file locks and may
// duplicate effort but will not corrupt the cache.
//
// However, it is NOT safe for multiple processes on different machines
// to share a cache directory (for example, if the directory were stored
// in a network file system). File locking is notoriously unreliable in
// network file systems and may not suffice to protect the cache.
//
func Open(dir string) (*diskCacheManager, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "open", Path: dir, Err: fmt.Errorf("not a directory")}
	}
	for i := 0; i < 256; i++ {
		name := filepath.Join(dir, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(name, 0744); err != nil {
			return nil, err
		}
	}
	c := &diskCacheManager{
		cacheDir: dir,
	}
	return c, nil
}

func (d diskCacheManager) Put(ctx context.Context, key string, value map[string]string) error {
	_, err := internal.WriteHashFile(d.getFilename(key), value)
	return err
}

func (d diskCacheManager) Get(ctx context.Context, key string) (map[string]string, error) {
	return internal.ReadHashFile(d.getFilename(key))
}

func (d diskCacheManager) getFilename(key string) string {
	return filepath.Join(d.cacheDir, key[0:2], key, "hashes.json")
}
