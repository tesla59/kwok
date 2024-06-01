package file_test

import (
	"errors"
	"os"
	"path/filepath"
	"sigs.k8s.io/kwok/pkg/utils/file"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("Create a file in test directory", func(t *testing.T) {
		// Create the file
		testFilePath := "../../../test/data"
		err := file.Create(testFilePath)
		if err != nil {
			t.Errorf("Failed to create file: %v", err)
		}

		// Check if the file is created
		if !file.Exists(testFilePath) {
			t.Errorf("File not created: %v", err)
		}

		// Clean up: Remove the test file
		err = file.Remove(testFilePath)
		if err != nil {
			t.Errorf("Failed to remove test file: %v", err)
		}
	})

	t.Run("Create a file in non-existing directory", func(t *testing.T) {
		nonExistentDir := "/non/existent/dir/test.txt"
		err := file.Create(nonExistentDir)
		if err == nil {
			t.Errorf("Expected error when creating file in non-existent directory, but got nil")
		} else {
			expectedErr := os.ErrNotExist
			if !errors.Is(err, expectedErr) {
				t.Errorf("Expected error when creating file in non-existent directory, but got %v", err)
			}
		}
	})
}

func TestCopy(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testcopy")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func(path string) {
		err := file.RemoveAll(path)
		if err != nil {
			t.Errorf("Failed to remove temporary directory: %v", err)
		}
	}(tempDir)

	t.Run("Copy file", func(t *testing.T) {
		oldPath := filepath.Join(tempDir, "oldfile.txt")
		newPath := filepath.Join(tempDir, "newfile.txt")

		err = file.Write(oldPath, []byte("Hello, World!"))
		if err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		err = file.Copy(oldPath, newPath)
		if err != nil {
			t.Errorf("Copy failed: %v", err)
		}

		// Check if the new file exists and has the same content as the old file
		newContent, err := file.Read(newPath)
		if err != nil {
			t.Errorf("Failed to read new file: %v", err)
		}
		if string(newContent) != "Hello, World!" {
			t.Errorf("New file content doesn't match the old file content")
		}
	})

	t.Run("Copy non-existing file", func(t *testing.T) {
		oldPath := filepath.Join(tempDir, "nonexistent.txt")
		newPath := filepath.Join(tempDir, "newfile.txt")

		err = file.Copy(oldPath, newPath)
		if err == nil {
			t.Errorf("Expected error when copying non-existing file, but got nil")
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Unexpected error when copying non-existing file: %v", err)
		}
	})

	t.Run("Copy to existing file", func(t *testing.T) {
		oldPath := filepath.Join(tempDir, "oldfile.txt")
		newPath := filepath.Join(tempDir, "newfile.txt")

		err = file.WriteWithMode(oldPath, []byte("Hello, World!"), 0644)
		if err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		err = file.WriteWithMode(newPath, []byte("Existing content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create new file: %v", err)
		}

		err = file.Copy(oldPath, newPath)
		if err != nil {
			t.Errorf("Copy failed: %v", err)
		}

		// Check if the new file content was overwritten by the old file content
		newContent, err := file.Read(newPath)
		if err != nil {
			t.Errorf("Failed to read new file: %v", err)
		}
		if string(newContent) != "Hello, World!" {
			t.Errorf("New file content doesn't match the old file content")
		}
	})
}

func TestAppend(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testappend")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Errorf("Failed to remove temporary directory: %v", err)
		}
	}(tempDir)

	t.Run("Append to new file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "newfile.txt")
		content := []byte("Hello, World!")

		err = file.Append(filePath, content)
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}

		// Check if the file was created and has the correct content
		fileContent, err := file.Read(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		if string(fileContent) != "Hello, World!" {
			t.Errorf("File content doesn't match the appended content")
		}
	})

	t.Run("Append to existing file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "existingfile.txt")
		initialContent := []byte("Initial content ")
		appendedContent := []byte("Appended content")

		err = file.WriteWithMode(filePath, initialContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		err = file.Append(filePath, appendedContent)
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}

		// Check if the file content was appended correctly
		fileContent, err := file.Read(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}
		expectedContent := string(initialContent) + string(appendedContent)
		if string(fileContent) != expectedContent {
			t.Errorf("File content doesn't match the expected content")
		}
	})

	t.Run("Append to non-existing directory", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non-existing-dir", "file.txt")
		content := []byte("Hello, World!")

		err = file.Append(filePath, content)
		if err == nil {
			t.Errorf("Expected error when appending to non-existing directory, but got nil")
		} else if !os.IsNotExist(err) {
			t.Errorf("Unexpected error when appending to non-existing directory: %v", err)
		}
	})
}
