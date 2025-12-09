package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeRelative(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		absPath     string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:        "simple subdirectory",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/me/project/dna",
			want:        "dna",
			wantErr:     false,
		},
		{
			name:        "nested subdirectory",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/me/project/shared/dna",
			want:        "shared/dna",
			wantErr:     false,
		},
		{
			name:        "same path as project root",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/me/project",
			want:        ".",
			wantErr:     false,
		},
		{
			name:        "outside project - parent directory",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/me/dna",
			wantErr:     true,
			errContains: "outside project root",
		},
		{
			name:        "outside project - different tree",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/other/dna",
			wantErr:     true,
			errContains: "outside project root",
		},
		{
			name:        "path with trailing slash",
			projectRoot: "/Users/me/project/",
			absPath:     "/Users/me/project/dna/",
			want:        "dna",
			wantErr:     false,
		},
		{
			name:        "path with dot segments cleaned",
			projectRoot: "/Users/me/project",
			absPath:     "/Users/me/project/./dna",
			want:        "dna",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeRelative(tt.projectRoot, tt.absPath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestResolveRelative(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		relPath     string
		wantSuffix  string // Check path ends with this
		wantErr     bool
		errContains string
	}{
		{
			name:        "simple relative path",
			projectRoot: "/Users/me/project",
			relPath:     "dna",
			wantSuffix:  filepath.Join("project", "dna"),
			wantErr:     false,
		},
		{
			name:        "nested relative path",
			projectRoot: "/Users/me/project",
			relPath:     "shared/dna",
			wantSuffix:  filepath.Join("project", "shared", "dna"),
			wantErr:     false,
		},
		{
			name:        "dot current directory",
			projectRoot: "/Users/me/project",
			relPath:     ".",
			wantSuffix:  "project",
			wantErr:     false,
		},
		{
			name:        "reject absolute path",
			projectRoot: "/Users/me/project",
			relPath:     "/Users/me/other",
			wantErr:     true,
			errContains: "expected relative path",
		},
		{
			name:        "path with ./prefix",
			projectRoot: "/Users/me/project",
			relPath:     "./dna",
			wantSuffix:  filepath.Join("project", "dna"),
			wantErr:     false,
		},
		{
			name:        "path trying to escape with ..",
			projectRoot: "/Users/me/project",
			relPath:     "../outside",
			wantErr:     true,
			errContains: "escapes project root",
		},
		{
			name:        "deeply nested escape attempt",
			projectRoot: "/Users/me/project",
			relPath:     "foo/../../outside",
			wantErr:     true,
			errContains: "escapes project root",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRelative(tt.projectRoot, tt.relPath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.True(t, filepath.IsAbs(got), "result should be absolute path")
				if tt.wantSuffix != "" {
					assert.True(t,
						strings.HasPrefix(got, tt.projectRoot) || strings.HasSuffix(got, tt.wantSuffix),
						"path %s should end with %s", got, tt.wantSuffix)
				}
			}
		})
	}
}

func TestValidateLocalPath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid relative path",
			projectRoot: "/Users/me/project",
			path:        "dna",
			wantErr:     false,
		},
		{
			name:        "valid nested relative path",
			projectRoot: "/Users/me/project",
			path:        "shared/dna",
			wantErr:     false,
		},
		{
			name:        "valid absolute path within project",
			projectRoot: "/Users/me/project",
			path:        "/Users/me/project/dna",
			wantErr:     false,
		},
		{
			name:        "invalid - escapes project",
			projectRoot: "/Users/me/project",
			path:        "../outside",
			wantErr:     true,
			errContains: "project",
		},
		{
			name:        "invalid - absolute path outside project",
			projectRoot: "/Users/me/project",
			path:        "/Users/other/dna",
			wantErr:     true,
			errContains: "outside project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLocalPath(tt.projectRoot, tt.path)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIsWithinProject(t *testing.T) {
	// Create a temp directory for real filesystem tests
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name        string
		projectRoot string
		path        string
		want        bool
		wantErr     bool
	}{
		{
			name:        "path within project",
			projectRoot: tmpDir,
			path:        subDir,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "path equals project root",
			projectRoot: tmpDir,
			path:        tmpDir,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "non-existent path",
			projectRoot: tmpDir,
			path:        filepath.Join(tmpDir, "nonexistent"),
			want:        false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsWithinProject(tt.projectRoot, tt.path)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIsWithinPath(t *testing.T) {
	tests := []struct {
		name       string
		parentPath string
		childPath  string
		want       bool
	}{
		{
			name:       "exact match",
			parentPath: "/Users/me/project",
			childPath:  "/Users/me/project",
			want:       true,
		},
		{
			name:       "child is subdirectory",
			parentPath: "/Users/me/project",
			childPath:  "/Users/me/project/dna",
			want:       true,
		},
		{
			name:       "child is nested subdirectory",
			parentPath: "/Users/me/project",
			childPath:  "/Users/me/project/shared/dna",
			want:       true,
		},
		{
			name:       "child is outside - sibling",
			parentPath: "/Users/me/project",
			childPath:  "/Users/me/other",
			want:       false,
		},
		{
			name:       "child is outside - parent",
			parentPath: "/Users/me/project",
			childPath:  "/Users/me",
			want:       false,
		},
		{
			name:       "child is similar prefix but not within",
			parentPath: "/Users/me/proj",
			childPath:  "/Users/me/project",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWithinPath(tt.parentPath, tt.childPath)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestWithSymlinks tests behavior with actual symlinks
func TestWithSymlinks(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create actual directories
	projectDir := filepath.Join(tmpDir, "project")
	require.NoError(t, os.Mkdir(projectDir, 0755))

	insideDir := filepath.Join(projectDir, "inside")
	require.NoError(t, os.Mkdir(insideDir, 0755))

	outsideDir := filepath.Join(tmpDir, "outside")
	require.NoError(t, os.Mkdir(outsideDir, 0755))

	// Create symlink inside project pointing to inside directory
	symlinkInside := filepath.Join(projectDir, "link-inside")
	require.NoError(t, os.Symlink(insideDir, symlinkInside))

	// Create symlink inside project pointing to outside directory
	symlinkOutside := filepath.Join(projectDir, "link-outside")
	require.NoError(t, os.Symlink(outsideDir, symlinkOutside))

	t.Run("symlink to inside directory - should be valid", func(t *testing.T) {
		err := ValidateLocalPath(projectDir, symlinkInside)
		assert.NoError(t, err)
	})

	t.Run("symlink to outside directory - should be invalid", func(t *testing.T) {
		err := ValidateLocalPath(projectDir, symlinkOutside)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "outside project")
	})

	t.Run("IsWithinProject with symlink to inside", func(t *testing.T) {
		within, err := IsWithinProject(projectDir, symlinkInside)
		require.NoError(t, err)
		assert.True(t, within)
	})

	t.Run("MakeRelative with symlink inside project", func(t *testing.T) {
		relPath, err := MakeRelative(projectDir, symlinkInside)
		require.NoError(t, err)
		// Should resolve to the actual target path
		assert.NotEmpty(t, relPath)
	})
}
