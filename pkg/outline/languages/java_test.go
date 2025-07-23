package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

func TestJavaOutlineWithImports(t *testing.T) {
	javaCode := `package com.example.demo;

import java.util.List;
import java.util.ArrayList;
import java.io.IOException;

/**
 * Demo class for testing
 */
public class Demo {
    private String name;
    private int age;

    /**
     * Constructor for Demo
     */
    public Demo(String name, int age) {
        this.name = name;
        this.age = age;
    }

    /**
     * Gets the name
     */
    public String getName() {
        return name;
    }

    /**
     * Processes data
     */
    public void processData() throws IOException {
        // Implementation
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(java.Language())); err != nil {
		t.Fatalf("Failed to set Java language: %v", err)
	}

	tree := parser.Parse([]byte(javaCode), nil)
	defer tree.Close()

	result := ExtractJavaOutline(tree.RootNode(), []byte(javaCode))

	// Check that package is included
	if !strings.Contains(result, "package com.example.demo;") {
		t.Error("Expected package declaration to be included")
	}

	// Check that imports are included
	if !strings.Contains(result, "import java.util.List;") {
		t.Error("Expected import declarations to be included")
	}

	// Check that class is included
	if !strings.Contains(result, "public class Demo") {
		t.Error("Expected class declaration to be included")
	}

	// Check that constructor is included
	if !strings.Contains(result, "Demo(String name, int age)") {
		t.Error("Expected constructor to be included")
	}

	// Check that methods are included
	if !strings.Contains(result, "public String getName()") {
		t.Error("Expected getter method to be included")
	}

	// Check that throws clause is included
	if !strings.Contains(result, "throws IOException") {
		t.Error("Expected throws clause to be included")
	}

	// Check that fields are included
	if !strings.Contains(result, "private String name;") {
		t.Error("Expected field declaration to be included")
	}

	t.Logf("Java outline result:\n%s", result)
}

func TestJavaInterface(t *testing.T) {
	javaCode := `package com.example;

/**
 * Repository interface
 */
public interface UserRepository extends BaseRepository<User> {
    /**
     * Finds user by email
     */
    User findByEmail(String email);

    /**
     * Finds users by age range
     */
    List<User> findByAgeRange(int minAge, int maxAge) throws DataAccessException;
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(java.Language())); err != nil {
		t.Fatalf("Failed to set Java language: %v", err)
	}

	tree := parser.Parse([]byte(javaCode), nil)
	defer tree.Close()

	result := ExtractJavaOutline(tree.RootNode(), []byte(javaCode))

	// Check that interface is included
	if !strings.Contains(result, "public interface UserRepository") {
		t.Error("Expected interface declaration to be included")
	}

	// Check that extends clause is included
	if !strings.Contains(result, "extends BaseRepository<User>") {
		t.Error("Expected extends clause to be included")
	}

	// Check that interface methods are included
	if !strings.Contains(result, "User findByEmail(String email)") {
		t.Error("Expected interface method to be included")
	}

	// Check that generic types are preserved
	if !strings.Contains(result, "List<User>") {
		t.Error("Expected generic type to be preserved")
	}

	t.Logf("Java interface outline result:\n%s", result)
}

func TestJavaEnum(t *testing.T) {
	javaCode := `package com.example;

/**
 * Status enumeration
 */
public enum Status {
    ACTIVE,
    INACTIVE,
    PENDING;

    /**
     * Gets display name
     */
    public String getDisplayName() {
        return name().toLowerCase();
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(java.Language())); err != nil {
		t.Fatalf("Failed to set Java language: %v", err)
	}

	tree := parser.Parse([]byte(javaCode), nil)
	defer tree.Close()

	result := ExtractJavaOutline(tree.RootNode(), []byte(javaCode))

	// Check that enum is included
	if !strings.Contains(result, "public enum Status") {
		t.Error("Expected enum declaration to be included")
	}

	// Check that enum constants are included
	if !strings.Contains(result, "ACTIVE,") {
		t.Error("Expected enum constant to be included")
	}

	// Check that enum methods are included
	if !strings.Contains(result, "public String getDisplayName()") {
		t.Error("Expected enum method to be included")
	}

	t.Logf("Java enum outline result:\n%s", result)
}

func TestJavaComplexClass(t *testing.T) {
	javaCode := `package com.example.service;

import java.util.*;
import javax.annotation.Nullable;

/**
 * User service implementation
 */
@Service
public class UserService extends BaseService implements UserManager {
    private static final String DEFAULT_ROLE = "USER";
    private final UserRepository repository;
    private Map<String, Object> cache;

    /**
     * Constructor with dependency injection
     */
    @Autowired
    public UserService(UserRepository repository) {
        this.repository = repository;
        this.cache = new HashMap<>();
    }

    /**
     * Creates a new user
     */
    @Override
    public User createUser(@Nullable String name, int age) throws ValidationException {
        // Implementation
        return new User(name, age);
    }

    /**
     * Inner utility class
     */
    private static class UserValidator {
        public boolean isValid(User user) {
            return user != null && user.getName() != null;
        }
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(java.Language())); err != nil {
		t.Fatalf("Failed to set Java language: %v", err)
	}

	tree := parser.Parse([]byte(javaCode), nil)
	defer tree.Close()

	result := ExtractJavaOutline(tree.RootNode(), []byte(javaCode))

	// Check that class with modifiers is included
	if !strings.Contains(result, "public class UserService") {
		t.Error("Expected class declaration to be included")
	}

	// Check that extends and implements are included
	if !strings.Contains(result, "extends BaseService") {
		t.Error("Expected extends clause to be included")
	}
	if !strings.Contains(result, "implements UserManager") {
		t.Error("Expected implements clause to be included")
	}

	// Check that constants are included
	if !strings.Contains(result, "private static final String DEFAULT_ROLE") {
		t.Error("Expected constant field to be included")
	}

	// Check that constructor is included
	if !strings.Contains(result, "UserService(UserRepository repository)") {
		t.Error("Expected constructor to be included")
	}

	// Check that method with annotations is included
	if !strings.Contains(result, "@Override") || !strings.Contains(result, "public User createUser") {
		t.Error("Expected annotated method to be included")
	}

	// Check that inner class is included
	if !strings.Contains(result, "private static class UserValidator") {
		t.Error("Expected inner class to be included")
	}

	t.Logf("Java complex class outline result:\n%s", result)
}

func TestJavaAbstractClass(t *testing.T) {
	javaCode := `package com.example;

/**
 * Abstract base class
 */
public abstract class Animal {
    protected String name;

    /**
     * Abstract method
     */
    public abstract void makeSound();

    /**
     * Concrete method
     */
    public final String getName() {
        return name;
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(java.Language())); err != nil {
		t.Fatalf("Failed to set Java language: %v", err)
	}

	tree := parser.Parse([]byte(javaCode), nil)
	defer tree.Close()

	result := ExtractJavaOutline(tree.RootNode(), []byte(javaCode))

	// Check that abstract class is included
	if !strings.Contains(result, "public abstract class Animal") {
		t.Error("Expected abstract class declaration to be included")
	}

	// Check that abstract method is included
	if !strings.Contains(result, "public abstract void makeSound()") {
		t.Error("Expected abstract method to be included")
	}

	// Check that final method is included
	if !strings.Contains(result, "public final String getName()") {
		t.Error("Expected final method to be included")
	}

	t.Logf("Java abstract class outline result:\n%s", result)
}